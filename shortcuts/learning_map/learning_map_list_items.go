package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapListItems = common.Shortcut{
	Service:     "learning_map",
	Command:     "+list-items",
	Description: "List learning map stage items",
	Risk:        "read",
	UserScopes:  []string{"learningMap:learningItem:readonly"},
	BotScopes:   []string{"learningMap:learningItem:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "map-id", Required: true, Desc: "Learning map ID"},
		{Name: "stage-id", Required: true, Desc: "Learning map stage ID"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/learningMap/learningItem/list").
			Desc("List learning map stage items").
			Params(map[string]interface{}{
				"map_id":    runtime.Str("map-id"),
				"stage_id":  runtime.Str("stage-id"),
				"page":      runtime.Int("page"),
				"page_size": runtime.Int("page-size"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"map_id":    runtime.Str("map-id"),
			"stage_id":  runtime.Str("stage-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}

		data, err := runtime.CallAPI("GET", "/learningMap/learningItem/list", params, nil)
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
				litem, _ := item.(map[string]interface{})
				if litem != nil {
					rows = append(rows, map[string]interface{}{
						"map_id":       litem["map_id"],
						"stage_id":     litem["stage_id"],
						"type":         litem["type"],
						"target_id":    litem["target_id"],
						"target_title": litem["target_title"],
						"display":      litem["display"],
						"learning_num": litem["learning_num"],
						"finished_num": litem["finished_num"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
