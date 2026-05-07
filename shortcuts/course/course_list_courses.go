package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseListCourses = common.Shortcut{
	Service:     "course",
	Command:     "+list-courses",
	Description: "List courses",
	Risk:        "read",
	UserScopes:  []string{"course:course:readonly"},
	BotScopes:   []string{"course:course:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "start-time", Required: true, Desc: "Start time (Unix timestamp in milliseconds)"},
		{Name: "end-time", Required: true, Desc: "End time (Unix timestamp in milliseconds)"},
		{Name: "category-id", Desc: "Course category ID"},
		{Name: "is-in", Type: "int", Desc: "0=purchased course, 1=self-built course"},
		{Name: "status", Type: "int", Desc: "Course status: 0=unpublished, 1=published, 2=closed"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if categoryID := runtime.Str("category-id"); categoryID != "" {
			params["category_id"] = categoryID
		}
		if isIn := runtime.Str("is-in"); isIn != "" {
			params["is_in"] = runtime.Int("is-in")
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = runtime.Int("status")
		}
		return common.NewDryRunAPI().
			GET(runtime.Config.APIBaseURL + "/course/course/list").
			Desc("List courses").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"start_time": runtime.Str("start-time"),
			"end_time":   runtime.Str("end-time"),
			"page":       runtime.Int("page"),
			"page_size":  runtime.Int("page-size"),
		}
		if categoryID := runtime.Str("category-id"); categoryID != "" {
			params["category_id"] = categoryID
		}
		if isIn := runtime.Str("is-in"); isIn != "" {
			params["is_in"] = runtime.Int("is-in")
		}
		if status := runtime.Str("status"); status != "" {
			params["status"] = runtime.Int("status")
		}

		data, err := runtime.CallAPI("GET", "/course/course/list", params, nil)
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
				course, _ := item.(map[string]interface{})
				if course != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":        course["uuid"],
						"title":       course["title"],
						"category_id": course["category_id"],
						"status":      course["status"],
						"lesson_num":  course["lesson_num"],
						"create_time": course["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
