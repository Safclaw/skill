package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Safclaw/skill/pkg/hook"
	"github.com/Safclaw/skill/pkg/manifest"
)

// Uninstaller 卸载器
type Uninstaller struct {
	hookExecutor *hook.Executor
}

// NewUninstaller 创建卸载器
func NewUninstaller() *Uninstaller {
	return &Uninstaller{
		hookExecutor: hook.NewExecutor(),
	}
}

// UninstallOptions 卸载选项
type UninstallOptions struct {
	InstallDir string // 安装目录
	ModuleDir  string // 模块相对路径
	RunHooks   bool   // 是否执行 hooks
}

// UninstallResult 卸载结果
type UninstallResult struct {
	Success bool
	Message string
}

// Uninstall 卸载 skill
func (u *Uninstaller) Uninstall(opts UninstallOptions) (*UninstallResult, error) {
	// 1. 确定技能路径
	skillPath := filepath.Join(opts.InstallDir, "reps", opts.ModuleDir)

	// 检查技能是否存在
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("skill not found: %s", opts.ModuleDir)
	}

	// 2. 如果启用 hooks，执行 pre_remove hooks
	if opts.RunHooks {
		if err := u.executePreRemoveHooks(skillPath); err != nil {
			// Hook 执行失败，但仍然继续卸载（记录警告）
			fmt.Printf("Warning: pre_remove hooks failed: %v\n", err)
		}
	}

	// 3. 删除技能目录
	if err := os.RemoveAll(skillPath); err != nil {
		return nil, fmt.Errorf("failed to remove skill directory: %w", err)
	}

	// 4. 清理空的父目录
	_ = cleanupEmptyDirs(filepath.Join(opts.InstallDir, "reps"))

	// 5. 更新清单文件
	manifestPath := filepath.Join(opts.InstallDir, ".skills.yaml")
	if err := u.updateManifest(manifestPath, opts.ModuleDir); err != nil {
		return nil, fmt.Errorf("failed to update manifest: %w", err)
	}

	return &UninstallResult{
		Success: true,
		Message: fmt.Sprintf("Successfully uninstalled %s", opts.ModuleDir),
	}, nil
}

// executePreRemoveHooks 执行 pre_remove hooks
func (u *Uninstaller) executePreRemoveHooks(skillPath string) error {
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
	_, err = u.hookExecutor.Execute(hook.PreRemove, skillPath, hooks)
	return err
}

// updateManifest 更新清单文件
func (u *Uninstaller) updateManifest(manifestPath string, moduleDir string) error {
	// 读取现有清单
	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return err
	}

	// 删除技能条目
	m.RemoveSkill(moduleDir)

	// 写回清单文件
	return manifest.WriteManifest(manifestPath, m)
}

// cleanupEmptyDirs 清理空目录
func cleanupEmptyDirs(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		// 检查目录是否为空
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			return os.Remove(path)
		}

		return nil
	})
}
