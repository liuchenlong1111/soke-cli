# 授客CLI Skills 安装和部署指南

本文档说明如何将授客CLI的Skills部署到GitHub，并通过 `npx skills add` 命令安装。

## 部署到GitHub

### 步骤1：确认当前仓库状态

```bash
cd /Users/edy/www/soke-lark-cli/soke-cli

# 查看当前git状态
git status

# 查看远程仓库
git remote -v
```

当前仓库地址：`https://codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli.git`

### 步骤2：创建GitHub镜像仓库

由于 `npx skills add` 需要从GitHub拉取skills，你需要：

**选项A：在GitHub上创建新仓库**

1. 访问 https://github.com/new
2. 创建仓库，例如：`sokeai/soke-cli` 或 `<your-org>/soke-cli`
3. 设置为公开仓库（Public）

**选项B：使用现有GitHub仓库**

如果已有GitHub仓库，跳过此步骤。

### 步骤3：添加GitHub远程仓库

```bash
# 添加GitHub远程仓库（如果还没有）
git remote add github https://github.com/<org>/soke-cli.git

# 或者修改现有的origin
# git remote set-url origin https://github.com/<org>/soke-cli.git
```

### 步骤4：提交skills文件

```bash
# 添加skills目录
git add skills/

# 提交
git commit -m "feat: add AI Agent skills for exam module

- Add soke-shared skill for common configuration and authentication
- Add soke-exam skill for exam management
- Support npx skills add installation"

# 推送到GitHub
git push github main
# 或 git push origin main
```

### 步骤5：验证GitHub访问

访问以下URL，确认skills目录可公开访问：

```
https://github.com/<org>/soke-cli/tree/main/skills
```

应该能看到：
- `soke-shared/SKILL.md`
- `soke-exam/SKILL.md`
- `README.md`

## 安装测试

### 方式1：从GitHub安装（推荐）

安装完成后，执行以下命令测试：

```bash
# 全局安装所有skills
npx skills add <org>/soke-cli -y -g

# 示例（替换为实际的GitHub仓库）
# npx skills add sokeai/soke-cli -y -g
```

**预期结果**：
- Skills安装到 `~/.claude/skills/soke-exam/` 和 `~/.claude/skills/soke-shared/`
- 可以看到安装成功的提示信息

### 方式2：本地测试（开发阶段）

在推送到GitHub之前，可以先本地测试：

```bash
# 复制到全局skills目录
cp -r skills/soke-shared ~/.claude/skills/
cp -r skills/soke-exam ~/.claude/skills/

# 验证安装
ls -la ~/.claude/skills/
```

### 方式3：从本地路径安装

```bash
# 使用绝对路径
npx skills add /Users/edy/www/soke-lark-cli/soke-cli -y -g

# 或使用相对路径
cd /Users/edy/www/soke-lark-cli
npx skills add ./soke-cli -y -g
```

## 验证功能

### 测试1：检查skills是否安装

```bash
# 查看已安装的skills
ls -la ~/.claude/skills/

# 应该看到：
# soke-shared/
# soke-exam/
```

### 测试2：在AI Agent中测试

在Claude Code或其他支持skills的AI Agent中测试：

**测试场景1：查询考试成绩**

用户输入：
```
帮我查询考试成绩，考试ID是exam123，用户ID是user456
```

AI Agent应该：
1. 识别到需要使用 `soke-exam` skill
2. 自动执行：`soke-cli exam +get-exam-user --exam-id exam123 --dept-user-id user456`
3. 展示成绩结果

**测试场景2：首次使用提示**

用户输入：
```
我想使用soke-cli
```

AI Agent应该：
1. 识别到需要使用 `soke-shared` skill
2. 提示执行：`soke-cli config init` 和 `soke-cli auth login`

**测试场景3：列出考试**

用户输入：
```
列出最近的考试
```

AI Agent应该：
1. 使用 `soke-exam` skill
2. 执行：`soke-cli exam +list-exams --start-time <timestamp> --end-time <timestamp>`

### 测试3：验证skill内容

```bash
# 查看skill内容
cat ~/.claude/skills/soke-exam/SKILL.md | head -20

# 应该看到正确的frontmatter
```

## 更新Skills

当你修改了skills内容后：

### 步骤1：更新版本号

编辑 `skills/soke-exam/SKILL.md`，更新version字段：

```yaml
---
name: soke-exam
version: 1.0.1  # 从1.0.0更新到1.0.1
...
---
```

### 步骤2：提交并推送

```bash
git add skills/
git commit -m "feat: update soke-exam skill to v1.0.1"
git push github main
```

### 步骤3：重新安装

```bash
# 用户需要重新安装以获取最新版本
npx skills add <org>/soke-cli -y -g
```

## 在README中添加安装说明

建议在项目的主README.md中添加skills安装说明：

```markdown
## AI Agent Skills

授客CLI提供了AI Agent技能，使AI能够自动发现和调用CLI功能。

### 安装

```bash
# 1. 安装CLI
npm install -g @sokeai/cli

# 2. 配置和认证
soke-cli config init
soke-cli auth login

# 3. 安装AI Agent Skills
npx skills add <org>/soke-cli -y -g
```

### 使用

安装后，在支持skills的AI Agent（如Claude Code）中，直接用自然语言描述需求：

- "查询考试成绩"
- "列出最近的考试"
- "查看考试分类"

AI Agent会自动调用相应的soke-cli命令。

详细文档：[skills/README.md](skills/README.md)
```

## 常见问题

### Q1: npx skills add 找不到仓库

**原因**: GitHub仓库不存在或不是公开仓库

**解决方案**:
1. 确认仓库URL正确
2. 确认仓库设置为Public
3. 确认skills目录已推送到main分支

### Q2: Skills安装后AI Agent无法识别

**原因**: Description字段不够明确

**解决方案**:
1. 检查SKILL.md的description字段
2. 确保包含"当用户需要...时使用"
3. 添加更多触发关键词

### Q3: 命令执行失败

**原因**: 用户未完成配置或认证

**解决方案**:
1. 确保soke-shared skill正确引导用户配置
2. 在SKILL.md中明确说明前置条件

### Q4: 如何调试skill

**方法1**: 查看AI Agent的执行日志

**方法2**: 手动测试命令
```bash
# 直接执行命令，验证是否正常工作
soke-cli exam +get-exam-user --exam-id exam123 --dept-user-id user456
```

**方法3**: 检查skill文件格式
```bash
# 验证YAML frontmatter格式
head -10 ~/.claude/skills/soke-exam/SKILL.md
```

## 下一步

1. ✅ 创建skills目录和文件
2. ✅ 编写soke-shared和soke-exam skills
3. ⏳ 推送到GitHub仓库
4. ⏳ 测试 `npx skills add` 安装
5. ⏳ 在AI Agent中验证功能
6. ⏳ 创建其他业务模块的skills（course, contact, training等）

## 相关资源

- Skills规范: https://github.com/agent-skills/spec
- npx skills文档: https://skywork.ai/clihub/keywords/skills.html
- Claude Code文档: https://code.claude.com/docs/en/skills
- 授客AI开放平台: https://opendev.soke.cn
