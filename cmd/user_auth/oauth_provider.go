// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	"github.com/golang-jwt/jwt/v5"
)

// oauthHTTPClient is a dedicated HTTP client for OAuth operations with
// explicit timeout and TLS configuration, replacing http.DefaultClient.
var oauthHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

// parseJWT parses and validates a JWT token signed with HS256.
// Returns the claims if valid, or an error if parsing/validation fails.
func parseJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid JWT token")
}

// OAuthProvider handles the DingTalk OAuth 2.0 authorization code flow.
type OAuthProvider struct {
	configDir  string
	clientID   string
	logger     *slog.Logger
	Output     io.Writer
	httpClient *http.Client
}

type openDevCreateAppRespData struct {
	AppID             string `json:"app_id"`
	AppKey            string `json:"app_key"`
	AppSecret         string `json:"app_secret"`
	SSOSecret         string `json:"sso_secret"`
	CallbackToken     string `json:"callback_token"`
	CallbackSecretKey string `json:"callback_secret_key"`
}

type openDevCreateAppResp struct {
	Code      string                   `json:"code"`
	Status    string                   `json:"status"`
	Message   string                   `json:"message"`
	Data      openDevCreateAppRespData `json:"data"`
	RequestID string                   `json:"request_id"`
}

// NewOAuthProvider creates a new OAuth provider.
func NewOAuthProvider(configDir string, logger *slog.Logger) *OAuthProvider {
	return &OAuthProvider{
		configDir:  configDir,
		clientID:   ClientID(),
		logger:     logger,
		Output:     os.Stderr,
		httpClient: oauthHTTPClient,
	}
}

// resetCredentialState clears any stale credential state inherited from
// previous login methods so that OAuth flow always starts fresh by
// fetching clientID from MCP.
func (p *OAuthProvider) resetCredentialState() {
	p.clientID = ""
	clientMu.Lock()
	clientIDFromMCP = false
	clientMu.Unlock()
}

func (p *OAuthProvider) output() io.Writer {
	if p != nil && p.Output != nil {
		return p.Output
	}
	return io.Discard
}

