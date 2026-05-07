package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapGetMapUser = common.Shortcut{
	Service:     "learning_map",
	Command:     "+get-map-user",
	Description: "Get learning map user details",
	Risk:        "read",
	UserScopes:  []string{"learningMap:learningMapUser:readonly"},
	BotScopes:   []string{"learningMap:learningMapUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "map-id", Required: true, Desc: "Learning map ID"},
		{Name: "dept-user-id", Required: true, Desc: "Department user ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/learningMap/learningMapUser/info").
			Desc("Get learning map user details").
			Params(map[string]interface{}{
				"map_id":       runtime.Str("map-id"),
				"dept_user_id": runtime.Str("dept-user-id"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"map_id":       runtime.Str("map-id"),
			"dept_user_id": runtime.Str("dept-user-id"),
		}

		data, err := runtime.CallAPI("GET", "/learningMap/learningMapUser/info", params, nil)
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
					"target_id":          dataObj["target_id"],
					"target_title":       dataObj["target_title"],
					"dept_user_id":       dataObj["dept_user_id"],
					"progress":           dataObj["progress"],
					"eligible_status":    dataObj["eligible_status"],
					"learn_status":       dataObj["learn_status"],
					"finish_stage_count": dataObj["finish_stage_count"],
					"miss_count":         dataObj["miss_count"],
					"finish_time":        dataObj["finish_time"],
					"use_days":           dataObj["use_days"],
					"point":              dataObj["point"],
					"credit":             dataObj["credit"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
