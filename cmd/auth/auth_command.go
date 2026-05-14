// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	authpkg "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/cmd/user_auth"
	authguard "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/auth"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	apperrors "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/config"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"
	"github.com/spf13/cobra"
)

type authLoginConfig struct {
	Token  string
	Force  bool
	Device bool
}

func NewUserAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "auth",
		Short:             "用户认证管理",
		Long:              "管理 CLI 的认证凭证。支持 OAuth 扫码登录。",
		Args:              cobra.NoArgs,
		TraverseChildren:  true,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	if !edition.Get().HideAuthLogin {
		cmd.AddCommand(newAuthLoginCommand())
	}
	cmd.AddCommand(
		newAuthLogoutCommand(),
	)
	return cmd
}

func newAuthLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "登录授客（自动获取token，需要扫码）",
		Long: `登录授客并获取认证凭证。

支持的登录方式:
  - OAuth 设备流 (默认): 通过授客扫码授权登录

不支持的登录方式:
  - 邮箱/密码登录
  - 手机号/验证码登录

示例:
   auth login              # 扫码登录
   auth login --force      # 强制重新登录 (忽略缓存 token)`,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := resolveAuthLoginConfig(cmd)
			if err != nil {
				return err
			}
			//configDir := defaultConfigDir()
			configDir := ""
			var tokenData *authpkg.TokenData

			switch {
			case strings.TrimSpace(cfg.Token) != "":
				tokenData = &authpkg.TokenData{
					AccessToken: cfg.Token,
					ExpiresAt:   time.Now().Add(config.ManualTokenExpiry),
				}
				if err := authpkg.SaveTokenData(configDir, tokenData); err != nil {
					return apperrors.NewInternal(fmt.Sprintf("failed to persist auth token: %v", err))
				}
			default:
				loginCtx, cancel := context.WithTimeout(cmd.Context(), config.OAuthFlowTimeout)
				defer cancel()

				provider := authpkg.NewOAuthProvider(configDir, nil)
				provider.Output = cmd.ErrOrStderr()
				//configureOAuthProviderCompatibility(provider, configDir)
				tokenData, err = provider.Login(loginCtx, cfg.Force)
				if err != nil {
					return apperrors.NewAuth(fmt.Sprintf(" login failed: %v", err))
				}

				// 更新 token 缓存
				if tokenData != nil && tokenData.AccessToken != "" {
					authguard.UpdateCachedToken(tokenData.AccessToken)
				}
			}

			clearCompatCache()

			w := cmd.OutOrStdout()

			// Check if JSON output is requested
			format, _ := cmd.Root().PersistentFlags().GetString("format")
			if strings.EqualFold(strings.TrimSpace(format), "json") {
				return writeAuthLoginJSON(w, tokenData, cfg.Force)
			}

			// Default table output - removed duplicate output
			return nil
		},
	}
	cmd.Flags().String("token", "", "Access token")
	_ = cmd.Flags().MarkHidden("token")
	cmd.Flags().Bool("device", false, "Use device authorization flow")
	_ = cmd.Flags().MarkHidden("device")
	cmd.Flags().Bool("force", false, "Force interactive login (ignore cached token)")
	_ = cmd.Flags().MarkHidden("force")
	// Hidden compatibility flags
	cmd.Flags().String("redirect-url", "", "Loopback redirect URL")
	cmd.Flags().String("scopes", "", "Space-separated DingTalk OAuth scopes")
	cmd.Flags().String("authorize-url", "", "Override DingTalk authorization URL")
	cmd.Flags().String("token-url", "", "Override DingTalk token exchange URL")
	cmd.Flags().String("refresh-url", "", "Override DingTalk refresh token URL")
	cmd.Flags().Int("login-timeout", 0, "Login timeout seconds")
	cmd.Flags().Bool("no-browser", false, "Suppress browser launch")
	_ = cmd.Flags().MarkHidden("redirect-url")
	_ = cmd.Flags().MarkHidden("scopes")
	_ = cmd.Flags().MarkHidden("authorize-url")
	_ = cmd.Flags().MarkHidden("token-url")
	_ = cmd.Flags().MarkHidden("refresh-url")
	_ = cmd.Flags().MarkHidden("login-timeout")
	_ = cmd.Flags().MarkHidden("no-browser")
	return cmd
}

func newAuthLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "logout",
		Short:             "清除认证信息",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := ""
			_, cancel := context.WithTimeout(cmd.Context(), 15*time.Second)
			defer cancel()

			var storedClientID string
			if tokenData, err := authpkg.LoadTokenData(configDir); err == nil && tokenData != nil {
				storedClientID = tokenData.ClientID
			}

			if err := authpkg.DeleteTokenData(configDir); err != nil {
				return apperrors.NewInternal(fmt.Sprintf("failed to clear token data: %v", err))
			}

			if err := clearCliConfig(); err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "警告: 清除配置文件失败: %v\n", err)
			}

			// 重置 token 缓存
			authguard.ResetTokenCache()

			// Clean up associated client secret and app token from keychain
			if storedClientID != "" {
				_ = authpkg.DeleteClientSecret(storedClientID)
				//_ = authpkg.DeleteAppTokenData(storedClientID)
			}
			
			_ = os.Remove(filepath.Join(configDir, "mcp_url"))
			_ = os.Remove(filepath.Join(configDir, "token"))
			_ = os.Remove(filepath.Join(configDir, "token.json"))
			//ResetRuntimeTokenCache()
			clearCompatCache()
			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "[OK] 已清除所有认证信息")
			if !edition.Get().IsEmbedded {
				fmt.Fprintln(w, "请运行  auth login 重新登录")
			}
			return nil
		},
	}
}



