# 授客AI官方CLI工具 (soke-cli)

`soke-cli` 是授客AI官方提供的命令行工具，旨在帮助开发者和系统管理员更便捷地通过命令行与授客AI开放平台进行交互。该工具使用 Go 语言开发，支持多平台，并提供了 NPM 包的安装方式。

## 目录

- [核心能力](#核心能力)
- [安装方法](#安装方法)
- [快速开始](#快速开始)
- [开发指南](#开发指南)
  - [从接口封装到发布的完整流程](#从接口封装到发布的完整流程)
  - [本地开发测试](#本地开发测试)
- [项目结构](#项目结构)
- [相关文档](#相关文档)

## 核心能力

`soke-cli` 提供了以下核心能力：

1. **认证与配置管理**
   - 支持通过 OAuth 设备码授权流程（Device Flow）进行快速登录。
   - 灵活的配置管理，支持多租户、多环境（AppKey, AppSecret, API地址, CorpID等）。

2. **通用 API 调用**
   - 提供 `api` 基础命令，支持直接发起任意 GET、POST、PUT、DELETE 等 HTTP 请求到开放平台接口。

3. **业务模块快捷命令**
   - 针对授客AI核心业务场景提供了丰富的快捷命令集：
     - **通讯录 (contact)**: 部门与用户查询、管理、搜索等。
     - **课程 (course)**: 课程列表、分类、学习记录、人脸识别记录等。
     - **考试 (exam)**: 考试列表、分类、考试成绩与记录查询。
     - **学习档案 (learning-profile)**: 学员学习档案查询、学习情况统计。
     - **学习地图 (learning-map)**: 学习地图、阶段、任务查询与分配。
     - **证书 (certificate)**: 证书发放记录、分类等。
     - **学分 (credit) & 积分 (point)**: 学分/积分日志及用户情况查询。
     - **培训 (training)**: 线下培训查询与分配。
     - **新闻公告 (news)**: 资讯列表与详情。

4. **AI Agent Skills**
   - 提供 AI Agent 技能（Skills），使 AI 助手能够自动发现和调用 CLI 功能。
   - 支持自然语言交互，无需记忆复杂的命令参数。

## 安装方法

### 方式一：通过 NPM 安装 (推荐)
如果你本地已安装 Node.js，可以直接使用 npm 全局安装：
```bash
npm install -g @sokeai/cli@latest
```

### 方式二：通过源码编译安装 (需要 Go 环境)
```bash
# 1. 克隆/下载本仓库代码
git clone <repository_url>
cd soke-cli

# 2. 编译并安装到 /usr/local/bin
make install
```

## AI Agent Skills

`soke-cli` 提供了 AI Agent 技能（Skills），使 AI 助手能够自动发现和调用 CLI 功能，无需记忆复杂的命令参数。

### 安装 Skills

```bash
# 安装 CLI 后，安装 AI Agent Skills
npx skills add liuchenlong1111/soke-cli -y -g
```

### 使用方式

安装 Skills 后，在支持 Skills 的 AI Agent（如 Claude Code）中，可以直接用自然语言描述需求：

- **"查询考试成绩"** - AI 会自动调用 `soke-cli exam +get-exam-user`
- **"列出最近的考试"** - AI 会自动调用 `soke-cli exam +list-exams`
- **"查看考试分类"** - AI 会自动调用 `soke-cli exam +list-categories`

AI Agent 会自动：
1. 识别你的意图
2. 选择合适的命令
3. 提示你提供必要的参数
4. 执行命令并展示结果

### 可用的 Skills

- **soke-shared**: 配置初始化、用户认证、权限处理等基础功能
- **soke-exam**: 考试管理（查询考试、考试成绩、考试分类）
- **soke-course**: 课程管理（查询课程、课程分类、学习记录）
- **soke-learning-profile**: 学习档案查询（查询学员学习档案、学习情况统计）
- 更多业务模块的 Skills 正在开发中...

详细文档：[skills/README.md](skills/README.md)

## 快速开始

### 1. 安装 CLI 和 Skills

```bash
# 安装 CLI
npm install -g @sokeai/cli@latest  

# 安装 AI Agent Skills（可选，推荐）
npx skills add liuchenlong1111/soke-cli -y -g
```

### 2. 初始化配置
首次使用前，需要配置你的开放平台凭证：
```bash
soke-cli config init
```
根据提示输入 `app_key`, `app_secret`, API地址(默认: `https://opendev.soke.cn`), `corpid` 和 `dept_user_id`。

可以使用 `soke-cli config show` 来查看当前配置。

### 3. 用户登录授权
配置完成后，进行设备授权登录：
```bash
soke-cli auth login
```
命令会返回一个设备码和链接，请在浏览器中打开链接并输入设备码完成授权。成功后即可调用需要用户授权的接口。
如果需要退出登录，可以执行 `soke-cli auth logout`。

### 4. 开始使用

**方式一：通过 AI Agent（推荐）**

如果已安装 Skills，在 AI Agent 中直接用自然语言：
```
"查询考试成绩"
"列出最近的考试"
```

**方式二：直接使用命令行**

```bash
# 查看课程列表
soke-cli course +list-courses --start-time 1672502400000 --end-time 1704038400000

# 查询考试成绩
soke-cli exam +get-exam-user --exam-id exam123 --dept-user-id user456
```

## 使用示例

### 调用快捷业务命令
以课程模块为例，你可以使用以下命令获取课程列表：
```bash
# 查看 course 模块的帮助和可用命令
soke-cli course --help

# 获取课程列表 (使用必要的参数)
soke-cli course +list-courses --start-time 1672502400000 --end-time 1704038400000 --page 1 --page-size 10
```

### 使用通用 API 命令
如果你需要调用的接口尚未包含在快捷命令中，可以直接使用 `api` 命令：
```bash
# 发起 GET 请求
soke-cli api GET /users/me

# 发起 POST 请求，携带 JSON 参数
soke-cli api POST /some/endpoint --data '{"key": "value"}'
```

## 开发指南

### 从接口封装到发布的完整流程

#### 1️⃣ 接口封装为 CLI 命令

**步骤 1: 在 shortcuts 目录下定义接口元数据**

```bash
# 创建新模块目录
mkdir -p shortcuts/your-module

# 创建命令定义文件
touch shortcuts/your-module/list_items.go
```

**示例：定义一个查询列表的命令**

```go
// shortcuts/your-module/list_items.go
package yourmodule

import "soke-cli/internal/client"

// ListItemsShortcut 定义查询列表命令
func ListItemsShortcut() client.Shortcut {
    return client.Shortcut{
        Name:        "list-items",           // 命令名称
        Description: "查询项目列表",          // 命令描述
        Method:      "GET",                  // HTTP 方法
        Path:        "/api/v1/items",        // API 路径
        Params: []client.Param{              // 参数定义
            {
                Name:        "page",
                Type:        "int",
                Description: "页码",
                Required:    false,
                Default:     "1",
            },
            {
                Name:        "page_size",
                Type:        "int",
                Description: "每页数量",
                Required:    false,
                Default:     "10",
            },
            {
                Name:        "keyword",
                Type:        "string",
                Description: "搜索关键词",
                Required:    false,
            },
        },
    }
}
```

**步骤 2: 注册命令到模块**

```go
// shortcuts/your-module/shortcuts.go
package yourmodule

import "soke-cli/internal/client"

// Shortcuts 返回该模块的所有快捷命令
func Shortcuts() []client.Shortcut {
    return []client.Shortcut{
        ListItemsShortcut(),
        // 添加更多命令...
    }
}
```

**步骤 3: 在 cmd 层注册模块**

```go
// cmd/your_module/your_module.go
package yourmodule

import (
    "soke-cli/internal/client"
    yourmodule "soke-cli/shortcuts/your-module"
    "github.com/spf13/cobra"
)

// YourModuleCmd 模块根命令
var YourModuleCmd = &cobra.Command{
    Use:   "your-module",
    Short: "你的模块管理",
    Long:  "管理你的模块相关功能",
}

func init() {
    // 自动注册所有快捷命令
    client.RegisterShortcuts(YourModuleCmd, yourmodule.Shortcuts())
}
```

**步骤 4: 在根命令中注册模块**

```go
// cmd/root.go
import (
    yourmodule "soke-cli/cmd/your_module"
)

func init() {
    // 注册模块命令
    rootCmd.AddCommand(yourmodule.YourModuleCmd)
}
```

#### 2️⃣ 创建 AI Agent Skill

**步骤 1: 创建 Skill 目录和文档**

```bash
# 创建 Skill 目录
mkdir -p skills/your-module

# 创建 Skill 定义文件
touch skills/your-module/SKILL.md
```

**步骤 2: 编写 Skill 文档**

```markdown
<!-- skills/your-module/SKILL.md -->
---
name: your-module
description: 你的模块管理：查询项目列表、创建项目、更新项目等
trigger: 当用户需要查询项目、管理项目时使用
---

# 你的模块管理 Skill

## 功能说明

提供项目管理相关功能，包括：
- 查询项目列表
- 创建新项目
- 更新项目信息
- 删除项目

## 使用场景

- 用户询问："查询所有项目"
- 用户询问："创建一个新项目"
- 用户询问："更新项目信息"

## 可用命令

### 查询项目列表

\`\`\`bash
soke-cli your-module +list-items [选项]
\`\`\`

**参数说明：**
- `--page <number>`: 页码（可选，默认：1）
- `--page-size <number>`: 每页数量（可选，默认：10）
- `--keyword <string>`: 搜索关键词（可选）

**使用示例：**
\`\`\`bash
# 查询第一页
soke-cli your-module +list-items

# 搜索包含"测试"的项目
soke-cli your-module +list-items --keyword "测试"
\`\`\`

## 交互流程

1. 识别用户意图（查询/创建/更新/删除）
2. 提示用户提供必要参数
3. 执行对应的 CLI 命令
4. 解析并展示结果

## 错误处理

- 如果用户未登录，提示执行 `soke-cli auth login`
- 如果参数缺失，提示用户补充必要参数
- 如果 API 返回错误，展示友好的错误信息
```

**步骤 3: 更新 skills/README.md**

在 `skills/README.md` 中添加新 Skill 的说明。

#### 3️⃣ 本地测试 Skill

**方法一：一键测试脚本（推荐）**

```bash
# 1. 编译、安装到全局、运行测试
bash ./scripts/local-test.sh
# 提示时输入 'y' 更新全局安装

# 2. 链接 Skills 到 AI Agent
bash ./scripts/link-skills.sh
# 选择 'all' 链接到所有目录

# 3. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询所有项目"
```

**方法二：手动测试**

```bash
# 1. 编译项目
go build -o soke-cli main.go

# 2. 测试命令是否正常工作
./soke-cli your-module +list-items --help

# 3. 安装到全局
sudo cp soke-cli /usr/local/bin/

# 4. 链接 Skill 到 Claude
mkdir -p ~/.claude/skills/your-module
cp skills/your-module/SKILL.md ~/.claude/skills/your-module/

# 5. 在 Claude Code 中测试
# 输入："查询所有项目"
```

**测试检查清单：**
- [ ] CLI 命令能正常执行
- [ ] 参数验证正确
- [ ] API 调用成功
- [ ] 输出格式正确
- [ ] AI Agent 能识别意图
- [ ] AI Agent 能正确调用命令
- [ ] 错误处理友好

#### 4️⃣ 用 Go 编译二进制文件

**编译单平台二进制**

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/soke-cli-darwin-amd64 main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/soke-cli-darwin-arm64 main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/soke-cli-linux-amd64 main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/soke-cli-windows-amd64.exe main.go
```

**一键编译所有平台**

```bash
# 使用编译脚本
bash ./scripts/build-binaries.sh

# 编译结果在 bin/ 目录下
ls -lh bin/
```

**编译优化选项**

```bash
# 减小二进制文件大小
go build -ldflags="-s -w" -o soke-cli main.go

# 添加版本信息
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.Version=$VERSION" -o soke-cli main.go
```

#### 5️⃣ 发布到 NPM 和 GitHub

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

**发布流程**

**方法一：一键发布（推荐）**

```bash
# 1. 更新版本号
npm version patch  # 1.0.0 -> 1.0.1
# 或
npm version minor  # 1.0.0 -> 1.1.0
# 或
npm version major  # 1.0.0 -> 2.0.0

# 2. 执行一键发布脚本
bash ./scripts/release.sh

# 脚本会自动完成：
# - 编译所有平台二进制文件
# - 创建 Git 标签
# - 推送标签到 GitHub
# - 创建 GitHub Release 并上传二进制文件
# - 发布到 NPM
```

**方法二：手动发布**

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

# 6. 验证发布
npm view @sokeai/cli version
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

### 本地开发测试

**🚀 快速测试流程（推荐）**

```bash
# 1. 一键测试：编译、安装到全局、运行测试（一条命令）
bash ./scripts/local-test.sh
# 提示时输入 'y' 更新全局安装

# 2. 链接 Skills 到 Claude（首次需要）
bash ./scripts/link-skills.sh
# 选择 'all' 链接到所有目录

# 3. 在 AI Agent 中测试
# 打开 Claude Code，输入："查询张三的学习档案"
```

**📖 详细指南：** [docs/LOCAL_TESTING.md](docs/LOCAL_TESTING.md)

**完整开发流程**

```bash
# 1. 修改代码后，运行本地测试（会提示是否更新全局安装）
bash ./scripts/local-test.sh

# 2. 运行完整 E2E 测试
bash ./scripts/e2e-test.sh

# 3. 测试 Skills（需要先更新全局安装）
npx skills add liuchenlong1111/soke-cli -y -g

# 4. 在 AI Agent 中测试自然语言交互
# 例如："查询张三的学习档案"
```

**运行完整 E2E 测试**

```bash
# 测试所有模块
bash ./scripts/e2e-test.sh

# 测试特定模块
bash ./scripts/e2e-test.sh learning-profile
bash ./scripts/e2e-test.sh contact
```

## 项目结构

```
soke-cli/
├── cmd/                    # CLI 子命令定义（按业务模块划分）
│   ├── root.go            # 根命令和命令注册
│   ├── auth/              # 认证相关命令
│   ├── config/            # 配置管理命令
│   ├── api/               # 通用 API 调用命令
│   └── your_module/       # 你的业务模块命令
├── internal/              # 核心逻辑（不对外暴露）
│   ├── auth/              # 认证逻辑（OAuth Device Flow）
│   ├── client/            # API 客户端封装
│   ├── core/              # 配置管理
│   └── output/            # 输出格式化
├── shortcuts/             # 业务模块快捷命令元数据定义
│   └── your-module/       # 各模块的 API 路径、参数、默认值等
├── skills/                # AI Agent Skills 定义
│   ├── README.md          # Skills 使用说明
│   └── your-module/       # 各模块的 Skill 文档
│       └── SKILL.md
├── scripts/               # 构建和发布脚本
│   ├── install.js         # NPM 安装时下载二进制文件
│   ├── run.js             # NPM 运行时入口
│   ├── build-binaries.sh  # 编译所有平台二进制
│   ├── release.sh         # 一键发布脚本
│   ├── local-test.sh      # 本地测试脚本
│   └── link-skills.sh     # 链接 Skills 到 AI Agent
├── main.go                # CLI 主入口
├── go.mod                 # Go 依赖管理
├── package.json           # NPM 包配置
├── Makefile               # 编译和安装脚本
└── README.md              # 项目说明文档
```

## 相关文档

- **[QUICKSTART.md](QUICKSTART.md)** - 快速开始指南
- **[CLAUDE.md](CLAUDE.md)** - 项目架构和开发规范
- **[docs/LOCAL_TESTING.md](docs/LOCAL_TESTING.md)** - 本地测试详细指南
- **[skills/README.md](skills/README.md)** - AI Agent Skills 使用说明
- **[npm.md](npm.md)** - NPM 包发布详细文档
