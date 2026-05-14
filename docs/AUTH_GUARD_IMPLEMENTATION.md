# 认证拦截机制实现文档

## 概述

本文档说明了 `soke-cli` 项目中实现的登录授权拦截机制，确保用户在执行需要认证的命令前必须先完成登录。

## 实现架构

```
用户命令输入
    ↓
main.go
    ↓
cmd/root.go - Execute()
    ↓
PersistentPreRunE (认证拦截点)
    ↓
internal/auth/guard.go - CheckAuth()
    ↓
【检查 token 是否存在和有效】
    ↓
如果未登录 → 返回错误并提示
如果已登录 → 继续执行命令
```

## 核心文件

### 1. internal/auth/guard.go

**职责**：认证拦截逻辑的核心实现

**关键函数**：

- `CheckAuth(ctx, commandName)` - 唯一的认证拦截点
  - 检查命令是否需要认证（跳过 auth、config、version、help）
  - 解析并验证 OAuth token
  - 检查 token 是否过期
  - 返回友好的错误提示

- `ResolveAuthToken(ctx)` - Token 解析
  - 使用进程缓存避免重复读取
  - 从 OAuth token 文件加载

- `IsTokenExpired(ctx)` - Token 过期检查
  - 使用 `TokenData.IsAccessTokenValid()` 方法
  - 自动处理 5 分钟缓冲期

- `ResetTokenCache()` - 重置缓存
  - 在登出时调用

- `UpdateCachedToken(token)` - 更新缓存
  - 在登录成功后调用

**跳过认证的命令**：
```go
var skipAuthCommands = map[string]bool{
    "auth":    true,
    "config":  true,
    "version": true,
    "help":    true,
}
```

### 2. cmd/root.go

**修改点**：在 `rootCmd` 的 `PersistentPreRunE` 中添加认证检查

```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    // 版本检查
    version.CheckForUpdates()

    // 认证检查（唯一拦截点）
    ctx := context.Background()
    commandName := cmd.Name()
    if cmd.Parent() != nil && cmd.Parent().Name() != "soke-cli" {
        commandName = cmd.Parent().Name()
    }

    if err := authpkg.CheckAuth(ctx, commandName); err != nil {
        errors.PrintHuman(os.Stderr, err)
        return err
    }

    return nil
},
```

### 3. cmd/auth/auth_command.go

**修改点**：

1. 导入 `authguard` 包
2. 在 `login` 成功后调用 `authguard.UpdateCachedToken()`
3. 在 `logout` 时调用 `authguard.ResetTokenCache()`

### 4. internal/errors/errors.go

**已存在**：完善的错误处理系统，支持：
- 错误分类（auth、api、validation 等）
- 友好的错误提示
- 建议操作
- 多级详细度输出

## 工作流程

### 未登录场景

```bash
$ soke-cli contact +list-departments

Error: [AUTH] 未登录，请先执行 soke-cli auth login
Hint: 运行 'soke-cli auth login' 完成登录后重试
Action: soke-cli auth login
```

**流程**：
1. 用户执行命令
2. `PersistentPreRunE` 触发
3. `CheckAuth()` 检查 token
4. Token 不存在 → 返回 `AuthError`
5. 错误被格式化输出
6. 命令执行被中断

### 已登录场景

```bash
$ soke-cli auth login
# 扫码登录成功

$ soke-cli contact +list-departments
# 正常执行，返回部门列表
```

**流程**：
1. 用户登录 → token 保存到 keychain
2. 用户执行命令
3. `PersistentPreRunE` 触发
4. `CheckAuth()` 检查 token
5. Token 存在且有效 → 通过检查
6. 命令正常执行

### Token 过期场景

```bash
$ soke-cli contact +list-departments

Error: [AUTH] 登录已过期，请重新登录
Hint: 运行 'soke-cli auth login' 重新登录
Action: soke-cli auth login
```

**流程**：
1. 用户执行命令
2. `CheckAuth()` 检查 token
3. Token 存在但已过期 → 返回 `AuthError`
4. 提示用户重新登录

## 设计优势

### 1. Fail-Fast 机制
- 在发起任何网络请求之前检查认证状态
- 避免无效的 API 调用
- 提供即时的错误反馈