// Login performs authentication with smart degradation:
// 1. If force=false, try silent token refresh first (refresh_token)
// 2. If all silent methods fail (or force=true), fall back to browser OAuth flow
func (p *OAuthProvider) Login(ctx context.Context, force bool) (*TokenData, error) {
	// Smart degradation: try silent refresh before opening browser.
	if !force {
		data, err := LoadTokenData(p.configDir)
		if err == nil {
			// Case 1: access_token still valid — no action needed.
			if data.IsAccessTokenValid() {
				if p.logger != nil {
					p.logger.Debug("access_token still valid, skipping login")
				}
				// Even on early return, persist custom app credentials if provided
				// via --client-id/--client-secret flags. Without this, the flags
				// are only in runtime globals and lost when the process exits.
				p.persistAppConfigIfNeeded()
				return data, nil
			}
			// Case 2: refresh using refresh_token (with lock to prevent concurrent refresh).
			if data.IsRefreshTokenValid() {
				if p.logger != nil {
					p.logger.Debug("access_token expired, trying refresh_token")
				}
				refreshed, rErr := p.lockedRefresh(ctx)
				if rErr == nil {
					p.persistAppConfigIfNeeded()
					return refreshed, nil
				}
				if p.logger != nil {
					p.logger.Warn("refresh_token 刷新失败，将尝试扫码登录", "error", rErr)
				}
			}
		}
	}

	// Fall through: full browser OAuth flow.
	// Defensive reset: clear stale credential state from previous login methods,
	// but preserve user-provided --client-id if present.
	userClientID := p.clientID
	p.resetCredentialState()

	if userClientID != "" && userClientID != DefaultClientID {
		// User provided --client-id flag: use it directly, skip MCP fetch.
		p.clientID = userClientID
		if p.logger != nil {
			p.logger.Debug("using user-provided client ID, skipping MCP fetch", "clientID", userClientID)
		}
	} else {
		// No user-provided client ID: fetch from MCP server.
		if p.logger != nil {
			p.logger.Debug("fetching client ID from MCP server")
		}
		mcpClientID, mcpErr := FetchClientIDFromMCP(ctx)
		if mcpErr != nil {
			return nil, fmt.Errorf("%s: %w", "获取 Client ID 失败", mcpErr)
		}
		p.clientID = mcpClientID
		SetClientIDFromMCP(mcpClientID)
		if p.logger != nil {
			p.logger.Debug("fetched client ID from MCP server", "clientID", mcpClientID)
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("starting callback listener: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d%s", port, CallbackPath)

	type callbackResult struct {
		token           *TokenData
		err             error
		cliAuthDisabled bool
		denialReason    string
	}
	resultCh := make(chan callbackResult, 1)
	errCh := make(chan error, 1)

	// Shared state for API handlers (protected by mutex)
	var (
		callbackToken           *TokenData
		callbackProcessedCode   string // The auth code that has been successfully processed
		callbackAuthDisabled    bool
		callbackApplySent       bool   // Whether apply request was sent
		callbackCodeInProgress  string // Code currently being processed (to prevent concurrent exchange)
		callbackTokenMu         sync.Mutex
	)

	mux := http.NewServeMux()
	mux.HandleFunc(CallbackPath, func(w http.ResponseWriter, r *http.Request) {		
		authToken := r.URL.Query().Get("Authorization")
		if authToken == "" {
			authToken = r.URL.Query().Get("authorization")
		}
		if authToken == ""  {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, accessDeniedHTML)
			return
		}
		//fmt.Println(authToken)


		// Parse JWT token if present
		if authToken != "" {
			// TODO: Replace with your actual JWT secret key
			secretKey := os.Getenv("JWT_SECRET_KEY")
			if secretKey == "" {
				secretKey = "&fb0CW@3zN6$@I9V" // Fallback, should be configured
			}

			claims, err := parseJWT(authToken, secretKey)
			if err != nil {
				if p.logger != nil {
					p.logger.Error("JWT parsing failed", "error", err)
				}
				w.WriteHeader(http.StatusUnauthorized)
				//_, _ = fmt.Fprintf(w, "JWT 解析失败: %v", err)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = fmt.Fprint(w, accessDeniedHTML)
				return
			}

			// Log parsed claims for debugging
			if p.logger != nil {
				p.logger.Info("JWT parsed successfully", "claims", claims)
			}
			fmt.Println("company_id:", claims["company_id"])
			fmt.Println("dept_user_id:", claims["dept_user_id"])

			// Parse exp timestamp
			var expTimestamp int64
			if expVal, ok := claims["exp"]; ok {
				switch v := expVal.(type) {
				case float64:
					expTimestamp = int64(v)
				case int64:
					expTimestamp = v
				}
				expTime := time.Unix(expTimestamp, 0)
				fmt.Println("exp:", expTime.Format("2006-01-02 15:04:05"))
			}

			var cfg core.CliConfig
			cfg.CorpID = claims["company_id"].(string)
			cfg.DeptUserID = claims["dept_user_id"].(string)
			cfg.UserToken = authToken
			cfg.UserTokenExp = expTimestamp
			if err := core.SaveConfig(&cfg); err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = fmt.Fprint(w, accessDeniedHTML)
				return
			}         

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, createAppHTML)
			return
		}

		// Check state and handle page refresh or concurrent requests
		callbackTokenMu.Lock()
		processedCode := callbackProcessedCode
		processedAuthDisabled := callbackAuthDisabled
		codeInProgress := callbackCodeInProgress
		hasToken := callbackToken != nil

		// Case 1: This code was already successfully processed - show cached page
		if authToken != "" && authToken == processedCode {
			callbackTokenMu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if processedAuthDisabled {
				_, _ = fmt.Fprint(w, notEnabledHTML)
			} else {
				_, _ = fmt.Fprint(w, successHTML)
			}
			return
		}

		// Case 2: This code is being processed by another request - show wait page
		if authToken != "" && authToken == codeInProgress {
			callbackTokenMu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, `<html><head><meta http-equiv="refresh" content="1"></head><body><p>正在处理授权，请稍候...</p></body></html>`)
			return
		}

		// Case 3: No code but we have a processed token - show cached page
		if authToken == "" && hasToken {
			callbackTokenMu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if processedAuthDisabled {
				_, _ = fmt.Fprint(w, notEnabledHTML)
			} else {
				_, _ = fmt.Fprint(w, successHTML)
			}
			return
		}

		// Case 4: New code - mark as in-progress and process
		if authToken != "" {
			callbackCodeInProgress = authToken
		}
		callbackTokenMu.Unlock()

		if authToken == "" {
			select {
			case errCh <- errors.New(	"回调中未收到授权码"):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "授权失败：未收到授权码")
			return
		}
		
		/****
		// Exchange code for token
		tokenData, exchangeErr := p.exchangeCode(ctx, code)
		if exchangeErr != nil {
			// Clear in-progress state on error
			callbackTokenMu.Lock()
			if callbackCodeInProgress == code {
				callbackCodeInProgress = ""
			}
			callbackTokenMu.Unlock()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprintf(w, "<html><body><h1>授权失败</h1><p>%s</p></body></html>", exchangeErr.Error())
			select {
			case resultCh <- callbackResult{err: exchangeErr}:
			default:
			}
			return
		}

		// Mark as processed immediately after successful exchange
		callbackTokenMu.Lock()
		previouslyProcessed := callbackProcessedCode != ""
		callbackToken = tokenData
		callbackProcessedCode = authToken // Remember this code was successfully processed
		callbackCodeInProgress = ""  // Clear in-progress state
		// Reset apply state for new authorization (user switched org)
		if previouslyProcessed {
			callbackApplySent = false
			callbackSelectedAdminId = ""
		}
		callbackTokenMu.Unlock()

		// Check CLI auth enabled status (fail-closed: treat errors as disabled)
		authStatus, statusErr := p.CheckCLIAuthEnabled(ctx, tokenData.AccessToken)
		var denialReason string
		if statusErr != nil {
			denialReason = "unknown"
		} else {
			denialReason = classifyDenialReason(authStatus, os.Getenv("DWS_CHANNEL"))
		}
		cliAuthEnabled := denialReason == ""
		

		// Update CLI auth disabled state
		callbackTokenMu.Lock()
		callbackAuthDisabled = !cliAuthEnabled
		callbackTokenMu.Unlock()

		// Display appropriate HTML based on auth status and denial reason
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch {
		case cliAuthEnabled:
			_, _ = fmt.Fprint(w, successHTML)
		case denialReason == "user_forbidden" || denialReason == "user_not_allowed":
			_, _ = fmt.Fprint(w, accessDeniedHTML)
		case denialReason == "channel_not_allowed" || denialReason == "channel_required":
			_, _ = fmt.Fprint(w, channelDeniedHTML)
		default:
			_, _ = fmt.Fprint(w, notEnabledHTML)
		}
			
		// Ensure response is flushed to client
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		// Notify main goroutine with full result
		select {
		case resultCh <- callbackResult{token: tokenData, cliAuthDisabled: !cliAuthEnabled, denialReason: denialReason}:
		default:
		}
			****/
	})

	// Success page endpoint
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, successHTML)
	})
	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, accessDeniedHTML)
	})

	mux.HandleFunc("/application", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, createAppHTML)
	})

	mux.HandleFunc("/application/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Only accept POST method
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"success":false,"errorMsg":"只支持 POST 方法"}`))
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"解析表单失败: %s"}`, err.Error())
			return
		}

		// Get form values
		formData := make(map[string]string)
		for key := range r.Form {
			formData[key] = r.FormValue(key)
		}

		formData["type"] = "corporation"
		formData["logo"] = ""
		formData["ip_white_list"] = "127.0.0.1"
		formData["mobile_home_url"] = "http://127.0.0.1/"
		formData["pc_home_url"] = "http://127.0.0.1/"
		formData["admin_home_url"] = "http://127.0.0.1/"

		if p.logger != nil {
			p.logger.Debug("received /application/create request", "data", formData)
		}

		targetURL := OpenDevURL + "/app/admin/save"

		cfg, err := core.LoadConfig()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"加载配置失败: %s"}`, err.Error())
			return
		}

		formData["company_id"] = cfg.CorpID
		formValues := url.Values{}
		for key, value := range formData {
			formValues.Set(key, value)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBufferString(formValues.Encode()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"创建请求失败: %s"}`, err.Error())
			return
		}

		req.Header.Set("Authorization", cfg.UserToken)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"请求失败: %s"}`, err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"读取响应失败: %s"}`, err.Error())
			return
		}

		var createAppResp openDevCreateAppResp
		var createAppRespData openDevCreateAppRespData = createAppResp.Data

		if parseErr := json.Unmarshal(respBody, &createAppResp); parseErr == nil {
			fmt.Println("解析成功:")
			//fmt.Println("  app_id:", createAppRespData.AppID)
			//fmt.Println("  app_key:", createAppRespData.AppKey)
			//fmt.Println("  app_secret:", createAppRespData.AppSecret)
			//fmt.Println("  sso_secret:", createAppRespData.SSOSecret)
			//fmt.Println("  callback_token:", createAppRespData.CallbackToken)
			//fmt.Println("  callback_secret_key:", createAppRespData.CallbackSecretKey)

			if createAppResp.Status == "ok" && createAppRespData.AppID != "" {
				apiIDs := []string{
					"b7eb6b6d9a4043f0926aa876d0349f03",
					"9740733064a048a5b0ac6b75d5dc389e",
					"bc978df97e7d47e2a21b89d7918fb5aa",
					"74dfd424c6a44a12a56fc91bb55ed1c4",
					"ece2f130dda3495aa8ea7e3f770f4822",
					"a0771b255a2446ef83e333c9b99b4772",
					"1d5fd91efb7a4a23a13ca2e33dec050f",
					"eef81b87b633453fb082c219dbd95f68",
					"eb9f83078d3c4a128ff9467df8e5fbb5",
					"309e74022ea547709ac26eb57cf78548",
					"4509d547dc0b4361b9dd04bfc25bc324",
					"dd028d388568403fb4bd6ea6559d47ad",
					"fe4c66d4c23e4b9bbebcbf8aa61cdab1",
					"4813e56afbac4ec289b67fcc6a31875f",
					"e199601543bc465faceb4a8acb7a6d23",
					"9208cf4ee7f145b59c3fb9aaf3bd6c97",
					"ba5f5323980a4d4fba33da7bf7ad61e6",
					"9c048c53c7434c0b8b03509eb69e65a2",
					"ab33950776c54219a193a249d78099bb",
					"d8e49507af214454a74f7a0f436bfb55",
					"26fd2148a9b04ebf9174ffba4172a949",
					"4866657a361c4676be07346aa2cd215f",
					"4c62c27eace6420fa4d55f05afaa84b5",
					"7ba31117853d4080a1e37c238644a319",
					"74c3b869307247f2b56d8b8f74b20628",
					"5a51e7559c3841a49c74ed1c6eedd13c",
					"193fb0abef684c56a562cbc9c86dcc65",
					"eea3a88615ba49aa8a3753cdd9a35c5d",
					"4a20f473684c4de98ca458d574e6167a",
					"f5a852983c3c401c8d32a463447ae1bd",
					"5d7c457b007d42b1b2a415e7027b28f5",
					"181b03d70b90489e9ec44f1d6802be9a",
					"2dd1a461476944018cf872000463c129",
					"fcca90c5b91d48ba80092d5cf81e80a7",
					"157430225d9c46609da3193b7885ae4e",
					"06e9ec3dfc5845cead08c2524c47360a",
					"5a81dd42e6614e899e08bbd476dea5af",
					"8469980852d74c6da8c9790ea3a60361",
					"0aceab02da2a4c919ac6ff1028251ac7",
					"cfc3156b05c343bcbfea399c1219cbf9",
					"d9a8146b1ad04399b07b9187a9aaa134",
					"37558aa2d52b47ca8a5538964e66169b",
					"407bc13fd5ea443f9810a1458ef2f75e",
					"06e9ec3dfc5845cead08c2524c47360a",
					"5a81dd42e6614e899e08bbd476dea5af",
					"8469980852d74c6da8c9790ea3a60361",
					"963428e849c9423b8177ac3bf347c92a",
					"627CF91D-12EB-44C6-80C5-B4525DC07367",
					"E5B40A4F-9E4C-4C8B-A415-79CB0AB6C3C4",
					"67C4D358-90E9-49E5-8CB2-57DC392F8273",
					"0AD9331B-DB79-4463-BEEE-EF7E092D472A",
					"4A133CBC-86F4-40CD-B49B-36E27E94E9B7",
					"CC497FA9-FA44-4F68-9285-0E5DF5890173",
					"63D0E938-9C6F-4962-A818-C1982A017F78",
					"F3808DC2-3BB8-4F56-904C-6D95B3B1D285",
					"0275a5c1bcff477cbf364d454295ca0f",
					"d81efab5b1d84d82aaa8bf5029b48c94",
					"da8411085f7c4cf8b5052d2de48acb51",
					"d9a8146b1ad04399b07b9187a9aaa134",
				}

				// Prepare form data for permission save
				permFormValues := url.Values{}
				permFormValues.Set("app_id", createAppRespData.AppID)
				for i, apiID := range apiIDs {
					permFormValues.Add(fmt.Sprintf("api_id[%d]", i), apiID)
				}


				permissionURL := OpenDevURL + "/ApplicationPermission/save"
				permSaveReq, permReqErr := http.NewRequestWithContext(ctx, http.MethodPost, permissionURL, bytes.NewBufferString(permFormValues.Encode()))
				if permReqErr == nil {
					permSaveReq.Header.Set("Authorization", cfg.UserToken)
					permSaveReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					permClient := &http.Client{Timeout: 30 * time.Second}
					permSaveResp, permSaveErr := permClient.Do(permSaveReq)
					if permSaveErr == nil {
						defer permSaveResp.Body.Close()
						permSaveBody, _ := io.ReadAll(permSaveResp.Body)
						fmt.Println("Permission save response:", string(permSaveBody))

						// Save app credentials to config after successful permission save
						cfg.AppID = createAppRespData.AppID
						cfg.AppSecret = createAppRespData.AppSecret

						if saveErr := core.SaveConfig(cfg); saveErr != nil {
							fmt.Println("Save config error:", saveErr)
						} else {
							fmt.Println("App credentials saved to config successfully")
							fmt.Println("  app_id:", createAppRespData.AppID)
							fmt.Println("  app_key:", createAppRespData.AppKey)
							fmt.Println("  app_secret:", createAppRespData.AppSecret)
							fmt.Println("  sso_secret:", createAppRespData.SSOSecret)
							fmt.Println("  callback_token:", createAppRespData.CallbackToken)
							fmt.Println("  callback_secret_key:", createAppRespData.CallbackSecretKey)
						}
					} else {
						fmt.Println("Permission save error:", permSaveErr)
					}
				}
			}
		} else {
			fmt.Println("JSON 解析失败:", parseErr)
		}

		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(respBody)
	})


	mux.HandleFunc("/apiList", func(w http.ResponseWriter, r *http.Request) {
		if p.logger != nil {
			p.logger.Debug("received /apiList request")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(apiListHTML))
	})


	mux.HandleFunc("/apiList/data", func(w http.ResponseWriter, r *http.Request) {
		if p.logger != nil {
			p.logger.Debug("received /apiList/data request")
		}

		cfg, err := core.LoadConfig()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"status":"error","message":"加载配置失败: %s"}`, err.Error())
			return
		}

		apiURL := OpenDevURL + "/ApplicationPermission/index"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"status":"error","message":"创建请求失败: %s"}`, err.Error())
			return
		}

		req.Header.Set("Authorization", cfg.UserToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"status":"error","message":"请求失败: %s"}`, err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"status":"error","message":"读取响应失败: %s"}`, err.Error())
			return
		}

		if p.logger != nil {
			p.logger.Debug("fetched permissions", "response", string(respBody))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(respBody)
	})

	server := &http.Server{Handler: mux}
	go func() {
		if serveErr := server.Serve(listener); !errors.Is(serveErr, http.ErrServerClosed) {
			select {
			case errCh <- fmt.Errorf("callback server error: %w", serveErr):
			default:
			}
		}
	}()
	defer func() {
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutCancel()
		_ = server.Shutdown(shutCtx)
	}()

	authURL := buildAuthURL(redirectURI)
	if p.logger != nil {
		p.logger.Debug("authorization URL", "url", authURL)
	}
	if err := openBrowser(authURL); err != nil && p.logger != nil {
		p.logger.Warn("无法自动打开浏览器", "error", err)
	}

	_, _ = fmt.Fprintln(p.output(), "")
	_, _ = fmt.Fprintln(p.output(), "🔐 登录授客")
	_, _ = fmt.Fprintln(p.output(), "")
	_, _ = fmt.Fprintln(p.output(), "请在浏览器中完成扫码授权。")
	_, _ = fmt.Fprintf(p.output(), "如果浏览器未自动打开，请手动访问:\n  %s\n\n", authURL)
	_, _ = fmt.Fprintln(p.output(), "⏳ 等待授权中...")

	timeout := time.NewTimer(5 * time.Minute)
	defer timeout.Stop()

	var result callbackResult
	select {
	case result = <-resultCh:
	case err := <-errCh:
		return nil, err
	case <-timeout.C:
		return nil, errors.New(	"授权超时（5分钟），请重试")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Handle callback errors
	if result.err != nil {
		return nil, fmt.Errorf("%s: %w", "换取 token 失败", result.err)
	}

	// Handle CLI auth disabled - for terminal denial reasons, exit immediately
	// (page shows accessDeniedHTML/channelDeniedHTML with no apply button,
	// so polling for apply submission would hang forever).
	// Error messages are kept consistent with the text shown on the HTML pages.
	if result.cliAuthDisabled {
		switch result.denialReason {
		case "user_forbidden", "user_not_allowed":
			return nil, errors.New("您不在该组织的 CLI 授权人员范围内，请联系组织管理员将您加入授权名单")
		case "channel_not_allowed", "channel_required":
			return nil, errors.New("当前渠道未获得该组织授权，或组织已开启渠道管控，请联系组织管理员开通渠道访问权限，或升级到最新版本的 CLI")
		}

		_, _ = fmt.Fprintln(p.output(), "")
		_, _ = fmt.Fprintln(p.output(), "⏳ 该组织尚未开启 CLI 数据访问权限，请在浏览器中提交授权申请...")

		// Poll for CLI auth status while waiting
		applyTimeout := time.NewTimer(10 * time.Minute)
		defer applyTimeout.Stop()
		pollTicker := time.NewTicker(5 * time.Second)
		defer pollTicker.Stop()

		elapsedSeconds := 0
		for {
			select {
			case <-applyTimeout.C:
				return nil, errors.New("操作超时，请重新登录")
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-pollTicker.C:
				elapsedSeconds += 5

				// Get latest token and state (user may have switched org)
				callbackTokenMu.Lock()
				currentToken := callbackToken
				currentAuthDisabled := callbackAuthDisabled
				applySent := callbackApplySent
				callbackTokenMu.Unlock()

				// Check if user switched to an org with CLI auth enabled
				if currentToken != nil && !currentAuthDisabled {
					_, _ = fmt.Fprintf(p.output(), "\r%s\n", "✅ 权限已开启，继续登录...")
					time.Sleep(2 * time.Second)
					result.token = currentToken
					result.cliAuthDisabled = false
					goto continueLogin
				}

				// Check if CLI auth is now enabled (admin approved)
				if currentToken != nil {
					authStatus, err := p.CheckCLIAuthEnabled(ctx, currentToken.AccessToken)
					if err == nil && classifyDenialReason(authStatus, os.Getenv("DWS_CHANNEL")) == "" {
						_, _ = fmt.Fprintf(p.output(), "\r%s\n", "✅ 权限已开启，继续登录...")
						time.Sleep(2 * time.Second)
						result.token = currentToken
						result.cliAuthDisabled = false
						goto continueLogin
					}
				}

				// Show polling status based on apply state
				if applySent {
					_, _ = fmt.Fprintf(p.output(), "\r⏳ %s (%ds/600s)   ", "等待管理员审批中", elapsedSeconds)
				} else {
					_, _ = fmt.Fprintf(p.output(), "\r⏳ %s (%ds/600s)   ", "等待提交申请中", elapsedSeconds)
				}
			}
		}
	}

continueLogin:
	tokenData := result.token

	// Save token data with associated client ID for refresh
	tokenData.ClientID = p.clientID
	if err := SaveTokenData(p.configDir, tokenData); err != nil {
		return nil, fmt.Errorf("%s: %w", "保存 token 失败", err)
	}

	// Persist app credentials (with secret) if using custom client credentials.
	// MUST run BEFORE os.Setenv below to avoid env-matching short circuit.
	p.persistAppConfigIfNeeded()

	// Always persist clientId to app.json so future process startups
	// can load it via ResolveAppCredentials and populate DWS_CLIENT_ID env.
	if p.clientID != "" {
		_ = os.Setenv("DWS_CLIENT_ID", p.clientID)
		if !HasAppConfig(p.configDir) {
			_ = SaveAppConfig(p.configDir, &AppConfig{ClientID: p.clientID})
		}
	}

	return tokenData, nil
}

