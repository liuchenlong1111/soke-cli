package contact

import (
	"context"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactGetDepartmentUser = common.Shortcut{
	Service:     "contact",
	Command:     "+get-department-user",
	Description: "Get department user details by user ID",
	Risk:        "read",
	UserScopes:  []string{"contact:user:readonly"},
	BotScopes:   []string{"contact:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Required: true, Desc: "Department user ID (employee ID)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		deptUserID := runtime.Str("dept-user-id")
		return common.NewDryRunAPI().
			GET("/oa/departmentUser/info").
			Desc("Get department user details").
			Params(map[string]interface{}{"dept_user_id": deptUserID})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		deptUserID := runtime.Str("dept-user-id")

		params := map[string]interface{}{"dept_user_id": deptUserID}
		data, err := runtime.CallAPI("GET", "/oa/departmentUser/info", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			user, _ := data["data"].(map[string]interface{})
			if user == nil {
				return
			}

			deptIDs := ""
			if deptIDsArr, ok := user["dept_ids"].([]interface{}); ok {
				for i, id := range deptIDsArr {
					if i > 0 {
						deptIDs += ", "
					}
					deptIDs += fmt.Sprintf("%v", id)
				}
			}

			output.PrintTable(w, []map[string]interface{}{{
				"dept_user_id":   user["dept_user_id"],
				"dept_user_name": user["dept_user_name"],
				"avatar":         user["avatar"],
				"position":       user["position"],
				"dept_ids":       deptIDs,
				"company_id":     user["company_id"],
			}})
		})
		return nil
	},
}
