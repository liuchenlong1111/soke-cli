package auth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	userauth "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/user_auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/errors"
)

// 需要跳过认证检查的命令列表
var skipAuthCommands = map[string]bool{
	"auth":    true,
	"config":  true,
	"version": true,
	"help":    true,
}

// 配置信息缓存
var (
	cachedConfig     *core.CliConfig
	configCachedAt   time.Time
	configCacheTTL   = 2 * time.Hour
	configMutex      sync.RWMutex
)

// CheckAuth 检查用户是否已登录授权
// 这是整个系统中唯一的认证拦截点，在命令执行前进行检查
// 参考 dingtalk-workspace-cli 的实现：internal/app/runner.go:339-348
func CheckAuth(ctx context.Context, commandName string) error {
	// 跳过不需要认证的命令
	if skipAuthCommands[commandName] {
		return nil
	}

	// 检查配置信息是否完整
	cfg, err := LoadAndValidateConfig(ctx)
	if err != nil {
		return err
	}

	// 检查 UserToken 是否存在且有效（验证用户是否已登录授权）
	if strings.TrimSpace(cfg.UserToken) == "" {
		return errors.NewAuth(
			"未登录授权，请先执行 soke-cli auth login",
			errors.WithReason("not_authenticated"),
			errors.WithHint("运行 'soke-cli auth login' 完成登录授权后重试"),
			errors.WithActions("soke-cli auth login"),
		)
	}

	// 检查 UserToken 是否过期
	if cfg.UserTokenExp > 0 && time.Now().Unix() > cfg.UserTokenExp {
		return errors.NewAuth(
			"登录已过期，请重新登录",
			errors.WithReason("token_expired"),
			errors.WithHint("运行 'soke-cli auth login' 重新登录"),
			errors.WithActions("soke-cli auth login"),
		)
	}

	// 检查应用配置信息是否完整（app_key、app_secret、corpid）
	if !IsConfigComplete(cfg) {
		return errors.NewAuth(
			"用户数据信息不完整，请先执行 soke-cli auth login 完成授权",
			errors.WithReason("incomplete_config"),
			errors.WithHint("运行 'soke-cli auth login' 完成授权后重试"),
			errors.WithActions("soke-cli auth login"),
		)
	}

	return nil
}

// LoadAndValidateConfig 加载并验证配置信息
// 优先从缓存读取，缓存有效期2小时
func LoadAndValidateConfig(ctx context.Context) (*core.CliConfig, error) {
	// 快速路径：检查缓存是否有效
	configMutex.RLock()
	if cachedConfig != nil && time.Since(configCachedAt) < configCacheTTL {
		cfg := cachedConfig
		configMutex.RUnlock()
		return cfg, nil
	}
	configMutex.RUnlock()

	// 慢速路径：从文件加载配置
	configMutex.Lock()
	defer configMutex.Unlock()

	// 双重检查：防止并发重复加载
	if cachedConfig != nil && time.Since(configCachedAt) < configCacheTTL {
		return cachedConfig, nil
	}

	// 从文件加载配置
	cfg, err := core.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NewAuth(
				"未登录授权，请先执行 soke-cli auth login",
				errors.WithReason("not_authenticated"),
				errors.WithHint("运行 'soke-cli auth login' 完成登录授权后重试"),
				errors.WithActions("soke-cli auth login"),
			)
		}
		return nil, errors.NewInternal("加载配置失败: " + err.Error())
	}

	// 缓存配置信息
	cachedConfig = cfg
	configCachedAt = time.Now()

	return cfg, nil
}

// IsConfigComplete 检查配置信息是否完整
// 需要包含：AppID、AppSecret、CorpID、APIBaseURL
func IsConfigComplete(cfg *core.CliConfig) bool {
	if cfg == nil {
		return false
	}
	return strings.TrimSpace(cfg.AppID) != "" &&
		strings.TrimSpace(cfg.AppKey) != "" &&
		strings.TrimSpace(cfg.AppSecret) != "" &&
		strings.TrimSpace(cfg.CorpID) != "" &&
		strings.TrimSpace(cfg.APIBaseURL) != ""
}

// ResetConfigCache 重置配置缓存
// 在登录/登出操作后调用，强制下次访问时重新加载
func ResetConfigCache() {
	configMutex.Lock()
	defer configMutex.Unlock()
	cachedConfig = nil
	configCachedAt = time.Time{}
}

// GetCachedConfig 获取缓存的配置
// 如果缓存不存在或已过期，返回 nil
func GetCachedConfig() *core.CliConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if cachedConfig != nil && time.Since(configCachedAt) < configCacheTTL {
		return cachedConfig
	}
	return nil
}

// UpdateCachedConfig 更新缓存的配置
// 在登录成功后调用
func UpdateCachedConfig(cfg *core.CliConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	cachedConfig = cfg
	configCachedAt = time.Now()
}

// ResolveAuthToken 解析认证 token（保留用于兼容性）
// 优先级顺序：
// 1. 环境变量 SOKE_TOKEN
// 2. 进程缓存的 token
// 3. OAuth token 文件
func ResolveAuthToken(ctx context.Context) string {
	// 1. 优先使用环境变量（支持 CI/CD 场景）
	envToken := os.Getenv("SOKE_TOKEN")
	if token := strings.TrimSpace(envToken); token != "" {
		return token
	}

	// 2. 从缓存的配置中获取 UserToken
	if cfg := GetCachedConfig(); cfg != nil && cfg.UserToken != "" {
		return cfg.UserToken
	}

	// 3. 从配置文件加载
	cfg, err := core.LoadConfig()
	if err == nil && cfg.UserToken != "" {
		return cfg.UserToken
	}

	// 4. 尝试从 OAuth token 文件加载（兼容旧版本）
	return loadTokenFromOAuth(ctx)
}

// loadTokenFromOAuth 从 OAuth token 文件加载 token（兼容旧版本）
func loadTokenFromOAuth(ctx context.Context) string {
	configDir := getDefaultConfigDir()
	tokenData, err := userauth.LoadTokenData(configDir)
	if err != nil {
		return ""
	}

	if tokenData == nil || !tokenData.IsAccessTokenValid() {
		return ""
	}

	return tokenData.AccessToken
}

// getDefaultConfigDir 获取默认配置目录
func getDefaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".soke-cli")
}
