package cmd

import (
	"fmt"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd 创建 version 命令
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示 soke-cli 版本信息",
		Long:  `显示当前安装的 soke-cli 版本号`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("soke-cli 版本: %s\n", version.GetVersion())
		},
	}
}
