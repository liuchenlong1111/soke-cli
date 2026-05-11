# soke-cli NPM 包发布与使用流程

## 架构设计

soke-cli 采用 **Go 编译 + NPM 分发** 的混合架构：
- **核心代码**：Go 语言编写，编译为跨平台二进制文件
- **分发方式**：通过 NPM 包管理器分发
- **安装机制**：postinstall 钩子根据用户平台自动下载对应二进制文件

## 用户使用流程

### 1. 安装

```bash
# 全局安装
npm install -g @sokeai/cli

# 或项目内安装
npm install --save-dev @sokeai/cli
```

### 2. 安装过程（自动执行）

```
npm install @sokeai/cli
    ↓
下载 NPM 包（只包含 scripts/ 目录，不含二进制文件）
    ↓
触发 postinstall 钩子
    ↓
执行 scripts/install.js
    ↓
检测用户操作系统和架构
    - macOS Intel: darwin-amd64
    - macOS Apple Silicon: darwin-arm64
    - Linux: linux-amd64
    - Windows: windows-amd64
    ↓
从 GitHub Releases 下载对应平台的二进制文件
    URL: https://github.com/sokeai/soke-cli/releases/download/v{version}/soke-cli-{platform}-{arch}
    ↓
保存到 node_modules/@sokeai/cli/bin/soke-cli
    ↓
设置可执行权限 (chmod +x)
    ↓
安装完成
```

### 3. 使用命令

```bash
# 全局安装后直接使用
soke-cli config init
soke-cli auth login
soke-cli course +list-courses --start-time 1672502400000 --end-time 1704038400000

# 项目内安装使用 npx
npx soke-cli config show

# 或在 package.json 中配置脚本
{
  "scripts": {
    "soke": "soke-cli"
  }
}
npm run soke -- course +list-courses --page 1
```

### 4. 命令执行流程

```
用户输入: soke-cli course +list-courses --page 1
    ↓
Shell 查找命令: /usr/local/bin/soke-cli (符号链接)
    ↓
指向: node_modules/@sokeai/cli/scripts/run.js
    ↓
run.js 调用真正的二进制文件: node_modules/@sokeai/cli/bin/soke-cli
    ↓
Go 二进制文件启动
    ↓
main.go → cmd/root.go → cmd/course/course.go
    ↓
遍历 shortcuts/course/Shortcuts()
    ↓
找到 "+list-courses" 命令
    ↓
解析参数: --page 1
    ↓
执行 shortcuts/course/course_list_courses.go 的 Execute 函数
    ↓
调用授客AI API: GET /course/course/list
    ↓
格式化输出结果
```

## 开发者发布流程

### 1. 编译多平台二进制文件

```bash
# 使用编译脚本
./scripts/build-binaries.sh

# 或手动编译
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/soke-cli-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/soke-cli-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/soke-cli-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/soke-cli-windows-amd64.exe .
```

### 2. 上传二进制文件到 GitHub Releases

```bash
# 创建 Git 标签
git tag v1.0.0
git push origin v1.0.0

# 使用 GitHub CLI 创建 Release 并上传文件
gh release create v1.0.0 \
  bin/soke-cli-darwin-amd64 \
  bin/soke-cli-darwin-arm64 \
  bin/soke-cli-linux-amd64 \
  bin/soke-cli-windows-amd64.exe \
  --title "v1.0.0" \
  --notes "Release notes here"
```

### 3. 更新 package.json 版本号

```bash
npm version patch  # 1.0.0 -> 1.0.1
# 或
npm version minor  # 1.0.0 -> 1.1.0
# 或
npm version major  # 1.0.0 -> 2.0.0
```

### 4. 发布到 NPM

```bash
# 登录 NPM（首次）
npm login

# 发布
npm publish --access public

# 如果是 scoped package (@sokeai/cli)，需要 --access public
```

### 5. 验证安装

```bash
# 在新环境测试安装
npm install -g @sokeai/cli

# 验证命令
soke-cli --version
soke-cli --help
```

## 文件结构

```
soke-cli/
├── main.go                      # Go 主入口
├── cmd/                         # 命令定义
├── internal/                    # 核心逻辑
├── shortcuts/                   # 快捷命令元数据
├── scripts/
│   ├── install.js              # NPM postinstall 钩子（下载二进制文件）
│   ├── run.js                  # NPM bin 入口（调用二进制文件）
│   └── build-binaries.sh       # 多平台编译脚本
├── bin/                        # 编译产物（不提交到 Git）
│   ├── soke-cli-darwin-amd64
│   ├── soke-cli-darwin-arm64
│   ├── soke-cli-linux-amd64
│   └── soke-cli-windows-amd64.exe
├── package.json                # NPM 包配置
├── go.mod                      # Go 依赖管理
└── Makefile                    # 本地开发编译脚本
```

## NPM 包配置说明

```json
{
  "name": "@sokeai/cli",
  "version": "1.0.0",
  "description": "授客AI官方CLI工具",
  "bin": {
    "soke-cli": "scripts/run.js"
  },
  "scripts": {
    "postinstall": "node scripts/install.js"
  },
  "files": [
    "scripts/"
  ],
  "repository": {
    "type": "git",
    "url": "https://codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli.git"
  },
  "keywords": ["soke", "cli", "授客AI"],
  "license": "MIT"
}
```

**关键配置**：
- `bin`: 定义全局命令入口为 `scripts/run.js`
- `postinstall`: 安装后自动执行 `scripts/install.js` 下载二进制文件
- `files`: 只打包 `scripts/` 目录到 NPM，不包含二进制文件（减小包体积）

## 二进制文件存储位置

| 环境 | 二进制文件位置 | 说明 |
|------|---------------|------|
| **开发环境** | `/Users/edy/www/soke-cli/soke-cli` | 本地编译产物 |
| **GitHub Releases** | `https://github.com/sokeai/soke-cli/releases/download/v1.0.0/soke-cli-darwin-arm64` | 发布的二进制文件 |
| **NPM 全局安装** | `/usr/local/lib/node_modules/@sokeai/cli/bin/soke-cli` | postinstall 下载后保存位置 |
| **NPM 项目安装** | `./node_modules/@sokeai/cli/bin/soke-cli` | 项目内安装位置 |

## 优势

1. **包体积小**：NPM 包只包含脚本文件（~10KB），不包含二进制文件
2. **按需下载**：用户安装时只下载自己平台的二进制文件（~7MB）
3. **跨平台支持**：自动适配 macOS、Linux、Windows
4. **开发体验好**：用户通过熟悉的 `npm install` 安装，无需手动下载
5. **版本管理**：利用 NPM 的版本管理机制

## 类似项目参考

- **esbuild**: Go 编写，通过 NPM 分发
- **swc**: Rust 编写，通过 NPM 分发
- **prisma**: Rust 编写，通过 NPM 分发

## 注意事项

1. **GitHub Releases URL 必须公开访问**，否则用户无法下载二进制文件
2. **版本号必须一致**：package.json 的 version 必须与 GitHub Release 的 tag 一致
3. **二进制文件命名规范**：`soke-cli-{platform}-{arch}` 格式
4. **权限问题**：Linux/macOS 需要设置可执行权限 `chmod +x`
5. **网络问题**：国内用户可能需要配置镜像或使用 CDN 加速
