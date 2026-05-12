package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	"github.com/spf13/cobra"
)                                                                                                                                                                                                                                              
                                                                                                                                                                                                                                                 
func NewAuthCmd() *cobra.Command {                                                                                                                                                                                                             
    cmd := &cobra.Command{                                                                                                                                                                                                                     
        Use:   "auth",                                                                                                                                                                                                                         
        Short: "认证管理",                                                                                                                                                                                                                     
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    cmd.AddCommand(newLoginCmd())                                                                                                                                                                                                              
    cmd.AddCommand(newLogoutCmd())
    return cmd                                                                                                                                                                                                                                 
}                                                                                                                                                                                                                                              
                                                                                                                                                                                                                                                 
func newLoginCmd() *cobra.Command {                                                                                                                                                                                                            
    var scopes []string
                                                                                                                                                                                                                                               
    cmd := &cobra.Command{                                                                                                                                                                                                                     
        Use:   "login",                                                                                                                                                                                                                        
        Short: "用户登录",                                                                                                                                                                                                                     
        RunE: func(cmd *cobra.Command, args []string) error {                                                                                                                                                                                  
            cfg, err := core.LoadConfig()                                                                                                                                                                                                      
            if err != nil {                                                                                                                                                                                                                    
                return fmt.Errorf("请先运行 'soke-cli config init'")                                                                                                                                                                         
            }                                                                                                                                                                     
                                                                                                                                                                                                                                               
            // 创建OAuth客户端                                                                                                                                                                                                                 
            oauthClient := &auth.OAuthClient{
                AppID:      cfg.AppID,                                                                                                                                                                                                         
                AppSecret:  cfg.AppSecret,                                                                                                                                                                                                     
                AuthURL:    cfg.APIBaseURL + "/oauth/device/code",                                                                                                                                                                             
                TokenURL:   cfg.APIBaseURL + "/oauth/token",                                                                                                                                                                                   
                HTTPClient: &http.Client{Timeout: 10 * time.Second},                                                                                                                                                                           
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            // 请求设备授权                                                                                                                                                                                                                    
            ctx := context.Background()                                                                                                                                                                                                        
            authResp, err := oauthClient.RequestDeviceAuth(ctx, scopes)                                                                                                                                                                        
            if err != nil {                                                                                                                                                                                                                    
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            // 显示授权信息                                                                                                                                                                                                                    
            fmt.Println("\n请在浏览器中打开以下链接完成授权:")
            fmt.Printf("\n  %s\n\n", authResp.VerificationURL)                                                                                                                                                                                 
            fmt.Printf("用户码: %s\n\n", authResp.UserCode)                                                                                                                                                                                    
            fmt.Println("等待授权...")                                                                                                                                                                                                         
                                                                                                                                                                                                                                               
            // 轮询获取token                                                                                                                                                                                                                   
            token, err := oauthClient.PollToken(ctx, authResp.DeviceCode, authResp.Interval)                                                                                                                                                   
            if err != nil {                                                                                                                                                                                                                    
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            // 保存token                                                                                                                                                                                                                       
            cfg.UserToken = token
            if err := core.SaveConfig(cfg); err != nil {                                                                                                                                                                                       
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            fmt.Println("\n✓ 登录成功!")                                                                                                                                                                                                       
            return nil
        },                                                                                                                                                                                                                                     
    }           
                                                                                                                                                                                                                                               
    cmd.Flags().StringSliceVar(&scopes, "scope", []string{"read", "write"}, "权限范围")                                                                                                                                                        
 
    return cmd                                                                                                                                                                                                                                 
}               
                                                                                                                                                                                                                                                 
func newLogoutCmd() *cobra.Command {                                                                                                                                                                                                           
    return &cobra.Command{
        Use:   "logout",                                                                                                                                                                                                                       
        Short: "退出登录",                                                                                                                                                                                                                     
        RunE: func(cmd *cobra.Command, args []string) error {                                                                                                                                                                                  
            cfg, err := core.LoadConfig()                                                                                                                                                                                                      
            if err != nil {                                                                                                                                                                                                                    
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            cfg.UserToken = ""                                                                                                                                                                                                                 
            cfg.BotToken = ""
                                                                                                                                                                                                                                               
            if err := core.SaveConfig(cfg); err != nil {                                                                                                                                                                                       
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            fmt.Println("✓ 已退出登录")                                                                                                                                                                                                        
            return nil                                                                                                                                                                                                                         
        },                                                                                                                                                                                                                                     
    }           
}         