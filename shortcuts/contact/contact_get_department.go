// Copyright (c) 2026 Soke Technologies.
// SPDX-License-Identifier: MIT

package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactGetDepartment = common.Shortcut{
	Service:     "contact",
	Command:     "+get-department",
	Description: "Get department details by department ID",
	Risk:        "read",
	UserScopes:  []string{"contact:department:readonly"},
	BotScopes:   []string{"contact:department:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-id", Required: true, Desc: "Department ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		deptID := runtime.Str("dept-id")
		return common.NewDryRunAPI().
			GET("https://oapi.soke.cn/oa/department/info").
			Desc("Get department details").
			Params(map[string]interface{}{"dept_id": deptID})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		deptID := runtime.Str("dept-id")

		params := map[string]interface{}{"dept_id": deptID}
		data, err := runtime.CallAPI("GET", "https://oapi.soke.cn/oa/department/info", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			dept, _ := data["data"].(map[string]interface{})
			if dept == nil {
				return
			}
			output.PrintTable(w, []map[string]interface{}{{
				"dept_id":    dept["dept_id"],
				"dept_name":  dept["dept_name"],
				"parent_id":  dept["parent_id"],
				"company_id": dept["company_id"],
			}})
		})
		return nil
	},
}
