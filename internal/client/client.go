package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
)

type Client struct {
    Config       *core.CliConfig
    HTTPClient   *http.Client
    TokenManager *auth.TokenManager
}

func NewClient(cfg *core.CliConfig) *Client {
  tokenManager := auth.NewTokenManager(cfg.AppKey, cfg.AppSecret, cfg.CorpID, cfg.APIBaseURL)
  return &Client{
      Config:       cfg,
      TokenManager: tokenManager,
      HTTPClient: &http.Client{
          Timeout: 30 * time.Second,
      },
  }
}

// DoRequest 执行API请求
// 1. 优先从内存缓存中获取 access_token
// 2. 如果缓存中没有，使用 app_key、app_secret、corpid 从 API 获取
// 3. 如果配置信息不完整，返回错误提示需要登录授权
func (c *Client) DoRequest(ctx context.Context, req *core.APIRequest) (interface{}, error) {
  // 从缓存或配置文件中获取配置信息
  cfg, err := auth.LoadAndValidateConfig(ctx)
  if err != nil {
      return nil, err
  }

  // 检查配置信息是否完整
  if !auth.IsConfigComplete(cfg) {
      return nil, fmt.Errorf("用户信息不完整，请先执行 soke-cli auth login 完成登录授权")
  }

  // 使用 TokenManager 获取 access_token
  // TokenManager 内部会自动处理缓存和续期
  token, err := c.TokenManager.GetAccessToken(ctx)
  if err != nil {
      return nil, fmt.Errorf("获取 access_token 失败: %w", err)
  }

  // 构建完整URL
  baseURL := cfg.APIBaseURL + req.Path

  // 添加access_token到query参数
  if req.Query == nil {
      req.Query = make(map[string]interface{})
  }
  req.Query["access_token"] = token

  // 构建URL参数
  params := url.Values{}
  for key, value := range req.Query {
      params.Add(key, fmt.Sprintf("%v", value))
  }

  fullURL := baseURL
  if len(params) > 0 {
      fullURL = baseURL + "?" + params.Encode()
  }

  // 构建HTTP请求
  var bodyReader io.Reader
  if req.Body != nil {
      data, err := json.Marshal(req.Body)
      if err != nil {
          return nil, err
      }
      bodyReader = bytes.NewReader(data)
  }

  httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
  if err != nil {
      return nil, err
  }

  httpReq.Header.Set("Content-Type", "application/json")

  // 发送请求
  resp, err := c.HTTPClient.Do(httpReq)
  if err != nil {
      return nil, err
  }
  defer resp.Body.Close()

  // 解析响应
  body, err := io.ReadAll(resp.Body)
  if err != nil {
      return nil, err
  }

  if resp.StatusCode >= 400 {
      return nil, fmt.Errorf("API error: %d %s", resp.StatusCode, string(body))
  }

  var result interface{}
  if err := json.Unmarshal(body, &result); err != nil {
      return nil, err
  }

  return result, nil
}