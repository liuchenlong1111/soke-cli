package training

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var TrainingCreate = common.Shortcut{
	Service:     "training",
	Command:     "+create",
	Description: "Create training",
	Risk:        "write",
	UserScopes:  []string{"training:training:write"},
	BotScopes:   []string{"training:training:write"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-id", Required: true, Desc: "Creator user ID"},
		{Name: "title", Required: true, Desc: "Training title (max 225 characters)"},
		{Name: "category-id", Required: true, Desc: "Training category ID"},
		{Name: "start-time", Required: true, Desc: "Training start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "Training end time (Unix timestamp in milliseconds)"},
		{Name: "train-place", Required: true, Desc: "Training place (max 100 characters)"},
		{Name: "total-length", Required: true, Desc: "Total training duration"},
		{Name: "lector-id", Required: true, Desc: "Lecturer ID"},
		{Name: "lector-title", Required: true, Desc: "Lecturer name (max 100 characters)"},
		{Name: "train-id", Desc: "Training ID (10-36 characters, auto-generated if not provided)"},
		{Name: "description", Desc: "Detailed description (max 65535 characters)"},
		{Name: "outline", Desc: "Training outline (max 65535 characters)"},
		{Name: "point", Desc: "Point reward (default 0.00)"},
		{Name: "credit", Desc: "Credit reward (default 0.00)"},
		{Name: "lector-point", Desc: "Lecturer point reward (default 0.00)"},
		{Name: "training-form", Desc: "Training form: 1=offline, 2=live (default 1)"},
		{Name: "remind-time", Desc: "Remind time before training starts (minutes, default 0)"},
		{Name: "valid-time-second", Desc: "QR code refresh time (seconds, default 600)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"dept_user_id":  runtime.Str("dept-user-id"),
			"title":         runtime.Str("title"),
			"category_id":   runtime.Str("category-id"),
			"start_time":    runtime.Str("start-time"),
			"end_time":      runtime.Str("end-time"),
			"train_place":   runtime.Str("train-place"),
			"total_length":  runtime.Str("total-length"),
			"lector_id":     runtime.Str("lector-id"),
			"lector_title":  runtime.Str("lector-title"),
		}
		if trainID := runtime.Str("train-id"); trainID != "" {
			params["train_id"] = trainID
		}
		if description := runtime.Str("description"); description != "" {
			params["description"] = description
		}
		if outline := runtime.Str("outline"); outline != "" {
			params["outline"] = outline
		}
		if point := runtime.Str("point"); point != "" {
			params["point"] = point
		}
		if credit := runtime.Str("credit"); credit != "" {
			params["credit"] = credit
		}
		if lectorPoint := runtime.Str("lector-point"); lectorPoint != "" {
			params["lector_point"] = lectorPoint
		}
		if trainingForm := runtime.Str("training-form"); trainingForm != "" {
			params["training_form"] = trainingForm
		}
		if remindTime := runtime.Str("remind-time"); remindTime != "" {
			params["remind_time"] = remindTime
		}
		if validTimeSecond := runtime.Str("valid-time-second"); validTimeSecond != "" {
			params["valid_time_second"] = validTimeSecond
		}
		return common.NewDryRunAPI().
			POST(runtime.Config.APIBaseURL + "/training/training/create").
			Desc("Create training").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"dept_user_id":  runtime.Str("dept-user-id"),
			"title":         runtime.Str("title"),
			"category_id":   runtime.Str("category-id"),
			"start_time":    runtime.Str("start-time"),
			"end_time":      runtime.Str("end-time"),
			"train_place":   runtime.Str("train-place"),
			"total_length":  runtime.Str("total-length"),
			"lector_id":     runtime.Str("lector-id"),
			"lector_title":  runtime.Str("lector-title"),
		}
		if trainID := runtime.Str("train-id"); trainID != "" {
			params["train_id"] = trainID
		}
		if description := runtime.Str("description"); description != "" {
			params["description"] = description
		}
		if outline := runtime.Str("outline"); outline != "" {
			params["outline"] = outline
		}
		if point := runtime.Str("point"); point != "" {
			params["point"] = point
		}
		if credit := runtime.Str("credit"); credit != "" {
			params["credit"] = credit
		}
		if lectorPoint := runtime.Str("lector-point"); lectorPoint != "" {
			params["lector_point"] = lectorPoint
		}
		if trainingForm := runtime.Str("training-form"); trainingForm != "" {
			params["training_form"] = trainingForm
		}
		if remindTime := runtime.Str("remind-time"); remindTime != "" {
			params["remind_time"] = remindTime
		}
		if validTimeSecond := runtime.Str("valid-time-second"); validTimeSecond != "" {
			params["valid_time_second"] = validTimeSecond
		}

		data, err := runtime.CallAPI("POST", "/training/training/create", nil, params)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintJSON(w, data)
		})
		return nil
	},
}
