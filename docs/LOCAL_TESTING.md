# 本地开发测试指南

本指南帮助开发者在本地测试 soke-cli 的新功能和 AI Agent Skills。

## 问题说明

当你修改了 soke-cli 代码后，需要同时满足两个条件才能在 AI Agent 中测试 Skills：

1. ✅ **全局 CLI 必须是最新版本** - AI Agent 调用的是全局安装的 `soke-cli`
2. ✅ **Skills 文件必须在 Claude 的 skills 目录** - AI Agent 才能发现和加载 Skills

## 快速开始

### 一键测试（推荐）

```bash
# 1. 编译、安装、测试（一条命令完成）
bash ./scripts/local-test.sh
# 当提示时，输入 'y' 更新全局安装

# 2. 链接 Skills 到 Claude
bash ./scripts/link-skills.sh
# 选择 'all' 链接到所有 Claude 目录，或选择特定序号

# 3. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询张三的学习档案"
```

## 详细步骤

### 步骤 1: 更新全局 CLI

**为什么需要这一步？**
- AI Agent 调用的是全局安装的 `soke-cli`（通常在 `/Users/xxx/.nvm/versions/node/xxx/bin/soke-cli`）
- 如果不更新全局版本，AI Agent 会调用旧版本，导致找不到新命令

**操作：**
```bash
bash ./scripts/local-test.sh
```

**脚本会做什么：**
1. 编译项目
2. 检查本地版本功能
3. 检测全局安装状态
4. 提示是否更新全局安装（会自动备份原文件）
5. 运行功能测试

**重要提示：**
- 当提示 "是否要用本地版本覆盖全局安装? (y/n)" 时，输入 `y`
- 脚本会自动备份原文件，安全可靠
- 如果安装失败，会自动恢复备份

### 步骤 2: 链接 Skills

**为什么需要这一步？**
- AI Agent 从特定目录加载 Skills（如 `~/.codex/skills`）
- 需要将本地的 `skills/` 目录链接到 Claude 的 skills 目录

**操作：**
```bash
bash ./scripts/link-skills.sh
```

**脚本会做什么：**
1. 自动查找所有 Claude skills 目录
2. 让你选择链接到哪个目录（或选择 'all' 链接到所有）
3. 创建符号链接（修改本地文件会立即生效）

**常见的 skills 目录：**
- `~/.codex/skills` - Claude Code
- `~/.openclaw/skills` - OpenClaw
- `~/.agents/skills` - Agents
- `~/.workclaw/skills` - WorkClaw

**提示：**
- 选择 `all` 可以链接到所有目录，确保在任何 AI Agent 中都能使用
- 符号链接的好处：修改本地 skills 文件会立即生效，无需重新链接

### 步骤 3: 验证安装

**验证全局 CLI：**
```bash
# 检查 learning-profile 模块
soke-cli learning-profile --help

# 检查 contact 新命令
soke-cli contact +search-dept --help
soke-cli contact +search-user --help
```

**验证 Skills 链接：**
```bash
# 检查链接是否存在
ls -la ~/.codex/skills/soke-*

# 应该看到类似输出：
# lrwxr-xr-x  soke-learning-profile -> /path/to/soke-cli/skills/soke-learning-profile
# lrwxr-xr-x  soke-exam -> /path/to/soke-cli/skills/soke-exam
# lrwxr-xr-x  soke-course -> /path/to/soke-cli/skills/soke-course
```

### 步骤 4: 在 AI Agent 中测试

打开 Claude Code（或其他支持 Skills 的 AI Agent），输入自然语言：

```
查询张三的学习档案
```

**AI Agent 应该：**
1. 自动识别需要使用 `soke-learning-profile` skill
2. 先调用 `soke-cli contact +search-user --dept-user-name "张三"` 获取 dept_user_id
3. 再调用 `soke-cli learning-profile +list --dept-user-ids <dept_user_id>` 查询学习档案
4. 展示结果

## 常见问题

### Q1: AI Agent 提示 "soke-cli 没有 learning-profile 模块"

**原因：** 全局 CLI 还是旧版本

**解决：**
```bash
# 重新运行安装脚本，并选择 'y' 更新全局
bash ./scripts/local-test.sh

# 验证全局版本
soke-cli learning-profile --help
```

### Q2: AI Agent 没有触发 skill

**原因：** Skills 文件没有链接到 Claude 的 skills 目录

**解决：**
```bash
# 运行链接脚本
bash ./scripts/link-skills.sh

# 验证链接
ls -la ~/.codex/skills/soke-*
```

### Q3: 修改了 skill 文档，但 AI Agent 还是用旧的

**原因：** 可能需要重启 AI Agent 或清除缓存

**解决：**
1. 重启 Claude Code
2. 或者在新的对话中测试

### Q4: 有多个 skills 目录，应该链接到哪个？

**建议：** 选择 `all` 链接到所有目录

**原因：**
- 不同的 AI Agent 可能使用不同的目录
- 链接到所有目录确保在任何环境都能使用

### Q5: 如何删除链接？

```bash
# 删除特定 skill
rm ~/.codex/skills/soke-learning-profile

# 删除所有 soke skills
rm ~/.codex/skills/soke-*

# 删除所有目录的链接
rm ~/.codex/skills/soke-*
rm ~/.openclaw/skills/soke-*
rm ~/.agents/skills/soke-*
rm ~/.workclaw/skills/soke-*
```

## 完整工作流

```bash
# 1. 修改代码
vim shortcuts/learning_profile/learning_profile_list.go

# 2. 测试和安装
bash ./scripts/local-test.sh
# 输入 'y' 更新全局安装

# 3. 链接 Skills（首次需要）
bash ./scripts/link-skills.sh
# 选择 'all'

# 4. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询张三的学习档案"

# 5. 如果修改了 skill 文档
vim skills/soke-learning-profile/SKILL.md
# 修改会立即生效（因为是符号链接）

# 6. 运行完整 E2E 测试
bash ./scripts/e2e-test.sh
```

## 脚本说明

### `local-test.sh`
- **功能：** 编译、安装、测试
- **使用：** `bash ./scripts/local-test.sh`
- **特点：** 自动备份、安全可靠

### `link-skills.sh`
- **功能：** 链接 Skills 到 Claude
- **使用：** `bash ./scripts/link-skills.sh`
- **参数：** 可选，传入序号自动选择目录
- **特点：** 支持链接到所有目录

### `e2e-test.sh`
- **功能：** 完整的端到端测试
- **使用：** `bash ./scripts/e2e-test.sh [module]`
- **示例：** `bash ./scripts/e2e-test.sh learning-profile`

## 注意事项

1. **权限问题：** 更新全局 CLI 需要 sudo 权限
2. **备份文件：** 脚本会自动备份，文件名格式：`soke-cli.backup.YYYYMMDD_HHMMSS`
3. **符号链接：** 修改本地文件会立即生效，无需重新链接
4. **多环境：** 如果使用多个 AI Agent，建议链接到所有 skills 目录

## 发布前检查清单

- [ ] 运行 `bash ./scripts/local-test.sh` 通过
- [ ] 运行 `bash ./scripts/e2e-test.sh` 通过
- [ ] 在 AI Agent 中测试 Skills 功能正常
- [ ] 更新 CHANGELOG.md
- [ ] 更新版本号
- [ ] 提交代码并打 tag
