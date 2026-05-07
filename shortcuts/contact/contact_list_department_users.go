package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactListDepartmentUsers = common.Shortcut{
	Service:     "contact",
	Command:     "+list-department-users",
	Description: "List users in a department",
	Risk:        "read",
	UserScopes:  []string{"contact:user:readonly"},
	BotScopes:   []string{"contact:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-id", Required: true, Desc: "Department ID"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		deptID := runtime.Str("dept-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/oa/departmentUser/list").
			Desc("List users in department").
			Params(map[string]interface{}{
				"dept_id":   deptID,
				"page":      page,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		deptID := runtime.Str("dept-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"dept_id":   deptID,
			"page":      page,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("GET", runtime.Config.APIBaseURL+"/oa/departmentUser/list", params, nil)
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
						"company_id":     user["company_id"],
						"email":          user["email"],
						"mobile":         user["mobile"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
