package clock

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ClockListLearnings = common.Shortcut{
	Service:     "clock",
	Command:     "list-learnings",
	Description: "List homework assignments",
	Risk:        "read",
	UserScopes:  []string{"clock:learning:readonly"},
	BotScopes:   []string{"clock:learning:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)", Required: true},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)", Required: true},
		{Name: "start-time", Type: "int", Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Type: "int", Desc: "End time (Unix timestamp in milliseconds)"},
		{Name: "target-id", Desc: "Source ID (e.g., course ID)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		}
		if startTime := runtime.Int("start-time"); startTime > 0 {
			params["start_time"] = startTime
		}
		if endTime := runtime.Int("end-time"); endTime > 0 {
			params["end_time"] = endTime
		}
		if targetID := runtime.Str("target-id"); targetID != "" {
			params["target_id"] = targetID
		}
		return common.NewDryRunAPI().
			GET("/clock/learning/list").
			Desc("List homework assignments").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		}
		if startTime := runtime.Int("start-time"); startTime > 0 {
			params["start_time"] = startTime
		}
		if endTime := runtime.Int("end-time"); endTime > 0 {
			params["end_time"] = endTime
		}
		if targetID := runtime.Str("target-id"); targetID != "" {
			params["target_id"] = targetID
		}

		data, err := runtime.CallAPI("GET", "/clock/learning/list", params, nil)
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
				learning, _ := item.(map[string]interface{})
				if learning != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":          learning["uuid"],
						"title":         learning["title"],
						"homework_type": learning["homework_type"],
						"target_id":     learning["target_id"],
						"target_title":  learning["target_title"],
						"total_score":   learning["total_score"],
						"pass_score":    learning["pass_score"],
						"create_time":   learning["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
