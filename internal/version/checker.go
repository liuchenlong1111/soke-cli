package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 当前版本（编译时注入）
var Version = "1.0.26"

// NPM 注册表响应结构
type npmRegistryResponse struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

// 版本检查缓存文件路径
func getCacheFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".soke-cli", "version_check_cache.json")
}

// 版本检查缓存结构
type versionCache struct {
	LastCheckTime int64  `json:"last_check_time"`
	LatestVersion string `json:"latest_version"`
}

// 读取缓存
func readCache() (*versionCache, error) {
	data, err := os.ReadFile(getCacheFilePath())
	if err != nil {
		return nil, err
	}

	var cache versionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// 写入缓存
func writeCache(cache *versionCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheDir := filepath.Dir(getCacheFilePath())
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(getCacheFilePath(), data, 0644)
}

// 从 NPM 获取最新版本
func fetchLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://registry.npmjs.org/@sokeai/cli")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NPM 注册表返回状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var npmResp npmRegistryResponse
	if err := json.Unmarshal(body, &npmResp); err != nil {
		return "", err
	}

	return npmResp.DistTags.Latest, nil
}

// 解析版本号为数字数组
func parseVersion(version string) ([]int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("无效的版本号格式: %s", version)
	}

	result := make([]int, 3)
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("无效的版本号部分: %s", part)
		}
		result[i] = num
	}

	return result, nil
}

// 比较版本号（语义化版本比较）
func isNewerVersion(latest, current string) bool {
	if latest == current {
		return false
	}

	latestParts, err := parseVersion(latest)
	if err != nil {
		return false
	}

	currentParts, err := parseVersion(current)
	if err != nil {
		return false
	}

	// 比较主版本号、次版本号、修订号
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

// CheckForUpdates 检查更新（带缓存机制，每24小时检查一次）
func CheckForUpdates() {
	// 读取缓存
	cache, err := readCache()
	now := time.Now().Unix()

	// 如果缓存存在且未过期（24小时内），使用缓存数据
	if err == nil && now-cache.LastCheckTime < 86400 {
		if isNewerVersion(cache.LatestVersion, Version) {
			printUpdateNotice(cache.LatestVersion)
		}
		return
	}

	// 从 NPM 获取最新版本
	latestVersion, err := fetchLatestVersion()
	if err != nil {
		// 静默失败，不影响用户使用
		return
	}

	// 更新缓存
	newCache := &versionCache{
		LastCheckTime: now,
		LatestVersion: latestVersion,
	}
	_ = writeCache(newCache)

	// 如果有新版本，显示提示
	if isNewerVersion(latestVersion, Version) {
		printUpdateNotice(latestVersion)
	}
}

// CheckForUpdatesAsync 异步检查更新（用于 PersistentPreRun）
func CheckForUpdatesAsync() {
	go CheckForUpdates()
}

// 打印更新提示
func printUpdateNotice(latestVersion string) {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "╭────────────────────────────────────────────────────────╮\n")
	fmt.Fprintf(os.Stderr, "│  🎉 发现新版本可用！                                   │\n")
	fmt.Fprintf(os.Stderr, "│                                                        │\n")
	fmt.Fprintf(os.Stderr, "│  当前版本: %s                                      │\n", Version)
	fmt.Fprintf(os.Stderr, "│  最新版本: %s                                      │\n", latestVersion)
	fmt.Fprintf(os.Stderr, "│                                                        │\n")
	fmt.Fprintf(os.Stderr, "│  更新命令:                                             │\n")
	fmt.Fprintf(os.Stderr, "│  npm update -g @sokeai/cli                             │\n")
	fmt.Fprintf(os.Stderr, "│                                                        │\n")
	fmt.Fprintf(os.Stderr, "╰────────────────────────────────────────────────────────╯\n")
	fmt.Fprintf(os.Stderr, "\n")
}

// GetVersion 获取当前版本
func GetVersion() string {
	return Version
}
