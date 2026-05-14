# 课程列表查询参考

## 命令说明

`soke-cli course +list-courses` 用于查询课程列表，支持按时间范围、分类、状态等条件筛选。

## 基本用法

### 查询指定时间范围的课程

```bash
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000
```

**说明**:
- `--start-time`: 课程创建开始时间（Unix时间戳，毫秒）
- `--end-time`: 课程创建结束时间（Unix时间戳，毫秒）
- 时间范围不能超过365天

## 高级筛选

### 按课程状态筛选

```bash
# 查询已发布的课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --status 1

# 查询未发布的课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --status 0

# 查询已关闭的课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --status 2
```

**状态说明**:
- `0`: 未发布
- `1`: 已发布
- `2`: 已关闭

### 按课程来源筛选

```bash
# 查询自建课
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --is-in 1

# 查询采购课
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --is-in 0
```

**来源说明**:
- `0`: 采购课（从外部采购的课程）
- `1`: 自建课（企业自己创建的课程）

### 按课程分类筛选

```bash
soke-cli course +list-categories

soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --category-id "category-uuid-here"
```

### 组合筛选

```bash
# 查询已发布的自建课，且属于特定分类
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --status 1 \
  --is-in 1 \
  --category-id "category123"
```

## 分页查询

```bash
# 第一页，每页10条
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --page 1 \
  --page-size 10

# 第二页，每页10条
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --page 2 \
  --page-size 10
```

**分页说明**:
- `--page`: 页码，从1开始，默认为1
- `--page-size`: 每页数量，最大100，默认为100

## 返回数据示例

```json
{
  "code": "200",
  "status": "ok",
  "message": "success",
  "data": {
    "list": [
      {
        "uuid": "course-uuid-123",
        "title": "Go语言入门教程",
        "category_id": "category-uuid-456",
        "certificate_id": "cert-uuid-789",
        "lector_id": "lector-uuid-101",
        "study_type": 1,
        "credit": 10.00,
        "point": 100.00,
        "status": 1,
        "lesson_num": 20,
        "total_length": 7200,
        "description": "这是一门Go语言入门课程",
        "pc_url": "https://example.com/course/123",
        "mobile_url": "https://m.example.com/course/123",
        "create_time": 1704038400000,
        "update_time": 1704124800000,
        "create_dept_user_id": "user-uuid-111",
        "create_dept_user_name": "张三"
      }
    ],
    "has_more": 0
  }
}
```

## 返回字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| uuid | String | 课程唯一ID |
| title | String | 课程标题 |
| category_id | String | 课程分类ID |
| certificate_id | String | 关联证书ID |
| lector_id | String | 关联讲师ID |
| study_type | Int | 学习模式（1=自由式, 2=解锁式） |
| credit | Decimal | 学分数量 |
| point | Decimal | 积分数量 |
| status | Int | 课程发布状态（-1=删除, 0=未发布, 1=已发布, 2=关闭） |
| lesson_num | Int | 课件数量 |
| total_length | Int | 学时长度（单位：秒） |
| description | String | 课程描述 |
| pc_url | String | PC端跳转链接 |
| mobile_url | String | 移动端跳转链接 |
| create_time | Int | 创建时间（Unix时间戳，毫秒） |
| update_time | Int | 更新时间（Unix时间戳，毫秒） |
| create_dept_user_id | String | 创建人ID |
| create_dept_user_name | String | 创建人姓名 |
| has_more | Int | 是否还有更多数据（1=有, 0=没有） |

## 时间戳转换

### JavaScript/Node.js
```javascript
// 获取当前时间戳（毫秒）
const now = Date.now();

// 获取指定日期的时间戳
const date = new Date('2024-01-01 00:00:00');
const timestamp = date.getTime();
```

### Python
```python
import time
from datetime import datetime

# 获取当前时间戳（毫秒）
now = int(time.time() * 1000)

# 获取指定日期的时间戳
dt = datetime(2024, 1, 1, 0, 0, 0)
timestamp = int(dt.timestamp() * 1000)
```

### Go
```go
import "time"

// 获取当前时间戳（毫秒）
now := time.Now().UnixMilli()

// 获取指定日期的时间戳
t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
timestamp := t.UnixMilli()
```

## 常见问题

### Q: 为什么时间范围不能超过365天？
A: 这是API的限制，为了防止一次查询返回过多数据。如果需要查询更长时间范围的数据，可以分多次查询。

### Q: 如何查询所有课程？
A: 可以按年份分批查询，例如：
```bash
# 查询2024年的课程
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000

# 查询2025年的课程
soke-cli course +list-courses \
  --start-time 1735660800000 \
  --end-time 1767196799000
```

### Q: 如何知道是否还有更多数据？
A: 查看返回数据中的 `has_more` 字段：
- `1`: 还有更多数据，需要继续分页查询
- `0`: 没有更多数据了

### Q: 学习模式的区别是什么？
A: 
- **自由式（1）**: 学员可以按任意顺序学习课件
- **解锁式（2）**: 学员必须按顺序完成课件，完成前一个才能解锁下一个

## 相关命令

- `soke-cli course +get-course`: 获取单个课程详情
- `soke-cli course +list-categories`: 查询课程分类
- `soke-cli course +list-lessons`: 查询课程的课件列表
- `soke-cli course +list-course-users`: 查询课程的学习记录

## 权限要求

- 权限范围: `course:course:readonly`
- 需要先完成认证: `soke-cli auth login`
