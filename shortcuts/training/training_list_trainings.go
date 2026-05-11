package training

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var TrainingListTrainings = common.Shortcut{
	Service:     "training",
	Command:     "+list-trainings",
	Description: "List trainings",
	Risk:        "read",
	UserScopes:  []string{"training:training:readonly"},
	BotScopes:   []string{"training:training:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "start-time", Required: true, Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds, max 365 days range)"},
		{Name: "status", Desc: "Status: 0=unpublished, 1=published, 2=closed"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}
		return common.NewDryRunAPI().
			GET("/training/training/list").
			Desc("List trainings").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}

		data, err := runtime.CallAPI("GET", "/training/training/list", params, nil)
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
				training, _ := item.(map[string]interface{})
				if training != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":              training["uuid"],
						"title":             training["title"],
						"category_id":       training["category_id"],
						"status":            training["status"],
						"enroll_num":        training["enroll_num"],
						"training_num":      training["training_num"],
						"train_start_time":  training["train_start_time"],
						"train_end_time":    training["train_end_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
