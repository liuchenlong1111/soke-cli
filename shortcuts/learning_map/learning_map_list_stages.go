package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapListStages = common.Shortcut{
	Service:     "learning_map",
	Command:     "+list-stages",
	Description: "List learning map stages",
	Risk:        "read",
	UserScopes:  []string{"learningMap:learningStage:readonly"},
	BotScopes:   []string{"learningMap:learningStage:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "map-id", Required: true, Desc: "Learning map ID"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/learningMap/learningStage/list").
			Desc("List learning map stages").
			Params(map[string]interface{}{
				"map_id":    runtime.Str("map-id"),
				"page":      runtime.Int("page"),
				"page_size": runtime.Int("page-size"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"map_id":    runtime.Str("map-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}

		data, err := runtime.CallAPI("GET", "/learningMap/learningStage/list", params, nil)
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
				stage, _ := item.(map[string]interface{})
				if stage != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":         stage["uuid"],
						"map_id":       stage["map_id"],
						"title":        stage["title"],
						"cycle":        stage["cycle"],
						"display":      stage["display"],
						"learning_num": stage["learning_num"],
						"finish_num":   stage["finish_num"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