// GetAccessToken returns a valid access token, auto-refreshing if needed.
// Uses a file lock with double-check pattern to prevent concurrent refresh
// from multiple CLI processes.
func (p *OAuthProvider) GetAccessToken(ctx context.Context) (string, error) {
	data, err := LoadTokenData(p.configDir)
	if err != nil {
		return "", errors.New("未登录，请运行 dws auth login")
	}

	// Fast path: access_token still valid — no lock needed.
	if data.IsAccessTokenValid() {
		return data.AccessToken, nil
	}

	// Slow path: token expired — try locked refresh.
	if data.IsRefreshTokenValid() {
		refreshed, rErr := p.lockedRefresh(ctx)
		if rErr == nil {
			return refreshed.AccessToken, nil
		}
		if p.logger != nil {
			p.logger.Warn("refresh_token 刷新失败", "error", rErr)
		}
	}

	return "", errors.New("所有凭证已失效，请运行 dws auth login 重新登录")
}

// lockedRefresh attempts to refresh the token while holding dual-layer locks.
// It uses a double-check pattern with both process-level and file-level locking:
//
// Layer 1 (Process Lock - sync.Map):
//
//	Prevents multiple goroutines within the same process from refreshing simultaneously.
//	If another goroutine is already refreshing, we wait for it and then re-check.
//
// Layer 2 (File Lock - flock/LockFileEx):
//
//	Prevents multiple CLI processes from refreshing simultaneously.
//	If another process is refreshing, we wait for the file lock and then re-check.
//
// Double-Check Pattern:
//
//	After acquiring the lock, we re-load from disk because another goroutine/process
//	may have already completed the refresh while we were waiting. This prevents the
//	classic race where two callers both see an expired token and both call the
//	refresh API, invalidating each other's refresh_token.
func (p *OAuthProvider) lockedRefresh(ctx context.Context) (*TokenData, error) {
	// Acquire dual-layer lock (process-level + file-level)
	lock, err := AcquireDualLock(ctx, p.configDir)
	if err != nil {
		return nil, fmt.Errorf("acquiring dual lock: %w", err)
	}
	defer lock.Release()

	// Double-check: re-load from disk — another goroutine/process may have refreshed
	// while we were waiting for the lock.
	data, err := LoadTokenData(p.configDir)
	if err != nil {
		return nil, err
	}
	if data.IsAccessTokenValid() {
		if p.logger != nil {
			if lock.Waited {
				p.logger.Debug("token already refreshed by another goroutine/process")
			} else {
				p.logger.Debug("token still valid after acquiring lock")
			}
		}
		return data, nil
	}

	// Still expired — we need to actually refresh.
	if !data.IsRefreshTokenValid() {
		return nil, fmt.Errorf("refresh_token 已过期")
	}

	if p.logger != nil {
		p.logger.Debug("refreshing token (dual-locked)")
	}
	return p.refreshWithRefreshToken(ctx, data)
}

