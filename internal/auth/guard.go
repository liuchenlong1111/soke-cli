package auth

import (
	"context"
	"strings"
	"sync"

	userauth "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/user_auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/errors"
)

// 需要跳过认证检查的命令列表
var skipAuthCommands = map[string]bool{
	"auth":    true,
	"config":  true,
	"version": true,
	"help":    true,
}

// Token 缓存
var (
	cachedToken     string
	cachedTokenOnce sync.Once
	tokenMutex      sync.RWMutex
)

// CheckAuth 检查用户是否已登录授权
// 这是整个系统中唯一的认证拦截点，在命令执行前进行检查
func CheckAuth(ctx context.Context, commandName string) error {
	// 跳过不需要认证的命令
	// 检查命令名称或其父命令是否在跳过列表中
	if skipAuthCommands[commandName] {
		return nil
	}

	// 解析 token
	token := ResolveAuthToken(ctx)

	// Fail-fast: 如果 token 为空，直接返回错误，不会发起任何请求
	if strings.TrimSpace(token) == "" {
		return errors.NewAuth(
			"未登录，请先执行 soke-cli auth login",
			errors.WithReason("not_authenticated"),
			errors.WithHint("运行 'soke-cli auth login' 完成登录后重试"),
			errors.WithActions("soke-cli auth login"),
		)
	}

	// 检查 token 是否过期
	if IsTokenExpired(ctx) {
		return errors.NewAuth(
			"登录已过期，请重新登录",
			errors.WithReason("token_expired"),
			errors.WithHint("运行 'soke-cli auth login' 重新登录"),
			errors.WithActions("soke-cli auth login"),
		)
	}

	return nil
}

// ResolveAuthToken 解析认证 token
// 优先级顺序：
// 1. 环境变量 SOKE_TOKEN
// 2. 进程缓存的 token
// 3. OAuth token 文件
func ResolveAuthToken(ctx context.Context) string {
	// 1. 优先使用环境变量
	// envToken := os.Getenv("SOKE_TOKEN")
	// if token := strings.TrimSpace(envToken); token != "" {
	// 	return token
	// }

	// 2. 使用进程缓存的 token（避免重复读取 keychain）
	return GetCachedToken(ctx)
}

// GetCachedToken 获取缓存的 token，只在进程生命周期内加载一次
func GetCachedToken(ctx context.Context) string {
	tokenMutex.RLock()
	if cachedToken != "" {
		token := cachedToken
		tokenMutex.RUnlock()
		return token
	}
	tokenMutex.RUnlock()

	// 双重检查锁
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if cachedToken != "" {
		return cachedToken
	}

	// 从 OAuth token 文件加载 token
	token := loadTokenFromOAuth(ctx)
	if token != "" {
		cachedToken = token
	}

	return cachedToken
}

// loadTokenFromOAuth 从 OAuth token 文件加载 token
func loadTokenFromOAuth(ctx context.Context) string {
	configDir := "" // 使用默认配置目录
	tokenData, err := userauth.LoadTokenData(configDir)
	if err != nil {
		return ""
	}

	// 检查 token 是否有效
	if tokenData == nil || !tokenData.IsAccessTokenValid() {
		return ""
	}

	return tokenData.AccessToken
}

// IsTokenExpired 检查 token 是否过期
func IsTokenExpired(ctx context.Context) bool {
	configDir := "" // 使用默认配置目录
	tokenData, err := userauth.LoadTokenData(configDir)
	if err != nil {
		return true
	}

	// 检查 token 是否有效
	return tokenData == nil || !tokenData.IsAccessTokenValid()
}

// ResetTokenCache 重置 token 缓存
// 在登录/登出操作后调用，强制下次访问时重新加载
func ResetTokenCache() {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	cachedToken = ""
	cachedTokenOnce = sync.Once{}
}

// UpdateCachedToken 更新缓存的 token
// 在登录成功后调用
func UpdateCachedToken(token string) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	cachedToken = token
}
