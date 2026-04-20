package point

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PointGetUserInfo = common.Shortcut{
	Service:     "point",
	Command:     "+get-user-info",
	Description: "Get user point info",
	Risk:        "read",
	UserScopes:  []string{"user:point:readonly"},
	BotScopes:   []string{"user:point:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Desc: "Department user ID"},
		{Name: "user-id", Desc: "User ID in Lark/Feishu platform"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{}
		if deptUserID := runtime.Str("dept-user-id"); deptUserID != "" {
			params["dept_user_id"] = deptUserID
		}
		if userID := runtime.Str("user-id"); userID != "" {
			params["user_id"] = userID
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/user/point/info").
			Desc("Get user point info").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{}
		if deptUserID := runtime.Str("dept-user-id"); deptUserID != "" {
			params["dept_user_id"] = deptUserID
		}
		if userID := runtime.Str("user-id"); userID != "" {
			params["user_id"] = userID
		}

		data, err := runtime.CallAPI("GET", "/user/point/info", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			dataObj, _ := data["data"].(map[string]interface{})
			if dataObj == nil {
				return
			}
			rows := []map[string]interface{}{
				{
					"dept_user_id":   dataObj["dept_user_id"],
					"dept_user_name": dataObj["dept_user_name"],
					"total_points":   dataObj["total_points"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
