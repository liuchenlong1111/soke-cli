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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/config"
)

const apiListHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>授权确认</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: #f5f5f5;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
            max-width: 480px;
            width: 100%;
            padding: 40px 32px;
            text-align: center;
        }

        .avatar {
            width: 80px;
            height: 80px;
            border-radius: 50%;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0 auto 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 36px;
            color: white;
        }

        .app-name {
            font-size: 14px;
            color: #666;
            margin-bottom: 32px;
        }

        .title {
            font-size: 20px;
            font-weight: 500;
            color: #333;
            margin-bottom: 32px;
        }

        .permissions {
            text-align: left;
            margin-bottom: 32px;
        }

        .permission-item {
            display: flex;
            align-items: center;
            padding: 8px 0;
            color: #666;
            font-size: 14px;
        }

        .permission-item::before {
            content: "•";
            margin-right: 12px;
            color: #999;
            font-size: 18px;
        }

        .consent-section {
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 32px;
            padding: 16px;
            background: #f8f9fa;
            border-radius: 8px;
        }

        .checkbox-wrapper {
            display: flex;
            align-items: center;
            cursor: pointer;
        }

        input[type="checkbox"] {
            width: 18px;
            height: 18px;
            margin-right: 10px;
            cursor: pointer;
            accent-color: #1a73e8;
        }

        .consent-text {
            font-size: 14px;
            color: #333;
            user-select: none;
        }

        .authorize-btn {
            width: 100%;
            padding: 14px 24px;
            background: #1a73e8;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }

        .authorize-btn:hover {
            background: #1557b0;
        }

        .authorize-btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="app-name">授客 CLI</div>
        <h1 class="title">确定开通并授权以下权限吗？</h1>

        <div class="permissions" id="permissionsList">
            <div class="permission-item">加载中...</div>
        </div>

        <div class="consent-section">
            <label class="checkbox-wrapper">
                <input type="checkbox" id="consentCheckbox" checked>
                <span class="consent-text">一并开通审批的常用权限</span>
            </label>
        </div>

        <button class="authorize-btn" id="authorizeBtn">开通并授权</button>
    </div>

    <script>
        const permissionsList = document.getElementById('permissionsList');
        const authorizeBtn = document.getElementById('authorizeBtn');

        // 加载权限列表
        async function loadPermissions() {
            try {
                const response = await fetch('/apiList/data', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                const result = await response.json();

                if (result.status === 'ok' && result.data && Array.isArray(result.data)) {
                    if (result.data.length === 0) {
                        permissionsList.innerHTML = '<div class="permission-item">暂无权限项</div>';
                        return;
                    }

                    // 渲染权限列表
                    permissionsList.innerHTML = result.data.map(permission => {
                        return '<div class="permission-item">' + (permission.name || permission.description || permission) + '</div>';
                    }).join('');
                } else {
                    permissionsList.innerHTML = '<div class="permission-item" style="color: #f44336;">加载失败: ' + (result.message || '未知错误') + '</div>';
                }
            } catch (error) {
                console.error('加载权限失败:', error);
                permissionsList.innerHTML = '<div class="permission-item" style="color: #f44336;">网络请求失败</div>';
            }
        }

        // 页面加载时获取权限列表
        loadPermissions();

        // 授权按钮点击事件
        authorizeBtn.addEventListener('click', function() {
            const consentChecked = document.getElementById('consentCheckbox').checked;
            console.log('授权确认', { consentChecked: consentChecked });

            alert('授权成功！');
        });
    </script>
</body>
</html>
`

const createAppHTML = `<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>创建飞书 CLI 应用</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }

      body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        background-color: #f5f5f5;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        padding: 20px;
      }

      .container {
        background: white;
        border-radius: 8px;
        padding: 40px;
        width: 100%;
        max-width: 500px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      }

      h1 {
        text-align: center;
        font-size: 24px;
        font-weight: 500;
        margin-bottom: 40px;
        color: #1f2329;
      }

      .form-group {
        margin-bottom: 24px;
      }

      label {
        display: block;
        font-size: 14px;
        color: #1f2329;
        margin-bottom: 8px;
      }

      label .required {
        color: #f53f3f;
      }

      input[type="text"] {
        width: 100%;
        padding: 10px 12px;
        border: 1px solid #e5e6eb;
        border-radius: 4px;
        font-size: 14px;
        transition: border-color 0.2s;
      }

      input[type="text"]:focus {
        outline: none;
        border-color: #3370ff;
      }

      .hint {
        font-size: 12px;
        color: #86909c;
        margin-top: 8px;
        cursor: pointer;
      }

      .hint:hover {
        color: #3370ff;
      }

      .btn-group {
        display: flex;
        flex-direction: column;
        gap: 12px;
        margin-top: 32px;
      }

      button {
        width: 100%;
        padding: 12px;
        border: none;
        border-radius: 4px;
        font-size: 14px;
        cursor: pointer;
        transition: all 0.2s;
      }

      .btn-primary {
        background: #3370ff;
        color: white;
      }

      .btn-primary:hover {
        background: #2b5dd4;
      }

      .btn-primary:disabled {
        background: #c9cdd4;
        cursor: not-allowed;
      }

      .btn-secondary {
        background: white;
        color: #1f2329;
        border: 1px solid #e5e6eb;
      }

      .btn-secondary:hover {
        background: #f7f8fa;
      }

      .error-message {
        color: #f53f3f;
        font-size: 12px;
        margin-top: 4px;
        display: none;
      }

      .success-message {
        background: #e8f7f0;
        color: #00b42a;
        padding: 12px;
        border-radius: 4px;
        margin-bottom: 20px;
        display: none;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="success-message" id="successMessage">
        应用创建成功！
      </div>

      <h1>创建 CLI 应用</h1>

      <form id="appForm">
        <div class="form-group">
          <label for="name">名称<span class="required">*</span></label>
          <input type="text" id="name" name="name" placeholder="授客 CLI应用" required>
          <div class="hint">创建后，自动完成所有配置 ></div>
          <div class="error-message" id="nameError">请输入应用名称</div>
        </div>

        <div class="form-group">
          <label for="description">描述</label>
          <input type="text" id="description" name="description" placeholder="应用描述">
        </div>

        <div class="btn-group">
          <button type="submit" class="btn-primary" id="submitBtn">创建</button>
          <button type="button" class="btn-secondary">选择已有应用</button>
        </div>
      </form>
    </div>

    <script>
      const form = document.getElementById('appForm');
      const submitBtn = document.getElementById('submitBtn');
      const successMessage = document.getElementById('successMessage');

      form.addEventListener('submit', async function(e) {
        e.preventDefault();

        document.querySelectorAll('.error-message').forEach(el => el.style.display = 'none');
        successMessage.style.display = 'none';

        const name = document.getElementById('name').value.trim();

        if (!name) {
          document.getElementById('nameError').style.display = 'block';
          return;
        }

        const formData = new FormData(form);
        const data = {};
        formData.forEach((value, key) => {
          data[key] = value;
        });

        submitBtn.disabled = true;
        submitBtn.textContent = '创建中...';

        try {
          const response = await fetch('/application/create', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams(data)
          });

          if (response.ok) {
            const result = await response.json();
            if (result.status === 'ok') {
              successMessage.style.display = 'block';
              form.reset();
              console.log('创建成功:', result);
              // 跳转到成功页面
              setTimeout(() => {
                window.location.href = '/success';
              }, 3000);
            } else {
              alert('创建失败: ' + (result.message || result.msg || '未知错误'));
            }
          } else {
            alert('创建失败，请稍后重试');
          }
        } catch (error) {
          console.error('请求错误:', error);
          alert('网络错误，请检查连接');
        } finally {
          submitBtn.disabled = false;
          submitBtn.textContent = '创建';
        }
      });
    </script>
  </body>
