# 本地测试脚本使用指南

## 📋 功能说明

`scripts/local-test.js` 是一个本地测试工具，用于在开发完 skill 后快速分发到本地的 AI Agent 环境进行测试，无需发布到 npm。

**重要说明**：
- ✅ 此脚本**只分发到 Agent skills 目录**（如 `~/.workclaw/skills/`）
- ❌ **不包括全局安装目录**（如 npm global 目录）
- 💡 全局安装测试需要使用 `npm install -g` 或其他方式

## 🎯 使用场景

1. **开发新 skill 后**：在发布前先在本地 Agent 中测试功能
2. **修改现有 skill**：快速验证修改效果
3. **调试问题**：在本地 Agent 环境中调试 skill 行为

## 🚀 使用方法

### 1. 分发所有 skills

```bash
node scripts/local-test.js
```

**效果**：
- 自动检测 `skills/` 目录下的所有 `soke-*` skills
- 分发到所有检测到的本地 AI Agent 环境
- 自动更新 workclaw 的 registry.json

### 2. 分发指定 skill

```bash
node scripts/local-test.js --skill=soke-course
```

**效果**：
- 只分发 `soke-course` 到本地环境
- 适合只测试单个 skill 的场景

### 3. 清理所有 skills

```bash
node scripts/local-test.js --clean
```

**效果**：
- 删除所有本地环境中的 `soke-*` skills
- 清理 workclaw 的 registry.json
- 恢复到干净状态

## 📊 支持的环境

脚本会自动检测以下本地 AI Agent 环境：

| 环境 | 路径 | 说明 |
|------|------|------|
| **workclaw** | `~/.workclaw/skills/` | Claude Code Agent skills |
| **claude** | `~/.claude/skills/` | Claude Desktop Agent skills |
| **sokeclaw** | `~/.sokeclaw/openai-agents/workspaces/main/skills/` | Sokeclaw Agent skills |
| **zev** | `~/.zev/openai-agents/workspaces/main/skills/` | Zev Agent skills |

**注意**：
- ✅ 这些是 **Agent skills 目录**，用于 AI Agent 调用
- ❌ **不包括全局安装目录**（如 `$(npm root -g)/@sokeai/cli/skills/`）
- 💡 如需测试全局安装，请使用 `npm install -g` 或 `npm link`

## 📝 完整测试流程

### 步骤1: 开发 skill

```bash
# 创建新 skill
mkdir -p skills/soke-newskill

# 创建 SKILL.md（包含 frontmatter）
vim skills/soke-newskill/SKILL.md
```

### 步骤2: 分发到本地

```bash
# 分发所有 skills
node scripts/local-test.js

# 或只分发新 skill
node scripts/local-test.js --skill=soke-newskill
```

**输出示例**：
```
============================================================
  本地 Skill 测试工具
============================================================

ℹ️  检测到 3 个 skills: soke-course, soke-exam, soke-shared

ℹ️  检测到 4 个本地环境:
ℹ️    • workclaw (Claude Code)
ℹ️      /Users/edy/.workclaw/skills
ℹ️    • claude (Claude Desktop)
ℹ️      /Users/edy/.claude/skills
ℹ️    • sokeclaw
ℹ️      /Users/edy/.sokeclaw/openai-agents/workspaces/main/skills
ℹ️    • zev
ℹ️      /Users/edy/.zev/openai-agents/workspaces/main/skills

============================================================
  分发 Skills 到本地环境
============================================================

ℹ️  分发 soke-course...

✅   → workclaw (Claude Code)
ℹ️       /Users/edy/.workclaw/skills/soke-course
ℹ️       已更新 registry.json

✅   → claude (Claude Desktop)
ℹ️       /Users/edy/.claude/skills/soke-course

✅   → sokeclaw
ℹ️       /Users/edy/.sokeclaw/openai-agents/workspaces/main/skills/soke-course

✅   → zev
ℹ️       /Users/edy/.zev/openai-agents/workspaces/main/skills/soke-course

============================================================
  分发完成
============================================================

ℹ️  总计: 3 个 skills
✅ 成功: 3 个
```

### 步骤3: 重启 AI Agent

**重要**：分发后需要重启 AI Agent 才能加载新的 skills。

- **Claude Code**: 重启 VS Code 或重新打开 Claude Code
- **Claude Desktop**: 重启 Claude Desktop 应用
- **Sokeclaw**: 重启 sokeclaw 进程
- **Zev**: 重启 zev 进程

### 步骤4: 测试功能

#### 方法1: 在对话中测试

在 AI Agent 对话中询问：
```
查询课程列表
查询考试成绩
```

AI Agent 应该会自动调用对应的 skill。

#### 方法2: 验证命令

```bash
# 验证命令是否可用
soke-cli course --help
soke-cli course +list-courses --help

# 测试实际查询
soke-cli course +list-courses \
  --start-time 1704038400000 \
  --end-time 1735660799000
```

### 步骤5: 清理（可选）

测试完成后，可以清理本地环境：

```bash
node scripts/local-test.js --clean
```

## 🔍 脚本功能详解

### 自动检测功能

