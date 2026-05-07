package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactGetUser = common.Shortcut{
	Service:     "contact",
	Command:     "+get-user",
	Description: "Get user info (omit user_id for self; provide user_id for specific user)",
	Risk:        "read",
	UserScopes:  []string{"contact:user.basic_profile:readonly"},
	BotScopes:   []string{"contact:user.base:readonly", "contact:contact.base:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "user-id", Desc: "user ID (omit to get current user)"},
		{Name: "user-id-type", Default: "open_id", Desc: "user ID type: open_id | union_id | user_id"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if runtime.Str("user-id") == "" && runtime.IsBot() {
			return common.FlagErrorf("bot identity cannot get current user info, specify --user-id")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		userId := runtime.Str("user-id")
		if userId == "" {
			return common.NewDryRunAPI().
				GET("https://oapi.soke.cn/authen/v1/user_info").
				Desc("(when --user-id omitted) Get current authenticated user info").
				Set("mode", "current_user")
		}
		userIdType := runtime.Str("user-id-type")
		if userIdType == "" {
			userIdType = "open_id"
		}
		return common.NewDryRunAPI().
			GET("https://oapi.soke.cn/contact/v3/users/:user_id").
			Desc("Get user info by user ID").
			Params(map[string]interface{}{"user_id_type": userIdType}).
			Set("user_id", userId).Set("user_id_type", userIdType)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		userId := runtime.Str("user-id")
		userIdType := runtime.Str("user-id-type")

		if userIdType == "" {
			userIdType = "open_id"
		}

		if userId == "" {
			data, err := runtime.CallAPI("GET", "https://oapi.soke.cn/authen/v1/user_info", nil, nil)
			if err != nil {
				return err
			}

			runtime.OutFormat(data, nil, func(w io.Writer) {
				user, _ := data["data"].(map[string]interface{})
				if user == nil {
					return
				}
				output.PrintTable(w, []map[string]interface{}{{
					"name":            user["name"],
					"open_id":        user["open_id"],
					"union_id":       user["union_id"],
					"user_id":        user["user_id"],
					"email":          user["email"],
					"mobile":         user["mobile"],
				}})
			})
			return nil
		}

		params := map[string]interface{}{"user_id_type": userIdType}
		data, err := runtime.CallAPI("GET", "https://oapi.soke.cn/contact/v3/users/"+userId, params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			user, _ := data["data"].(map[string]interface{})
			if user == nil {
				return
			}
			output.PrintTable(w, []map[string]interface{}{{
				"name":            user["name"],
				"open_id":        user["open_id"],
				"union_id":       user["union_id"],
				"user_id":        user["user_id"],
				"email":          user["email"],
				"mobile":         user["mobile"],
				"position":       user["position"],
			}})
		})
		return nil
	},
}
