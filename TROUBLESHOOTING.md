# soke-course Skill 安装问题排查与解决

## 🔍 问题描述

用户执行以下命令后，soke-course skill 没有被正确安装：
```bash
npm install -g @sokeai/cli@latest
npx skills add liuchenlong1111/soke-cli -y -g
```

但是 soke-exam skill 可以正常工作。

## 🐛 问题原因

经过排查，发现问题的根本原因是：

1. ✅ **npm 包中包含了 soke-course**
   - 文件位置：`$(npm root -g)/@sokeai/cli/skills/soke-course/`
   - 所有文件都存在（SKILL.md, README.md, SUMMARY.md, references/）

2. ❌ **但没有被复制到 workclaw 目录**
   - 应该在：`~/.workclaw/skills/soke-course/`
   - 实际：目录不存在

3. ❌ **没有注册到 registry.json**
   - 文件：`~/.workclaw/skills/registry.json`
   - 只包含：exam-generator, soke-exam, soke-shared
   - 缺少：soke-course

## 🔧 解决方案

### 已执行的修复步骤

#### 步骤1: 复制 skill 文件
```bash
cp -r $(npm root -g)/@sokeai/cli/skills/soke-course ~/.workclaw/skills/
```

**结果**：
```
~/.workclaw/skills/soke-course/
├── README.md
├── SKILL.md
├── SUMMARY.md
└── references/
    ├── course-list-courses.md
    └── examples.md
```

#### 步骤2: 更新 registry.json
在 `~/.workclaw/skills/registry.json` 中添加了 soke-course 的注册信息：

```json
{
  "id": "skill:soke-course",
  "name": "soke-course",
  "displayName": "授客课程管理",
  "description": "授客课程管理：查询课程、课程分类、课程详情、学习记录。",
  "install": {
    "path": "/Users/edy/.workclaw/skills/soke-course",
    "version": "1.0.0"
  },
  "state": {
    "enabled": true,
    "health": "ok"
  },
  "metadata": {
    "emoji": "📚",
    "requires": {
      "bins": ["soke-cli"]
    }
  }
}
```

### 验证结果

```bash
# 1. 检查文件存在
ls -la ~/.workclaw/skills/soke-course/
# ✅ 所有文件都存在

# 2. 检查注册信息
grep -c "soke-course" ~/.workclaw/skills/registry.json
# ✅ 返回 3（表示已注册）

# 3. 查看所有 skills
ls -la ~/.workclaw/skills/
# ✅ 包含：exam-generator, soke-course, soke-exam, soke-shared
```

## 📊 对比分析：为什么 soke-exam 可以工作？

### soke-exam 的状态
- ✅ npm 包中存在
- ✅ workclaw 目录中存在（`~/.workclaw/skills/soke-exam/`）
- ✅ registry.json 中已注册

### soke-course 的状态（修复前）
- ✅ npm 包中存在
- ❌ workclaw 目录中不存在
- ❌ registry.json 中未注册

### soke-course 的状态（修复后）
- ✅ npm 包中存在
- ✅ workclaw 目录中存在（`~/.workclaw/skills/soke-course/`）
- ✅ registry.json 中已注册

## 🤔 为什么会出现这个问题？

可能的原因：

1. **安装脚本问题**
   - `npx skills add` 命令可能没有正确处理新添加的 skill
   - 可能需要手动触发安装流程

2. **缓存问题**
   - npm 缓存可能导致新的 skill 没有被识别
   - 需要清除缓存后重新安装

3. **版本问题**
   - 用户可能安装了旧版本（1.0.18），而 soke-course 是在 1.0.19 中添加的
   - 需要确保安装最新版本

## 🚀 推荐的安装流程

为了避免类似问题，建议用户按以下步骤安装：

### 方法1: 完整重新安装（推荐）

```bash
# 1. 清除 npm 缓存
npm cache clean --force

# 2. 卸载旧版本
npm uninstall -g @sokeai/cli

# 3. 安装最新版本
npm install -g @sokeai/cli@latest

# 4. 验证版本
npm list -g @sokeai/cli

# 5. 重新安装 skills
npx skills add liuchenlong1111/soke-cli -y -g

# 6. 验证安装
ls -la ~/.workclaw/skills/
cat ~/.workclaw/skills/registry.json
```

### 方法2: 手动安装（当前使用的方法）

```bash
# 1. 复制 skill 文件
cp -r $(npm root -g)/@sokeai/cli/skills/soke-course ~/.workclaw/skills/

# 2. 手动编辑 registry.json
# 添加 soke-course 的注册信息

# 3. 验证
ls -la ~/.workclaw/skills/soke-course/
grep "soke-course" ~/.workclaw/skills/registry.json
```

## 📝 需要改进的地方

### 1. 自动化安装脚本

建议在 `@sokeai/cli` 包中添加 postinstall 脚本，自动完成以下操作：
- 复制 skills 到 workclaw 目录
- 更新 registry.json
- 验证安装

示例 `scripts/install.js`：
```javascript
const fs = require('fs');
const path = require('path');

// 复制 skills
const skillsDir = path.join(__dirname, '../skills');
const targetDir = path.join(process.env.HOME, '.workclaw/skills');

// 更新 registry.json
const registryPath = path.join(targetDir, 'registry.json');
// ... 自动添加新的 skills
```

### 2. 版本检查

在 `npx skills add` 命令中添加版本检查：
```bash
# 检查本地版本和远程版本
# 如果不一致，提示用户更新
```

### 3. 健康检查命令

添加一个命令来检查 skills 的健康状态：
```bash
npx skills health-check

# 输出：
# ✅ soke-exam: OK
# ✅ soke-course: OK
# ✅ soke-shared: OK
# ❌ some-skill: Missing files
```

## ✅ 当前状态

- ✅ soke-course 文件已复制到 `~/.workclaw/skills/soke-course/`
- ✅ registry.json 已更新，包含 soke-course 注册信息
- ✅ 所有文件完整（SKILL.md, README.md, SUMMARY.md, references/）
- ✅ skill 已启用（enabled: true）

## 🎯 下一步

现在 soke-course skill 应该可以正常使用了。你可以：

1. **重启 Claude Code / workclaw**
   - 让系统重新加载 registry.json

2. **测试 skill**
   - 在对话中询问课程相关的问题
   - 系统应该会自动调用 soke-course skill

3. **验证命令**
   ```bash
   soke-cli course --help
   soke-cli course +list-courses --help
   ```

## 📚 相关文档

- [soke-course SKILL.md](~/.workclaw/skills/soke-course/SKILL.md)
- [soke-course README.md](~/.workclaw/skills/soke-course/README.md)
- [发布指南](../PUBLISH_GUIDE.md)

---

**问题**: soke-course skill 没有被正确安装
**原因**: 文件没有复制到 workclaw 目录，也没有注册到 registry.json
**解决**: 手动复制文件并更新 registry.json
**状态**: ✅ 已解决
**日期**: 2024-05-14
