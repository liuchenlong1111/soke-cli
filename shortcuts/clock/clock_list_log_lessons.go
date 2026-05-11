package clock

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ClockListLogLessons = common.Shortcut{
	Service:     "clock",
	Command:     "list-log-lessons",
	Description: "List student homework submission materials",
	Risk:        "read",
	UserScopes:  []string{"clock:log:readonly"},
	BotScopes:   []string{"clock:log:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)", Required: true},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)", Required: true},
		{Name: "target-id", Desc: "Student submission ID (uuid from list-logs)", Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		targetID := runtime.Str("target-id")
		return common.NewDryRunAPI().
			GET("/clock/logLesson/list").
			Desc("List student homework submission materials").
			Params(map[string]interface{}{
				"page":      page,
				"page_size": pageSize,
				"target_id": targetID,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		targetID := runtime.Str("target-id")

		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"target_id": targetID,
		}

		data, err := runtime.CallAPI("GET", "/clock/logLesson/list", params, nil)
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
				lesson, _ := item.(map[string]interface{})
				if lesson != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        lesson["uuid"],
						"title":       lesson["title"],
						"type":        lesson["type"],
						"media_id":    lesson["media_id"],
						"media_name":  lesson["media_name"],
						"status":      lesson["status"],
						"create_time": lesson["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
