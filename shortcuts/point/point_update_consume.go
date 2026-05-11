package point

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PointUpdateConsume = common.Shortcut{
	Service:     "point",
	Command:     "+update-consume",
	Description: "Add or reduce user points",
	Risk:        "write",
	UserScopes:  []string{"user:point:write"},
	BotScopes:   []string{"user:point:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "trade-no", Required: true, Desc: "Third-party order trade number (unique)"},
		{Name: "dept-user-id", Required: true, Desc: "Department user ID"},
		{Name: "title", Required: true, Desc: "Order description (max 32 characters)"},
		{Name: "point", Required: true, Desc: "Point amount (negative to consume, positive to add)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			POST("/user/point/consumeUpdate").
			Desc("Add or reduce user points").
			Params(map[string]interface{}{
				"trade_no":     runtime.Str("trade-no"),
				"dept_user_id": runtime.Str("dept-user-id"),
				"title":        runtime.Str("title"),
				"point":        runtime.Str("point"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"trade_no":     runtime.Str("trade-no"),
			"dept_user_id": runtime.Str("dept-user-id"),
			"title":        runtime.Str("title"),
			"point":        runtime.Str("point"),
		}

		data, err := runtime.CallAPI("POST", "/user/point/consumeUpdate", nil, params)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintJSON(w, data)
		})
		return nil
	},
}
