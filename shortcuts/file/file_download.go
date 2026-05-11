package file

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var FileDownload = common.Shortcut{
	Service:     "file",
	Command:     "download",
	Description: "Get file download URL",
	Risk:        "read",
	UserScopes:  []string{"file:file:readonly"},
	BotScopes:   []string{"file:file:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "uuid", Desc: "File UUID", Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		uuid := runtime.Str("uuid")
		return common.NewDryRunAPI().
			GET("/file/download").
			Desc("Get file download URL").
			Params(map[string]interface{}{
				"uuid": uuid,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		uuid := runtime.Str("uuid")

		params := map[string]interface{}{
			"uuid": uuid,
		}

		data, err := runtime.CallAPI("GET", "/file/download", params, nil)
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
					"url": dataObj["url"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
