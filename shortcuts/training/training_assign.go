package training

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var TrainingAssign = common.Shortcut{
	Service:     "training",
	Command:     "+assign",
	Description: "Assign training to users",
	Risk:        "write",
	UserScopes:  []string{"training:training:write"},
	BotScopes:   []string{"training:training:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Required: true, Desc: "Operator user ID"},
		{Name: "training-id", Required: true, Desc: "Training ID"},
		{Name: "assign-user-ids", Desc: "Assign user IDs, comma separated, max 1000"},
		{Name: "assign-dept-ids", Desc: "Assign department IDs, comma separated, max 1000"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"training_id":  runtime.Str("training-id"),
		}
		if assignUserIDs := runtime.Str("assign-user-ids"); assignUserIDs != "" {
			params["assign_user_ids"] = assignUserIDs
		}
		if assignDeptIDs := runtime.Str("assign-dept-ids"); assignDeptIDs != "" {
			params["assign_dept_ids"] = assignDeptIDs
		}
		return common.NewDryRunAPI().
			POST(runtime.Config.APIBaseURL + "/training/training/assign").
			Desc("Assign training to users").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"training_id":  runtime.Str("training-id"),
		}
		if assignUserIDs := runtime.Str("assign-user-ids"); assignUserIDs != "" {
			params["assign_user_ids"] = assignUserIDs
		}
		if assignDeptIDs := runtime.Str("assign-dept-ids"); assignDeptIDs != "" {
			params["assign_dept_ids"] = assignDeptIDs
		}

		data, err := runtime.CallAPI("POST", "/training/training/assign", nil, params)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintJSON(w, data)
		})
		return nil
	},
}
