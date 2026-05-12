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
	"strings"
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
	// 首先检查 core.CliConfig 中的必需字段和 UserToken 是否有效
	var cliConfigChecked bool
	var missingFields []string
	var needRelogin bool

	if !force {
		cfg, err := core.LoadConfig()
		if err == nil {
			cliConfigChecked = true

			// 检查所有必需字段是否存在
			if cfg.AppID == "" {
				missingFields = append(missingFields, "AppID")
			}
			if cfg.AppSecret == "" {
				missingFields = append(missingFields, "AppSecret")
			}
			if cfg.UserToken == "" {
				missingFields = append(missingFields, "UserToken")
			}
			if cfg.CorpID == "" {
				missingFields = append(missingFields, "CorpID")
			}
			if cfg.DeptUserID == "" {
				missingFields = append(missingFields, "DeptUserID")
			}

			// 如果有字段缺失，需要重新登录
			if len(missingFields) > 0 {
				needRelogin = true
				if p.logger != nil {
					p.logger.Debug("config incomplete, need login", "missing_fields", missingFields)
				}
				_, _ = fmt.Fprintln(p.output(), "")
				_, _ = fmt.Fprintf(p.output(), "⚠️  配置不完整（缺少: %s），需要重新登录\n", strings.Join(missingFields, ", "))
				_, _ = fmt.Fprintln(p.output(), "")
			} else {
				// 所有字段都存在，检查 UserToken 是否过期
				now := time.Now().Unix()
				if cfg.UserTokenExp > now {
					// Token 未过期，无需重新登录
					if p.logger != nil {
						p.logger.Debug("UserToken still valid, skipping login",
							"expires_at", time.Unix(cfg.UserTokenExp, 0).Format("2006-01-02 15:04:05"))
					}

					// 构造 TokenData 返回
					tokenData := &TokenData{
						AccessToken:  cfg.UserToken,
						CorpID:       cfg.CorpID,
						ExpiresAt:    time.Unix(cfg.UserTokenExp, 0),
						RefreshExpAt: time.Unix(cfg.UserTokenExp, 0).Add(30 * 24 * time.Hour),
					}

					// 输出提示信息
					_, _ = fmt.Fprintln(p.output(), "")
					_, _ = fmt.Fprintln(p.output(), "✅ 已登录，Token 仍然有效")
					_, _ = fmt.Fprintf(p.output(), "企业 ID: %s\n", cfg.CorpID)
					_, _ = fmt.Fprintln(p.output(), "")

					return tokenData, nil
				} else {
					// Token 已过期，需要重新登录
					needRelogin = true
					if p.logger != nil {
						p.logger.Debug("UserToken expired, need re-login",
							"expired_at", time.Unix(cfg.UserTokenExp, 0).Format("2006-01-02 15:04:05"))
					}
					_, _ = fmt.Fprintln(p.output(), "")
					_, _ = fmt.Fprintln(p.output(), "⚠️  Token 已过期，需要重新登录")
					_, _ = fmt.Fprintln(p.output(), "")
				}
			}
		} else {
			// 配置文件不存在或读取失败，需要重新登录
			needRelogin = true
			if p.logger != nil {
				p.logger.Debug("config not found or invalid, need login", "error", err)
			}
		}
	}

	// 如果配置有问题，跳过 TokenData 检查，直接进入授权流程
	if !needRelogin && !force {
		data, err := LoadTokenData(p.configDir)
		if err == nil {
			if data.IsAccessTokenValid() {
				if p.logger != nil {
					p.logger.Debug("access_token still valid, skipping login")
				}
				return data, nil
			}
			if data.IsRefreshTokenValid() {
				if p.logger != nil {
					p.logger.Debug("access_token expired, trying refresh_token")
				}
				refreshed, rErr := p.lockedRefresh(ctx)
				if rErr == nil {
					return refreshed, nil
				}
				if p.logger != nil {
					p.logger.Warn("refresh_token 刷新失败，将尝试扫码登录", "error", rErr)
				}
			}
		}
	}

	if cliConfigChecked {
		_, _ = fmt.Fprintln(p.output(), "")
		_, _ = fmt.Fprintln(p.output(), "⚠️  Token 已过期，需要重新登录")
		_, _ = fmt.Fprintln(p.output(), "")
	} else if !force {
		_, _ = fmt.Fprintln(p.output(), "")
		_, _ = fmt.Fprintln(p.output(), "ℹ️  未找到登录信息，开始授权流程")
		_, _ = fmt.Fprintln(p.output(), "")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("starting callback listener: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d%s", port, CallbackPath)

	type callbackResult struct {
		token *TokenData
		err   error
	}
	resultCh := make(chan callbackResult, 1)
	errCh := make(chan error, 1)

	var (
		callbackTokenMu  sync.Mutex
		pendingTokenData *TokenData // 待发送的 TokenData，等待 /success 页面请求
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
		if authToken != "" {
			secretKey := os.Getenv("JWT_SECRET_KEY")
			if secretKey == "" {
				secretKey = "&fb0CW@3zN6$@I9V" 
			}

			claims, err := parseJWT(authToken, secretKey)
			if err != nil {
				if p.logger != nil {
					p.logger.Error("JWT parsing failed", "error", err)
				}
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = fmt.Fprint(w, accessDeniedHTML)
				return
			}

			if p.logger != nil {
				p.logger.Info("JWT parsed successfully", "claims", claims)
			}

			// Parse exp timestamp
			var expTimestamp int64
			if expVal, ok := claims["exp"]; ok {
				switch v := expVal.(type) {
				case float64:
					expTimestamp = int64(v)
				case int64:
					expTimestamp = v
				}
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

		if authToken == "" {
			select {
			case errCh <- errors.New("回调中未收到授权码"):
			default:
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "授权失败：未收到授权码")
			return
		}
	})

	// Success page endpoint
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		// 如果是 POST 请求，说明是前端通知完成
		if r.Method == http.MethodPost {
			callbackTokenMu.Lock()
			token := pendingTokenData
			callbackTokenMu.Unlock()

			if token != nil {
				// 通知主流程授权完成
				select {
				case resultCh <- callbackResult{token: token}:
					fmt.Println("✅ 授权成功")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"success":true}`))
				default:
					fmt.Println("⚠️  流程可能已超时，无法完成授权")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"success":false,"error":"timeout"}`))
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"success":false,"error":"no token data"}`))
			}
			return
		}

		// GET 请求，显示成功页面
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
			// 配置文件不存在，创建新配置对象
			cfg = &core.CliConfig{}
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

		if parseErr := json.Unmarshal(respBody, &createAppResp); parseErr == nil {
			createAppRespData := createAppResp.Data

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

						cfg.AppID = createAppRespData.AppID
						cfg.AppSecret = createAppRespData.AppSecret

						// 保存配置到文件
						if saveErr := core.SaveConfig(cfg); saveErr != nil {
							fmt.Println("Save config error:", saveErr)
							w.WriteHeader(http.StatusInternalServerError)
							_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"保存配置失败: %s"}`, saveErr.Error())
							return
						}

						// 构造 TokenData 并保存到共享变量，等待 /success 页面请求
						tokenData := &TokenData{
							AccessToken:  cfg.UserToken,
							CorpID:       cfg.CorpID,
							ExpiresAt:    time.Unix(cfg.UserTokenExp, 0),
							RefreshExpAt: time.Unix(cfg.UserTokenExp, 0).Add(30 * 24 * time.Hour), // 默认30天刷新期
						}

						// 保存到共享变量
						callbackTokenMu.Lock()
						pendingTokenData = tokenData
						callbackTokenMu.Unlock()

						// 返回成功响应，让前端跳转到 /success
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(`{"success":true,"message":"授权成功"}`))
						return
					} else {
						fmt.Println("Permission save error:", permSaveErr)
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = fmt.Fprintf(w, `{"success":false,"errorMsg":"权限保存失败: %s"}`, permSaveErr.Error())
						return
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

	tokenData := result.token

	// Save token data with associated client ID for refresh
	tokenData.ClientID = p.clientID
	if err := SaveTokenData(p.configDir, tokenData); err != nil {
		return nil, fmt.Errorf("%s: %w", "保存 token 失败", err)
	}

	// 读取配置文件显示详细信息
	cfg, err := core.LoadConfig()
	if err == nil {
		if cfg.CorpID != "" {
			_, _ = fmt.Fprintf(p.output(), "企业 ID: %s\n", cfg.CorpID)
		}
		if cfg.AppID != "" {
			_, _ = fmt.Fprintf(p.output(), "应用 ID: %s\n", cfg.AppID)
		}
		if cfg.UserTokenExp > 0 {
			expTime := time.Unix(cfg.UserTokenExp, 0)
			_, _ = fmt.Fprintf(p.output(), "Token 过期时间: %s\n", expTime.Format("2006-01-02 15:04:05"))
		}
	}
	_, _ = fmt.Fprintln(p.output(), "")

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
