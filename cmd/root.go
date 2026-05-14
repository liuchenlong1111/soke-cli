package cmd

import (
	"context"
	"fmt"
	"os"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/api"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/certificate"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/clock"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/config"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/contact"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/course"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/credit"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/exam"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/file"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/learning_map"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/news"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/point"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/training"
	versionCmd "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/version"
	authpkg "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/errors"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "soke-cli",
	Short:   "授客AI官方CLI工具",
	Version: version.GetVersion(),
	Long: `授客AI CLI - 命令行工具
	   使用示例:
	   soke-cli auth login
	   soke-cli api GET /users/me
	   soke-cli calendar list`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 在每次命令执行前检查更新（同步执行，使用缓存快速返回）
		version.CheckForUpdates()

		// 认证检查：这是整个系统中唯一的认证拦截点
		// 在发起任何 API 请求之前进行检查，提供 fail-fast 机制
		ctx := context.Background()

		// 获取根命令名称（第一级命令）
		commandName := cmd.Name()
		if cmd.Parent() != nil && cmd.Parent().Name() != "soke-cli" {
			// 如果有父命令且父命令不是根命令，使用父命令名称
			commandName = cmd.Parent().Name()
		}

		if err := authpkg.CheckAuth(ctx, commandName); err != nil {
			// 格式化错误输出
			errors.PrintHuman(os.Stderr, err)
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	//rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(auth.NewUserAuthCommand())
	rootCmd.AddCommand(api.NewAPICmd())
	rootCmd.AddCommand(config.NewConfigCmd())
	rootCmd.AddCommand(contact.NewContactCmd())
	rootCmd.AddCommand(course.NewCourseCmd())
	rootCmd.AddCommand(exam.NewExamCmd())
	rootCmd.AddCommand(certificate.NewCertificateCmd())
	rootCmd.AddCommand(credit.NewCreditCmd())
	rootCmd.AddCommand(point.NewPointCmd())
	rootCmd.AddCommand(learning_map.NewLearningMapCmd())
	rootCmd.AddCommand(training.NewTrainingCmd())
	rootCmd.AddCommand(news.NewNewsCmd())
	rootCmd.AddCommand(clock.NewClockCmd())
	rootCmd.AddCommand(file.NewFileCmd())
	rootCmd.AddCommand(versionCmd.NewVersionCmd())
}
