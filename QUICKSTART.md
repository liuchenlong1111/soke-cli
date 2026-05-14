# soke-cli 快速使用指南

本文档提供 soke-cli 的快速参考指南，涵盖从开发到发布的完整流程。

## 目录

- [用户使用指南](#用户使用指南)
- [开发者快速参考](#开发者快速参考)
- [完整开发流程](#完整开发流程)
- [故障排查](#故障排查)

---

## 用户使用指南

### 安装

```bash
# 1. 安装 CLI
npm install -g @sokeai/cli

# 2. 安装 AI Agent Skills（可选，推荐）
npx skills add liuchenlong1111/soke-cli -y -g
```

#### 在 sokeclaw 客户端中使用 Skills（重要）

从 `@sokeai/cli@1.0.7` 起，执行 `npm install -g @sokeai/cli` 时会自动检测本机是否已安装 sokeclaw（`~/.sokeclaw/`），并将内置的 Skills（`soke-shared`、`soke-exam`）同步到当前工作区的 skills 目录，默认位置为：

`~/.sokeclaw/openai-agents/workspaces/main/skills/`

因此用户只需要执行两条命令即可在 sokeclaw 工作区里看到并调用（例如 `/skill soke-exam`）：

- `npm install -g @sokeai/cli`
- `npx skills add liuchenlong1111/soke-cli -y -g`（可选但推荐：用于把同一套技能安装到 `~/.agents/skills/`，供其它支持 skills 的 Agent 使用）

如果你安装后没有立刻看到技能，重启/刷新 sokeclaw 即可触发重新扫描。

如果 sokeclaw 内执行技能时提示 `command not found: soke-cli`，说明 sokeclaw 运行时的 GUI 环境没有读取到你的 NVM/Node 环境变量。

**一键修复命令：**
```bash
soke-cli setup-gui-env
```
执行此命令后，它会自动将你的 Node/NVM 环境变量注入给 GUI 环境，并重启 SokeClaw 客户端，之后即可正常使用。

### 初始化配置

```bash
soke-cli config init
```

### 登录

```bash
soke-cli auth login
```

### 使用示例

**方式一：通过 AI Agent（推荐）**

如果已安装 Skills，在支持 Skills 的 AI Agent（如 Claude Code）中，直接用自然语言：

```
"查询考试成绩"
"列出最近的考试"
"查看考试分类"
```

AI Agent 会自动：
- 识别你的意图
- 选择合适的命令
- 提示你提供必要的参数
- 执行命令并展示结果

**方式二：直接使用命令行**

```bash
# 查看帮助
soke-cli --help

# 查看课程列表
soke-cli course +list-courses \
  --start-time 1672502400000 \
  --end-time 1704038400000 \
  --page 1 \
  --page-size 10

# 查看考试列表
soke-cli exam +list-exams \
  --start-time 1672502400000 \
  --end-time 1704038400000

# 查询考试成绩
soke-cli exam +get-exam-user \
  --exam-id exam123 \
  --dept-user-id user456

# 查看通讯录
soke-cli contact +list-users --page 1
```

### AI Agent Skills 说明

已安装的 Skills：
- **soke-shared**: 配置初始化、用户认证、权限处理
- **soke-exam**: 考试管理（查询考试、考试成绩、考试分类）
- **soke-course**: 课程管理（查询课程、课程分类、学习记录）
- **soke-learning-profile**: 学习档案查询（查询学员学习档案、学习情况统计）

详细文档：[skills/README.md](skills/README.md)

---

## 开发者快速参考

### 本地开发

```bash
# 编译
go build -o soke-cli .

# 或使用 Makefile
make build

# 安装到系统
make install

# 运行
./soke-cli --help
```

### 添加新命令

参考 `CLAUDE.md` 中的快捷命令开发流程：

1. 在 `shortcuts/<module>/` 创建命令定义文件
2. 在 `shortcuts/<module>/shortcuts.go` 的 `Shortcuts()` 函数中注册
3. 在 `cmd/<module>/` 中创建模块命令
4. 在 `cmd/root.go` 中注册模块
5. 创建对应的 Skill 文档：`skills/<module>/SKILL.md`

详细步骤请参考 [README.md - 从接口封装到发布的完整流程](README.md#从接口封装到发布的完整流程)

### 测试 NPM 包安装流程

```bash
# 本地测试安装脚本
node scripts/install.js

# 测试运行脚本
node scripts/run.js --help

# 本地打包测试
npm pack
# 会生成 sokeai-cli-1.0.0.tgz

# 在其他目录测试安装
npm install -g /path/to/sokeai-cli-1.0.0.tgz
```

### 发布新版本

```bash
# 1. 更新版本号
npm version patch  # 1.0.0 -> 1.0.1
# 或
npm version minor  # 1.0.0 -> 1.1.0

# 2. 一键发布（推荐）
./scripts/release.sh

# 或手动执行各步骤：

# 2.1 编译所有平台
./scripts/build-binaries.sh

# 2.2 创建标签
git tag v1.0.0
git push origin v1.0.0

# 2.3 上传到 GitHub Releases
gh release create v1.0.0 \
  bin/soke-cli-darwin-amd64 \
  bin/soke-cli-darwin-arm64 \
  bin/soke-cli-linux-amd64 \
  bin/soke-cli-windows-amd64.exe \
  --title "v1.0.0" \
  --notes "Release v1.0.0"

# 2.4 发布到 NPM
npm publish --access public
```

---

## 完整开发流程

### 1. 接口封装为 CLI 命令

```bash
# 创建新模块目录
mkdir -p shortcuts/your-module
mkdir -p cmd/your_module

# 创建命令定义文件
touch shortcuts/your-module/list_items.go
touch shortcuts/your-module/shortcuts.go
touch cmd/your_module/your_module.go
```

详细代码示例请参考 [README.md - 接口封装为 CLI 命令](README.md#1️⃣-接口封装为-cli-命令)

### 2. 创建 AI Agent Skill

```bash
# 创建 Skill 目录和文档
mkdir -p skills/your-module
touch skills/your-module/SKILL.md
```

详细文档格式请参考 [README.md - 创建 AI Agent Skill](README.md#2️⃣-创建-ai-agent-skill)

### 3. 本地测试

**一键测试（推荐）**

```bash
# 1. 编译、安装到全局、运行测试
bash ./scripts/local-test.sh

# 2. 链接 Skills 到 AI Agent
bash ./scripts/link-skills.sh

# 3. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询所有项目"
```

**手动测试**

```bash
# 1. 编译项目
go build -o soke-cli main.go

# 2. 测试命令
./soke-cli your-module +list-items --help

# 3. 安装到全局
sudo cp soke-cli /usr/local/bin/

# 4. 测试 Skill
mkdir -p ~/.claude/skills/your-module
cp skills/your-module/SKILL.md ~/.claude/skills/your-module/
```

**测试检查清单：**
- [ ] CLI 命令能正常执行
- [ ] 参数验证正确
- [ ] API 调用成功
- [ ] 输出格式正确
- [ ] AI Agent 能识别意图
- [ ] AI Agent 能正确调用命令
- [ ] 错误处理友好

### 4. 编译二进制文件

```bash
# 一键编译所有平台
bash ./scripts/build-binaries.sh

# 或手动编译单个平台
GOOS=darwin GOARCH=amd64 go build -o bin/soke-cli-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o bin/soke-cli-darwin-arm64 main.go
GOOS=linux GOARCH=amd64 go build -o bin/soke-cli-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o bin/soke-cli-windows-amd64.exe main.go
```

### 5. 发布到 NPM 和 GitHub

**准备工作**

```bash
# 1. 确保已登录 NPM
npm login

# 2. 确保已登录 GitHub CLI
gh auth login

# 3. 确保代码已提交
git add .
git commit -m "feat: 添加新模块"
git push
```

**一键发布**

```bash
# 1. 更新版本号
npm version patch

# 2. 执行发布脚本
bash ./scripts/release.sh
```

**手动发布**

```bash
# 1. 更新版本号
npm version patch

# 2. 编译所有平台
bash ./scripts/build-binaries.sh

# 3. 创建并推送标签
VERSION=$(node -p "require('./package.json').version")
git tag v$VERSION
git push origin v$VERSION

# 4. 创建 GitHub Release
gh release create v$VERSION \
  bin/soke-cli-darwin-amd64 \
  bin/soke-cli-darwin-arm64 \
  bin/soke-cli-linux-amd64 \
  bin/soke-cli-windows-amd64.exe \
  --title "v$VERSION" \
  --notes "Release v$VERSION"

# 5. 发布到 NPM
npm publish --access public
```

**发布后验证**

```bash
# 1. 测试 NPM 安装
npm install -g @sokeai/cli@latest

# 2. 验证版本
soke-cli --version

# 3. 测试命令
soke-cli your-module +list-items --help

# 4. 测试 Skills 安装
npx skills add liuchenlong1111/soke-cli -y -g

# 5. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询所有项目"
```

---

## 故障排查

### 安装失败

```bash
# 检查网络连接
curl -I https://github.com

# 手动下载二进制文件
# 访问 https://github.com/liuchenlong1111/soke-cli/releases
# 下载对应平台的文件到 node_modules/@sokeai/cli/bin/

# 设置权限（macOS/Linux）
chmod +x node_modules/@sokeai/cli/bin/soke-cli
```

### 命令找不到

```bash
# 检查 NPM 全局路径
npm config get prefix

# 检查 PATH 环境变量
echo $PATH

# 重新安装
npm uninstall -g @sokeai/cli
npm install -g @sokeai/cli
```

### Skills 无法识别

```bash
# 检查 Skills 目录
ls -la ~/.claude/skills/

# 重新链接 Skills
bash ./scripts/link-skills.sh

# 重启 AI Agent
```

### 编译失败

```bash
# 检查 Go 版本
go version

# 更新依赖
go mod tidy

# 清理缓存
go clean -cache

# 重新编译
go build -o soke-cli main.go
```

---

## 文件说明

| 文件 | 说明 |
|------|------|
| `scripts/run.js` | NPM bin 入口，调用真正的二进制文件 |
| `scripts/install.js` | postinstall 钩子，从 GitHub 下载对应平台的二进制文件 |
| `scripts/build-binaries.sh` | 编译所有平台的二进制文件 |
| `scripts/release.sh` | 一键发布脚本（编译+标签+GitHub+NPM） |
| `scripts/local-test.sh` | 本地测试脚本（编译+安装+测试） |
| `scripts/link-skills.sh` | 链接 Skills 到 AI Agent |
| `npm.md` | NPM 包发布与使用的详细文档 |

## 注意事项

1. **GitHub 仓库地址**: 需要在 `scripts/install.js` 中配置正确的仓库地址
2. **版本一致性**: `package.json` 的 version 必须与 Git tag 一致
3. **权限问题**: 发布到 NPM 需要先 `npm login`
4. **GitHub CLI**: 需要安装 `gh` 并登录 `gh auth login`
5. **依赖管理**: 项目使用 Go Modules，运行 `go mod tidy` 整理依赖

## 相关链接

- **NPM 包**: https://www.npmjs.com/package/@sokeai/cli
- **GitHub 仓库**: https://github.com/liuchenlong1111/soke-cli
- **AI Agent Skills**: [skills/README.md](skills/README.md)
- **授客AI 开放平台**: https://opendev.soke.cn
- **详细开发指南**: [README.md](README.md)
- **本地测试指南**: [docs/LOCAL_TESTING.md](docs/LOCAL_TESTING.md)
