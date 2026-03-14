package config

import (
	"os"
	"strings"
)

// Config 配置信息
type Config struct {
	Proxy        string // SKILLPROXY - 代理服务器列表，逗号分隔
	Private      string // SKILLPRIVATE - 私有仓库模式
	NoProxy      string // SKILLNOPROXY - 不使用代理的模式
	CacheDir     string // 缓存目录
	GlobalDir    string // 全局安装目录
	WorkspaceDir string // Workspace 安装目录
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		Proxy:        getEnvOrDefault("SKILLPROXY", "direct"),
		Private:      getEnvOrDefault("SKILLPRIVATE", ""),
		NoProxy:      getEnvOrDefault("SKILLNOPROXY", ""),
		CacheDir:     getEnvOrDefault("SKILL_CACHE", homeDir+"/.safeclaw/skill/cache"),
		GlobalDir:    getEnvOrDefault("SKILL_GLOBAL_DIR", homeDir+"/.safclaw/skills"),
		WorkspaceDir: getEnvOrDefault("SKILL_WORKSPACE_DIR", homeDir+"/.safclaw/workspace"),
	}
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return DefaultConfig()
}

// ProxyList 返回代理服务器列表
func (c *Config) ProxyList() []string {
	if c.Proxy == "" {
		return []string{"direct"}
	}

	parts := strings.Split(c.Proxy, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		return []string{"direct"}
	}

	return result
}

// UseProxy 判断是否对指定 host 使用代理
func (c *Config) UseProxy(host string) bool {
	// 检查 NoProxy
	if c.NoProxy != "" {
		noProxyParts := strings.Split(c.NoProxy, ",")
		for _, np := range noProxyParts {
			np = strings.TrimSpace(np)
			if np == "" {
				continue
			}

			// 支持通配符 *.example.com
			if strings.HasPrefix(np, "*.") {
				suffix := np[1:]
				if strings.HasSuffix(host, suffix) {
					return false
				}
			} else if host == np {
				return false
			}
		}
	}

	// 检查是否在 Private 列表中
	if c.Private != "" {
		privateParts := strings.Split(c.Private, ",")
		for _, pp := range privateParts {
			pp = strings.TrimSpace(pp)
			if pp == "" {
				continue
			}

			// 支持通配符 github.com/myorg/*
			if strings.HasSuffix(pp, "/*") {
				prefix := pp[:len(pp)-2]
				if strings.HasPrefix(host, prefix) {
					return false
				}
			}
		}
	}

	return true
}

// IsDirect 是否直接下载（不使用代理）
func (c *Config) IsDirect() bool {
	proxies := c.ProxyList()
	for _, p := range proxies {
		if p == "direct" {
			return true
		}
	}
	return false
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
