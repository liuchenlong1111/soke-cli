package config                                                                                                                                                                                                                                 
                  
import (                                                                                                                                                                                                                                       
  "fmt"       
  "github.com/spf13/cobra"                                                                                                                                                                                                                   
  "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"                                                                                                                                                                                             
)                                                                                                                                                                                                                                              
                  
func NewConfigCmd() *cobra.Command {                                                                                                                                                                                                           
    cmd := &cobra.Command{
        Use:   "config",                                                                                                                                                                                                                       
        Short: "管理配置",                                                                                                                                                                                                                     
    }                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                               
    cmd.AddCommand(newInitCmd())                                                                                                                                                                                                               
    cmd.AddCommand(newShowCmd())                                                                                                                                                                                                               
                                                                                                                                                                                                                                               
    return cmd                                                                                                                                                                                                                                 
}                                                                                                                                                                                                                                              
                                                                                                                                                                                                                                                 
func newInitCmd() *cobra.Command {                                                                                                                                                                                                             
    return &cobra.Command{
        Use:   "init",                                                                                                                                                                                                                         
        Short: "初始化配置",                                                                                                                                                                                                                   
        RunE: func(cmd *cobra.Command, args []string) error {                                                                                                                                                                                  
            // 交互式配置                                                                                                                                                                                                                      
            var cfg core.CliConfig                                                                                                                                                                                                             
                                                                                                                                                                                                                                               
            fmt.Print("请输入app_key: ")                                                                                                                                                                                                        
            fmt.Scanln(&cfg.AppID)                                                                                                                                                                                                             
                                                                                                                                                                                                                                               
            fmt.Print("请输入app_secret: ")                                                                                                                                                                                                    
            fmt.Scanln(&cfg.AppSecret)                                                                                                                                                                                                         
                                                                                                                                                                                                                                               
            fmt.Print("请输入API地址 (默认: https://opendev.soke.cn): ")
            fmt.Scanln(&cfg.APIBaseURL)
            if cfg.APIBaseURL == "" {
                cfg.APIBaseURL = "https://opendev.soke.cn"
            }

            fmt.Print("请输入corpid: ")
            fmt.Scanln(&cfg.CorpID)

            fmt.Print("请输入dept_user_id: ")
            fmt.Scanln(&cfg.DeptUserID)

            if err := core.SaveConfig(&cfg); err != nil {                                                                                                                                                                                      
                return err
            }                                                                                                                                                                                                                                  
                
            fmt.Println("✓ 配置已保存")                                                                                                                                                                                                        
            return nil
        },                                                                                                                                                                                                                                     
    }                                                                                                                                                                                                                                          
}
                                                                                                                                                                                                                                                 
func newShowCmd() *cobra.Command {                                                                                                                                                                                                             
    return &cobra.Command{
        Use:   "show",                                                                                                                                                                                                                         
        Short: "显示当前配置",                                                                                                                                                                                                                 
        RunE: func(cmd *cobra.Command, args []string) error {                                                                                                                                                                                  
            cfg, err := core.LoadConfig()                                                                                                                                                                                                      
            if err != nil {                                                                                                                                                                                                                    
                return err                                                                                                                                                                                                                     
            }                                                                                                                                                                                                                                  
                                                                                                                                                                                                                                               
            fmt.Printf("app_key: %s\n", cfg.AppID)
            fmt.Printf("API地址: %s\n", cfg.APIBaseURL)
            fmt.Printf("corpid: %s\n", cfg.CorpID)
            fmt.Printf("dept_user_id: %s\n", cfg.DeptUserID)                                                                                                                                                                            
                                                                                                                                                                                                                                               
            return nil                                                                                                                                                                                                                         
        },                                                                                                                                                                                                                                     
    }                                                                                                                                                                                                                                          
}                                                                                                                                                                                                                                              
   
func maskToken(token string) string {                                                                                                                                                                                                          
    if token == "" {
        return "(未设置)"                                                                                                                                                                                                                      
    }                                                                                                                                                                                                                                          
    if len(token) < 10 {                                                                                                                                                                                                                       
        return "***"                                                                                                                                                                                                                           
    }                                                                                                                                                                                                                                          
    return token[:4] + "..." + token[len(token)-4:]                                                                                                                                                                                            
}                                                                                                                                                                                                                                              
                