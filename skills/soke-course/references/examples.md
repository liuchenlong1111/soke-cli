# 课程查询使用示例

本文档提供了使用 `soke-course` skill 的实际示例。

## 示例1: 查询本月的课程

### 需求
查询2024年5月创建的所有已发布课程。

### 步骤

**1. 计算时间戳**

```javascript
// 2024年5月1日 00:00:00
const startTime = new Date('2024-05-01 00:00:00').getTime();
console.log(startTime); // 1714492800000

// 2024年5月31日 23:59:59
const endTime = new Date('2024-05-31 23:59:59').getTime();
console.log(endTime); // 1717171199000
```

**2. 执行查询**

```bash
soke-cli course +list-courses \
  --start-time 1714492800000 \
  --end-time 1717171199000 \
  --status 1
```

**3. 预期输出**

```
uuid                                  title              category_id  status  lesson_num  create_time
------------------------------------  -----------------  -----------  ------  ----------  -------------
abc123-def456-ghi789                  Go语言入门         cat001       1       20          1714492800000
xyz789-uvw456-rst123                  Python基础教程     cat002       1       15          1714579200000
```

---

## 示例2: 统计课程完成率

### 需求
统计某个课程的学习完成情况，包括完成人数、平均进度等。

### 步骤

**1. 获取课程学习记录**

```bash
soke-cli course +list-course-users \
  --course-id "abc123-def456-ghi789" \
  --page-size 100
```

**2. 分析返回数据**

假设返回数据如下：

```json
{
  "data": {
    "list": [
      {
        "dept_user_id": "user001",
        "dept_user_name": "张三",
        "study_progress": 100,
        "finish_status": 1,
        "study_duration": 7200
      },
      {
        "dept_user_id": "user002",
        "dept_user_name": "李四",
        "study_progress": 50,
        "finish_status": 0,
        "study_duration": 3600
      },
      {
        "dept_user_id": "user003",
        "dept_user_name": "王五",
        "study_progress": 100,
        "finish_status": 1,
        "study_duration": 6800
      }
    ]
  }
}
```

**3. 计算统计数据**

```javascript
const data = response.data.list;

// 总人数
const totalUsers = data.length; // 3

// 完成人数
const completedUsers = data.filter(u => u.finish_status === 1).length; // 2

// 完成率
const completionRate = (completedUsers / totalUsers * 100).toFixed(2); // 66.67%

// 平均进度
const avgProgress = (data.reduce((sum, u) => sum + u.study_progress, 0) / totalUsers).toFixed(2); // 83.33%

// 平均学习时长（小时）
const avgDuration = (data.reduce((sum, u) => sum + u.study_duration, 0) / totalUsers / 3600).toFixed(2); // 1.62小时

console.log(`总人数: ${totalUsers}`);
console.log(`完成人数: ${completedUsers}`);
console.log(`完成率: ${completionRate}%`);
console.log(`平均进度: ${avgProgress}%`);
console.log(`平均学习时长: ${avgDuration}小时`);
```

---

## 示例3: 查询用户的学习情况

### 需求
查询"张三"在某个课程的学习情况。

### 步骤

**1. 查询用户ID**

```bash
soke-cli contact +search-user --name "张三"
```

假设返回：
```json
{
  "data": {
    "dept_user_id": "user001",
    "dept_user_name": "张三"
  }
}
```

**2. 查询课程学习记录**

```bash
soke-cli course +get-course-user \
  --course-id "abc123-def456-ghi789" \
  --dept-user-id "user001"
```

**3. 预期输出**

```json
{
  "data": {
    "target_id": "record001",
    "dept_user_id": "user001",
    "dept_user_name": "张三",
    "study_progress": 100,
    "finish_status": 1,
    "study_duration": 7200,
    "create_time": 1714492800000,
    "update_time": 1714579200000
  }
}
```

**4. 解读结果**

- 学习进度: 100%（已完成）
- 完成状态: 1（已完成）
- 学习时长: 7200秒（2小时）
- 开始时间: 2024-05-01
- 最后更新: 2024-05-02

---

## 示例4: 查询课程的课件学习详情

### 需求
查询某个课程的所有课件，并统计每个课件的学习情况。

### 步骤

**1. 获取课程课件列表**

```bash
soke-cli course +list-lessons \
  --course-id "abc123-def456-ghi789"
```

**2. 预期输出**

```json
{
  "data": {
    "list": [
      {
        "uuid": "lesson001",
        "title": "第1课：Go语言简介",
        "type": "video",
        "duration": 1800,
        "sort": 1
      },
      {
        "uuid": "lesson002",
        "title": "第2课：环境搭建",
        "type": "video",
        "duration": 2400,
        "sort": 2
      }
    ]
  }
}
```

**3. 查询每个课件的学习记录**

```bash
# 查询第1课的学习记录
soke-cli course +list-lesson-learns --lesson-id "lesson001"

# 查询第2课的学习记录
soke-cli course +list-lesson-learns --lesson-id "lesson002"
```

**4. 统计分析**

```javascript
// 假设lesson001有50人学习，40人完成
// 假设lesson002有50人学习，30人完成

const lessonStats = [
  { title: "第1课：Go语言简介", total: 50, completed: 40, rate: "80%" },
  { title: "第2课：环境搭建", total: 50, completed: 30, rate: "60%" }
];

console.table(lessonStats);
```

