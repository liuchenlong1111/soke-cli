# +get-exam-user - 获取考试用户详细成绩

## 概述

获取单个用户在特定考试中的详细成绩信息，包括分数、状态、答题时间等。

## 命令格式

```bash
soke-cli exam +get-exam-user \
  --exam-id <exam_id> \
  --dept-user-id <dept_user_id>
```

## 参数说明

### 必需参数

| 参数 | 类型 | 说明 |
|------|------|------|
| `--exam-id` | string | 考试ID（uuid） |
| `--dept-user-id` | string | 部门用户ID |

### 可选参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--format` | string | json | 输出格式（json/table） |

## 返回数据

### JSON格式

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "target_id": "exam_user_123",
    "target_title": "2024年度安全培训考试",
    "dept_user_id": "user456",
    "score": 85,
    "exam_status": "已完成",
    "start_time": 1704038400000,
    "submit_time": 1704042000000,
    "question_count": 20,
    "create_time": 1704038400000
  }
}
```

### 表格格式

```
target_id      | target_title           | dept_user_id | score | exam_status | start_time    | submit_time   | question_count | create_time
exam_user_123  | 2024年度安全培训考试    | user456      | 85    | 已完成      | 1704038400000 | 1704042000000 | 20             | 1704038400000
```

## 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `target_id` | string | 考试用户记录ID |
| `target_title` | string | 考试标题 |
| `dept_user_id` | string | 部门用户ID |
| `score` | number | 考试成绩（分数） |
| `exam_status` | string | 考试状态（如：已完成、进行中、未开始） |
| `start_time` | number | 开始答题时间（Unix时间戳，毫秒） |
| `submit_time` | number | 提交时间（Unix时间戳，毫秒） |
| `question_count` | number | 题目总数 |
| `create_time` | number | 记录创建时间（Unix时间戳，毫秒） |

## 使用示例

### 示例1: 查询单个用户成绩

```bash
soke-cli exam +get-exam-user \
  --exam-id exam123 \
  --dept-user-id user456
```

### 示例2: 以表格格式输出

```bash
soke-cli exam +get-exam-user \
  --exam-id exam123 \
  --dept-user-id user456 \
  --format table
```

## 常见场景

### 场景1: 用户询问自己的成绩

**用户输入**: "我的考试成绩是多少？"

**处理步骤**:
1. 获取当前用户的 `dept_user_id`（通过 `soke-cli api GET /users/me`）
2. 确认考试ID（可能需要先列出考试）
3. 执行查询命令

```bash
# 步骤1: 获取当前用户信息
soke-cli api GET /users/me

# 步骤2: 查询成绩
soke-cli exam +get-exam-user \
  --exam-id exam123 \
  --dept-user-id <从步骤1获取的user_id>
```

### 场景2: 管理员查询员工成绩

**用户输入**: "查询张三的考试成绩"

**处理步骤**:
1. 通过姓名查询用户ID（使用 `soke-cli contact +search-user`）
2. 确认考试ID
3. 执行查询命令

```bash
# 步骤1: 查询用户ID
soke-cli contact +search-user --name "张三"

# 步骤2: 查询成绩
soke-cli exam +get-exam-user \
  --exam-id exam123 \
  --dept-user-id <从步骤1获取的dept_user_id>
```

### 场景3: 批量查询多个用户成绩

**用户输入**: "查询所有人的考试成绩"

**处理步骤**:
使用 `+list-exam-users` 更合适，可以一次获取所有用户的成绩列表。

```bash
soke-cli exam +list-exam-users --exam-id exam123
```

## 权限要求

- **所需权限**: `exam:examUser:readonly`
- **认证方式**: 需要先执行 `soke-cli auth login` 完成用户认证

## 错误处理

### 错误1: 考试不存在

```json
{
  "code": 404,
  "msg": "考试不存在"
}
```

**解决方案**: 检查 `exam-id` 是否正确

### 错误2: 用户未参加考试

```json
{
  "code": 404,
  "msg": "用户未参加该考试"
}
```

**解决方案**: 确认用户是否已参加该考试

### 错误3: 权限不足

```json
{
  "code": 403,
  "msg": "权限不足"
}
```

**解决方案**: 
1. 确认已执行 `soke-cli auth login`
2. 联系管理员开通 `exam:examUser:readonly` 权限

### 错误4: 参数缺失

```bash
Error: required flag(s) "exam-id", "dept-user-id" not set
```

**解决方案**: 检查是否提供了所有必需参数

## API详情

- **HTTP方法**: GET
- **API路径**: `/exam/user/info`
- **请求参数**:
  - `exam_id`: 考试ID
  - `dept_user_id`: 部门用户ID

## 相关命令

- `+list-exam-users`: 列出考试用户成绩列表
- `+list-exams`: 列出考试列表
- `soke-cli contact +search-user`: 查询用户信息

## 注意事项

1. **时间戳格式**: 所有时间字段都是Unix时间戳（毫秒），不是秒
2. **用户ID**: 必须使用 `dept_user_id`，不能使用用户名或其他标识
3. **考试状态**: 状态值可能因系统配置而异，常见值包括：已完成、进行中、未开始、已过期
4. **成绩计算**: 成绩字段可能为null（如果考试未完成或未提交）
