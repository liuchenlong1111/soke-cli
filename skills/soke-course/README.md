# soke-course Skill

授客课程管理 Skill，用于查询课程列表、课程分类、课程详情、学习记录等。

## 功能概述

这个 skill 提供了完整的课程管理查询功能，包括：

- 📚 **课程查询**: 查询课程列表、课程详情
- 🏷️ **分类管理**: 查询课程分类
- 📝 **课件管理**: 查询课程下的课件列表
- 👥 **学习记录**: 查询用户的课程学习记录
- 📊 **学习统计**: 查询课件学习记录和人脸识别记录

## 快速开始

### 前置条件

1. 安装 soke-cli
```bash
npm install -g @sokeai/cli
```

2. 配置认证
```bash
soke-cli config init
soke-cli auth login
```

### 基本使用

#### 查询课程列表

```bash
# 查询2024年的所有课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000

# 查询已发布的课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --status 1

# 查询自建课
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --is-in 1
```

#### 查询课程详情

```bash
soke-cli course +get-course --uuid "course-uuid-here"
```

#### 查询课程分类

```bash
soke-cli course +list-categories
```

#### 查询课程学习记录

```bash
# 查询某个课程的所有学习记录
soke-cli course +list-course-users --course-id "course-uuid-here"

# 查询特定用户的学习记录
soke-cli course +get-course-user \
  --course-id "course-uuid-here" \
  --dept-user-id "user-uuid-here"
```

## 主要命令

| 命令 | 说明 |
|------|------|
| `+list-courses` | 查询课程列表 |
| `+get-course` | 查询课程详情 |
| `+list-categories` | 查询课程分类 |
| `+list-lessons` | 查询课程课件列表 |
| `+list-course-users` | 查询课程学习记录 |
| `+get-course-user` | 查询用户课程学习详情 |
| `+list-lesson-learns` | 查询课件学习记录 |
| `+list-lesson-faces` | 查询课件人脸识别记录 |

## 使用场景

### 场景1: 统计课程完成情况

```bash
# 1. 获取课程列表
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000

# 2. 查询某个课程的学习记录
soke-cli course +list-course-users --course-id "course-uuid"

# 3. 分析返回数据中的 finish_status 和 study_progress 字段
```

### 场景2: 查询用户学习情况

```bash
# 1. 查询用户ID（如果只知道用户名）
soke-cli contact +search-user --name "张三"

# 2. 查询该用户的课程学习记录
soke-cli course +get-course-user \
  --course-id "course-uuid" \
  --dept-user-id "user-uuid"
```

### 场景3: 课程内容分析

```bash
# 1. 获取课程详情
soke-cli course +get-course --uuid "course-uuid"

# 2. 获取课程的课件列表
soke-cli course +list-lessons --course-id "course-uuid"

# 3. 查询每个课件的学习记录
soke-cli course +list-lesson-learns --lesson-id "lesson-uuid"
```

## 参数说明

### 时间参数

所有时间参数使用 Unix 时间戳（毫秒）格式：

```javascript
// JavaScript 示例
const startTime = new Date('2024-01-01').getTime(); // 1704038400000
const endTime = new Date('2024-12-31').getTime();   // 1735660799000
```

### 状态参数

**课程状态 (status)**:
- `0`: 未发布
- `1`: 已发布
- `2`: 已关闭

**课程来源 (is-in)**:
- `0`: 采购课
- `1`: 自建课

**学习模式 (study_type)**:
- `1`: 自由式（可任意顺序学习）
- `2`: 解锁式（必须按顺序学习）

## 注意事项

1. **时间范围限制**: `+list-courses` 的时间范围不能超过365天
2. **分页查询**: 默认每页100条，最大100条
3. **权限要求**: 需要相应的只读权限（`course:*:readonly`）
4. **认证要求**: 所有命令都需要先完成 `soke-cli auth login`

## 错误处理

### 权限不足

```bash
# 错误信息: Permission denied
# 解决方案: 检查是否已登录，并确认账号有相应权限
soke-cli auth login
```

### 时间范围超限

```bash
# 错误信息: Time range exceeds 365 days
# 解决方案: 缩小时间范围或分批查询
```

### 数据不存在

```bash
# 错误信息: Course not found
# 解决方案: 检查课程ID是否正确
```

## 相关文档

- [SKILL.md](./SKILL.md) - 完整的 Skill 文档
- [course-list-courses.md](./references/course-list-courses.md) - 课程列表查询详细参考
- [授客开放平台API文档](../../cli-doc/授客开放平台API接口文档.md) - API 接口文档

## 技术支持

如有问题，请参考：
- 授客开放平台: https://opendev.soke.cn
- soke-cli 文档: [../../cli-doc/soke-cli.md](../../cli-doc/soke-cli.md)

## 版本历史

- v1.0.0 (2024-05-14): 初始版本，支持课程查询、分类、学习记录等功能
