---
name: soke-learning-profile
summary: 学员学习档案查询
version: 1.0.0
description: "学员学习档案查询：查询学员的完整学习记录，包括课程学习、考试成绩、证书获取、学分积分等。当用户需要查询学员学习档案、查看学员学习情况、统计学员学习数据、查询学员综合学习信息时使用。"
metadata:
  requires:
    bins: ["soke-cli"]
  cliHelp: "soke-cli learning-profile --help"
---

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../soke-shared/SKILL.md`](../soke-shared/SKILL.md)，其中包含认证、配置、权限处理**

## 核心概念

- **LearningProfile（学习档案）**: 学员的完整学习记录，包含课程、考试、证书、培训等多维度数据
- **DeptUser（部门用户）**: 企业内的学员，通过 `dept_user_id` 标识
- **Department（部门）**: 组织架构单元，通过 `dept_id` 标识

## 资源关系

```
LearningProfile (学习档案)
├── DeptUser (学员)
│   ├── dept_user_id (学员ID)
│   ├── dept_user_name (学员姓名)
│   └── dept_names (所属部门)
├── CourseData (课程数据)
│   ├── required_learning/finished (必修课程)
│   ├── optional_learning/finished (选修课程)
│   └── learn_time (学习时长)
├── ExamData (考试数据)
│   ├── not_attempt (未开始)
│   ├── passed (通过)
│   ├── notpassed (未通过)
│   └── reviewing (批阅中)
├── CertificateData (证书数据)
│   └── certificate_number (证书数量)
├── PointsAndCredits (积分学分)
│   ├── points (积分)
│   └── credits (学分)
└── TrainingData (培训数据)
    ├── training_finished (线下培训)
    ├── learning_map_finished (学习地图)
    └── training_class_finished (培训班)
```

## Shortcuts（推荐优先使用）

| Shortcut | 说明 |
|----------|------|
| [`+list`](#list) | 查询学员学习档案列表 |

## 辅助命令（contact 模块）

| Shortcut | 说明 |
|----------|------|
| `contact +search-user` | 根据姓名搜索学员 |
| `contact +search-dept` | 根据名称搜索部门 |

---

## 命令详解

### +list

查询学员学习档案列表，支持按部门、学员ID筛选。

**命令格式**:
```bash
soke-cli learning-profile +list \
  [--dept-user-ids <user_id1,user_id2,...>] \
  [--dept-ids <dept_id1,dept_id2,...>] \
  [--is-new <0|1>] \
  [--offset <offset>] \
  [--page-size <size>]
```

**参数说明**:
- `--dept-user-ids`: 学员ID列表，多个ID用逗号分隔（可选）
- `--dept-ids`: 部门ID列表，多个ID用逗号分隔（可选）
- `--is-new`: 是否新员工，0-否，1-是（可选）
- `--offset`: 偏移量，默认从0开始（可选）
- `--page-size`: 每页条数，最大100，默认10（可选）

**返回字段**（表格显示关键字段）:
- `姓名`: 学员姓名
- `部门`: 所属部门
- `职位`: 职位
- `必修完成`: 必修课程完成情况（完成数/总数）
- `选修完成`: 选修课程完成情况（完成数/总数）
- `考试通过`: 考试通过数
- `学习时长`: 学习时长（格式化为"X小时Y分钟"）
- `证书数`: 获得证书数量
- `学分`: 学分
- `积分`: 积分

**JSON 输出包含完整字段**（30+ 字段）:
- 基础信息：`dept_user_id`, `dept_user_name`, `dept_names`, `position`, `job_number`, `avatar`, `is_leave`, `hired_date`
- 课程学习：`optional_learning`, `optional_finished`, `required_learning`, `required_finished`, `learn_time`
- 考试情况：`not_attempt`, `passed`, `notpassed`, `reviewing`
- 证书积分：`certificate_number`, `points`, `credits`
- 培训情况：`training_not_attempt`, `training_finished`, `class_length`
- 学习地图：`leaning_map_attempt`, `learning_map_finished`
- 培训班：`training_class_attempt`, `training_class_finished`
- 其他：`live_learn_time`, `external_training_time`, `knowledge_number`, `evaluation_score`

**示例**:
```bash
# 查询所有学员学习档案
soke-cli learning-profile +list --offset 0 --page-size 20

