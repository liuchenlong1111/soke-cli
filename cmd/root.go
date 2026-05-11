package cmd

import (
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 在每次命令执行前检查更新（同步执行，使用缓存快速返回）
		version.CheckForUpdates()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(auth.NewAuthCmd())
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
