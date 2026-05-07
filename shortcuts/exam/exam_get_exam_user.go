package exam

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var ExamGetExamUser = common.Shortcut{
	Service:     "exam",
	Command:     "+get-exam-user",
	Description: "Get exam user details",
	Risk:        "read",
	UserScopes:  []string{"exam:examUser:readonly"},
	BotScopes:   []string{"exam:examUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "exam-id", Required: true, Desc: "Exam ID"},
		{Name: "dept-user-id", Required: true, Desc: "Department user ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		examID := runtime.Str("exam-id")
		deptUserID := runtime.Str("dept-user-id")
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/exam/user/info").
			Desc("Get exam user details").
			Params(map[string]interface{}{
				"exam_id":      examID,
				"dept_user_id": deptUserID,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		examID := runtime.Str("exam-id")
		deptUserID := runtime.Str("dept-user-id")

		params := map[string]interface{}{
			"exam_id":      examID,
			"dept_user_id": deptUserID,
		}

		data, err := runtime.CallAPI("GET", "/exam/user/info", params, nil)
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
					"target_id":      dataObj["target_id"],
					"target_title":   dataObj["target_title"],
					"dept_user_id":   dataObj["dept_user_id"],
					"score":          dataObj["score"],
					"exam_status":    dataObj["exam_status"],
					"start_time":     dataObj["start_time"],
					"submit_time":    dataObj["submit_time"],
					"question_count": dataObj["question_count"],
					"create_time":    dataObj["create_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
