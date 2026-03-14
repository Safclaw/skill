package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Safclaw/skill/pkg/hook"
	"github.com/Safclaw/skill/pkg/manifest"
)

// Installer 安装器
type Installer struct {
	hookExecutor *hook.Executor
}

// NewInstaller 创建安装器
func NewInstaller() *Installer {
	return &Installer{
		hookExecutor: hook.NewExecutor(),
	}
}

// InstallOptions 安装选项
type InstallOptions struct {
	InstallDir string // 安装目录
	SkillPath  string // skill 源路径（从缓存或下载）
	ModuleDir  string // 模块相对路径（如：github.com/Safclaw/skills/read-json）
	Version    string // 版本号
	Signature  string // SHA256 签名
	RunHooks   bool   // 是否执行 hooks
}

// InstallResult 安装结果
type InstallResult struct {
	Success bool
	Message string
}

// Install 安装 skill
func (i *Installer) Install(opts InstallOptions) (*InstallResult, error) {
	// 1. 确定目标路径
	targetPath := filepath.Join(opts.InstallDir, "reps", opts.ModuleDir)

	// 2. 确保目标目录存在
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 3. 复制 skill 文件到目标目录
	if err := copyDirectory(opts.SkillPath, targetPath); err != nil {
		return nil, fmt.Errorf("failed to copy files: %w", err)
	}

	// 4. 如果启用 hooks，执行 post_add hooks
	if opts.RunHooks {
		if err := i.executePostAddHooks(targetPath); err != nil {
			// Hook 执行失败，回滚
			_ = os.RemoveAll(targetPath)
			return nil, fmt.Errorf("post_add hooks failed: %w", err)
		}
	}

	// 5. 更新清单文件
	manifestPath := filepath.Join(opts.InstallDir, ".skills.yaml")
	if err := i.updateManifest(manifestPath, opts); err != nil {
		// 清单更新失败，回滚
		_ = os.RemoveAll(targetPath)
		return nil, fmt.Errorf("failed to update manifest: %w", err)
	}

	return &InstallResult{
		Success: true,
		Message: fmt.Sprintf("Successfully installed %s@%s", opts.ModuleDir, opts.Version),
	}, nil
}

// executePostAddHooks 执行 post_add hooks
func (i *Installer) executePostAddHooks(skillPath string) error {
	// 读取 skill.yaml
	skillYamlPath := filepath.Join(skillPath, "skill.yaml")
	if _, err := os.Stat(skillYamlPath); os.IsNotExist(err) {
		return nil // 没有 skill.yaml，跳过 hooks
	}

	// 解析 hooks
	parser := hook.NewParser()
	hooks, err := parser.Parse(skillYamlPath)
	if err != nil {
		return fmt.Errorf("failed to parse hooks: %w", err)
	}

	if len(hooks) == 0 {
		return nil // 没有 hooks 定义
	}

	// 执行 hooks
	_, err = i.hookExecutor.Execute(hook.PostAdd, skillPath, hooks)
	return err
}

// updateManifest 更新清单文件
func (i *Installer) updateManifest(manifestPath string, opts InstallOptions) error {
	// 读取现有清单
	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return err
	}

	// 添加或更新技能条目
	entry := manifest.SkillEntry{
		Name:    opts.ModuleDir, // TODO: 从 skill.yaml 读取可读名称
		Dir:     opts.ModuleDir,
		Version: opts.Version,
		Sig:     opts.Signature,
	}
	m.AddSkill(entry)

	// 写回清单文件
	return manifest.WriteManifest(manifestPath, m)
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