1. **检测 skills**
   - 扫描 `skills/` 目录
   - 只包含 `soke-*` 开头的目录
   - 按字母顺序排序

2. **检测本地环境**
   - 检查 `~/.workclaw/`
   - 检查 `~/.claude/`
   - 检查 `~/.sokeclaw/`
   - 检查 `~/.zev/`

3. **解析元数据**
   - 从 `SKILL.md` frontmatter 读取
   - 提取 name, summary, description, version, bins

### 分发功能

1. **复制文件**
   - 递归复制整个 skill 目录
   - 保留符号链接
   - 保留文件权限

2. **更新 registry**（仅 workclaw）
   - 自动更新 `registry.json`
   - 添加或更新 skill 条目
   - 设置正确的元数据

### 清理功能

1. **删除文件**
   - 删除所有 `soke-*` skill 目录
   - 递归删除，包括子目录

2. **清理 registry**（仅 workclaw）
   - 从 `registry.json` 中移除对应条目
   - 保持文件格式

## 💡 使用技巧

### 1. 快速迭代开发

```bash
# 修改 skill 代码
vim skills/soke-course/SKILL.md

# 重新分发
node scripts/local-test.js --skill=soke-course

# 重启 AI Agent 并测试
```

### 2. 对比测试

```bash
# 分发到所有环境
node scripts/local-test.js

# 在不同环境中测试相同功能
# 验证跨环境兼容性
```

### 3. 调试问题

```bash
# 清理旧版本
node scripts/local-test.js --clean

# 分发新版本
node scripts/local-test.js

# 重启并测试
```

## ⚠️ 注意事项

### 1. 必须重启 AI Agent

分发后**必须重启** AI Agent 才能加载新的 skills，否则看不到效果。

### 2. 不影响 npm 包

本地测试脚本只修改本地 Agent 环境，**不会影响**已发布的 npm 包。

### 3. 不影响全局安装

本地测试脚本**不会修改**全局安装目录（如 `$(npm root -g)/@sokeai/cli/skills/`）。

如需测试全局安装，请使用：
```bash
# 方法1: 使用 npm link（推荐）
npm link

# 方法2: 本地安装
npm install -g .

# 方法3: 手动复制
cp -r skills/* $(npm root -g)/@sokeai/cli/skills/
```

### 4. workclaw registry

只有 workclaw 环境会自动更新 `registry.json`，其他环境可能需要手动配置。

### 5. 文件覆盖

分发会**覆盖**已存在的同名 skill，请确保没有未保存的修改。

### 6. 清理操作

`--clean` 会删除**所有** `soke-*` skills，包括之前手动安装的。

## 🐛 常见问题

### Q1: 分发后看不到 skill

**原因**: 没有重启 AI Agent

**解决**: 重启 AI Agent（Claude Code / Claude Desktop / sokeclaw）

### Q2: registry.json 没有更新

**原因**: 只有 workclaw 环境会自动更新 registry

**解决**: 其他环境可能需要手动配置或使用各自的安装方式

### Q3: 分发失败

**原因**: 目标目录权限不足或不存在

**解决**: 
```bash
# 检查目录权限
ls -la ~/.workclaw/skills/

# 手动创建目录
mkdir -p ~/.workclaw/skills/
```

### Q4: 清理后 skill 还在

**原因**: 没有重启 AI Agent

**解决**: 重启 AI Agent 清除缓存

## 📚 相关文档

- [AUTO_DETECT_OPTIMIZATION.md](../AUTO_DETECT_OPTIMIZATION.md) - 自动检测机制说明
- [PUBLISH_GUIDE.md](../PUBLISH_GUIDE.md) - 发布到 npm 的流程
- [TROUBLESHOOTING.md](../TROUBLESHOOTING.md) - 问题排查指南

## 🎯 最佳实践

### 开发流程

```bash
# 1. 创建新 skill
mkdir -p skills/soke-newskill
vim skills/soke-newskill/SKILL.md

# 2. 本地测试
node scripts/local-test.js --skill=soke-newskill

# 3. 重启 AI Agent 并测试

# 4. 修改和迭代
vim skills/soke-newskill/SKILL.md
node scripts/local-test.js --skill=soke-newskill
# 重启并测试

# 5. 测试通过后发布
npm version patch -m "feat: add soke-newskill"
npm publish --access public

# 6. 清理本地测试环境
node scripts/local-test.js --clean
```

### 测试检查清单

- [ ] skill 文件已创建（SKILL.md, README.md 等）
- [ ] SKILL.md frontmatter 格式正确
- [ ] 运行 `node scripts/test-auto-detect.js` 验证元数据
- [ ] 运行 `node scripts/local-test.js` 分发到本地
- [ ] 重启 AI Agent
- [ ] 在对话中测试 skill 功能
- [ ] 验证命令行工具可用
- [ ] 测试各种参数组合
- [ ] 测试错误处理
- [ ] 清理测试环境

---

**脚本位置**: `scripts/local-test.js`
**创建日期**: 2024-05-14
**用途**: 本地 skill 测试和分发
**支持环境**: workclaw, claude, sokeclaw, zev
