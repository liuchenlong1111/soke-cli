package credit

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CreditListLogs = common.Shortcut{
	Service:     "credit",
	Command:     "+list-logs",
	Description: "List credit logs",
	Risk:        "read",
	UserScopes:  []string{"credit:logMeta:readonly"},
	BotScopes:   []string{"credit:logMeta:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "start-time", Desc: "Start time (Unix timestamp in milliseconds, max 365 days from now)"},
		{Name: "end-time", Desc: "End time (Unix timestamp in milliseconds, max 7 days range)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if startTime := runtime.Str("start-time"); startTime != "" {
			params["start_time"] = startTime
		}
		if endTime := runtime.Str("end-time"); endTime != "" {
			params["end_time"] = endTime
		}
		return common.NewDryRunAPI().
			GET("/credit/logMeta/list").
			Desc("List credit logs").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if startTime := runtime.Str("start-time"); startTime != "" {
			params["start_time"] = startTime
		}
		if endTime := runtime.Str("end-time"); endTime != "" {
			params["end_time"] = endTime
		}

		data, err := runtime.CallAPI("GET", "/credit/logMeta/list", params, nil)
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
				log, _ := item.(map[string]interface{})
				if log != nil {
					rows = append(rows, map[string]interface{}{
						"dept_user_id":  log["dept_user_id"],
						"target_id":     log["target_id"],
						"target_title":  log["target_title"],
						"module":        log["module"],
						"credits":       log["credits"],
						"remark":        log["remark"],
						"create_time":   log["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
