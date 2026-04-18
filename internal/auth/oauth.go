package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"                                                                                                                                                                                                                                  
    "strings"
    "time"                                                                                                                                                                                                                                     
)               

type OAuthClient struct {
    AppID        string
    AppSecret    string
    AuthURL      string
    TokenURL     string
    HTTPClient   *http.Client
}

// DeviceAuthResponse 设备授权响应
type DeviceAuthResponse struct {
    DeviceCode      string `json:"device_code"`
    UserCode        string `json:"user_code"`
    VerificationURL string `json:"verification_uri"`
    ExpiresIn       int    `json:"expires_in"`
    Interval        int    `json:"interval"`
}

// RequestDeviceAuth 请求设备授权
func (c *OAuthClient) RequestDeviceAuth(ctx context.Context, scopes []string) (*DeviceAuthResponse, error) {
    form := url.Values{}
    form.Set("client_id", c.AppID)
    form.Set("scope", strings.Join(scopes, " "))

    req, err := http.NewRequestWithContext(ctx, "POST", c.AuthURL, strings.NewReader(form.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
                                                                                                                                                                                                                                               
    resp, err := c.HTTPClient.Do(req)                                                                                                                                                                                                          
    if err != nil {                                                                                                                                                                                                                            
        return nil, err                                                                                                                                                                                                                        
    }                                                                                                                                                                                                                                          
    defer resp.Body.Close()                                                                                                                                                                                                                    
                                                                                                                                                                                                                                               
    body, _ := io.ReadAll(resp.Body)                                                                                                                                                                                                           
                                                                                                                                                                                                                                               
    var result DeviceAuthResponse                                                                                                                                                                                                              
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err                                                                                                                                                                                                                        
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    return &result, nil                                                                                                                                                                                                                        
}               
                                                                                                                                                                                                                                                 
// PollToken 轮询获取token                                                                                                                                                                                                                     
func (c *OAuthClient) PollToken(ctx context.Context, deviceCode string, interval int) (string, error) {
    ticker := time.NewTicker(time.Duration(interval) * time.Second)                                                                                                                                                                            
    defer ticker.Stop()                                                                                                                                                                                                                        
                                                                                                                                                                                                                                               
    for {                                                                                                                                                                                                                                      
        select {                                                                                                                                                                                                                               
        case <-ctx.Done():                                                                                                                                                                                                                     
            return "", ctx.Err()                                                                                                                                                                                                               
        case <-ticker.C:                                                                                                                                                                                                                       
            token, err := c.tryGetToken(ctx, deviceCode)                                                                                                                                                                                       
            if err == nil {                                                                                                                                                                                                                    
                return token, nil                                                                                                                                                                                                              
            }                                                                                                                                                                                                                                  
            // 继续轮询                                                                                                                                                                                                                        
        }                                                                                                                                                                                                                                      
    }                                                                                                                                                                                                                                          
}                                                                                                                                                                                                                                              
                                                                                                                                                                                                                                                 
func (c *OAuthClient) tryGetToken(ctx context.Context, deviceCode string) (string, error) {                                                                                                                                                    
    form := url.Values{}
    form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")                                                                                                                                                                     
    form.Set("device_code", deviceCode)                                                                                                                                                                                                        
    form.Set("client_id", c.AppID)                                                                                                                                                                                                             
    form.Set("client_secret", c.AppSecret)                                                                                                                                                                                                     
                                                                                                                                                                                                                                               
    req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(form.Encode()))                                                                                                                                          
    if err != nil {                                                                                                                                                                                                                            
        return "", err                                                                                                                                                                                                                         
    }           
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")                                                                                                                                                                        
                                                                                                                                                                                                                                               
    resp, err := c.HTTPClient.Do(req)                                                                                                                                                                                                          
    if err != nil {                                                                                                                                                                                                                            
        return "", err                                                                                                                                                                                                                         
    }           
    defer resp.Body.Close()                                                                                                                                                                                                                    
                                                                                                                                                                                                                                               
    body, _ := io.ReadAll(resp.Body)                                                                                                                                                                                                           
                                                                                                                                                                                                                                               
    var result map[string]interface{}                                                                                                                                                                                                          
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err                                                                                                                                                                                                                         
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    if errCode, ok := result["error"].(string); ok {                                                                                                                                                                                           
        if errCode == "authorization_pending" {
            return "", fmt.Errorf("pending")                                                                                                                                                                                                   
        }                                                                                                                                                                                                                                      
        return "", fmt.Errorf("error: %s", errCode)                                                                                                                                                                                            
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    if token, ok := result["access_token"].(string); ok {                                                                                                                                                                                      
        return token, nil                                                                                                                                                                                                                      
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    return "", fmt.Errorf("no token in response")                                                                                                                                                                                              
}    