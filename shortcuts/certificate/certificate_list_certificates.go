package certificate

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CertificateListCertificates = common.Shortcut{
	Service:     "certificate",
	Command:     "+list-certificates",
	Description: "List certificates",
	Risk:        "read",
	UserScopes:  []string{"certificate:certificate:readonly"},
	BotScopes:   []string{"certificate:certificate:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "start-time", Required: true, Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds, max 365 days range)"},
		{Name: "status", Desc: "Certificate status: 0=unpublished, 1=published, 2=closed, -1=deleted"},
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
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/certificate/certificate/list").
			Desc("List certificates").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = status
		}

		data, err := runtime.CallAPI("GET", "/certificate/certificate/list", params, nil)
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
				cert, _ := item.(map[string]interface{})
				if cert != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        cert["uuid"],
						"title":       cert["title"],
						"status":      cert["status"],
						"code":        cert["code"],
						"category_id": cert["category_id"],
						"create_time": cert["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
