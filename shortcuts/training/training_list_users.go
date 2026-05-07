package training

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var TrainingListUsers = common.Shortcut{
	Service:     "training",
	Command:     "+list-users",
	Description: "List training enrolled users",
	Risk:        "read",
	UserScopes:  []string{"training:user:readonly"},
	BotScopes:   []string{"training:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "training-id", Required: true, Desc: "Training ID"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "start-time", Desc: "Start time (Unix timestamp in milliseconds, max 365 days from now)"},
		{Name: "end-time", Desc: "End time (Unix timestamp in milliseconds, max 7 days range)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"training_id": runtime.Str("training-id"),
			"page":        runtime.Int("page"),
			"page_size":   runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if startTime := runtime.Str("start-time"); startTime != "" {
			params["start_time"] = startTime
		}
		if endTime := runtime.Str("end-time"); endTime != "" {
			params["end_time"] = endTime
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/training/user/list").
			Desc("List training enrolled users").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"training_id": runtime.Str("training-id"),
			"page":        runtime.Int("page"),
			"page_size":   runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if startTime := runtime.Str("start-time"); startTime != "" {
			params["start_time"] = startTime
		}
		if endTime := runtime.Str("end-time"); endTime != "" {
			params["end_time"] = endTime
		}

		data, err := runtime.CallAPI("GET", "/training/user/list", params, nil)
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
						"dept_user_id":     user["dept_user_id"],
						"training_status":  user["training_status"],
						"enroll_status":    user["enroll_status"],
						"enroll_time":      user["enroll_time"],
						"sign_status":      user["sign_status"],
						"evaluate_status":  user["evaluate_status"],
						"point":            user["point"],
						"credit":           user["credit"],
						"completed_time":   user["completed_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
