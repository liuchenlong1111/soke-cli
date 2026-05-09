package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseListLessons = common.Shortcut{
	Service:     "course",
	Command:     "+list-lessons",
	Description: "List course lessons",
	Risk:        "read",
	UserScopes:  []string{"course:lesson:readonly"},
	BotScopes:   []string{"course:lesson:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "course-id", Required: true, Desc: "Course ID"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		courseID := runtime.Str("course-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")
		return common.NewDryRunAPI().
			GET("/course/lesson/list").
			Desc("List course lessons").
			Params(map[string]interface{}{
				"course_id": courseID,
				"page":      page,
				"page_size": pageSize,
			})
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		courseID := runtime.Str("course-id")
		page := runtime.Int("page")
		pageSize := runtime.Int("page-size")

		params := map[string]interface{}{
			"course_id": courseID,
			"page":      page,
			"page_size": pageSize,
		}

		data, err := runtime.CallAPI("GET", "/course/lesson/list", params, nil)
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
				lesson, _ := item.(map[string]interface{})
				if lesson != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        lesson["uuid"],
						"title":       lesson["title"],
						"course_id":   lesson["course_id"],
						"sort":        lesson["sort"],
						"create_time": lesson["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