</html>
`

func (p *OAuthProvider) exchangeCode(ctx context.Context, code string) (*TokenData, error) {
	// Use MCP mode if clientID is from MCP server
	if IsClientIDFromMCP() {
		return p.exchangeCodeViaMCP(ctx, code)
	}
	// Direct mode with client secret
	clientID := ClientID()
	clientSecret := ClientSecret()
	body := map[string]string{
		"clientId":     clientID,
		"clientSecret": clientSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}
	resp, err := p.postJSON(ctx, UserAccessTokenURL, body)
	if err != nil {
		return nil, err
	}
	data, err := p.parseTokenResponse(resp)
	if err != nil {
		return nil, err
	}
	// Snapshot credentials used for this token (for refresh)
	data.ClientID = clientID
	data.Source = resolveCredentialSource()
	// Save clientSecret for future refresh (even if env changes)
	if err := SaveClientSecret(clientID, clientSecret); err != nil {
		// Log warning but don't fail login
		fmt.Fprintf(p.Output, "Warning: failed to save client secret: %v\n", err)
	}
	return data, nil
}

// ExchangeCodeForToken exchanges an authorization code for token data using
// the currently configured client credentials.  This is a convenience wrapper
// around OAuthProvider.exchangeCode for callers outside the auth package.
func ExchangeCodeForToken(ctx context.Context, configDir, code string) (*TokenData, error) {
	p := &OAuthProvider{
		configDir:  configDir,
		clientID:   ClientID(),
		Output:     io.Discard,
		httpClient: oauthHTTPClient,
	}
	return p.exchangeCode(ctx, code)
}

// exchangeCodeViaMCP exchanges auth code for token via MCP proxy.
// This is used when client secret is not available (server-side secret management).
func (p *OAuthProvider) exchangeCodeViaMCP(ctx context.Context, code string) (*TokenData, error) {
	clientID := ClientID()
	url := GetMCPBaseURL() + MCPOAuthTokenPath
	body := map[string]string{
		"clientId":  clientID,
		"authCode":  code,
		"grantType": "authorization_code",
	}
	resp, err := p.postJSON(ctx, url, body)
	if err != nil {
		return nil, err
	}
	data, err := p.parseMCPTokenResponse(resp)
	if err != nil {
		return nil, err
	}
	// Snapshot credentials used for this token (for refresh)
	data.ClientID = clientID
	data.Source = "mcp"
	// MCP mode doesn't need to save clientSecret (server-side managed)
	return data, nil
}

func (p *OAuthProvider) refreshWithRefreshToken(ctx context.Context, data *TokenData) (*TokenData, error) {
	// Use stored Source to determine refresh path (not current runtime state)
	// This ensures refresh works even if environment variables changed since login
	if data.Source == "mcp" {
		return p.refreshViaMCP(ctx, data)
	}

	// Direct mode: use stored clientId and load saved clientSecret
	clientID := data.ClientID
	if clientID == "" {
		// Fallback for legacy tokens without stored clientId
		clientID = ClientID()
	}
	clientSecret := LoadClientSecret(clientID)
	if clientSecret == "" {
		// Fallback: try current environment
		clientSecret = ClientSecret()
	}

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("无法刷新 token: 缺少 clientId 或 clientSecret，请重新登录")
	}

	body := map[string]string{
		"clientId":     clientID,
		"clientSecret": clientSecret,
		"refreshToken": data.RefreshToken,
		"grantType":    "refresh_token",
	}
	resp, err := p.postJSON(ctx, UserAccessTokenURL, body)
	if err != nil {
		return nil, err
	}
	updated, err := p.parseTokenResponse(resp)
	if err != nil {
		return nil, err
	}
	// Preserve original credentials info
	updated.ClientID = data.ClientID
	updated.Source = data.Source
	updated.PersistentCode = data.PersistentCode
	updated.CorpID = data.CorpID
	updated.UserID = data.UserID
	updated.UserName = data.UserName
	updated.CorpName = data.CorpName

	if err := SaveTokenData(p.configDir, updated); err != nil {
		return nil, fmt.Errorf("保存刷新后的 token 失败（旧 refresh_token 已失效，请重新登录）: %w", err)
	}
	return updated, nil
}

// refreshViaMCP refreshes token via MCP proxy.
func (p *OAuthProvider) refreshViaMCP(ctx context.Context, data *TokenData) (*TokenData, error) {
	// Use stored clientId from token data
	clientID := data.ClientID
	if clientID == "" {
		// Fallback for legacy tokens
		clientID = ClientID()
	}

	if clientID == "" {
		return nil, fmt.Errorf("无法刷新 token: 缺少 clientId，请重新登录")
	}

	url := GetMCPBaseURL() + MCPRefreshTokenPath
	body := map[string]string{
		"clientId":     clientID,
		"refreshToken": data.RefreshToken,
		"grantType":    "refresh_token",
	}
	resp, err := p.postJSON(ctx, url, body)
	if err != nil {
		return nil, err
	}
	updated, err := p.parseMCPTokenResponse(resp)
	if err != nil {
		return nil, err
	}
	// Preserve original credentials info
	updated.ClientID = data.ClientID
	updated.Source = data.Source
	updated.PersistentCode = data.PersistentCode
	updated.CorpID = data.CorpID
	updated.UserID = data.UserID
	updated.UserName = data.UserName
	updated.CorpName = data.CorpName

	if err := SaveTokenData(p.configDir, updated); err != nil {
		return nil, fmt.Errorf("保存刷新后的 token 失败（旧 refresh_token 已失效，请重新登录）: %w", err)
	}
	return updated, nil
}

func (p *OAuthProvider) postJSON(ctx context.Context, endpoint string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := p.httpClient
	if client == nil {
		client = oauthHTTPClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateBody(data, 200))
	}
	return data, nil
}

func (p *OAuthProvider) parseTokenResponse(body []byte) (*TokenData, error) {
	var resp struct {
		AccessToken    string `json:"accessToken"`
		RefreshToken   string `json:"refreshToken"`
		PersistentCode string `json:"persistentCode"`
		ExpiresIn      int64  `json:"expiresIn"`
		CorpID         string `json:"corpId"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	if resp.AccessToken == "" {
		return nil, fmt.Errorf("token response missing accessToken")
	}

	now := time.Now()
	expiresIn := resp.ExpiresIn
	if expiresIn <= 0 {
		// 默认 2 小时有效期（钉钉 access_token 标准有效期）
		expiresIn = config.DefaultAccessTokenExpiry
	}
	data := &TokenData{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    now.Add(time.Duration(expiresIn) * time.Second),
		RefreshExpAt: now.Add(config.DefaultRefreshTokenLifetime),
		CorpID:       resp.CorpID,
	}
	if resp.PersistentCode != "" {
		data.PersistentCode = resp.PersistentCode
	}
	return data, nil
}

