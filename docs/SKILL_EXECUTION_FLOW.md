# DingTalk Workspace CLI - Skill 执行流程与登录授权拦截机制

## 概述

本文档详细说明了 `dingtalk-workspace-cli` 项目中 skill 的执行流程，以及用户未登录时的拦截和提示机制。

## 整体架构

```
用户命令输入
    ↓
main.go (cmd/main.go)
    ↓
app.Execute() (internal/app/root.go)
    ↓
NewRootCommandWithEngine() - 构建命令树
    ↓
newMCPCommand() - 构建 MCP 动态命令
    ↓
BuildDynamicCommands() - 从 discovery catalog 生成 cobra 命令
    ↓
NewDirectCommand() - 为每个 tool 创建可执行命令
    ↓
RunE 执行 (internal/compat/registry.go:214)
    ↓
runner.Run() (internal/app/runner.go)
    ↓
executeInvocation() - 执行具体调用
    ↓
【认证检查点】resolveAuthToken() + 空 token 拦截
    ↓
transport.CallTool() - 发起 HTTP 请求
```

## 关键流程详解

### 1. 命令树构建阶段

**文件**: `internal/app/root.go`

```go
// Execute() 是整个 CLI 的入口点
func Execute() (exitCode int) {
    // 1. 创建 pipeline engine（处理命令的各个阶段）
    engine := newPipelineEngine()
    
    // 2. 构建根命令树
    root := NewRootCommandWithEngine(ctx, engine)
    
    // 3. 执行命令
    executed, err := root.ExecuteC()
}

// NewRootCommandWithEngine 构建完整的命令树
func NewRootCommandWithEngine(ctx context.Context, engine *pipeline.Engine) *cobra.Command {
    // 创建 runner（负责实际执行 MCP 调用）
    runner := newCommandRunnerWithFlags(loader, flags)
    
    // 创建 MCP 动态命令（从 discovery catalog 生成）
    mcpCmd := newMCPCommand(rootCtx, loader, runner, engine)
    
    // 添加到根命令
    root.AddCommand(mcpCmd)
}
```

### 2. 动态命令生成阶段

**文件**: `internal/compat/dynamic_commands.go`

```go
// BuildDynamicCommands 从 discovery catalog 生成 cobra 命令
func BuildDynamicCommands(servers []market.ServerDescriptor, runner executor.Runner, detailsByID map[string][]market.DetailTool) []*cobra.Command {
    // 遍历每个 server（产品）
    for _, server := range servers {
        // 为每个 tool 创建命令
        for _, toolName := range toolNames {
            // 创建可执行的命令
            cmd := NewDirectCommand(route, runner)
        }
    }
}
```

### 3. 命令执行阶段

**文件**: `internal/compat/registry.go`

```go
// NewDirectCommand 创建一个可执行的 cobra 命令
func NewDirectCommand(route Route, runner executor.Runner) *cobra.Command {
    cmd := &cobra.Command{
        RunE: func(cmd *cobra.Command, args []string) error {
            // 1. 收集参数
            params := collectParameters(cmd, args, route.Bindings)
            
            // 2. 创建调用对象
            invocation := executor.NewCompatibilityInvocation(
                route.Target.CanonicalProduct,
                route.Target.Tool,
                params,
            )
            
            // 3. 执行调用（这里会触发认证检查）
            result, err := runner.Run(cmd.Context(), invocation)
            
            // 4. 输出结果
            return output.WriteCommandPayload(cmd, result, output.FormatJSON)
        },
    }
    return cmd
}
```

## 🔐 登录授权拦截机制

### 拦截点位置

**文件**: `internal/app/runner.go:339-348`

这是整个系统中**唯一的认证拦截点**，在发起网络请求之前进行检查。

```go
func (r *runtimeRunner) executeInvocation(ctx context.Context, endpoint string, invocation executor.Invocation) (result executor.Result, retErr error) {
    // ... 前置处理 ...
    
    // 解析认证 token
    authToken := r.resolveAuthToken(ctx)
    
    // 【关键拦截点】Fail-fast: 在发起网络请求前检查认证状态
    // 如果 token 为空，直接返回错误，不会发起 HTTP 请求
    if strings.TrimSpace(authToken) == "" {
        return executor.Result{}, apperrors.NewAuth(
            "未登录，请先执行 dws auth login",
            apperrors.WithReason("not_authenticated"),
            apperrors.WithHint("运行 'dws auth login' 完成登录后重试"),
            apperrors.WithActions("dws auth login"),
        )
    }
    
    // 只有 token 不为空时，才会继续执行网络请求
    callResult, err := tc.CallTool(callCtx, endpoint, invocation.Tool, invocation.Params)
    // ...
}
```

