package exam

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ExamListExamUsers = common.Shortcut{
	Service:     "exam",
	Command:     "+list-exam-users",
	Description: "List exam user results",
	Risk:        "read",
	UserScopes:  []string{"exam:examUser:readonly"},
	BotScopes:   []string{"exam:examUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "exam-id", Required: true, Desc: "Exam ID"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "finish-start-time", Desc: "Finish start time (Unix timestamp in milliseconds)"},
		{Name: "finish-end-time", Desc: "Finish end time (Unix timestamp in milliseconds)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"exam_id":   runtime.Str("exam-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if finishStartTime := runtime.Str("finish-start-time"); finishStartTime != "" {
			params["finish_start_time"] = finishStartTime
		}
		if finishEndTime := runtime.Str("finish-end-time"); finishEndTime != "" {
			params["finish_end_time"] = finishEndTime
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/exam/user/list").
			Desc("List exam user results").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"exam_id":   runtime.Str("exam-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		if finishStartTime := runtime.Str("finish-start-time"); finishStartTime != "" {
			params["finish_start_time"] = finishStartTime
		}
		if finishEndTime := runtime.Str("finish-end-time"); finishEndTime != "" {
			params["finish_end_time"] = finishEndTime
		}

		data, err := runtime.CallAPI("GET", "/exam/user/list", params, nil)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			dataObj, _ := data["data"].(map[string]interface{})
			if dataObj == nil {
				return
			}
			list, _ := dataObj["list"].([]interface{})
			var rows []map[string]interface{}
			for _, item := range list {
				examUser, _ := item.(map[string]interface{})
				if examUser != nil {
					rows = append(rows, map[string]interface{}{
						"target_id":    examUser["target_id"],
						"dept_user_id": examUser["dept_user_id"],
						"score":        examUser["score"],
						"exam_status":  examUser["exam_status"],
						"create_time":  examUser["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
