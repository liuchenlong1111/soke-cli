package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactSearchDepartment = common.Shortcut{
	Service:     "contact",
	Command:     "+search-department",
	Description: "Search departments by name",
	Risk:        "read",
	UserScopes:  []string{"contact:department:readonly"},
	BotScopes:   []string{"contact:department:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-name", Required: true, Desc: "Department name to search"},
		{Name: "page-size", Type: "int", Default: "10", Desc: "Page size (max 100)"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if len(runtime.Str("dept-name")) == 0 {
			return common.FlagErrorf("department name cannot be empty")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		deptName := runtime.Str("dept-name")
		pageSize := runtime.Int("page-size")

		return common.NewDryRunAPI().
			POST(runtime.Config.APIBaseURL + "/oa/department/searchDepartmentByName").
			Desc("Search departments by name").
			Body(map[string]interface{}{
				"dept_name": deptName,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		deptName := runtime.Str("dept-name")
		pageSize := runtime.Int("page-size")

		body := map[string]interface{}{
			"dept_name": deptName,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("POST", runtime.Config.APIBaseURL+"/oa/department/searchDepartmentByName", nil, body)
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
						"dept_id":   dept["dept_id"],
						"dept_name": dept["dept_name"],
						"parent_id": dept["parent_id"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
