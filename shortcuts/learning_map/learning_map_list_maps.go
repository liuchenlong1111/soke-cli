package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapListMaps = common.Shortcut{
	Service:     "learning_map",
	Command:     "+list-maps",
	Description: "List learning maps",
	Risk:        "read",
	UserScopes:  []string{"learningMap:learningMap:readonly"},
	BotScopes:   []string{"learningMap:learningMap:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "start-time", Required: true, Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds, max 365 days range)"},
		{Name: "certificate-id", Desc: "Certificate ID"},
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
		if certID := runtime.Str("certificate-id"); certID != "" {
			params["certificate_id"] = certID
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/learningMap/learningMap/list").
			Desc("List learning maps").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if certID := runtime.Str("certificate-id"); certID != "" {
			params["certificate_id"] = certID
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}

		data, err := runtime.CallAPI("GET", "/learningMap/learningMap/list", params, nil)
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
				lmap, _ := item.(map[string]interface{})
				if lmap != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":          lmap["uuid"],
						"title":         lmap["title"],
						"category_id":   lmap["category_id"],
						"module":        lmap["module"],
						"is_new":        lmap["is_new"],
						"stage_number":  lmap["stage_number"],
						"learning_num":  lmap["learning_num"],
						"finish_num":    lmap["finish_num"],
						"publish_time":  lmap["publish_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
