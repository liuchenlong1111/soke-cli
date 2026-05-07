package point

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var PointGetConsume = common.Shortcut{
	Service:     "point",
	Command:     "+get-consume",
	Description: "Get point consume order details",
	Risk:        "read",
	UserScopes:  []string{"user:point:readonly"},
	BotScopes:   []string{"user:point:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "trade-no", Required: true, Desc: "Third-party order trade number"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/user/point/consume").
			Desc("Get point consume order details").
			Params(map[string]interface{}{
				"trade_no": runtime.Str("trade-no"),
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"trade_no": runtime.Str("trade-no"),
		}

		data, err := runtime.CallAPI("GET", "/user/point/consume", params, nil)
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
					"order_no":     dataObj["order_no"],
					"dept_user_id": dataObj["dept_user_id"],
					"trade_no":     dataObj["trade_no"],
					"title":        dataObj["title"],
					"point":        dataObj["point"],
					"module":       dataObj["module"],
					"create_time":  dataObj["create_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
