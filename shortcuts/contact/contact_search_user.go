package contact

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactSearchUser = common.Shortcut{
	Service:     "contact",
	Command:     "+search-user",
	Description: "Search users (results sorted by relevance)",
	Risk:        "read",
	UserScopes:  []string{"contact:user:search"},
	AuthTypes:   []string{"user"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "query", Required: true, Desc: "search keyword"},
		{Name: "page-size", Default: "20", Desc: "page size"},
		{Name: "page-token", Desc: "page token"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if len(runtime.Str("query")) == 0 {
			return common.FlagErrorf("search keyword empty")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		pageSizeStr := runtime.Str("page-size")
		pageToken := runtime.Str("page-token")

		pageSize := 20
		if n, err := strconv.Atoi(pageSizeStr); err == nil {
			if n < 1 {
				pageSize = 1
			} else if n > 200 {
				pageSize = 200
			} else {
				pageSize = n
			}
		}

		params := map[string]interface{}{
			"query":     runtime.Str("query"),
			"page_size": pageSize,
		}
		if pageToken != "" {
			params["page_token"] = pageToken
		}

		return common.NewDryRunAPI().
			GET("https://oapi.soke.cn/search/v1/user").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		query := runtime.Str("query")
		pageSizeStr := runtime.Str("page-size")
		pageToken := runtime.Str("page-token")

		pageSize := 20
		if n, err := strconv.Atoi(pageSizeStr); err == nil {
			if n < 1 {
				pageSize = 1
			} else if n > 200 {
				pageSize = 200
			} else {
				pageSize = n
			}
		}

		params := map[string]interface{}{
			"query":     query,
			"page_size": pageSize,
		}
		if pageToken != "" {
			params["page_token"] = pageToken
		}

		data, err := runtime.CallAPI("GET", "https://oapi.soke.cn/search/v1/user", params, nil)
		if err != nil {
			return err
		}

		users, _ := data["users"].([]interface{})

		for _, u := range users {
			if m, ok := u.(map[string]interface{}); ok {
				if av, ok := m["avatar"].(map[string]interface{}); ok {
					m["avatar"] = map[string]interface{}{"avatar_origin": av["avatar_origin"]}
				}
			}
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			if len(users) == 0 {
				fmt.Fprintln(w, "No matching users found.")
				return
			}

			var rows []map[string]interface{}
			for _, u := range users {
				if user, ok := u.(map[string]interface{}); ok {
					rows = append(rows, map[string]interface{}{
						"user_id":   user["user_id"],
						"name":      user["name"],
						"open_id":   user["open_id"],
						"union_id":  user["union_id"],
						"email":     user["email"],
						"mobile":    user["mobile"],
						"position":  user["position"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
