package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseListCourseUsers = common.Shortcut{
	Service:     "course",
	Command:     "+list-course-users",
	Description: "List course student learning records",
	Risk:        "read",
	UserScopes:  []string{"course:courseUser:readonly"},
	BotScopes:   []string{"course:courseUser:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "course-id", Required: true, Desc: "Course ID"},
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
			GET(runtime.Config.APIBaseURL + "/course/user/list").
			Desc("List course student learning records").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"course_id": runtime.Str("course-id"),
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

		data, err := runtime.CallAPI("GET", "/course/user/list", params, nil)
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
				courseUser, _ := item.(map[string]interface{})
				if courseUser != nil {
					rows = append(rows, map[string]interface{}{
						"target_id":         courseUser["target_id"],
						"dept_user_id":      courseUser["dept_user_id"],
						"lesson_num":        courseUser["lesson_num"],
						"lesson_finish_num": courseUser["lesson_finish_num"],
						"progress":          courseUser["progress"],
						"learn_status":      courseUser["learn_status"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
