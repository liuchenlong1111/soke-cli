package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseGetCourse = common.Shortcut{
	Service:     "course",
	Command:     "+get-course",
	Description: "Get course details",
	Risk:        "read",
	UserScopes:  []string{"course:course:readonly"},
	BotScopes:   []string{"course:course:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "course-id", Required: true, Desc: "Course ID"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		courseID := runtime.Str("course-id")
		return common.NewDryRunAPI().
			GET("/course/course/info").
			Desc("Get course details").
			Params(map[string]interface{}{
				"course_id": courseID,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		courseID := runtime.Str("course-id")

		params := map[string]interface{}{
			"course_id": courseID,
		}

		data, err := runtime.CallAPI("GET", "/course/course/info", params, nil)
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
					"uuid":        dataObj["uuid"],
					"title":       dataObj["title"],
					"category_id": dataObj["category_id"],
					"status":      dataObj["status"],
					"lesson_num":  dataObj["lesson_num"],
					"create_time": dataObj["create_time"],
				},
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
