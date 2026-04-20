package training

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var TrainingGetUser = common.Shortcut{
	Service:     "training",
	Command:     "+get-user",
	Description: "Get training user details",
	Risk:        "read",
	UserScopes:  []string{"training:user:readonly"},
	BotScopes:   []string{"training:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "training-id", Required: true, Desc: "Training ID"},
		{Name: "dept-user-id", Required: true, Desc: "Department user ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/training/user/info").
			Desc("Get training user details").
			Params(map[string]interface{}{
				"training_id":  runtime.Str("training-id"),
				"dept_user_id": runtime.Str("dept-user-id"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"training_id":  runtime.Str("training-id"),
			"dept_user_id": runtime.Str("dept-user-id"),
		}

		data, err := runtime.CallAPI("GET", "/training/user/info", params, nil)
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
					"target_id":        dataObj["target_id"],
					"target_title":     dataObj["target_title"],
					"dept_user_id":     dataObj["dept_user_id"],
					"training_status":  dataObj["training_status"],
					"enroll_status":    dataObj["enroll_status"],
					"enroll_time":      dataObj["enroll_time"],
					"sign_status":      dataObj["sign_status"],
					"sign_time":        dataObj["sign_time"],
					"evaluate_status":  dataObj["evaluate_status"],
					"point":            dataObj["point"],
					"credit":           dataObj["credit"],
					"completed_time":   dataObj["completed_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
