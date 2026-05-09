package contact

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PositionGetInfo = common.Shortcut{
	Service:     "contact",
	Command:     "+get-position",
	Description: "Get position details",
	Risk:        "read",
	UserScopes:  []string{"contact:position:readonly"},
	BotScopes:   []string{"contact:position:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "position-id", Required: true, Desc: "Position ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/position/external/info").
			Desc("Get position details").
			Params(map[string]interface{}{
				"position_id": runtime.Str("position-id"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"position_id": runtime.Str("position-id"),
		}

		data, err := runtime.CallAPI("GET", "/position/external/info", params, nil)
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