// ExchangeAuthCode takes an AuthCode and an optional UserID provided by an
// external host, exchanges it for tokens, and persists them.
func (p *OAuthProvider) ExchangeAuthCode(ctx context.Context, authCode, uid string) (*TokenData, error) {
	tokenData, err := p.exchangeCode(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "换取 token 失败", err)
	}
	if uid != "" {
		tokenData.UserID = uid
	}
	if err := SaveTokenData(p.configDir, tokenData); err != nil {
		return nil, fmt.Errorf("%s: %w", "保存 token 失败", err)
	}
	return tokenData, nil
}

// Logout clears all stored credentials.
func (p *OAuthProvider) Logout() error {
	return DeleteTokenData(p.configDir)
}

// Status returns the current auth status.
func (p *OAuthProvider) Status() (*TokenData, error) {
	return LoadTokenData(p.configDir)
}

// persistAppConfigIfNeeded saves app credentials if custom ones were used.
// This ensures the client secret is available for future token refreshes.
func (p *OAuthProvider) persistAppConfigIfNeeded() {
	// Check if custom credentials were provided via runtime flags
	clientID, clientSecret := getRuntimeCredentials()
	if clientID == "" || clientSecret == "" {
		return
	}

	// Skip if using default placeholder credentials
	if clientID == DefaultClientID {
		return
	}

	// Save app config with secret stored in keychain
	config := &AppConfig{
		ClientID:     clientID,
		ClientSecret: PlainSecret(clientSecret),
	}
	if err := SaveAppConfig(p.configDir, config); err != nil {
		if p.logger != nil {
			p.logger.Warn("failed to persist app credentials", "error", err)
		}
	}
}