// parseMCPTokenResponse parses token response from MCP proxy.
// MCP OAuth response format: {"accessToken": "...", "refreshToken": "...", "expiresIn": 7200, "corpId": "..."}
func (p *OAuthProvider) parseMCPTokenResponse(body []byte) (*TokenData, error) {
	var resp struct {
		AccessToken    string `json:"accessToken"`
		RefreshToken   string `json:"refreshToken"`
		PersistentCode string `json:"persistentCode"`
		ExpiresIn      int64  `json:"expiresIn"`
		CorpID         string `json:"corpId"`
		// Error fields (when request fails)
		ErrorCode string `json:"errorCode,omitempty"`
		ErrorMsg  string `json:"errorMsg,omitempty"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing MCP token response: %w (body: %s)", err, string(body))
	}
	// Check for error response
	if resp.ErrorCode != "" || resp.ErrorMsg != "" {
		return nil, fmt.Errorf("MCP token exchange failed: %s - %s", resp.ErrorCode, resp.ErrorMsg)
	}
	if resp.AccessToken == "" {
		return nil, fmt.Errorf("MCP token response missing accessToken (body: %s)", string(body))
	}

	now := time.Now()
	expiresIn := resp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = config.DefaultAccessTokenExpiry
	}
	data := &TokenData{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    now.Add(time.Duration(expiresIn) * time.Second),
		RefreshExpAt: now.Add(config.DefaultRefreshTokenLifetime),
		CorpID:       resp.CorpID,
	}
	if resp.PersistentCode != "" {
		data.PersistentCode = resp.PersistentCode
	}
	return data, nil
}

func buildAuthURL(redirectURI string) string {
	params := url.Values{
		"redirect_uri":  {redirectURI},
		//"response_type": {"code"},
		//"scope":         {DefaultScopes},
		//"prompt":        {"consent"},
	}
	return AuthorizeURL + "?" + params.Encode()
}

const successHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>授客 CLI</title>
    <script>
      if (window.location.pathname !== "/success") {
        setTimeout(function() {
          window.location.href = "/success";
        }, 500);
      }
    </script>
    <style>
      body {
        font-family:
          -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          "Helvetica Neue", Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        margin: 0;
        background: #f5f5f5;
        padding: 20px;
      }
      .card {
        height: 600px;
        width: 480px;
        border-radius: 16px;
        background: #ffffff;
        box-sizing: border-box;
        border: 1px solid #f2f2f6;
        box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.12);
        padding: 32px 24px 24px;
        text-align: center;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
      }
      .lock-icon {
        width: 120px;
        height: 120px;
        margin: 0 auto;
        object-fit: contain;
        display: block;
      }
      h1 {
        margin: 8px 0 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 18px;
        font-weight: 600;
        line-height: 44px;
        text-align: center;
        letter-spacing: normal;
        color: #181c1f;
      }
      p {
        margin: 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 14px;
        font-weight: normal;
        line-height: 21px;
        text-align: center;
        letter-spacing: normal;
        color: rgba(24, 28, 31, 0.6);
      }
    </style>
  </head>
  <body>
    <div class="card">
      <img
        class="lock-icon"
        src="https://img.alicdn.com/imgextra/i4/O1CN01fS3xxz1vbzZSGjbe0_!!6000000006192-2-tps-480-480.png"
        alt="lock icon"
      />
      <h1>授权成功</h1>
      <p>请返回终端继续操作。此页面可以关闭。</p>
    </div>
  </body>
</html>`

const notEnabledHTML = `<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>授客 CLI</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }
      body {
        font-family:
          -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          "Helvetica Neue", Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        background: #fff;
        padding: 20px;
      }
      .container {
        text-align: center;
        width: 100%;
        max-width: 400px;
        border-radius: 16px;
        background: #ffffff;
        border: 1px solid #f2f2f6;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.12);
        padding: 32px 24px 24px;
      }
      .lock-icon {
        width: 120px;
        height: 120px;
        margin: 0 auto;
        object-fit: contain;
        display: block;
      }
      h1 {
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 18px;
        font-weight: 600;
        line-height: 44px;
        text-align: center;
        color: #181c1f;
      }
      p {
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 21px;
        text-align: center;
        color: rgba(24, 28, 31, 0.6);
        margin-bottom: 24px;
      }
      .form-group {
        text-align: left;
        margin-bottom: 24px;
      }
      .form-label {
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 14px;
        color: rgba(24, 28, 31, 0.6);
        margin-top: 38px;
        margin-bottom: 8px;
        display: block;
      }
      .select-wrapper {
        position: relative;
      }
      .custom-select-trigger {
        width: 100%;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 16px;
        gap: 12px;
        border-radius: 8px;
        background: #ffffff;
        border: 1px solid rgba(126, 134, 142, 0.24);
        cursor: pointer;
        color: #181c1f;
        text-align: left;
      }
      .custom-select-text {
        flex: 1;
        font-size: 14px;
        line-height: 24px;
        color: rgba(24, 28, 31, 0.48);
      }
      .custom-select.has-value .custom-select-text {
        color: #181c1f;
      }
      .custom-select-arrow {
        width: 16px;
        height: 16px;
        flex-shrink: 0;
        background-image: url("https://img.alicdn.com/imgextra/i1/O1CN01MzGSB21oZ3iyQ8H5e_!!6000000005238-55-tps-16-16.svg");
        background-repeat: no-repeat;
        background-size: 16px 16px;
        background-position: center;
        opacity: 0.5;
      }
      .custom-select-menu {
        position: absolute;
        left: 0;
        right: 0;
        top: calc(100% + 8px);
        background: #ffffff;
        border-radius: 10px;
        padding: 8px 12px;
        list-style: none;
        margin: 0;
        box-shadow: 0 6px 18px rgba(0, 0, 0, 0.12);
        display: none;
        z-index: 20;
      }
      .custom-select.open .custom-select-menu {
        display: block;
      }
      .custom-select-option {
        width: 100%;
        height: 40px;
        border: none;
        background: transparent;
        text-align: left;
        padding: 8px 12px;
        border-radius: 8px;
        font-size: 14px;
        line-height: 24px;
        color: #181c1f;
        cursor: pointer;
      }
      .custom-select-option:hover {
        background: rgba(126, 134, 142, 0.16);
      }
      .custom-select-option.is-active {
        background: #e8eaee;
      }
      .btn {
        width: 100%;
        height: 40px;
        border-radius: 1000px;
        background: #0066ff;
        border: none;
        cursor: pointer;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 16px;
        font-weight: 500;
        line-height: 24px;
        color: #ffffff;
      }
      .btn:hover {
        background: #0b6cff;
      }
      .btn:disabled {
        background: #d9d9d9;
        cursor: not-allowed;
      }
      .link {
        color: #1890ff;
        font-size: 14px;
        text-decoration: none;
        margin-top: 16px;
        display: inline-block;
      }
      .success-msg {
        display: none;
        width: 100%;
        height: 60px;
        gap: 12px;
        padding: 16px 20px;
        margin-top: 50px;
        margin-bottom: 16px;
        background: #eaf1ff;
        border-radius: 12px;
        align-items: center;
      }
      .success-msg-icon {
        width: 24px;
        height: 24px;
        flex-shrink: 0;
      }
      .success-msg-text {
        flex: 1;
        font-size: 14px;
        line-height: 22px;
        color: #181c1f;
      }
      .error-msg {
        color: #ff4d4f;
        font-size: 14px;
        margin-top: 8px;
        display: none;
      }
      .loading {
        display: inline-block;
        width: 16px;
        height: 16px;
        border: 2px solid #fff;
        border-top-color: transparent;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
        margin-right: 8px;
        vertical-align: middle;
      }
      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }
    </style>
  </head>
  <body>
    <div class="container">
      <img
        class="lock-icon"
        src="https://img.alicdn.com/imgextra/i4/O1CN01fS3xxz1vbzZSGjbe0_!!6000000006192-2-tps-480-480.png"
        alt="lock icon"
      />
      <h1>该组织尚未开启 CLI 数据访问权限</h1>
      <p>
        你所选择的组织管理员尚未开启<br />「允许成员通过 CLI
        访问其个人数据」的权限。
      </p>

      <div class="form-group">
        <label class="form-label">选择一位主管理员发送开通申请</label>
        <div class="select-wrapper">
          <div class="custom-select" id="adminSelect">
            <button
              type="button"
              class="custom-select-trigger"
              aria-haspopup="listbox"
              aria-expanded="false"
            >
              <span class="custom-select-text">加载中...</span>
              <span class="custom-select-arrow"></span>
            </button>
            <ul class="custom-select-menu" role="listbox" id="adminMenu"></ul>
            <input type="hidden" name="adminStaffId" value="" />
          </div>
        </div>
        <div id="errorMsg" class="error-msg"></div>
      </div>

      <div id="successMsg" class="success-msg">
        <svg
          class="success-msg-icon"
          viewBox="0 0 16 16"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M8 1.33333C4.32 1.33333 1.33333 4.32 1.33333 8C1.33333 11.68 4.32 14.6667 8 14.6667C11.68 14.6667 14.6667 11.68 14.6667 8C14.6667 4.32 11.68 1.33333 8 1.33333ZM8 13.3333C5.05333 13.3333 2.66667 10.9467 2.66667 8C2.66667 5.05333 5.05333 2.66667 8 2.66667C10.9467 2.66667 13.3333 5.05333 13.3333 8C13.3333 10.9467 10.9467 13.3333 8 13.3333ZM7.33333 9.33333H8.66667V10.6667H7.33333V9.33333ZM7.33333 5.33333H8.66667V8H7.33333V5.33333Z"
            fill="#0066FF"
          />
        </svg>
        <span class="success-msg-text"
          >已向管理员发送权限申请，正在等待审核<br />审核通过后，请返回终端继续操作</span
        >
      </div>

      <button id="applyBtn" class="btn" disabled>立即申请</button>
      <a id="backLink" class="link" href="#">返回选择其他组织</a>
    </div>

    <script>
      const adminSelect = document.getElementById("adminSelect");
      const trigger = adminSelect.querySelector(".custom-select-trigger");
      const text = adminSelect.querySelector(".custom-select-text");
      const menu = document.getElementById("adminMenu");
      const hiddenInput = adminSelect.querySelector('input[name="adminStaffId"]');
      const btn = document.getElementById("applyBtn");
      const successMsg = document.getElementById("successMsg");
      const errorMsg = document.getElementById("errorMsg");
      const backLink = document.getElementById("backLink");

      let admins = [];
      let clientId = "";
      let applySent = false;
      let selectedAdminId = "";
      let pollTimer = null;

      function closeMenu() {
        adminSelect.classList.remove("open");
        trigger.setAttribute("aria-expanded", "false");
      }

      function openMenu() {
        adminSelect.classList.add("open");
        trigger.setAttribute("aria-expanded", "true");
      }

      function showError(msg) {
        errorMsg.textContent = msg;
        errorMsg.style.display = "block";
      }

      function hideError() {
        errorMsg.style.display = "none";
      }

      function setSelected(staffId, name) {
        hiddenInput.value = staffId || "";
        text.textContent = name || "请选择";
        if (staffId) {
          adminSelect.classList.add("has-value");
        } else {
          adminSelect.classList.remove("has-value");
        }
        btn.disabled = applySent || !staffId;
      }

      function renderAdminOptions(list) {
        menu.innerHTML = "";
        list.forEach(function (admin) {
          const li = document.createElement("li");
          const option = document.createElement("button");
          option.type = "button";
          option.className = "custom-select-option";
          option.setAttribute("data-value", admin.staffId);
          option.textContent = admin.name;

          if (selectedAdminId && admin.staffId === selectedAdminId) {
            option.classList.add("is-active");
          }

          option.addEventListener("click", function () {
            selectedAdminId = admin.staffId;
            const all = menu.querySelectorAll(".custom-select-option");
            all.forEach(function (item) {
              item.classList.remove("is-active");
            });
            option.classList.add("is-active");
            setSelected(admin.staffId, admin.name);
            closeMenu();
            hideError();
          });

          li.appendChild(option);
          menu.appendChild(li);
        });
      }

      function setAppliedState() {
        btn.disabled = true;
        btn.textContent = "立即申请";
        trigger.disabled = true;
        adminSelect.classList.remove("open");
        successMsg.style.display = "flex";
        backLink.style.pointerEvents = "none";
        backLink.style.color = "#999";
        backLink.onclick = function (e) {
          e.preventDefault();
          return false;
        };
        startPolling();
      }

      function startPolling() {
        if (pollTimer) return;
        pollTimer = setInterval(checkAuthStatus, 5000);
        checkAuthStatus();
      }

      function stopPolling() {
        if (pollTimer) {
          clearInterval(pollTimer);
          pollTimer = null;
        }
      }

      async function checkAuthStatus() {
        try {
          const res = await fetch("/api/cliAuthEnabled");
          const data = await res.json();
          if (data.success && data.result && data.result.cliAuthEnabled) {
            stopPolling();
            location.href = "/success";
          }
        } catch (e) {
          console.error("Poll error", e);
        }
      }

      async function loadAdmins() {
        try {
          const res = await fetch("/api/superAdmin");
          const data = await res.json();
          if (data.success && data.result && data.result.length > 0) {
            admins = data.result;
            renderAdminOptions(admins);

            if (selectedAdminId) {
              const selected = admins.find(function (a) {
                return a.staffId === selectedAdminId;
              });
              if (selected) {
                setSelected(selected.staffId, selected.name);
              } else {
                setSelected("", "请选择");
              }
            } else {
              setSelected("", "请选择");
            }
          } else {
            setSelected("", "暂无可选管理员");
            trigger.disabled = true;
            showError((data && data.errorMsg) || "获取管理员列表失败");
          }
        } catch (e) {
          setSelected("", "加载失败");
          trigger.disabled = true;
          showError("网络错误，请重试");
        }
      }

      async function init() {
        try {
          const statusRes = await fetch("/api/status");
          const status = await statusRes.json();
          clientId = status.clientId || "";
          applySent = status.applySent || false;
          selectedAdminId = status.selectedAdminId || "";

          if (clientId) {
            const port = location.port;
            const redirectUri = encodeURIComponent(
              "http://127.0.0.1:" + port + "/callback"
            );
            backLink.href =
              "https://login.dingtalk.com/oauth2/auth?client_id=" +
              clientId +
              "&prompt=consent&redirect_uri=" +
              redirectUri +
              "&response_type=code&scope=openid+corpid";
          }

          if (applySent) {
            setAppliedState();
          }
        } catch (e) {
          console.error("Failed to load status", e);
        }

        await loadAdmins();
      }

      trigger.addEventListener("click", function () {
        if (trigger.disabled) return;
        if (adminSelect.classList.contains("open")) {
          closeMenu();
        } else {
          openMenu();
        }
      });

      document.addEventListener("click", function (event) {
        if (!adminSelect.contains(event.target)) {
          closeMenu();
        }
      });

      btn.onclick = async function () {
        const value = hiddenInput.value;
        if (!value) return;

        btn.disabled = true;
        btn.innerHTML = '<span class="loading"></span>申请中...';
        hideError();
        try {
          const res = await fetch(
            "/api/sendApply?adminStaffId=" + encodeURIComponent(value)
          );
          const data = await res.json();
          if (data.success && data.result) {
            setAppliedState();
          } else {
            showError(data.errorMsg || "申请失败，请重试");
            btn.disabled = false;
            btn.textContent = "立即申请";
          }
        } catch (e) {
          showError("网络错误，请重试");
          btn.disabled = false;
          btn.textContent = "立即申请";
        }
      };

      init();
    </script>
  </body>
</html>`

const accessDeniedHTML = `<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>授客 CLI</title>
    <script>
    if (window.location.pathname !== "/fail") {
      setTimeout(function() {
        window.location.href = "/fail";
      }, 500);
    }
  </script>
    <style>
      body {
        font-family:
          -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          "Helvetica Neue", Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        margin: 0;
        background: #f5f5f5;
        padding: 20px;
      }
      .card {
        height: 600px;
        width: 480px;
        border-radius: 16px;
        background: #ffffff;
        box-sizing: border-box;
        border: 1px solid #f2f2f6;
        box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.12);
        padding: 32px 24px 24px;
        text-align: center;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
      }
      .lock-icon {
        width: 120px;
        height: 120px;
        margin: 0 auto;
        object-fit: contain;
        display: block;
      }
      h1 {
        margin: 8px 0 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 18px;
        font-weight: 600;
        line-height: 44px;
        text-align: center;
        letter-spacing: normal;
        color: #181c1f;
      }
      p {
        margin: 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 14px;
        font-weight: normal;
        line-height: 21px;
        text-align: center;
        letter-spacing: normal;
        color: rgba(24, 28, 31, 0.6);
      }
    </style>
  </head>
  <body>
    <div class="card">
      <img
        class="lock-icon"
        src="https://img.alicdn.com/imgextra/i4/O1CN01fS3xxz1vbzZSGjbe0_!!6000000006192-2-tps-480-480.png"
        alt="lock icon"
      />
      <h1>无权限访问</h1>
      <p>您不在该组织的 CLI 授权人员范围内。请联系组织管理员将您加入授权名单。此页面可以关闭。</p>
    </div>
  </body>
</html>`

const channelDeniedHTML = `<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>钉钉 CLI</title>
    <style>
      body {
        font-family:
          -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          "Helvetica Neue", Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        margin: 0;
        background: #f5f5f5;
        padding: 20px;
      }
      .card {
        height: 600px;
        width: 480px;
        border-radius: 16px;
        background: #ffffff;
        box-sizing: border-box;
        border: 1px solid #f2f2f6;
        box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.12);
        padding: 32px 24px 24px;
        text-align: center;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
      }
      .lock-icon {
        width: 120px;
        height: 120px;
        margin: 0 auto;
        object-fit: contain;
        display: block;
      }
      h1 {
        margin: 8px 0 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 18px;
        font-weight: 600;
        line-height: 44px;
        text-align: center;
        letter-spacing: normal;
        color: #181c1f;
      }
      p {
        margin: 0;
        font-family:
          "PingFang SC",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          "Helvetica Neue",
          Arial,
          sans-serif;
        font-size: 14px;
        font-weight: normal;
        line-height: 21px;
        text-align: center;
        letter-spacing: normal;
        color: rgba(24, 28, 31, 0.6);
      }
    </style>
  </head>
  <body>
    <div class="card">
      <img
        class="lock-icon"
        src="https://img.alicdn.com/imgextra/i4/O1CN01fS3xxz1vbzZSGjbe0_!!6000000006192-2-tps-480-480.png"
        alt="lock icon"
      />
      <h1>渠道未授权</h1>
      <p>当前渠道未获得该组织授权，或组织已开启渠道管控。请联系组织管理员开通渠道访问权限，或升级到最新版本的 CLI。此页面可以关闭。</p>
    </div>
  </body>
</html>`

// CLIAuthStatus represents the response from /cli/cliAuthEnabled API.
type CLIAuthStatus struct {
	Success   bool           `json:"success"`
	ErrorCode string         `json:"errorCode,omitempty"`
	ErrorMsg  string         `json:"errorMsg,omitempty"`
	Result    *CLIAuthResult `json:"result"`
}

// CLIAuthResult holds the business data returned by /cli/cliAuthEnabled.
// The server computes cliAuthEnabled by considering the org switch, userScope,
// and channelScope together; the CLI uses it as-is.
type CLIAuthResult struct {
	CLIAuthEnabled       bool     `json:"cliAuthEnabled"`
	UserScope            string   `json:"userScope,omitempty"`            // "all" | "specified" | "forbidden"
	AllowedUsers         []string `json:"allowedUsers,omitempty"`         // staffId list when userScope="specified"
	ChannelScope         string   `json:"channelScope,omitempty"`         // "all" | "specified"
	AllowedChannels      []string `json:"allowedChannels,omitempty"`      // channelCode list when channelScope="specified"
	ChannelConfigEnabled bool     `json:"channelConfigEnabled,omitempty"` // whether org has any channel restriction configured
}

// classifyDenialReason inspects a CLIAuthStatus response and returns a machine-readable
// denial reason string. Returns "" when access is granted.
//
// Priority rationale:
//  1. Explicit org-wide ban (userScope=forbidden) always wins.
//  2. Channel scope is evaluated BEFORE user scope because the CLI has
//     authoritative knowledge of DWS_CHANNEL and can verify membership against
//     allowedChannels. This avoids falsely blaming the user when the real
//     denial cause is a channel mismatch (e.g. user is in allowedUsers but the
//     current channel is not in allowedChannels).
//  3. Only when the channel is unrestricted or matches do we attribute the
//     denial to the user scope.
func classifyDenialReason(status *CLIAuthStatus, currentChannel string) string {
	if status.ErrorCode == "CHANNEL_REQUIRED" {
		return "channel_required"
	}
	if status.ErrorCode == "NO_AUTH" {
		return "no_auth"
	}
	if status.Result == nil || !status.Success {
		return "unknown"
	}
	r := status.Result
	if r.CLIAuthEnabled {
		return ""
	}

	if r.UserScope == "forbidden" {
		return "user_forbidden"
	}

	if r.ChannelScope == "specified" {
		if currentChannel == "" {
			return "channel_required"
		}
		if !slices.Contains(r.AllowedChannels, currentChannel) {
			return "channel_not_allowed"
		}
	}

	if r.UserScope == "specified" {
		return "user_not_allowed"
	}
	return "cli_not_enabled"
}

// SuperAdmin represents a corp super admin.
type SuperAdmin struct {
	StaffID string `json:"staffId"`
	Name    string `json:"name"`
}

// SuperAdminResponse represents the response from /cli/superAdmin API.
type SuperAdminResponse struct {
	Success   bool         `json:"success"`
	ErrorCode string       `json:"errorCode,omitempty"`
	ErrorMsg  string       `json:"errorMsg,omitempty"`
	Result    []SuperAdmin `json:"result"`
}

// SendApplyResponse represents the response from /cli/sendCliAuthApply API.
type SendApplyResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode,omitempty"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
	Result    bool   `json:"result"`
}

// mcpRequestMaxRetries is the maximum number of attempts for MCP API calls
// (e.g. /cli/cliAuthEnabled, /cli/clientId, /cli/superAdmin, /cli/sendCliAuthApply)
// to tolerate transient network errors before propagating the failure.
const mcpRequestMaxRetries = 3

// CheckCLIAuthEnabled checks if CLI authorization is enabled for the current corp.
// It retries up to mcpRequestMaxRetries times on transient errors to avoid
// false negatives caused by momentary network issues.
func (p *OAuthProvider) CheckCLIAuthEnabled(ctx context.Context, accessToken string) (*CLIAuthStatus, error) {
	var lastErr error
	for attempt := 0; attempt < mcpRequestMaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		status, err := p.doCheckCLIAuthEnabled(ctx, accessToken)
		if err == nil {
			return status, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("check CLI auth status failed after %d attempts: %w", mcpRequestMaxRetries, lastErr)
}

func (p *OAuthProvider) doCheckCLIAuthEnabled(ctx context.Context, accessToken string) (*CLIAuthStatus, error) {
	url := GetMCPBaseURL() + CLIAuthEnabledPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("x-user-access-token", accessToken)
	if ch := os.Getenv("DWS_CHANNEL"); ch != "" {
		req.Header.Set("x-dws-channel", ch)
	}

	client := p.httpClient
	if client == nil {
		client = oauthHTTPClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var status CLIAuthStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &status, nil
}

// GetSuperAdmins fetches the list of corp super admins.
// It retries up to mcpRequestMaxRetries times on transient errors.
func GetSuperAdmins(ctx context.Context, accessToken string) (*SuperAdminResponse, error) {
	var lastErr error
	for attempt := 0; attempt < mcpRequestMaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		result, err := doGetSuperAdmins(ctx, accessToken)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("get super admins failed after %d attempts: %w", mcpRequestMaxRetries, lastErr)
}

func doGetSuperAdmins(ctx context.Context, accessToken string) (*SuperAdminResponse, error) {
	url := GetMCPBaseURL() + SuperAdminPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("x-user-access-token", accessToken)

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result SuperAdminResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &result, nil
}

// SendCliAuthApply sends a CLI auth apply request to the specified admin.
// It retries up to mcpRequestMaxRetries times on transient errors.
func SendCliAuthApply(ctx context.Context, accessToken, adminStaffID string) (*SendApplyResponse, error) {
	var lastErr error
	for attempt := 0; attempt < mcpRequestMaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		result, err := doSendCliAuthApply(ctx, accessToken, adminStaffID)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("send CLI auth apply failed after %d attempts: %w", mcpRequestMaxRetries, lastErr)
}

func doSendCliAuthApply(ctx context.Context, accessToken, adminStaffID string) (*SendApplyResponse, error) {
	url := GetMCPBaseURL() + SendCliAuthApplyPath + "?adminStaffId=" + adminStaffID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("x-user-access-token", accessToken)

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result SendApplyResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &result, nil
}

// ClientIDResponse represents the response from /cli/clientId API.
type ClientIDResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode,omitempty"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
	Result    string `json:"result"`
}

// FetchClientIDFromMCP fetches the CLI client ID from MCP server.
// This is used when no client ID is provided via flags, config, or env vars.
// It retries up to mcpRequestMaxRetries times on transient errors.
func FetchClientIDFromMCP(ctx context.Context) (string, error) {
	var lastErr error
	for attempt := 0; attempt < mcpRequestMaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		id, err := doFetchClientIDFromMCP(ctx)
		if err == nil {
			return id, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("fetch client ID failed after %d attempts: %w", mcpRequestMaxRetries, lastErr)
}

func doFetchClientIDFromMCP(ctx context.Context) (string, error) {
	url := GetMCPBaseURL() + ClientIDPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBodySize))
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var result ClientIDResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	if !result.Success {
		return "", fmt.Errorf("%s: %s", result.ErrorCode, result.ErrorMsg)
	}
	return result.Result, nil
}
