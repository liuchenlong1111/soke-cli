package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var NewsGetInfo = common.Shortcut{
	Service:     "news",
	Command:     "+get-info",
	Description: "Get news details",
	Risk:        "read",
	UserScopes:  []string{"news:news:readonly"},
	BotScopes:   []string{"news:news:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "uuid", Required: true, Desc: "News UUID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/news/news/info").
			Desc("Get news details").
			Params(map[string]interface{}{
				"uuid": runtime.Str("uuid"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"uuid": runtime.Str("uuid"),
		}

		data, err := runtime.CallAPI("GET", "/news/news/info", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			jsonData, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(w, string(jsonData))
		})
		return nil
	},
}
