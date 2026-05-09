package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapListStageUsers = common.Shortcut{
	Service:     "learning_map",
	Command:     "+list-stage-users",
	Description: "List learning map stage user results",
	Risk:        "read",
	UserScopes:  []string{"learningMap:learningStageUser:readonly"},
	BotScopes:   []string{"learningMap:learningStageUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "map-id", Required: true, Desc: "Learning map ID"},
		{Name: "stage-id", Required: true, Desc: "Learning map stage ID"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "start-time", Desc: "Start time (Unix timestamp in milliseconds, max 365 days from now)"},
		{Name: "end-time", Desc: "End time (Unix timestamp in milliseconds, max 7 days range)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"map_id":    runtime.Str("map-id"),
			"stage_id":  runtime.Str("stage-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
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
			GET("/learningMap/learningStageUser/list").
			Desc("List learning map stage user results").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"map_id":    runtime.Str("map-id"),
			"stage_id":  runtime.Str("stage-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
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

		data, err := runtime.CallAPI("GET", "/learningMap/learningStageUser/list", params, nil)
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
						"progress":       user["progress"],
						"learn_status":   user["learn_status"],
						"finish_count":   user["finish_count"],
						"miss_item_count": user["miss_item_count"],
						"finish_time":    user["finish_time"],
						"use_days":       user["use_days"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
