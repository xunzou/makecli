/**
 * [INPUT]: 依赖 internal/config（Load/LoadConfig）、internal/api（New/WithDebug/WithHeaders）、fmt；从 root.go 读取全局 Profile / ServerURL / DebugMode
 * [OUTPUT]: 对外提供 newClientFromProfile 函数
 * [POS]: cmd 模块的公共 helper，统一「全局命令行入参 → API 客户端」的构建逻辑——profile / server / debug 三态全部由 root PersistentFlag 注入，子命令零参数调用
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/qfeius/makecli/internal/api"
	"github.com/qfeius/makecli/internal/config"
)

// newClientFromProfile 按当前全局 Profile / ServerURL / DebugMode 构建 API 客户端。
// 三个全局态都来自 rootCmd 的 PersistentFlag，子命令无需也不应再传参。
func newClientFromProfile() (*api.Client, error) {
	creds, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("加载凭证失败: %w", err)
	}

	p, ok := creds[Profile]
	if !ok || p.AccessToken == "" {
		return nil, fmt.Errorf("profile '%s' 未配置，请先运行: makecli configure --profile %s", Profile, Profile)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 解析 server URL: --server-url > config > default
	server := defaultMetaServer
	headers := map[string]string{}
	if cp, ok := cfg[Profile]; ok {
		if cp.ServerURL != "" {
			server = cp.ServerURL
		}
		if cp.XTenantID != "" {
			headers["X-Tenant-ID"] = cp.XTenantID
		}
		if cp.OperatorID != "" {
			headers["X-Operator-ID"] = cp.OperatorID
		}
	}
	if ServerURL != "" {
		server = ServerURL
	}

	return api.New(server, p.AccessToken, api.WithDebug(DebugMode), api.WithHeaders(headers)), nil
}
