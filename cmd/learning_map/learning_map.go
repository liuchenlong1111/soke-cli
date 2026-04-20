package learning_map

import (
	"context"

	"github.com/spf13/cobra"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	learningMapShortcuts "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/learning_map"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

func NewLearningMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "learning-map",
		Short: "学习地图相关接口",
	}

	shortcuts := learningMapShortcuts.Shortcuts()
	for _, shortcut := range shortcuts {
		cmd.AddCommand(createShortcutCommand(shortcut))
	}

	return cmd
}

func createShortcutCommand(shortcut common.Shortcut) *cobra.Command {
	cmd := &cobra.Command{
		Use:   shortcut.Command,
		Short: shortcut.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := core.LoadConfig()
			if err != nil {
				return err
			}

			ctx := context.Background()
			runtime := &common.RuntimeContext{
				Config: cfg,
				Cmd:    cmd,
			}
			runtime.SetContext(ctx)

			return shortcut.Execute(ctx, runtime)
		},
	}

	for _, flag := range shortcut.Flags {
		switch flag.Type {
		case "int":
			defaultVal := 0
			if flag.Default != "" {
				defaultVal = parseInt(flag.Default)
			}
			cmd.Flags().Int(flag.Name, defaultVal, flag.Desc)
		case "bool":
			defaultVal := false
			if flag.Default == "true" {
				defaultVal = true
			}
			cmd.Flags().Bool(flag.Name, defaultVal, flag.Desc)
		default:
			cmd.Flags().String(flag.Name, flag.Default, flag.Desc)
		}

		if flag.Required {
			cmd.MarkFlagRequired(flag.Name)
		}
	}

	return cmd
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