---

## 示例5: 按分类查询课程

### 需求
查询"技术培训"分类下的所有已发布课程。

### 步骤

**1. 获取课程分类列表**

```bash
soke-cli course +list-categories
```

**2. 找到目标分类**

假设返回：
```json
{
  "data": {
    "list": [
      {
        "uuid": "cat001",
        "title": "技术培训",
        "parent_id": "0"
      },
      {
        "uuid": "cat002",
        "title": "管理培训",
        "parent_id": "0"
      }
    ]
  }
}
```

**3. 查询该分类下的课程**

```bash
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000 \
  --category-id "cat001" \
  --status 1
```

---

## 示例6: 查询人脸识别记录

### 需求
查询某个课件的人脸识别记录，确保学员本人学习。

### 步骤

**1. 查询课件的人脸识别记录**

```bash
soke-cli course +list-lesson-faces \
  --lesson-id "lesson001" \
  --start-time 1714492800000 \
  --end-time 1717171199000
```

**2. 预期输出**

```json
{
  "data": {
    "list": [
      {
        "uuid": "face001",
        "dept_user_id": "user001",
        "dept_user_name": "张三",
        "face_status": 1,
        "face_time": 1714492800000
      },
      {
        "uuid": "face002",
        "dept_user_id": "user002",
        "dept_user_name": "李四",
        "face_status": 2,
        "face_time": 1714493400000
      }
    ]
  }
}
```

**3. 分析识别结果**

- `face_status = 1`: 识别成功（张三）
- `face_status = 2`: 识别失败（李四）

---

## 示例7: 批量查询多个用户的学习记录

### 需求
查询多个用户在某个课程的学习情况。

### 步骤

**1. 准备用户ID列表**

```bash
# 用户ID列表（逗号分隔，最多100个）
USER_IDS="user001,user002,user003,user004,user005"
```

**2. 批量查询**

```bash
soke-cli course +list-course-users \
  --course-id "abc123-def456-ghi789" \
  --userid-list "$USER_IDS"
```

**3. 预期输出**

```json
{
  "data": {
    "list": [
      {
        "dept_user_id": "user001",
        "dept_user_name": "张三",
        "study_progress": 100,
        "finish_status": 1
      },
      {
        "dept_user_id": "user002",
        "dept_user_name": "李四",
        "study_progress": 50,
        "finish_status": 0
      }
      // ... 其他用户
    ]
  }
}
```

---

## 示例8: 查询某个时间段完成的学习记录

### 需求
查询2024年5月完成课程学习的所有用户。

### 步骤

**1. 计算时间范围**

```javascript
const finishStartTime = new Date('2024-05-01 00:00:00').getTime(); // 1714492800000
const finishEndTime = new Date('2024-05-31 23:59:59').getTime();   // 1717171199000
```

**2. 查询完成记录**

```bash
soke-cli course +list-course-users \
  --course-id "abc123-def456-ghi789" \
  --finish-start-time 1714492800000 \
  --finish-end-time 1717171199000
```

**3. 筛选已完成的记录**

```javascript
const data = response.data.list;
const completedUsers = data.filter(u => u.finish_status === 1);

console.log(`5月完成学习的用户数: ${completedUsers.length}`);
completedUsers.forEach(u => {
  console.log(`- ${u.dept_user_name}: ${u.study_progress}%`);
});
```

---

## 常用脚本

### 脚本1: 课程完成率统计脚本

```bash
#!/bin/bash

COURSE_ID="abc123-def456-ghi789"

# 获取学习记录
RESULT=$(soke-cli course +list-course-users --course-id "$COURSE_ID" --page-size 100)

# 使用jq解析JSON（需要安装jq）
TOTAL=$(echo "$RESULT" | jq '.data.list | length')
COMPLETED=$(echo "$RESULT" | jq '[.data.list[] | select(.finish_status == 1)] | length')
RATE=$(echo "scale=2; $COMPLETED * 100 / $TOTAL" | bc)

echo "课程ID: $COURSE_ID"
echo "总人数: $TOTAL"
echo "完成人数: $COMPLETED"
echo "完成率: $RATE%"
```

### 脚本2: 批量导出课程数据

```bash
#!/bin/bash

START_TIME=1704038400000
END_TIME=1735660799000
OUTPUT_FILE="courses.json"

# 导出所有课程
soke-cli course +list-courses \
  --start-time "$START_TIME" \
  --end-time "$END_TIME" \
  --page-size 100 > "$OUTPUT_FILE"

echo "课程数据已导出到: $OUTPUT_FILE"
```

---

## 注意事项

1. **时间戳格式**: 所有时间参数必须使用毫秒级Unix时间戳
2. **分页处理**: 如果数据量大，需要循环分页查询
3. **权限检查**: 确保已完成认证并有相应权限
4. **错误处理**: 建议在脚本中添加错误处理逻辑
5. **API限流**: 注意API调用频率限制

## 相关资源

- [SKILL.md](./SKILL.md) - 完整文档
- [README.md](./README.md) - 快速开始
- [course-list-courses.md](./references/course-list-courses.md) - 详细参考
