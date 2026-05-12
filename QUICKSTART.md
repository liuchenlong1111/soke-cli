# soke-cli 快速使用指南

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

### 添加新命令

参考 `CLAUDE.md` 中的快捷命令开发流程：

1. 在 `shortcuts/<module>/` 创建命令定义文件
2. 在 `shortcuts/<module>/shortcuts.go` 的 `Shortcuts()` 函数中注册
3. 完成！`cmd/` 层会自动注册

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

如果 sokeclaw 内仍提示 `command not found: soke-cli`，说明 GUI/运行时没有加载你的 shell 环境（常见于 NVM）。可用以下方式确认并修复：

```bash
which soke-cli
soke-cli -v
```

修复方式二选一：
- 在 sokeclaw 执行命令时使用 `which soke-cli` 输出的绝对路径
- 将 `soke-cli` 放到更通用的 PATH 目录（如 `~/.local/bin` 或 `/usr/local/bin`），确保 GUI 也能找到

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

详细文档：[skills/README.md](skills/README.md)

## 文件说明

| 文件 | 说明 |
|------|------|
| `scripts/run.js` | NPM bin 入口，调用真正的二进制文件 |
| `scripts/install.js` | postinstall 钩子，从 GitHub 下载对应平台的二进制文件 |
| `scripts/build-binaries.sh` | 编译所有平台的二进制文件 |
| `scripts/release.sh` | 一键发布脚本（编译+标签+GitHub+NPM） |
| `npm.md` | NPM 包发布与使用的详细文档 |

## 注意事项

1. **GitHub 仓库地址**: 需要在 `scripts/install.js` 中配置正确的仓库地址
2. **版本一致性**: `package.json` 的 version 必须与 Git tag 一致
3. **权限问题**: 发布到 NPM 需要先 `npm login`
4. **GitHub CLI**: 需要安装 `gh` 并登录 `gh auth login`

## 故障排查

### 安装失败

```bash
# 检查网络连接
curl -I https://github.com

# 手动下载二进制文件
# 访问 https://github.com/sokeai/soke-cli/releases
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

## 相关链接

- NPM 包: https://www.npmjs.com/package/@sokeai/cli
- GitHub 仓库: https://github.com/liuchenlong1111/soke-cli
- AI Agent Skills: [skills/README.md](skills/README.md)
- 授客AI 开放平台: https://opendev.soke.cn
