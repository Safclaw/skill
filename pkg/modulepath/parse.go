package modulepath

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// ParsedPath 解析后的模块路径
type ParsedPath struct {
	Host      string // github.com, gitlab.com, etc.
	Namespace string // Safclaw/skills
	Name      string // read-json
	Version   string // v1.0.0, latest, main
	SubDir    string // subdirectory within the repo (e.g., "empty" in Safclaw/skill/empty)
	Raw       string // 原始输入
}

// Parse 解析技能路径，格式：{host}/{namespace}/{name}@{version}
func Parse(path string) (*ParsedPath, error) {
	if path == "" {
		return nil, fmt.Errorf("empty module path")
	}

	// 分离版本部分
	var version string
	namePart := path
	if idx := strings.LastIndex(path, "@"); idx != -1 {
		namePart = path[:idx]
		version = path[idx+1:]
	}

	// 如果没有指定版本，使用 latest
	if version == "" {
		version = "latest"
	}

	// 解析 host/namespace/name
	parts := strings.Split(namePart, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid module path format: %s, expected: {host}/{namespace}/{name}", path)
	}

	host := parts[0]

	// Check for subdirectory (e.g., github.com/Safclaw/skill/empty -> namespace=Safclaw/skill, name=empty)
	// But if there are more than 3 parts after host, treat the extra as subdirectory
	var subDir string
	var namespace, name string

	if len(parts) > 4 {
		// More than 4 parts: github.com/org/repo/subdir1/subdir2
		// namespace = org, name = repo, subDir = subdir1/subdir2
		namespace = parts[1]
		name = parts[2]
		subDir = strings.Join(parts[3:], "/")
	} else if len(parts) == 4 {
		// Exactly 4 parts: github.com/org/repo/subdir
		// namespace = org, name = repo, subDir = subdir
		namespace = parts[1]
		name = parts[2]
		subDir = parts[3]
	} else {
		// Exactly 3 parts: github.com/org/repo
		// namespace = org, name = repo
		namespace = strings.Join(parts[1:len(parts)-1], "/")
		name = parts[len(parts)-1]
	}

	// 验证 host 格式
	if !isValidHost(host) {
		return nil, fmt.Errorf("invalid host: %s", host)
	}

	// 验证 namespace
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	// 验证 name
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// 验证 version
	if !isValidVersion(version) {
		return nil, fmt.Errorf("invalid version: %s", version)
	}

	return &ParsedPath{
		Host:      host,
		Namespace: namespace,
		Name:      name,
		Version:   version,
		SubDir:    subDir,
		Raw:       path,
	}, nil
}

// String 返回标准格式的字符串
func (p *ParsedPath) String() string {
	return fmt.Sprintf("%s/%s/%s@%s", p.Host, p.Namespace, p.Name, p.Version)
}

// DirPath 返回目录路径（不带版本）
func (p *ParsedPath) DirPath() string {
	if p.SubDir != "" {
		return filepath.Join(p.Host, p.Namespace, p.Name, p.SubDir)
	}
	return filepath.Join(p.Host, p.Namespace, p.Name)
}

// FullPath 返回完整路径（带版本）
func (p *ParsedPath) FullPath() string {
	return filepath.Join(p.DirPath(), p.Version)
}

// EscapePath 转义路径用于文件系统（Windows 兼容性）
func (p *ParsedPath) EscapePath() string {
	// Windows 不区分大小写，使用 ! 前缀避免冲突
	escapedHost := "!" + strings.ToLower(p.Host)
	return filepath.Join(escapedHost, strings.ToLower(p.Namespace), strings.ToLower(p.Name))
}

// CachePath 返回缓存目录路径
func (p *ParsedPath) CachePath(cacheRoot string) string {
	return filepath.Join(cacheRoot, p.EscapePath()+"@"+p.Version)
}

// InstallPath 返回安装目录路径
func (p *ParsedPath) InstallPath(installRoot string) string {
	return filepath.Join(installRoot, "reps", p.DirPath())
}

// DownloadURL 生成下载 URL
func (p *ParsedPath) DownloadURL() (string, error) {
	switch p.Host {
	case "github.com":
		return fmt.Sprintf("https://codeload.github.com/%s/%s/zip/%s",
			p.Namespace, p.Name, p.Version), nil
	case "gitlab.com":
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/archive/%s/%s-%s.zip",
			p.Namespace, p.Name, p.Version, p.Namespace, p.Version), nil
	case "gitee.com":
		return fmt.Sprintf("https://gitee.com/%s/%s/repository/archive/%s.zip",
			p.Namespace, p.Name, p.Version), nil
	default:
		// 尝试通用 Git 下载
		return "", fmt.Errorf("unsupported host: %s", p.Host)
	}
}

// Signature 计算签名
func (p *ParsedPath) Signature(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// isValidHost 验证 host 是否合法
func isValidHost(host string) bool {
	if host == "" {
		return false
	}

	// 简单的域名验证
	pattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	matched, _ := regexp.MatchString(pattern, host)
	return matched
}

// isValidVersion 验证版本号是否合法
func isValidVersion(version string) bool {
	if version == "" {
		return false
	}

	// 支持特殊标签
	if version == "latest" || version == "stable" {
		return true
	}

	// 支持语义化版本
	semverPattern := `^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?(\+[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$`
	matched, _ := regexp.MatchString(semverPattern, version)
	if matched {
		return true
	}

	// 支持 Git commit hash（短格式和长格式）
	hashPattern := `^[a-fA-F0-9]{7,40}$`
	matched, _ = regexp.MatchString(hashPattern, version)
	if matched {
		return true
	}

	// 支持分支名
	branchPattern := `^[a-zA-Z][a-zA-Z0-9\-_./]+$`
	matched, _ = regexp.MatchString(branchPattern, version)
	return matched
}

// IsAbsolute 判断是否是绝对路径
func IsAbsolute(path string) bool {
	return strings.HasPrefix(path, "/") || (len(path) > 1 && path[1] == ':')
}

// SanitizePath 清理路径，防止路径注入
func SanitizePath(path string) (string, error) {
	// 解码 URL
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return "", err
	}

	// 清理路径
	cleaned := filepath.Clean(decoded)

	// 检查是否包含 ..
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains '..': %s", cleaned)
	}

	// 检查是否是绝对路径
	if IsAbsolute(cleaned) {
		return "", fmt.Errorf("absolute path not allowed: %s", cleaned)
	}

	return cleaned, nil
}
