package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactListGroupUsers = common.Shortcut{
	Service:     "contact",
	Command:     "+list-group-users",
	Description: "List users in a user group",
	Risk:        "read",
	UserScopes:  []string{"contact:group:readonly"},
	BotScopes:   []string{"contact:group:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "group-id", Required: true, Desc: "User group ID"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		groupID := runtime.Str("group-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/oa/group/userList").
			Desc("List users in user group").
			Params(map[string]interface{}{
				"group_id":  groupID,
				"page":      page,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		groupID := runtime.Str("group-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"group_id":  groupID,
			"page":      page,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("GET", runtime.Config.APIBaseURL+"/oa/group/userList", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			dataObj, _ := data["data"].(map[string]interface{})
			if dataObj == nil {
				return
			}
			list, _ := dataObj["list"].([]interface{})
			var rows []map[string]interface{}
			for _, item := range list {
				user, _ := item.(map[string]interface{})
				if user != nil {
					rows = append(rows, map[string]interface{}{
						"group_id":    user["group_id"],
						"user_id":     user["user_id"],
						"user_name":   user["user_name"],
						"create_time": user["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
