package contact

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PositionDelete = common.Shortcut{
	Service:     "contact",
	Command:     "+delete-position",
	Description: "Delete position",
	Risk:        "write",
	UserScopes:  []string{"contact:position:write"},
	BotScopes:   []string{"contact:position:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "position-id", Required: true, Desc: "Position ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			POST(runtime.Config.APIBaseURL + "/position/external/delete").
			Desc("Delete position").
			Params(map[string]interface{}{
				"position_id": runtime.Str("position-id"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"position_id": runtime.Str("position-id"),
		}

		data, err := runtime.CallAPI("POST", "/position/external/delete", params, nil)
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
