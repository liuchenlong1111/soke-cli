---
name: soke-business-training-report
description: 生成某场业务培训考试的整体分析报告，包括参考人数、通过率、平均分、最高分及未通过学员名单。
---

# 培训考试整体分析报告 Skill

## 目标
根据用户提供的培训/考试名称，通过调用开放平台的考试相关 API，自动生成一份结构化的培训考试分析报告。

## 触发条件
当用户提出类似以下问题时，应使用此 Skill：
- "帮我分析一下本次业务培训的考试情况"
- "统计一下这次培训的通过率"
- "输出一份『新员工入职培训』的考试报告"

## 执行步骤

### 第一步：确认目标考试
1. 调用 `GET /exam/exam/list` 接口，获取考试列表。
   - **必填参数**：`start_time` 和 `end_time`（Unix 时间戳，单位毫秒），时间差不超过 365 天。
   - 若用户未提供时间范围，自动设置为近 3 个月。
2. 根据用户提供的培训名称关键词，在返回的 `data.list` 中模糊匹配 `title` 字段。
   - 如果匹配到多个，列出候选供用户选择。
3. 记录目标考试的 `uuid`、`title`、`total_score`、`pass_score`。

### 第二步：获取考试学员成绩列表
1. 调用 `GET /exam/user/list` 接口，参数 `exam_id` 为考试 UUID。
2. 提取每个学员的 `pass_status`（passed/unpassed/not_attempt）、`last_score`。

### 第三步：计算整体统计指标
- 参考人数（排除 not_attempt）
- 通过人数（passed）
- 通过率、平均分、最高分
- 未通过学员列表

### 第四步：输出报告
按照以下格式输出：

📊 【培训名称】考试分析报告
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
👥 参考人数：XX 人
✅ 通过人数：XX 人
📈 通过率：XX%
📊 平均分：XX 分（满分 XX 分）
🏆 最高分：XX 分
❌ 未通过学员：XX 人
  - 学员A（得分：XX）
  ...

## API 调用与鉴权
- 基础 URL：`https://oapi.soke.cn`
- 接口需要 `access_token`，在 Agent 环境中自动获取，无需手动处理。

## 注意事项
- `start_time` 和 `end_time` 必填，Skill 自动生成。
- 注意区分 `not_attempt`（未参考）、`reviewing`（阅卷中）、`unpassed`（未通过）和 `passed`（通过）。
