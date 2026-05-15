# 授客AI CLI 工具 - 快速上手指南

本文档将指导你从零开始安装 `soke-cli`，完成用户登录授权，并开始使用授客AI开放平台的各项功能。

---

## 📋 目录

- [前置要求](#前置要求)
- [安装步骤](#安装步骤)
- [用户登录授权](#用户登录授权)
- [验证安装](#验证安装)
- [常见问题](#常见问题)
- [下一步](#下一步)

---

## 前置要求

在开始之前，请确保你的系统满足以下要求：

### 必需条件
- **Node.js**: 版本 >= 14.0.0（用于 NPM 安装）
- **网络连接**: 能够访问授客AI开放平台 API（默认：`https://opendev.soke.cn`）
- **浏览器**: 用于完成 OAuth 授权流程

> 💡 **提示**: 使用 `soke-cli` 无需预先配置任何凭证，所有配置将在登录授权过程中自动完成。

---

## 安装步骤

### 通过 NPM 安装

这是推荐的安装方式，简单快捷，适合所有用户。

#### macOS / Linux 安装

打开终端，执行以下命令：

```bash
npm install -g @sokeai/cli
```

#### Windows 安装

**方式一：使用 PowerShell（推荐）**

1. 以管理员身份打开 PowerShell：
   - 按 `Win + X`，选择"Windows PowerShell (管理员)"
   - 或在开始菜单搜索"PowerShell"，右键选择"以管理员身份运行"

2. 执行安装命令：
   ```powershell
   npm install -g @sokeai/cli
   ```

**方式二：使用 CMD 命令提示符**

1. 以管理员身份打开命令提示符：
   - 按 `Win + R`，输入 `cmd`，按 `Ctrl + Shift + Enter`
   - 或在开始菜单搜索"cmd"，右键选择"以管理员身份运行"

2. 执行安装命令：
   ```cmd
   npm install -g @sokeai/cli
   ```

> 💡 **Windows 提示**: 
> - 如果提示"npm 不是内部或外部命令"，请先安装 [Node.js](https://nodejs.org/)
> - 安装过程中可能需要管理员权限
> - 如果遇到权限问题，确保以管理员身份运行终端

---

### 安装过程说明

无论使用哪个平台，安装过程都是自动化的：

1. NPM 会自动下载 `@sokeai/cli` 包
2. 安装脚本会根据你的操作系统（macOS/Linux/Windows）和架构（x64/arm64）自动下载对应的二进制文件
3. 二进制文件会被放置在对应的目录：
   - **macOS/Linux**: `$(npm root -g)/@sokeai/cli/bin/soke-cli`
   - **Windows**: `%APPDATA%\npm\node_modules\@sokeai\cli\bin\soke-cli.exe`
4. 如果你安装了 AI Agent（如 Claude Code），相关的 Skills 会自动同步到对应目录

---

### 验证安装

安装完成后，验证是否安装成功：

**macOS / Linux**:
```bash
soke-cli --version
```

**Windows (PowerShell / CMD)**:
```cmd
soke-cli --version
```

如果显示版本号（如 `1.0.15`），说明安装成功。

**如果提示"命令未找到"**，请参考[常见问题](#常见问题)章节。

---

## 用户登录授权

安装完成后，直接执行登录命令即可开始使用。无需预先配置，所有必要的配置信息将在授权过程中自动获取并保存。

### OAuth 浏览器授权流程

`soke-cli` 使用 OAuth 2.0 授权码流程（Authorization Code Flow），通过浏览器完成安全授权。

### 执行登录命令

```bash
soke-cli auth login
```

### 授权步骤

1. **自动打开浏览器**
   
   命令执行后，会自动在默认浏览器中打开授权页面。如果浏览器没有自动打开，会显示授权链接：
   ```
   ℹ️  未找到登录信息，开始授权流程
   
   请在浏览器中打开以下链接进行授权:
   http://127.0.0.1:8080/callback?redirect_uri=...
   
   等待授权中...
   ```

2. **浏览器授权**
   
   在浏览器授权页面中：
   - 如果未登录授客AI平台，先完成登录
   - 查看应用请求的权限列表
   - 勾选"一并开通审批的常用权限"（可选）
   - 点击"开通并授权"按钮

3. **授权成功**
   
   授权完成后，浏览器会显示成功页面，CLI 终端会显示：
   ```
   ✅ 已登录，Token 仍然有效
   企业 ID: corp_123456
   ```

4. **配置自动保存**
   
   授权成功后，以下信息会自动保存到 `~/.soke-cli/config.json`：
   - 用户访问令牌（`user_token`）
   - 企业 ID（`corpid`）
   - 部门用户 ID（`dept_user_id`）
   - Token 过期时间

### Token 有效期

- **访问令牌（Access Token）**: 默认有效期根据平台设置，通常为数小时到数天
- **自动检查**: 每次执行命令时，CLI 会自动检查 Token 是否有效
- **自动续期**: Token 过期后，再次执行 `soke-cli auth login` 即可重新授权

### 强制重新登录

如果需要切换账号或强制重新授权，可以使用 `--force` 参数：

```bash
soke-cli auth login --force
```

### 退出登录

如果需要退出登录，清除本地保存的令牌：

```bash
soke-cli auth logout
```

输出：
```
✓ 已退出登录
```

---

## 验证安装

完成登录授权后，让我们验证一切是否正常工作。

### 1. 查看配置信息

```bash
soke-cli config show
```

确认配置信息已自动保存：
```
当前配置:
  app_key: (未设置)
  app_secret: (未设置)
  api_base_url: 
  corpid: corp_123456
  dept_user_id: user_789012
  user_token: eyJh...xyz
  bot_token: (未设置)
```

> 💡 **说明**: `app_key` 和 `app_secret` 在当前版本中不需要手动配置，授权流程会自动处理所有必要的认证信息。

### 2. 测试 API 调用

尝试调用一个简单的 API：

```bash
soke-cli api GET /users/me
```

如果返回你的用户信息（JSON 格式），说明一切配置正确！

### 3. 使用业务命令

尝试使用一个业务快捷命令，例如查看考试列表：

```bash
soke-cli exam +list-exams --page 1 --page-size 10
```

---

## 常见问题

### Q1: 安装时提示 "网络连接失败"

**原因**: 无法从 GitHub Releases 下载二进制文件。

**解决方案**:
1. 检查网络连接
2. 如果在国内，可能需要配置代理
3. 尝试重新安装：`npm install -g @sokeai/cli`

### Q2: 执行 `soke-cli` 提示 "command not found" 或 "不是内部或外部命令"

**原因**: 二进制文件不在系统 PATH 中。

**解决方案**:

**macOS/Linux**:
```bash
# 检查安装位置
which soke-cli

# 如果没有输出，手动创建软链接
sudo ln -s $(npm root -g)/@sokeai/cli/bin/soke-cli /usr/local/bin/soke-cli
```

**Windows**:

1. **检查 NPM 全局路径是否在 PATH 中**
   
   打开 PowerShell，执行：
   ```powershell
   npm config get prefix
   ```
   
   记下输出的路径（例如：`C:\Users\YourName\AppData\Roaming\npm`）

2. **添加到系统 PATH**
   
   - 按 `Win + X`，选择"系统"
   - 点击"高级系统设置"
   - 点击"环境变量"
   - 在"用户变量"或"系统变量"中找到 `Path`
   - 点击"编辑"，添加上面的 NPM 路径
   - 点击"确定"保存

3. **重启终端**
   
   关闭并重新打开 PowerShell 或 CMD，再次尝试：
   ```cmd
   soke-cli --version
   ```

**或者使用完整路径运行**:
```cmd
# Windows 示例
%APPDATA%\npm\soke-cli --version

# 或者
C:\Users\YourName\AppData\Roaming\npm\soke-cli.exe --version
```

### Q3: 登录时浏览器没有自动打开

**原因**: 系统无法自动打开默认浏览器。

**解决方案**:
1. 复制终端显示的授权链接
2. 手动在浏览器中打开该链接
3. 完成授权流程

### Q4: 浏览器授权后，CLI 仍显示 "等待授权中..."

**原因**: 
- 浏览器授权未成功完成
- 本地回调服务器端口被占用
- 网络连接问题

**解决方案**:
1. 确认浏览器页面显示"授权成功"
2. 检查防火墙是否阻止了本地端口（127.0.0.1）
3. 按 `Ctrl+C` 取消，重新执行 `soke-cli auth login`

### Q5: API 调用返回 "401 Unauthorized"

**原因**: Token 已过期或无效。

**解决方案**:
```bash
# 重新登录
soke-cli auth logout
soke-cli auth login
```

### Q6: 如何在 AI Agent 中使用？

**解决方案**:

安装 AI Agent Skills：
```bash
npx skills add liuchenlong1111/soke-cli -y -g
```

然后在 AI Agent（如 Claude Code）中直接用自然语言：
- "查询考试成绩"
- "列出最近的考试"

AI 会自动调用相应的 `soke-cli` 命令。

---

## 下一步

恭喜！你已经成功安装并配置了 `soke-cli`。接下来你可以：

### 📚 学习更多命令

查看所有可用命令：
```bash
soke-cli --help
```

查看特定模块的命令：
```bash
soke-cli course --help
soke-cli exam --help
soke-cli contact --help
```

### 🚀 使用业务功能

- **课程管理**: `soke-cli course +list-courses`
- **考试管理**: `soke-cli exam +list-exams`
- **通讯录**: `soke-cli contact +list-users`
- **学习地图**: `soke-cli learning-map +list-maps`

### 🤖 安装 AI Agent Skills（推荐）

```bash
npx skills add liuchenlong1111/soke-cli -y -g
```

安装后，在 AI Agent 中可以用自然语言操作授客AI平台，无需记忆复杂命令。

### 📖 阅读详细文档

- [README.md](../README.md) - 项目概览
- [CLAUDE.md](../CLAUDE.md) - 项目架构和开发规范
- [skills/README.md](../skills/README.md) - AI Agent Skills 详细说明

---

## 获取帮助

如果遇到问题：

1. **查看帮助文档**: `soke-cli --help`
2. **查看配置**: `soke-cli config show`
3. **访问开放平台**: https://opendev.soke.cn
4. **提交 Issue**: https://github.com/liuchenlong1111/soke-cli/issues

---

**祝你使用愉快！** 🎉
