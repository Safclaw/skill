package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Safclaw/skill/pkg/cache"
	"github.com/Safclaw/skill/pkg/config"
	"github.com/Safclaw/skill/pkg/downloader"
	"github.com/Safclaw/skill/pkg/installer"
	"github.com/Safclaw/skill/pkg/modulepath"
	"github.com/spf13/cobra"
)

var (
	globalFlag    bool
	workspaceFlag bool
	workspacePath string
	noHooksFlag   bool
)

func initAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <module>@<version>",
		Short: "Add a skill",
		Long: `Add a skill to the installation.

Examples:
  # Add to global installation
  skill add github.com/Safclaw/skills/read-json@v1.0.0 -g
  
  # Add to workspace
  skill add github.com/Safclaw/skills/read-json@latest -w
  
  # Add to specified directory
  skill add github.com/Safclaw/skills/read-json --workspace /path/to/dir`,
		Args: cobra.ExactArgs(1),
		RunE: runAdd,
	}

	cmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Install globally")
	cmd.Flags().BoolVarP(&workspaceFlag, "workspace", "w", false, "Install to workspace")
	cmd.Flags().StringVar(&workspacePath, "workspace-path", "", "Install to specified directory")
	cmd.Flags().BoolVar(&noHooksFlag, "no-hooks", false, "Skip running hooks")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	// 1. 解析模块路径
	modPath, err := modulepath.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid module path: %w", err)
	}

	// 2. 确定安装目录
	cfg := config.LoadConfig()
	installDir, err := getInstallDir(cfg, globalFlag, workspaceFlag, workspacePath)
	if err != nil {
		return err
	}

	fmt.Printf("Installing %s to %s...\n", modPath.String(), installDir)

	// 3. 检查缓存
	cacheManager := cache.NewManager(cfg.CacheDir)
	cacheInfo, err := cacheManager.Get(modPath.Host, modPath.Namespace, modPath.Name, modPath.Version)
	if err != nil {
		return fmt.Errorf("cache check failed: %w", err)
	}

	var skillPath string
	var checksum string
	var extractedPath string // Path after extracting subdirectory

	if cacheInfo != nil {
		// 缓存命中
		fmt.Println("Using cached version...")
		skillPath = cacheInfo.Path
		checksum = cacheInfo.Checksum
		extractedPath = skillPath
	} else {
		// 缓存未命中，需要下载
		fmt.Println("Downloading...")
		downloader := downloader.NewGitHubDownloader()
		downloadResult, err := downloader.Download(cmd.Context(), modPath.Host, modPath.Namespace, modPath.Name, modPath.Version)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer os.Remove(downloadResult.Path) // 清理临时 zip 文件

		// 保存到缓存
		skillPath, err = cacheManager.Put(modPath.Host, modPath.Namespace, modPath.Name, modPath.Version, downloadResult.Path, downloadResult.Checksum)
		if err != nil {
			return fmt.Errorf("failed to save to cache: %w", err)
		}
		checksum = downloadResult.Checksum
		extractedPath = skillPath

		fmt.Printf("Downloaded version: %s (size: %d bytes, checksum: %s)\n",
			downloadResult.Version, downloadResult.Size, downloadResult.Checksum[:16]+"...")
	}

	// Handle subdirectory extraction
	if modPath.SubDir != "" {
		fmt.Printf("Extracting subdirectory: %s\n", modPath.SubDir)
		extractedPath, err = extractorSubdirectory(extractedPath, modPath.SubDir)
		if err != nil {
			return fmt.Errorf("failed to extract subdirectory: %w", err)
		}
	}

	// 4. 安装 skill
	inst := installer.NewInstaller()
	result, err := inst.Install(installer.InstallOptions{
		InstallDir: installDir,
		SkillPath:  extractedPath,
		ModuleDir:  modPath.DirPath(),
		Version:    modPath.Version,
		Signature:  checksum,
		RunHooks:   !noHooksFlag,
	})
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("✓ %s\n", result.Message)
	return nil
}

// extractorSubdirectory extracts a subdirectory from the downloaded content
func extractorSubdirectory(basePath, subDir string) (string, error) {
	subDirPath := filepath.Join(basePath, subDir)

	// Check if subdirectory exists
	info, err := os.Stat(subDirPath)
	if err != nil {
		return "", fmt.Errorf("subdirectory not found: %s", subDir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", subDir)
	}

	return subDirPath, nil
}

// getInstallDir 根据标志确定安装目录
func getInstallDir(cfg *config.Config, global, workspace bool, customPath string) (string, error) {
	count := 0
	if global {
		count++
	}
	if workspace {
		count++
	}
	if customPath != "" {
		count++
	}

	if count > 1 {
		return "", fmt.Errorf("cannot specify multiple installation modes")
	}

	if global {
		return cfg.GlobalDir, nil
	}

	if workspace {
		return cfg.WorkspaceDir, nil
	}

	if customPath != "" {
		// 确保目录存在
		if err := os.MkdirAll(customPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return customPath, nil
	}

	// 默认使用全局安装
	return cfg.GlobalDir, nil
}