### 2. 集中式拦截
- **唯一拦截点**：`cmd/root.go` 的 `PersistentPreRunE`
- 所有子命令自动继承认证检查
- 易于维护和调试

### 3. 进程级缓存
- Token 只在进程生命周期内加载一次
- 避免重复访问 keychain（~70ms/次）
- 提升命令执行速度

### 4. 友好的错误提示
- 清晰的错误信息
- 具体的操作建议
- 分级详细度（normal/verbose/debug）

### 5. 灵活的跳过机制
- 配置类命令无需认证
- 帮助命令无需认证
- 认证命令本身无需认证

## Token 管理

### Token 存储
- 使用系统 keychain 安全存储
- 支持跨平台（macOS/Linux/Windows）
- 自动加密保护

### Token 生命周期
1. **登录**：OAuth 扫码 → 获取 token → 保存到 keychain
2. **使用**：从 keychain 加载 → 缓存到进程 → 验证有效性
3. **刷新**：自动检测过期 → 提示重新登录
4. **登出**：删除 keychain 数据 → 清除进程缓存

### Token 验证
```go
// 5 分钟缓冲期，提前提示用户
func (t *TokenData) IsAccessTokenValid() bool {
    if t == nil || t.AccessToken == "" {
        return false
    }
    return time.Now().Before(t.ExpiresAt.Add(-5 * time.Minute))
}
```

## 测试验证

### 测试场景

1. ✅ **未登录拦截**
   ```bash
   $ soke-cli contact +list-departments
   Error: [AUTH] 未登录，请先执行 soke-cli auth login
   ```

2. ✅ **跳过认证命令**
   ```bash
   $ soke-cli version
   soke-cli 版本: 1.0.15
   
   $ soke-cli config show
   # 正常显示配置
   ```

3. ✅ **登录后正常执行**
   ```bash
   $ soke-cli auth login
   # 扫码成功
   
   $ soke-cli contact +list-departments
   # 正常执行
   ```

4. ✅ **登出清除缓存**
   ```bash
   $ soke-cli auth logout
   [OK] 已清除所有认证信息
   
   $ soke-cli contact +list-departments
   Error: [AUTH] 未登录，请先执行 soke-cli auth login
   ```

## 与现有系统的集成

### OAuth 用户登录系统
- 位置：`cmd/user_auth/`
- 功能：扫码登录、token 管理、keychain 存储
- 用途：**用户身份认证**

### AppKey/AppSecret 系统
- 位置：`internal/auth/token.go`
- 功能：企业 access_token 获取
- 用途：**API 调用认证**

### 认证拦截的作用
- 确保用户已完成 OAuth 登录
- 不干预 AppKey/AppSecret 的 token 获取
- 两套系统独立运行，互不影响

## 注意事项

1. **命令名称识别**
   - 使用父命令名称判断（如 `contact`）
   - 而不是子命令名称（如 `+list-departments`）

2. **错误处理**
   - 使用 `errors.PrintHuman()` 格式化输出
   - 返回错误以中断命令执行

3. **缓存管理**
   - 登录后必须调用 `UpdateCachedToken()`
   - 登出后必须调用 `ResetTokenCache()`

4. **跨命令一致性**
   - 所有业务命令自动继承认证检查
   - 无需在每个命令中单独实现

## 未来扩展

### 可能的改进方向

1. **自动刷新 Token**
   - 检测到过期时自动刷新
   - 无需用户手动重新登录

2. **多租户支持**
   - 支持切换不同的企业账号
   - 管理多个 token

3. **离线模式**
   - 缓存部分数据
   - 支持离线查询

4. **Token 过期提醒**
   - 提前通知用户 token 即将过期
   - 避免命令执行时才发现

## 总结

本实现完全遵循了参考项目（dingtalk-workspace-cli）的设计模式：

- ✅ 唯一的认证拦截点
- ✅ Fail-fast 机制
- ✅ 进程级 token 缓存
- ✅ 友好的错误提示
- ✅ 清晰的操作建议
- ✅ 灵活的跳过机制

认证拦截机制已成功集成到 soke-cli 项目中，为所有需要认证的命令提供了统一、可靠的保护。
