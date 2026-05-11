package course

import (
	"context"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var CourseListLessonFaces = common.Shortcut{
	Service:     "course",
	Command:     "+list-lesson-faces",
	Description: "List lesson face recognition records",
	Risk:        "read",
	UserScopes:  []string{"course:lessonFace:readonly"},
	BotScopes:   []string{"course:lessonFace:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "target-id", Required: true, Desc: "Target ID (lesson uuid)"},
		{Name: "userid-list", Desc: "User ID list, comma separated, max 100"},
		{Name: "page", Type: "int", Default: "1", Desc: "Page number (starts from 1)"},
		{Name: "page-size", Type: "int", Default: "100", Desc: "Page size (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		params := map[string]interface{}{
			"target_id": runtime.Str("target-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}
		return common.NewDryRunAPI().
			GET("/course/lessonFace/list").
			Desc("List lesson face recognition records").
			Params(params)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		params := map[string]interface{}{
			"target_id": runtime.Str("target-id"),
			"page":      runtime.Int("page"),
			"page_size": runtime.Int("page-size"),
		}
		if useridList := runtime.Str("userid-list"); useridList != "" {
			params["userid_list"] = useridList
		}

		data, err := runtime.CallAPI("GET", "/course/lessonFace/list", params, nil)
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
				face, _ := item.(map[string]interface{})
				if face != nil {
					rows = append(rows, map[string]interface{}{
						"uuid":          face["uuid"],
						"target_id":     face["target_id"],
						"dept_user_id":  face["dept_user_id"],
						"qualified":     face["qualified"],
						"create_time":   face["create_time"],
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}
