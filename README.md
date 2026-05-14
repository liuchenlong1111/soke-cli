# 授客AI官方CLI工具 (soke-cli)

`soke-cli` 是授客AI官方提供的命令行工具，旨在帮助开发者和系统管理员更便捷地通过命令行与授客AI开放平台进行交互。该工具使用 Go 语言开发，支持多平台，并提供了 NPM 包的安装方式。

## 核心能力

`soke-cli` 提供了以下核心能力：

1. **认证与配置管理**
   - 支持通过 OAuth 设备码授权流程（Device Flow）进行快速登录。
   - 灵活的配置管理，支持多租户、多环境（AppKey, AppSecret, API地址, CorpID等）。

2. **通用 API 调用**
   - 提供 `api` 基础命令，支持直接发起任意 GET、POST、PUT、DELETE 等 HTTP 请求到开放平台接口。

3. **业务模块快捷命令**
   - 针对授客AI核心业务场景提供了丰富的快捷命令集：
     - **通讯录 (contact)**: 部门与用户查询、管理等。
     - **课程 (course)**: 课程列表、分类、学习记录、人脸识别记录等。
     - **考试 (exam)**: 考试列表、分类、考试成绩与记录查询。
     - **学习地图 (learning-map)**: 学习地图、阶段、任务查询与分配。
     - **证书 (certificate)**: 证书发放记录、分类等。
     - **学分 (credit) & 积分 (point)**: 学分/积分日志及用户情况查询。
     - **培训 (training)**: 线下培训查询与分配。
     - **新闻公告 (news)**: 资讯列表与详情。

## 安装方法

### 方式一：通过 NPM 安装 (推荐)
如果你本地已安装 Node.js，可以直接使用 npm 全局安装：
```bash
npm install -g @sokeai/cli
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

## 项目结构说明

- `cmd/`: 存放所有 CLI 子命令的定义（按业务模块划分，如 `auth`, `course`, `exam` 等）。
- `internal/`: 存放核心逻辑，包括认证 (`auth`)、API 客户端封装 (`client`)、配置管理 (`core`) 以及输出格式化 (`output`) 等。
- `shortcuts/`: 定义了各个业务模块的快捷命令元数据（API路径、参数、默认值等），供 `cmd/` 动态注册命令使用。
- `scripts/`: NPM 包相关的安装与运行脚本。
- `main.go`: CLI 的主入口文件。

## 开发与贡献

1. 依赖管理：项目使用 Go Modules，可以运行 `go mod tidy` 整理依赖。
2. 添加新命令：
   - 业务接口建议在 `shortcuts/` 目录下添加对应的结构定义。
   - 基础功能可在 `cmd/` 下新建对应包并在 `cmd/root.go` 中注册。
3. 测试：运行 `make test` 进行单元测试。
