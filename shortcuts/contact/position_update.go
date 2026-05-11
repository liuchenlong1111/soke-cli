package contact

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PositionUpdate = common.Shortcut{
	Service:     "contact",
	Command:     "+update-position",
	Description: "Update position",
	Risk:        "write",
	UserScopes:  []string{"contact:position:write"},
	BotScopes:   []string{"contact:position:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "position-id", Required: true, Desc: "Position ID"},
		{Name: "position", Required: true, Desc: "New position name"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			POST("/position/external/update").
			Desc("Update position").
			Params(map[string]interface{}{
				"position_id": runtime.Str("position-id"),
				"position":    runtime.Str("position"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"position_id": runtime.Str("position-id"),
			"position":    runtime.Str("position"),
		}

		data, err := runtime.CallAPI("POST", "/position/external/update", params, nil)
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