### Token 解析流程

**文件**: `internal/app/runner.go:526-575`

```go
// resolveAuthToken 解析认证 token 的优先级顺序
func (r *runtimeRunner) resolveAuthToken(ctx context.Context) string {
    // 1. 优先使用命令行显式传入的 --token
    if token := strings.TrimSpace(explicitToken); token != "" {
        return token
    }
    
    // 2. 使用 edition 提供的 TokenProvider（如果有）
    if tp := edition.Get().TokenProvider; tp != nil {
        token, _ := tp(ctx, func() (string, error) {
            return resolveAccessTokenFromDir(ctx, defaultConfigDir())
        })
        return token
    }
    
    // 3. 使用进程缓存的 token（避免重复访问 Keychain）
    return getCachedRuntimeToken(ctx)
}

// getCachedRuntimeToken 从缓存或存储中加载 token
func getCachedRuntimeToken(ctx context.Context) string {
    cachedRuntimeTokenOnce.Do(func() {
        // 从配置目录加载 token
        token, tokenErr := resolveAccessTokenFromDir(ctx, configDir)
        if token != "" {
            cachedRuntimeToken = token
        }
    })
    return cachedRuntimeToken
}
```

**文件**: `internal/app/access_token_resolve.go:32-49`

```go
// resolveAccessTokenFromDir 从配置目录加载 token
func resolveAccessTokenFromDir(ctx context.Context, configDir string) (string, error) {
    // 1. 尝试加载 OAuth token（新版认证方式）
    provider := authpkg.NewOAuthProvider(configDir, disc)
    token, tokenErr := provider.GetAccessToken(ctx)
    if tokenErr == nil && strings.TrimSpace(token) != "" {
        return strings.TrimSpace(token), nil
    }
    
    // 2. 如果 OAuth token 不存在，尝试加载 legacy token（旧版认证）
    manager := authpkg.NewManager(configDir, nil)
    if leg, _, err := manager.GetToken(); err == nil && strings.TrimSpace(leg) != "" {
        return strings.TrimSpace(leg), nil
    }
    
    // 3. 都没有，返回空字符串
    return "", nil
}
```

**文件**: `internal/auth/oauth_provider.go:532-555`

```go
// GetAccessToken 获取有效的 access token，必要时自动刷新
func (p *OAuthProvider) GetAccessToken(ctx context.Context) (string, error) {
    // 1. 加载已保存的 token 数据
    data, err := LoadTokenData(p.configDir)
    if err != nil {
        // 如果加载失败，说明未登录
        return "", errors.New(i18n.T("未登录，请运行 dws auth login"))
    }
    
    // 2. 快速路径：access_token 仍然有效
    if data.IsAccessTokenValid() {
        return data.AccessToken, nil
    }
    
    // 3. 慢速路径：token 已过期，尝试使用 refresh_token 刷新
    if data.IsRefreshTokenValid() {
        refreshed, rErr := p.lockedRefresh(ctx)
        if rErr == nil {
            return refreshed.AccessToken, nil
        }
    }
    
    // 4. 所有凭证都失效，需要重新登录
    return "", errors.New(i18n.T("所有凭证已失效，请运行 dws auth login 重新登录"))
}
```

## 认证状态检查的三层防护

### 第一层：Token 文件检查
- **位置**: `internal/auth/oauth_provider.go:533-536`
- **时机**: 尝试加载 token 文件时
- **错误**: `"未登录，请运行 dws auth login"`

### 第二层：Token 有效性检查
- **位置**: `internal/auth/oauth_provider.go:539-555`
- **时机**: Token 文件存在但已过期，且无法刷新时
- **错误**: `"所有凭证已失效，请运行 dws auth login 重新登录"`

### 第三层：空 Token 拦截（最终防线）
- **位置**: `internal/app/runner.go:341-348`
- **时机**: 在发起 HTTP 请求前
- **错误**: `"未登录，请先执行 dws auth login"`
- **特点**: 这是 fail-fast 机制，避免无意义的网络请求

