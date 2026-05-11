package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactSearchUser = common.Shortcut{
	Service:     "contact",
	Command:     "+search-user",
	Description: "Search users by name",
	Risk:        "read",
	UserScopes:  []string{"contact:user:readonly"},
	BotScopes:   []string{"contact:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-name", Required: true, Desc: "User name to search"},
		{Name: "page-size", Type: "int", Default: "10", Desc: "Page size (max 100)"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if len(runtime.Str("dept-user-name")) == 0 {
			return common.FlagErrorf("user name cannot be empty")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		deptUserName := runtime.Str("dept-user-name")
		pageSize := runtime.Int("page-size")

		return common.NewDryRunAPI().
			POST("/oa/departmentUser/searchDepartmentUserByName").
			Desc("Search users by name").
			Body(map[string]interface{}{
				"dept_user_name": deptUserName,
				"page_size":      pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		deptUserName := runtime.Str("dept-user-name")
		pageSize := runtime.Int("page-size")

		body := map[string]interface{}{
			"dept_user_name": deptUserName,
			"page_size":      pageSize,
		}

		data, err := runtime.CallAPI("POST", "/oa/departmentUser/searchDepartmentUserByName", nil, body)
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
						"dept_user_id":   user["dept_user_id"],
						"dept_user_name": user["dept_user_name"],
						"avatar":         user["avatar"],
						"position":       user["position"],
						"company_id":     user["company_id"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
