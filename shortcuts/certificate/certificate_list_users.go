package certificate

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CertificateListUsers = common.Shortcut{
	Service:     "certificate",
	Command:     "+list-users",
	Description: "List certificate users",
	Risk:        "read",
	UserScopes:  []string{"certificate:user:readonly"},
	BotScopes:   []string{"certificate:user:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "target-id", Required: true, Desc: "Source UUID (exam, course, training, learning map)"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "start-time", Desc: "Certificate obtain start time (Unix timestamp in milliseconds, max 365 days from now)"},
		{Name: "end-time", Desc: "Certificate obtain end time (Unix timestamp in milliseconds, max 7 days range)"},
		{Name: "update-start-time", Desc: "Certificate update start time (Unix timestamp in milliseconds, max 365 days from now)"},
		{Name: "update-end-time", Desc: "Certificate update end time (Unix timestamp in milliseconds, max 7 days range)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"target_id": runtime.Str("target-id"),
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
		if updateStartTime := runtime.Str("update-start-time"); updateStartTime != "" {
			params["update_start_time"] = updateStartTime
		}
		if updateEndTime := runtime.Str("update-end-time"); updateEndTime != "" {
			params["update_end_time"] = updateEndTime
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/certificate/user/list").
			Desc("List certificate users").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"target_id": runtime.Str("target-id"),
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
		if updateStartTime := runtime.Str("update-start-time"); updateStartTime != "" {
			params["update_start_time"] = updateStartTime
		}
		if updateEndTime := runtime.Str("update-end-time"); updateEndTime != "" {
			params["update_end_time"] = updateEndTime
		}

		data, err := runtime.CallAPI("GET", "/certificate/user/list", params, nil)
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
				user, _ := item.(map[string]interface{})
				if user != nil {
					rows = append(rows, map[string]interface{}{
						"dept_user_id":      user["dept_user_id"],
						"certificate_id":    user["certificate_id"],
						"certificate_title": user["certificate_title"],
						"certificate_num":   user["certificate_num"],
						"module":            user["module"],
						"status":            user["status"],
						"target_id":         user["target_id"],
						"create_time":       user["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