## 完整执行流程图

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 用户执行命令: dws calendar list                          │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. main.go → app.Execute()                                  │
│    - 初始化 pipeline engine                                  │
│    - 构建命令树                                              │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. NewRootCommandWithEngine()                               │
│    - 创建 runner (runtimeRunner)                            │
│    - 加载 discovery catalog                                 │
│    - 调用 newMCPCommand()                                   │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. BuildDynamicCommands()                                   │
│    - 遍历 servers.json 中的产品定义                         │
│    - 为每个 tool 生成 cobra.Command                         │
│    - 调用 NewDirectCommand()                                │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. cobra 解析命令并执行 RunE                                │
│    - 收集命令行参数和 flags                                 │
│    - 创建 executor.Invocation                               │
│    - 调用 runner.Run()                                      │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. runtimeRunner.Run()                                      │
│    - 查找产品和工具定义                                     │
│    - 解析 endpoint                                          │
│    - 调用 executeInvocation()                               │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 7. executeInvocation()                                      │
│    ┌──────────────────────────────────────────────────────┐│
│    │ 🔐 认证检查流程                                       ││
│    │                                                        ││
│    │ Step 1: resolveAuthToken()                            ││
│    │   ├─ 检查 --token flag                               ││
│    │   ├─ 检查 edition.TokenProvider                      ││
│    │   └─ getCachedRuntimeToken()                         ││
│    │       └─ resolveAccessTokenFromDir()                 ││
│    │           ├─ OAuthProvider.GetAccessToken()          ││
│    │           │   ├─ LoadTokenData() ← 第一层检查        ││
│    │           │   ├─ IsAccessTokenValid()                ││
│    │           │   └─ lockedRefresh() ← 第二层检查        ││
│    │           └─ Manager.GetToken() (legacy)             ││
│    │                                                        ││
│    │ Step 2: 空 token 拦截 ← 第三层检查（最终防线）       ││
│    │   if strings.TrimSpace(authToken) == "" {            ││
│    │       return apperrors.NewAuth(                      ││
│    │           "未登录，请先执行 dws auth login"          ││
│    │       )                                               ││
│    │   }                                                   ││
│    └──────────────────────────────────────────────────────┘│
│                                                             │
│    认证通过后：                                             │
│    - 创建带认证头的 HTTP client                            │
│    - 调用 transport.CallTool()                             │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 8. transport.CallTool()                                     │
│    - 构造 JSON-RPC 请求                                     │
│    - 发起 HTTP POST 请求                                    │
│    - 解析响应                                               │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 9. 返回结果并输出                                           │
│    - output.WriteCommandPayload()                           │
│    - 格式化为 JSON 输出到 stdout                            │
└─────────────────────────────────────────────────────────────┘
```

## 关键文件索引

| 文件路径 | 主要职责 |
|---------|---------|
| `cmd/main.go` | 程序入口 |
| `internal/app/root.go` | 命令树构建、根命令定义 |
| `internal/app/runner.go` | **认证拦截点**、MCP 调用执行 |
| `internal/app/access_token_resolve.go` | Token 解析逻辑 |
| `internal/compat/dynamic_commands.go` | 从 discovery catalog 生成命令 |
| `internal/compat/registry.go` | 命令 RunE 实现 |
| `internal/auth/oauth_provider.go` | OAuth 认证、Token 管理 |
| `internal/auth/manager.go` | Legacy 认证方式 |
| `internal/transport/client.go` | HTTP 请求发送 |

## 未登录时的错误提示

当用户未登录时，会在 `runner.go:341-348` 处被拦截，返回以下错误：

```
Error: 未登录，请先执行 dws auth login

Hint: 运行 'dws auth login' 完成登录后重试
Actions: dws auth login
```

这个错误会被 `apperrors` 包格式化后输出到 stderr，并返回特定的退出码。

## 总结

1. **Skill 执行流程**：
   - 命令树在启动时从 discovery catalog 动态生成
   - 每个 tool 对应一个 cobra.Command
   - RunE 中调用 runner.Run() 执行实际的 MCP 调用

2. **登录授权拦截**：
   - **唯一拦截点**：`internal/app/runner.go:341-348`
   - **拦截时机**：在发起 HTTP 请求之前（fail-fast）
   - **拦截条件**：`authToken` 为空字符串
   - **错误提示**：清晰的错误信息 + 操作建议

3. **Token 解析优先级**：
   1. 命令行 `--token` flag
   2. Edition TokenProvider
   3. 进程缓存的 token
   4. OAuth token 文件
   5. Legacy token 文件

4. **设计优势**：
   - 集中式认证检查，易于维护
   - Fail-fast 机制，避免无效请求
   - 清晰的错误提示，用户体验好
   - 支持 token 自动刷新
