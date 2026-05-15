package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TokenManager 管理access_token的获取和缓存
type TokenManager struct {
	appKey    string
	appSecret string
	corpID    string
	baseURL   string

	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time

	httpClient *http.Client
}

// TokenResponse gettoken接口的响应结构
type TokenResponse struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// NewTokenManager 创建TokenManager实例
func NewTokenManager(appKey, appSecret, corpID, baseURL string) *TokenManager {
	return &TokenManager{
		appKey:     appKey,
		appSecret:  appSecret,
		corpID:     corpID,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetAccessToken 获取access_token，自动处理缓存和续期
func (tm *TokenManager) GetAccessToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	if tm.accessToken != "" && time.Now().Before(tm.expiresAt) {
		token := tm.accessToken
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 双重检查，防止并发重复请求
	if tm.accessToken != "" && time.Now().Before(tm.expiresAt) {
		return tm.accessToken, nil
	}

	// 请求新的access_token
	token, err := tm.fetchAccessToken(ctx)
	if err != nil {
		return "", err
	}

	tm.accessToken = token
	// 设置过期时间为1小时50分钟（留10分钟buffer）
	tm.expiresAt = time.Now().Add(110 * time.Minute)

	return token, nil
}

// fetchAccessToken 从API获取access_token
func (tm *TokenManager) fetchAccessToken(ctx context.Context) (string, error) {
	params := url.Values{}
	params.Set("app_key", tm.appKey)
	params.Set("app_secret", tm.appSecret)
	params.Set("corpid", tm.corpID)

	reqURL := fmt.Sprintf("%s/service/corp/gettoken?%s", tm.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parse response failed: %w", err)
	}

	if tokenResp.Code != "200" {
		return "", fmt.Errorf("get token failed: code=%s, message=%s", tokenResp.Code, tokenResp.Message)
	}

	if tokenResp.Data == "" {
		return "", fmt.Errorf("empty access_token in response")
	}

	return tokenResp.Data, nil
}
