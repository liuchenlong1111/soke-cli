package file

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var FileListFiles = common.Shortcut{
	Service:     "file",
	Command:     "list-files",
	Description: "List files in file library",
	Risk:        "read",
	UserScopes:  []string{"file:file:readonly"},
	BotScopes:   []string{"file:file:readonly"},
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
			GET("/file/file/list").
			Desc("List files in file library").
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

		data, err := runtime.CallAPI("GET", "/file/file/list", params, nil)
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
				file, _ := item.(map[string]interface{})
				if file != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":           file["uuid"],
						"filename":       file["filename"],
						"filesize":       file["filesize"],
						"ext":            file["ext"],
						"type":           file["type"],
						"length":         file["length"],
						"status":         file["status"],
						"create_time":    file["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
