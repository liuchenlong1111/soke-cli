package clock

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ClockListLogs = common.Shortcut{
	Service:     "clock",
	Command:     "list-logs",
	Description: "List student homework submissions",
	Risk:        "read",
	UserScopes:  []string{"clock:log:readonly"},
	BotScopes:   []string{"clock:log:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)", Required: true},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)", Required: true},
		{Name: "clock-id", Desc: "Homework ID (uuid from list-learnings)", Required: true},
		{Name: "start-time", Type: "int", Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Type: "int", Desc: "End time (Unix timestamp in milliseconds)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		clockID := runtime.Str("clock-id")
		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"clock_id":  clockID,
		}
		if startTime := runtime.Int("start-time"); startTime > 0 {
			params["start_time"] = startTime
		}
		if endTime := runtime.Int("end-time"); endTime > 0 {
			params["end_time"] = endTime
		}
		return common.NewDryRunAPI().
			GET("/clock/log/list").
			Desc("List student homework submissions").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		clockID := runtime.Str("clock-id")
		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"clock_id":  clockID,
		}
		if startTime := runtime.Int("start-time"); startTime > 0 {
			params["start_time"] = startTime
		}
		if endTime := runtime.Int("end-time"); endTime > 0 {
			params["end_time"] = endTime
		}

		data, err := runtime.CallAPI("GET", "/clock/log/list", params, nil)
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
				log, _ := item.(map[string]interface{})
				if log != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":                  log["uuid"],
						"check_status":          log["check_status"],
						"content":               log["content"],
						"dept_user_id":          log["dept_user_id"],
						"dept_user_name":        log["dept_user_name"],
						"check_dept_user_id":    log["check_dept_user_id"],
						"check_dept_user_name":  log["check_dept_user_name"],
						"create_time":           log["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
