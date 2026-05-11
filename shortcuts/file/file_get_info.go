package file

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var FileGetInfo = common.Shortcut{
	Service:     "file",
	Command:     "get-info",
	Description: "Get file details",
	Risk:        "read",
	UserScopes:  []string{"file:file:readonly"},
	BotScopes:   []string{"file:file:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "file-id", Desc: "File ID", Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		fileID := runtime.Str("file-id")
		return common.NewDryRunAPI().
			GET("/file/file/info").
			Desc("Get file details").
			Params(map[string]interface{}{
				"file_id": fileID,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		fileID := runtime.Str("file-id")

		params := map[string]interface{}{
			"file_id": fileID,
		}

		data, err := runtime.CallAPI("GET", "/file/file/info", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			dataObj, _ := data["data"].(map[string]interface{})
			if dataObj == nil {
				return
			}
			rows := []map[string]interface{}{
				{
					"uuid":           dataObj["uuid"],
					"filename":       dataObj["filename"],
					"filesize":       dataObj["filesize"],
					"ext":            dataObj["ext"],
					"type":           dataObj["type"],
					"length":         dataObj["length"],
					"object":         dataObj["object"],
					"convert_object": dataObj["convert_object"],
					"storage":        dataObj["storage"],
					"status":         dataObj["status"],
					"create_time":    dataObj["create_time"],
					"update_time":    dataObj["update_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