func clearCompatCache() {
	//store := cacheStoreFromEnv()
	//if store != nil {
	//	_ = os.RemoveAll(store.Root)
	//}
}

func resolveAuthLoginConfig(cmd *cobra.Command) (authLoginConfig, error) {
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return authLoginConfig{}, apperrors.NewInternal("failed to read --token")
	}
	device, err := cmd.Flags().GetBool("device")
	if err != nil {
		return authLoginConfig{}, apperrors.NewInternal("failed to read --device")
	}
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return authLoginConfig{}, apperrors.NewInternal("failed to read --force")
	}
	return authLoginConfig{
		Token:  strings.TrimSpace(token),
		Force:  force,
		Device: device,
	}, nil
}

// authStatusResponse is the JSON response for auth status command.
type authStatusResponse struct {
	Success           bool   `json:"success"`
	Authenticated     bool   `json:"authenticated"`
	Message           string `json:"message,omitempty"`
	Refreshed         bool   `json:"refreshed,omitempty"`
	TokenValid        bool   `json:"token_valid,omitempty"`
	RefreshTokenValid bool   `json:"refresh_token_valid,omitempty"`
	ExpiresAt         string `json:"expires_at,omitempty"`
	RefreshExpiresAt  string `json:"refresh_expires_at,omitempty"`
	CorpID            string `json:"corp_id,omitempty"`
	CorpName          string `json:"corp_name,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	UserName          string `json:"user_name,omitempty"`
}

func writeAuthStatusJSON(w io.Writer, authenticated, refreshed bool, data *authpkg.TokenData) error {
	resp := authStatusResponse{
		Success:       true,
		Authenticated: authenticated,
	}

	if !authenticated {
		resp.Message = "未登录"
	} else if data != nil {
		resp.Refreshed = refreshed
		resp.TokenValid = data.IsAccessTokenValid()
		resp.RefreshTokenValid = data.IsRefreshTokenValid()
		if !data.ExpiresAt.IsZero() {
			resp.ExpiresAt = data.ExpiresAt.Format(time.RFC3339Nano)
		}
		if !data.RefreshExpAt.IsZero() {
			resp.RefreshExpiresAt = data.RefreshExpAt.Format(time.RFC3339Nano)
		}
		resp.CorpID = data.CorpID
		resp.CorpName = data.CorpName
		resp.UserID = data.UserID
		resp.UserName = data.UserName
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
}

// authLoginResponse is the JSON response for auth login command.
type authLoginResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	TokenValid        bool   `json:"token_valid,omitempty"`
	RefreshTokenValid bool   `json:"refresh_token_valid,omitempty"`
	ExpiresAt         string `json:"expires_at,omitempty"`
	RefreshExpiresAt  string `json:"refresh_expires_at,omitempty"`
	CorpID            string `json:"corp_id,omitempty"`
	CorpName          string `json:"corp_name,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	UserName          string `json:"user_name,omitempty"`
}

func writeAuthLoginJSON(w io.Writer, data *authpkg.TokenData, forced bool) error {
	resp := authLoginResponse{
		Success: true,
		Message: "登录成功",
	}

	if data != nil {
		if data.IsAccessTokenValid() && !forced {
			resp.Message = "Token 有效，无需重新登录"
		}
		resp.TokenValid = data.IsAccessTokenValid()
		resp.RefreshTokenValid = data.IsRefreshTokenValid()
		if !data.ExpiresAt.IsZero() {
			resp.ExpiresAt = data.ExpiresAt.Format(time.RFC3339Nano)
		}
		if !data.RefreshExpAt.IsZero() {
			resp.RefreshExpiresAt = data.RefreshExpAt.Format(time.RFC3339Nano)
		}
		resp.CorpID = data.CorpID
		resp.CorpName = data.CorpName
		resp.UserID = data.UserID
		resp.UserName = data.UserName
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
}

// clearCliConfig 清除 CliConfig 中的认证相关信息
func clearCliConfig() error {
	cfg, err := core.LoadConfig()
	if err != nil {
		// 如果配置文件不存在，不算错误
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 清除认证相关字段
	cfg.AppID = ""
	cfg.AppSecret = ""
	cfg.UserToken = ""
	cfg.UserTokenExp = 0
	cfg.BotToken = ""
	cfg.CorpID = ""
	cfg.DeptUserID = ""

	// 保存清空后的配置
	if err := core.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	return nil
}
