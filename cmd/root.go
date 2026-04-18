package cmd
import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
)

var rootCmd = &cobra.Command{
    Use:   "soke-cli",
    Short: "授客AI官方CLI工具",
    Long: `授客AI CLI - 命令行工具
	   使用示例:
	   soke-cli auth login
	   soke-cli api GET /users/me
	   soke-cli calendar list`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {   
    // 添加子命令
    rootCmd.AddCommand(newAuthCmd())
    rootCmd.AddCommand(newAPICmd())
    rootCmd.AddCommand(newConfigCmd())
}                                                                                                                                                                                                                                              
