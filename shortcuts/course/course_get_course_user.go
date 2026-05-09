package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseGetCourseUser = common.Shortcut{
	Service:     "course",
	Command:     "+get-course-user",
	Description: "Get course student learning details",
	Risk:        "read",
	UserScopes:  []string{"course:courseUser:readonly"},
	BotScopes:   []string{"course:courseUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "course-id", Required: true, Desc: "Course ID"},
		{Name: "dept-user-id", Required: true, Desc: "Department user ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		courseID := runtime.Str("course-id")
		deptUserID := runtime.Str("dept-user-id")
		return common.NewDryRunAPI().
			GET("/course/user/info").
			Desc("Get course student learning details").
			Params(map[string]interface{}{
				"course_id":    courseID,
				"dept_user_id": deptUserID,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		courseID := runtime.Str("course-id")
		deptUserID := runtime.Str("dept-user-id")

		params := map[string]interface{}{
			"course_id":    courseID,
			"dept_user_id": deptUserID,
		}

		data, err := runtime.CallAPI("GET", "/course/user/info", params, nil)
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
					"target_id":         dataObj["target_id"],
					"target_title":      dataObj["target_title"],
					"dept_user_id":      dataObj["dept_user_id"],
					"lesson_num":        dataObj["lesson_num"],
					"lesson_finish_num": dataObj["lesson_finish_num"],
					"progress":          dataObj["progress"],
					"learn_status":      dataObj["learn_status"],
					"learn_type":        dataObj["learn_type"],
					"create_time":       dataObj["create_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
