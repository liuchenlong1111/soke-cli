package contact

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ContactListGroups = common.Shortcut{
	Service:     "contact",
	Command:     "+list-groups",
	Description: "List user groups",
	Risk:        "read",
	UserScopes:  []string{"contact:group:readonly"},
	BotScopes:   []string{"contact:group:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "start-time", Required: true, Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		startTime := runtime.Str("start-time")
		endTime := runtime.Str("end-time")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/oa/group/list").
			Desc("List user groups").
			Params(map[string]interface{}{
				"start_time": startTime,
				"end_time":   endTime,
				"page":       page,
				"page_size":  pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		startTime := runtime.Str("start-time")
		endTime := runtime.Str("end-time")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
			"page":       page,
			"page_size":  pageSize,
		}

		data, err := runtime.CallAPI("GET", runtime.Config.APIBaseURL+"/oa/group/list", params, nil)
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
				group, _ := item.(map[string]interface{})
				if group != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":         group["uuid"],
						"parent_id":    group["parent_id"],
						"name":         group["name"],
						"member_count": group["member_count"],
						"create_time":  group["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
