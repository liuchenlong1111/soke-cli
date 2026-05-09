package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactListDepartments = common.Shortcut{
	Service:     "contact",
	Command:     "+list-departments",
	Description: "List departments under a parent department",
	Risk:        "read",
	UserScopes:  []string{"contact:department:readonly"},
	BotScopes:   []string{"contact:department:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "parent-id", Default: "1", Desc: "Parent department ID (default: 1 for root level)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		parentID := runtime.Str("parent-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET("/oa/department/list").
			Desc("List departments under parent department").
			Params(map[string]interface{}{
				"parent_id": parentID,
				"page":      page,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		parentID := runtime.Str("parent-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"parent_id": parentID,
			"page":      page,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("GET", "/oa/department/list", params, nil)
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
				dept, _ := item.(map[string]interface{})
				if dept != nil {
					rows = append(rows, map[string]interface{}{
						"dept_id":    dept["dept_id"],
						"dept_name":  dept["dept_name"],
						"parent_id":  dept["parent_id"],
						"company_id": dept["company_id"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
