package learning_map

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningMapCreate = common.Shortcut{
	Service:     "learning_map",
	Command:     "+create",
	Description: "Create learning map",
	Risk:        "write",
	UserScopes:  []string{"learningMap:learningMap:write"},
	BotScopes:   []string{"learningMap:learningMap:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Required: true, Desc: "Creator user ID"},
		{Name: "title", Required: true, Desc: "Learning map title"},
		{Name: "category-id", Required: true, Desc: "Learning map category ID"},
		{Name: "open-time", Required: true, Desc: "Open time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds)"},
		{Name: "map-id", Desc: "Learning map ID (10-36 characters, auto-generated if not provided)"},
		{Name: "type", Desc: "Type: fixed=fixed time, cycle=cycle type (default: fixed)"},
		{Name: "template", Desc: "Background template: 1-4 (default styles)"},
		{Name: "picture", Desc: "Cover image URL (https:// prefix)"},
		{Name: "description", Desc: "Detailed description (max 65535 characters)"},
		{Name: "is-new", Desc: "Is new employee map: 0=no, 1=yes (default: 0)"},
		{Name: "unlock-type", Desc: "Unlock type: 0=free, 3=by stage (default: 0)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"title":        runtime.Str("title"),
			"category_id":  runtime.Str("category-id"),
			"open_time":    runtime.Str("open-time"),
			"end_time":     runtime.Str("end-time"),
		}
		if mapID := runtime.Str("map-id"); mapID != "" {
			params["map_id"] = mapID
		}
		if typ := runtime.Str("type"); typ != "" {
			params["type"] = typ
		}
		if template := runtime.Str("template"); template != "" {
			params["template"] = template
		}
		if picture := runtime.Str("picture"); picture != "" {
			params["picture"] = picture
		}
		if description := runtime.Str("description"); description != "" {
			params["description"] = description
		}
		if isNew := runtime.Str("is-new"); isNew != "" {
			params["is_new"] = isNew
		}
		if unlockType := runtime.Str("unlock-type"); unlockType != "" {
			params["unlock_type"] = unlockType
		}
		return common.NewDryRunAPI().
			POST(runtime.Config.APIBaseURL + "/learningMap/learningMap/create").
			Desc("Create learning map").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"dept_user_id": runtime.Str("dept-user-id"),
			"title":        runtime.Str("title"),
			"category_id":  runtime.Str("category-id"),
			"open_time":    runtime.Str("open-time"),
			"end_time":     runtime.Str("end-time"),
		}
		if mapID := runtime.Str("map-id"); mapID != "" {
			params["map_id"] = mapID
		}
		if typ := runtime.Str("type"); typ != "" {
			params["type"] = typ
		}
		if template := runtime.Str("template"); template != "" {
			params["template"] = template
		}
		if picture := runtime.Str("picture"); picture != "" {
			params["picture"] = picture
		}
		if description := runtime.Str("description"); description != "" {
			params["description"] = description
		}
		if isNew := runtime.Str("is-new"); isNew != "" {
			params["is_new"] = isNew
		}
		if unlockType := runtime.Str("unlock-type"); unlockType != "" {
			params["unlock_type"] = unlockType
		}

		data, err := runtime.CallAPI("POST", "/learningMap/learningMap/create", nil, params)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintJSON(w, data)
		})
		return nil
	},
}
