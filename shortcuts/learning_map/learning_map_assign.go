package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapAssign = common.Shortcut{
	Service:     "learning_map",
	Command:     "+assign",
	Description: "Assign learning map to users",
	Risk:        "write",
	UserScopes:  []string{"learningMap:learningMap:write"},
	BotScopes:   []string{"learningMap:learningMap:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Required: true, Desc: "Operator user ID"},
		{Name: "map-id", Required: true, Desc: "Learning map ID"},
		{Name: "assign-user-ids", Desc: "Assign user IDs, comma separated, max 1000"},
		{Name: "assign-dept-ids", Desc: "Assign department IDs, comma separated, max 1000"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"map_id":       runtime.Str("map-id"),
		}
		if assignUserIDs := runtime.Str("assign-user-ids"); assignUserIDs != "" {
			params["assign_user_ids"] = assignUserIDs
		}
		if assignDeptIDs := runtime.Str("assign-dept-ids"); assignDeptIDs != "" {
			params["assign_dept_ids"] = assignDeptIDs
		}
		return common.NewDryRunAPI().
			POST("/learningMap/learningMap/assign").
			Desc("Assign learning map to users").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"map_id":       runtime.Str("map-id"),
		}
		if assignUserIDs := runtime.Str("assign-user-ids"); assignUserIDs != "" {
			params["assign_user_ids"] = assignUserIDs
		}
		if assignDeptIDs := runtime.Str("assign-dept-ids"); assignDeptIDs != "" {
			params["assign_dept_ids"] = assignDeptIDs
		}

		data, err := runtime.CallAPI("POST", "/learningMap/learningMap/assign", nil, params)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintJSON(w, data)
		})
		return nil
	},
}
