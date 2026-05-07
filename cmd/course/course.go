package course

import (
	"context"

	"github.com/spf13/cobra"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
	courseShortcuts "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/course"
)

func NewCourseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "course",
		Short: "课程相关接口",
	}

	shortcuts := courseShortcuts.Shortcuts()
	for _, shortcut := range shortcuts {
		cmd.AddCommand(createShortcutCommand(shortcut))
	}

	return cmd
}

func createShortcutCommand(s common.Shortcut) *cobra.Command {
	cmd := &cobra.Command{
		Use:   s.Command,
		Short: s.Description,
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

			if s.Execute != nil {
				return s.Execute(ctx, runtime)
			}

			return nil
		},
	}

	for _, flag := range s.Flags {
		switch flag.Type {
		case "int":
			cmd.Flags().Int(flag.Name, 0, flag.Desc)
		default:
			cmd.Flags().String(flag.Name, flag.Default, flag.Desc)
		}
		if flag.Required {
			cmd.MarkFlagRequired(flag.Name)
		}
	}

	return cmd
}
