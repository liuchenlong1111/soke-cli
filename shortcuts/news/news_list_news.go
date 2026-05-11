package news

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var NewsListNews = common.Shortcut{
	Service:     "news",
	Command:     "+list-news",
	Description: "List news",
	Risk:        "read",
	UserScopes:  []string{"news:news:readonly"},
	BotScopes:   []string{"news:news:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "category-id", Desc: "News category ID"},
		{Name: "status", Desc: "Status: 0=unpublished, 1=published, 2=closed"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if categoryID := runtime.Str("category-id"); categoryID != "" {
			params["category_id"] = categoryID
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}
		return common.NewDryRunAPI().
			GET("/news/news/list").
			Desc("List news").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if categoryID := runtime.Str("category-id"); categoryID != "" {
			params["category_id"] = categoryID
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}

		data, err := runtime.CallAPI("GET", "/news/news/list", params, nil)
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
				newsItem, _ := item.(map[string]interface{})
				if newsItem != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        newsItem["uuid"],
						"title":       newsItem["title"],
						"category_id": newsItem["category_id"],
						"status":      newsItem["status"],
						"is_top":      newsItem["is_top"],
						"pic_count":   newsItem["pic_count"],
						"target_type": newsItem["target_type"],
						"create_time": newsItem["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
