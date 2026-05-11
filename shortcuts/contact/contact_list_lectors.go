package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactListLectors = common.Shortcut{
	Service:     "contact",
	Command:     "+list-lectors",
	Description: "List lectors (instructors)",
	Risk:        "read",
	UserScopes:  []string{"contact:lector:readonly"},
	BotScopes:   []string{"contact:lector:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET("/lector/lector/list").
			Desc("List lectors").
			Params(map[string]interface{}{
				"page":      page,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("GET", "/lector/lector/list", params, nil)
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
				lector, _ := item.(map[string]interface{})
				if lector != nil {
					rows = append(rows, map[string]interface{}{
						"lector_id":            lector["lector_id"],
						"title":                lector["title"],
						"contact_information":  lector["contact_information"],
						"is_in":                lector["is_in"],
						"employer":             lector["employer"],
						"grade_title":          lector["grade_title"],
						"lector_score":         lector["lector_score"],
						"course_num":           lector["course_num"],
						"training_num":         lector["training_num"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
