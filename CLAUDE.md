# 授客AI CLI工具 (soke-cli)

## 项目概述

`soke-cli` 是授客AI官方提供的命令行工具，使用 Go 语言开发，通过 NPM 包分发。该工具帮助开发者和系统管理员通过命令行与授客AI开放平台进行交互。

**技术栈**:
- 语言: Go 1.25.0
- CLI框架: Cobra (github.com/spf13/cobra)
- 分发方式: NPM 包 (@sokeai/cli)
- 平台支持: 跨平台 (通过 Go 编译)

## 项目结构

```
soke-cli/
├── cmd/                    # CLI 子命令定义（按业务模块划分）
│   ├── root.go            # 根命令和命令注册
│   ├── auth/              # 认证相关命令
│   ├── config/            # 配置管理命令
│   ├── api/               # 通用 API 调用命令
│   ├── contact/           # 通讯录模块
│   ├── course/            # 课程模块
│   ├── exam/              # 考试模块
│   ├── learning_map/      # 学习地图模块
│   ├── certificate/       # 证书模块
│   ├── credit/            # 学分模块
│   ├── point/             # 积分模块
│   ├── training/          # 培训模块
│   └── news/              # 新闻公告模块
├── internal/              # 核心逻辑（不对外暴露）
│   ├── auth/              # 认证逻辑（OAuth Device Flow）
│   ├── client/            # API 客户端封装
│   ├── core/              # 配置管理
│   └── output/            # 输出格式化
├── shortcuts/             # 业务模块快捷命令元数据定义
│   └── contact/           # 各模块的 API 路径、参数、默认值等
├── scripts/               # NPM 包安装与运行脚本
│   ├── install.js         # 安装时下载对应平台的二进制文件
│   └── run.js             # 运行时入口
├── main.go                # CLI 主入口
├── go.mod                 # Go 依赖管理
├── package.json           # NPM 包配置
└── Makefile               # 编译和安装脚本
```

## 核心功能模块

### 1. 认证与配置 (auth & config)
- **OAuth Device Flow**: 通过设备码授权流程进行用户登录
- **配置管理**: 支持多租户、多环境配置（AppKey, AppSecret, API地址, CorpID等）
- 命令: `soke-cli auth login/logout`, `soke-cli config init/show`

### 2. 通用 API 调用 (api)
- 支持直接发起 GET、POST、PUT、DELETE 等 HTTP 请求
- 命令: `soke-cli api <METHOD> <PATH> [--data <JSON>]`

### 3. 业务模块快捷命令
针对授客AI核心业务场景提供的快捷命令集：
- **通讯录 (contact)**: 部门与用户查询、管理
- **课程 (course)**: 课程列表、分类、学习记录、人脸识别记录
- **考试 (exam)**: 考试列表、分类、成绩与记录查询
- **学习地图 (learning-map)**: 学习地图、阶段、任务查询与分配
- **证书 (certificate)**: 证书发放记录、分类
- **学分 (credit) & 积分 (point)**: 学分/积分日志及用户情况查询
- **培训 (training)**: 线下培训查询与分配
- **新闻公告 (news)**: 资讯列表与详情

## 开发规范

### 代码规范
- **语言**: 所有代码注释使用中文
- **命名**: 遵循 Go 语言命名规范（驼峰命名）
- **错误处理**: 所有错误信息使用中文提示
- **包结构**: 
  - `cmd/` 下按业务模块划分子包
  - `internal/` 下按功能职责划分（auth, client, core, output）
  - 使用 `internal/` 防止外部包引用内部实现

### 添加新命令
1. **业务接口快捷命令**:
   - 在 `shortcuts/<module>/` 目录下添加对应的结构定义
   - 定义 API 路径、参数、默认值等元数据
   - 在 `cmd/<module>/` 中注册命令

2. **基础功能命令**:
   - 在 `cmd/` 下新建对应包
   - 在 `cmd/root.go` 中注册命令

### 构建与测试
```bash
# 整理依赖
go mod tidy

# 本地编译
go build -o soke-cli main.go

# 编译并安装到 /usr/local/bin
make install

# 运行测试
make test
```

### NPM 包发布流程
1. 编译各平台二进制文件（Linux, macOS, Windows）
2. 上传二进制文件到 CDN 或 GitHub Releases
3. 更新 `package.json` 版本号
4. 发布到 NPM: `npm publish`
5. 安装时 `scripts/install.js` 会根据平台下载对应二进制文件

## 配置文件

配置文件位置: `~/.soke-cli/config.json`

配置项:
- `app_key`: 开放平台应用 Key
- `app_secret`: 开放平台应用 Secret
- `api_base_url`: API 地址（默认: https://opendev.soke.cn）
- `corpid`: 企业 ID
- `dept_user_id`: 部门用户 ID
- `access_token`: 用户访问令牌（登录后自动保存）
- `refresh_token`: 刷新令牌

## 常用命令示例

```bash
# 初始化配置
soke-cli config init

# 查看当前配置
soke-cli config show

# 用户登录
soke-cli auth login

# 获取课程列表
soke-cli course +list-courses --start-time 1672502400000 --end-time 1704038400000 --page 1 --page-size 10

# 通用 API 调用
soke-cli api GET /users/me
soke-cli api POST /some/endpoint --data '{"key": "value"}'
```

## 注意事项

1. **认证流程**: 首次使用需要先 `config init` 配置凭证，然后 `auth login` 进行用户授权
2. **Token 管理**: access_token 会自动保存到配置文件，过期后需要重新登录
3. **错误处理**: 所有 API 调用失败时会输出中文错误信息
4. **跨平台兼容**: 使用 Go 编译确保跨平台兼容性，NPM 安装时自动下载对应平台二进制文件

## 相关链接

- 授客AI开放平台: https://opendev.soke.cn
- NPM 包: @sokeai/cli
- 代码仓库: codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli
