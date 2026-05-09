# 授客CLI AI Agent Skills

本目录包含授客CLI的AI Agent技能定义，使AI Agent能够自动发现和调用授客CLI的各种功能。

## 安装方法

### 前提条件

1. 已安装授客CLI：
```bash
npm install -g @sokeai/cli
```

2. 已完成配置和认证：
```bash
soke-cli config init
soke-cli auth login
```

### 安装Skills

**方式1：从GitHub安装（推荐）**

```bash
# 全局安装所有skills
npx skills add <org>/soke-cli -y -g

# 或安装特定skill
npx skills add <org>/soke-cli --skill soke-exam -y -g
```

> 注意：需要先将本项目推送到GitHub公开仓库，然后将 `<org>` 替换为实际的GitHub组织或用户名。

**方式2：本地安装（开发测试）**

```bash
# 复制到全局skills目录
cp -r skills/soke-shared ~/.claude/skills/
cp -r skills/soke-exam ~/.claude/skills/

# 或复制到项目级skills目录
mkdir -p .claude/skills
cp -r skills/soke-shared .claude/skills/
cp -r skills/soke-exam .claude/skills/
```

### 验证安装

安装完成后，AI Agent会自动加载这些skills。你可以通过以下方式验证：

1. 在AI Agent对话中询问："查询考试成绩"
2. AI Agent应该能够识别并使用 `soke-exam` skill
3. AI Agent会自动执行 `soke-cli exam +get-exam-user` 命令

## Skills列表

### soke-shared

**功能**: 共享基础规则，包含配置、认证、权限处理

**触发条件**: 
- 首次使用soke-cli
- 需要配置初始化
- 需要用户登录
- 遇到权限错误

**关键内容**:
- 配置初始化（`soke-cli config init`）
- 用户认证（`soke-cli auth login`）
- 权限不足处理
- 错误处理规范
- 安全规则

### soke-exam

**功能**: 考试管理，查询考试、考试用户和成绩

**触发条件**:
- 查询考试成绩
- 查看考试列表
- 查询考试用户信息
- 查看考试分类

**支持的命令**:
- `+list-exams`: 列出考试列表
- `+list-exam-users`: 列出考试用户成绩列表
- `+get-exam-user`: 获取单个考试用户详细成绩
- `+list-categories`: 列出考试分类

**使用示例**:
```bash
# 查询考试成绩
soke-cli exam +get-exam-user --exam-id exam123 --dept-user-id user456

# 列出考试
soke-cli exam +list-exams --start-time 1672502400000 --end-time 1704038400000
```

## 目录结构

```
skills/
├── soke-shared/              # 共享基础skill
│   └── SKILL.md
├── soke-exam/                # 考试管理skill
│   ├── SKILL.md
│   └── references/           # 详细文档
│       └── exam-get-exam-user.md
└── README.md                 # 本文件
```

## Skill文件格式

每个skill目录包含一个 `SKILL.md` 文件，格式如下：

```markdown
---
name: skill-name              # 技能名称
version: 1.0.0                # 版本号
description: "简短描述。当用户需要...时使用。"  # 触发条件（关键）
metadata:
  requires:
    bins: ["soke-cli"]        # 依赖的CLI工具
  cliHelp: "soke-cli exam --help"  # 帮助命令
---

# Skill标题

[Markdown格式的详细说明]
```

**关键字段说明**:
- `name`: 技能标识符，小写字母、数字、连字符
- `version`: 语义版本号
- `description`: **最重要**，AI用它判断何时触发此skill，必须包含"当用户需要...时使用"
- `metadata.requires.bins`: 依赖的CLI工具列表

## 开发新的Skill

### 步骤1：创建目录

```bash
mkdir -p skills/soke-<module>/references
```

### 步骤2：编写SKILL.md

参考 `soke-exam/SKILL.md` 的格式，包含：
1. Frontmatter（name, version, description, metadata）
2. 核心概念说明
3. Shortcuts列表
4. 命令详解
5. 权限表
6. 常见工作流

### 步骤3：编写详细文档（可选）

在 `references/` 目录下为每个重要命令创建详细文档。

### 步骤4：测试

```bash
# 本地安装测试
cp -r skills/soke-<module> ~/.claude/skills/

# 在AI Agent中测试触发条件
```

## 后续扩展

可以按相同模式创建其他业务模块的skills：

- `soke-course` - 课程管理
- `soke-contact` - 组织架构（部门、用户、讲师）
- `soke-training` - 培训管理
- `soke-credit` - 学分管理
- `soke-certificate` - 证书管理
- 等等...

## 注意事项

1. **Description字段至关重要**: 这是AI匹配skill的唯一依据，必须清晰描述触发场景
2. **命令示例要完整**: 包含所有必需参数，避免AI猜测
3. **引用共享规则**: 每个业务skill都应引用soke-shared，避免重复
4. **保持简洁**: SKILL.md应该是快速参考，详细文档放在references/目录
5. **版本管理**: 更新skill时记得更新version字段

## 相关链接

- 授客AI开放平台: https://opendev.soke.cn
- NPM包: @sokeai/cli
- Skills规范: https://github.com/agent-skills/spec

## 许可证

MIT
