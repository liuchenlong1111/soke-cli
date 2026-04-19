package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseListLessonLearns = common.Shortcut{
	Service:     "course",
	Command:     "+list-lesson-learns",
	Description: "List lesson learning records",
	Risk:        "read",
	UserScopes:  []string{"course:lessonLearn:readonly"},
	BotScopes:   []string{"course:lessonLearn:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "course-id", Required: true, Desc: "Course ID"},
		{Name: "lesson-id", Desc: "Lesson ID"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "finish-start-time", Desc: "Finish start time (Unix timestamp in milliseconds)"},
		{Name: "finish-end-time", Desc: "Finish end time (Unix timestamp in milliseconds)"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"course_id": runtime.Str("course-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if lessonID := runtime.Str("lesson-id"); lessonID != "" {
			params["lesson_id"] = lessonID
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
			GET(runtime.Config.APIBaseURL + "/course/lessonLearn/list").
			Desc("List lesson learning records").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"course_id": runtime.Str("course-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if lessonID := runtime.Str("lesson-id"); lessonID != "" {
			params["lesson_id"] = lessonID
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

		data, err := runtime.CallAPI("GET", "/course/lessonLearn/list", params, nil)
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
				learn, _ := item.(map[string]interface{})
				if learn != nil {
					rows = append(rows, map[string]interface{}{
						"dept_user_id": learn["dept_user_id"],
						"lesson_id":    learn["lesson_id"],
						"course_id":    learn["course_id"],
						"learn_status": learn["learn_status"],
						"length":       learn["length"],
						"progress":     learn["progress"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
