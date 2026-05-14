package learning_profile

import (
	"context"
	"fmt"
	"io"
	"strings"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/output"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
)

var LearningProfileList = common.Shortcut{
	Service:     "learning_profile",
	Command:     "+list",
	Description: "Get student learning profile list",
	Risk:        "read",
	UserScopes:  []string{"learningProfile:readonly"},
	BotScopes:   []string{"learningProfile:readonly"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "dept-user-ids", Desc: "Student ID list (comma separated)"},
		{Name: "dept-ids", Desc: "Department ID list (comma separated)"},
		{Name: "is-new", Desc: "Is new employee: 0-no, 1-yes"},
		{Name: "offset", Type: "int", Default: "0", Desc: "Offset, default from 0"},
		{Name: "page-size", Type: "int", Default: "10", Desc: "Page size (max 100, default 10)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		body := map[string]interface{}{
			"offset":    runtime.Int("offset"),
			"page_size": runtime.Int("page-size"),
		}

		if userIDs := runtime.Str("dept-user-ids"); userIDs != "" {
			body["dept_user_ids"] = strings.Split(userIDs, ",")
		}

		if deptIDs := runtime.Str("dept-ids"); deptIDs != "" {
			body["dept_ids"] = strings.Split(deptIDs, ",")
		}

		if isNew := runtime.Str("is-new"); isNew != "" {
			body["is_new"] = isNew
		}

		return common.NewDryRunAPI().
			POST("/learningProfile/list").
			Desc("Get student learning profile list").
			Body(body)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		body := map[string]interface{}{
			"offset":    runtime.Int("offset"),
			"page_size": runtime.Int("page-size"),
		}

		// dept_user_ids 是数组类型
		if userIDs := runtime.Str("dept-user-ids"); userIDs != "" {
			body["dept_user_ids"] = strings.Split(userIDs, ",")
		}

		// dept_ids 是数组类型
		if deptIDs := runtime.Str("dept-ids"); deptIDs != "" {
			body["dept_ids"] = strings.Split(deptIDs, ",")
		}

		// is_new 是字符串类型
		if isNew := runtime.Str("is-new"); isNew != "" {
			body["is_new"] = isNew
		}

		data, err := runtime.CallAPI("POST", "/learningProfile/list", nil, body)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			list, _ := data["data"].([]interface{})
			if list == nil {
				return
			}
			var rows []map[string]interface{}
			for _, item := range list {
				profile, _ := item.(map[string]interface{})
				if profile != nil {
					// 格式化必修/选修完成情况
					requiredFinished := getStringValue(profile, "required_finished")
					requiredLearning := getStringValue(profile, "required_learning")
					optionalFinished := getStringValue(profile, "optional_finished")
					optionalLearning := getStringValue(profile, "optional_learning")

					rows = append(rows, map[string]interface{}{
						"姓名":   getStringValue(profile, "dept_user_name"),
						"部门":   getStringValue(profile, "dept_names"),
						"职位":   getStringValue(profile, "position"),
						"必修完成": fmt.Sprintf("%s/%s", requiredFinished, requiredLearning),
						"选修完成": fmt.Sprintf("%s/%s", optionalFinished, optionalLearning),
						"考试通过": getStringValue(profile, "passed"),
						"学习时长": formatSeconds(getStringValue(profile, "learn_time")),
						"证书数":  getStringValue(profile, "certificate_number"),
						"学分":   getStringValue(profile, "credits"),
						"积分":   getStringValue(profile, "points"),
					})
				}
			}
			output.PrintTable(w, rows)
		})
		return nil
	},
}

// getStringValue 安全获取字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// formatSeconds 格式化秒数为可读时长
func formatSeconds(secondsStr string) string {
	if secondsStr == "" || secondsStr == "0" {
		return "0"
	}

	var seconds int
	fmt.Sscanf(secondsStr, "%d", &seconds)

	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		return fmt.Sprintf("%d分钟", minutes)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		}
		return fmt.Sprintf("%d小时", hours)
	}
}
