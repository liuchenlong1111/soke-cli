package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
)

type Client struct {
    Config     *core.CliConfig
    HTTPClient *http.Client
}

func NewClient(cfg *core.CliConfig) *Client {
  return &Client{
      Config: cfg,
      HTTPClient: &http.Client{
          Timeout: 30 * time.Second,
      },
  }
}

// DoRequest 执行API请求
func (c *Client) DoRequest(ctx context.Context, req *core.APIRequest) (interface{}, error) {
  // 构建完整URL
  url := c.Config.APIBaseURL + req.Path

  // 添加query参数
  if len(req.Query) > 0 {
      // ... URL编码逻辑
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

  httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
  if err != nil {
      return nil, err
  }

  // 设置认证头
  token := c.getToken(req.As)
  httpReq.Header.Set("Authorization", "Bearer "+token)
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


func (c *Client) getToken(as core.Identity) string {
    if as == core.AsBot {
        return c.Config.BotToken
    }
    return c.Config.UserToken
}