# 查询特定学员（单个）
soke-cli learning-profile +list --dept-user-ids user123

# 查询多个学员
soke-cli learning-profile +list --dept-user-ids user123,user456,user789

# 查询特定部门的学员
soke-cli learning-profile +list --dept-ids dept456

# 查询多个部门的学员
soke-cli learning-profile +list --dept-ids dept456,dept789

# 查询新员工的学习档案
soke-cli learning-profile +list --is-new 1 --page-size 50

# 分页查询（使用 offset）
soke-cli learning-profile +list --offset 0 --page-size 10   # 第1页
soke-cli learning-profile +list --offset 10 --page-size 10  # 第2页
soke-cli learning-profile +list --offset 20 --page-size 10  # 第3页
```

**权限要求**: `learningProfile:readonly`

**使用场景**:
- 当用户询问"查询学员学习档案"时使用
- 当用户需要查看学员学习情况、统计学习数据时使用
- 当用户需要导出学员学习报表时使用

---

## 权限表

| 命令 | 所需权限 | 说明 |
|------|----------|------|
| `+list` | `learningProfile:readonly` | 查询学员学习档案 |
| `contact +search-user` | `contact:user:readonly` | 搜索学员 |
| `contact +search-dept` | `contact:department:readonly` | 搜索部门 |

---

## 常见工作流

### 工作流1: 根据姓名查询学员学习档案

当用户询问"查询张三的学习档案"时：

**步骤1**: 搜索学员获取 dept_user_id
```bash
soke-cli contact +search-user --dept-user-name "张三"
```

**步骤2**: 使用 dept_user_id 查询学习档案
```bash
soke-cli learning-profile +list --dept-user-ids <从步骤1获取的dept_user_id>
```

**注意**: 学习档案接口不支持按姓名直接查询，必须先通过 contact 模块获取 dept_user_id。

### 工作流2: 查询部门所有学员的学习档案

当用户询问"查询技术部所有学员的学习情况"时：

**步骤1**: 搜索部门获取 dept_id
```bash
soke-cli contact +search-dept --dept-name "技术部"
```

**步骤2**: 使用 dept_id 查询该部门学员学习档案
```bash
soke-cli learning-profile +list --dept-ids <从步骤1获取的dept_id> --page-size 100
```

### 工作流3: 批量导出学习档案

当用户需要导出学习档案数据时：

```bash
# 使用 JSON 格式输出并保存到文件
soke-cli learning-profile +list --page-size 100 --format json > learning_profiles.json
```

---

## 注意事项

1. **时间格式**: `learn_time`, `class_length` 等时长字段单位为秒，表格输出时会自动格式化为"X小时Y分钟"
2. **分页**: 使用 `offset` 和 `page_size` 进行分页，默认每页10条，最大100条
3. **数组参数**: `dept_user_ids` 和 `dept_ids` 是数组类型，多个ID用逗号分隔
4. **姓名查询**: 学习档案接口不支持按姓名查询，需要先通过 `contact +search-user` 获取 dept_user_id
5. **权限**: 需要 `learningProfile:readonly` 权限，如遇权限错误参考 soke-shared
6. **数据完整性**: 表格输出仅显示关键字段，完整数据请使用 `--format json`

---

## 错误处理

### 权限不足
如果遇到权限错误，参考 [`../soke-shared/SKILL.md`](../soke-shared/SKILL.md)

### 参数错误
使用 `--help` 查看命令参数说明：
```bash
soke-cli learning-profile +list --help
```

### 数据为空
如果查询结果为空，检查：
1. 筛选条件是否正确（dept_ids, dept_user_ids, is_new）
2. 学员是否有学习记录
3. 是否有权限查看该部门/学员的数据
4. 注意：学习档案接口不支持按姓名查询，需要先获取 dept_user_id
