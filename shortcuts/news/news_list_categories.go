package news

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var NewsListCategories = common.Shortcut{
	Service:     "news",
	Command:     "+list-categories",
	Description: "List news categories",
	Risk:        "read",
	UserScopes:  []string{"news:category:readonly"},
	BotScopes:   []string{"news:category:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/news/category/list").
			Desc("List news categories").
			Params(map[string]interface{}{
				"page":      runtime.Int("page"),
				"page_size": runtime.Int("page-size"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}

		data, err := runtime.CallAPI("GET", "/news/category/list", params, nil)
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
				category, _ := item.(map[string]interface{})
				if category != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        category["uuid"],
						"title":       category["title"],
						"parent_id":   category["parent_id"],
						"create_time": category["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
