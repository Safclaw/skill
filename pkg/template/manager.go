package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Safclaw/skill/pkg/downloader"
	"github.com/Safclaw/skill/pkg/modulepath"
)

// TemplateManager 模板管理器
type TemplateManager struct {
	cacheDir string
}

// NewTemplateManager 创建模板管理器
func NewTemplateManager(cacheDir string) *TemplateManager {
	return &TemplateManager{
		cacheDir: cacheDir,
	}
}

// CopyTemplate 从模板创建 skill
func (tm *TemplateManager) CopyTemplate(ctx context.Context, templatePath, destDir, moduleName string) error {
	// 判断是本地路径还是远程仓库
	if isLocalPath(templatePath) {
		return tm.copyLocalTemplate(templatePath, destDir, moduleName)
	}

	// 远程仓库模板
	return tm.copyRemoteTemplate(ctx, templatePath, destDir, moduleName)
}

// copyLocalTemplate 复制本地模板
func (tm *TemplateManager) copyLocalTemplate(srcPath, destDir, moduleName string) error {
	// 确保源路径存在
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("template path does not exist: %s", srcPath)
	}

	// 验证模板
	if err := validateTemplate(srcPath); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// 复制文件
	if err := copyDirectory(srcPath, destDir); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// 更新 skill.yaml 中的 name 字段
	if err := updateModuleName(destDir, moduleName); err != nil {
		return fmt.Errorf("failed to update module name: %w", err)
	}

	return nil
}

// copyRemoteTemplate 复制远程模板
func (tm *TemplateManager) copyRemoteTemplate(ctx context.Context, templatePath, destDir, moduleName string) error {
	// 解析模块路径
	parsed, err := modulepath.Parse(templatePath)
	if err != nil {
		return fmt.Errorf("invalid template path: %w", err)
	}

	// 检查是否有子目录路径（如：github.com/Safclaw/skill/empty 中的 empty）
	var subDirPath string
	pathParts := strings.Split(parsed.Raw, "/")
	if len(pathParts) > 3 {
		if len(pathParts) > 4 {
			// github.com/Safclaw/skill/empty/subdir
			subDirPath = strings.Join(pathParts[4:], "/")
		} else {
			// github.com/Safclaw/skill/empty
			subDirPath = pathParts[3]
		}
		
		// 重新解析为实际的仓库路径
		actualRepoPath := strings.Join(pathParts[:3], "/") // github.com/Safclaw/skill
		if idx := strings.LastIndex(versionPart(templatePath), "@"); idx != -1 {
			actualRepoPath += "@" + versionPart(templatePath)[idx+1:]
		}
		parsed, err = modulepath.Parse(actualRepoPath)
		if err != nil {
			return fmt.Errorf("invalid repository path: %w", err)
		}
	}

	// 使用 downloader 下载
	dl := downloader.NewGitHubDownloader()
	result, err := dl.Download(ctx, parsed.Host, parsed.Namespace, parsed.Name, parsed.Version)
	if err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}
	defer os.Remove(result.Path)

	// 创建临时目录解压
	tmpDir, err := os.MkdirTemp("", "skill-template-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 解压 zip 文件
	if err := downloader.ExtractZip(result.Path, tmpDir); err != nil {
		return fmt.Errorf("failed to extract template: %w", err)
	}

	// 规范化路径（移除顶层目录）
	normalizedDir := filepath.Join(tmpDir, "normalized")
	if err := downloader.NormalizeZipPath(tmpDir, normalizedDir); err != nil {
		return fmt.Errorf("failed to normalize path: %w", err)
	}

	// 确定模板源目录
	templateSourceDir := normalizedDir
	if subDirPath != "" {
		potentialPath := filepath.Join(normalizedDir, subDirPath)
		if info, err := os.Stat(potentialPath); err == nil && info.IsDir() {
			templateSourceDir = potentialPath
		} else {
			return fmt.Errorf("subdirectory not found: %s", subDirPath)
		}
	}

	// 验证模板
	if err := validateTemplate(templateSourceDir); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// 复制文件到目标目录
	if err := copyDirectory(templateSourceDir, destDir); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// 更新 skill.yaml 中的 name 字段
	if err := updateModuleName(destDir, moduleName); err != nil {
		return fmt.Errorf("failed to update module name: %w", err)
	}

	return nil
}

// versionPart 提取路径中的版本部分
func versionPart(path string) string {
	if idx := strings.LastIndex(path, "@"); idx != -1 {
		return path[idx+1:]
	}
	return "latest"
}

// isLocalPath 判断是否是本地路径
func isLocalPath(path string) bool {
	// 绝对路径或相对路径
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "../") {
		return true
	}

	// 检查是否包含常见的 Git host
	commonHosts := []string{"github.com", "gitlab.com", "gitee.com"}
	for _, host := range commonHosts {
		if strings.HasPrefix(path, host+"/") {
			return false
		}
	}

	// 如果路径中包含多个 / 且第一个部分看起来像域名，则认为是远程路径
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		// 检查第一部分是否包含点（可能是域名）
		if strings.Contains(parts[0], ".") {
			return false
		}
	}

	// 默认认为是本地路径
	return true
}

// validateTemplate 验证目录是否是有效的 skill 模板
func validateTemplate(dirPath string) error {
	// 必须包含 skill.yaml
	skillYamlPath := filepath.Join(dirPath, "skill.yaml")
	if _, err := os.Stat(skillYamlPath); os.IsNotExist(err) {
		return fmt.Errorf("missing skill.yaml")
	}

	// 可以包含 skill.md（可选）
	skillMdPath := filepath.Join(dirPath, "skill.md")
	if _, err := os.Stat(skillMdPath); os.IsNotExist(err) {
		// skill.md 不存在时，检查是否有其他入口文件
		// 这里暂时不强制要求
	}

	// 验证 skill.yaml 格式
	data, err := os.ReadFile(skillYamlPath)
	if err != nil {
		return fmt.Errorf("failed to read skill.yaml: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, "name:") {
		return fmt.Errorf("skill.yaml must contain 'name' field")
	}

	return nil
}

// copyDirectory 复制目录
func copyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// 跳过 .git 目录
		if entry.IsDir() && entry.Name() == ".git" {
			continue
		}

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateModuleName 更新 skill.yaml 中的模块名
func updateModuleName(dirPath, moduleName string) error {
	skillYamlPath := filepath.Join(dirPath, "skill.yaml")

	data, err := os.ReadFile(skillYamlPath)
	if err != nil {
		return err
	}

	content := string(data)

	// 替换 name 字段
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "name:") {
			indent := ""
			if idx := strings.Index(line, "name:"); idx > 0 {
				indent = line[:idx]
			}
			lines[i] = fmt.Sprintf("%sname: %s", indent, moduleName)
			break
		}
	}

	newContent := strings.Join(lines, "\n")
	return os.WriteFile(skillYamlPath, []byte(newContent), 0644)
}

// GetDefaultTemplate 获取默认模板路径
func GetDefaultTemplate() string {
	// 尝试多个可能的路径
	possiblePaths := []string{
		"./empty",
		"../empty",
		"../../empty",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		// 回退到当前目录
		return "./empty"
	}

	execDir := filepath.Dir(execPath)
	templatePath := filepath.Join(execDir, "empty")

	// 检查是否存在
	if _, err := os.Stat(templatePath); err == nil {
		return templatePath
	}

	// 最终回退
	return "./empty"
}
