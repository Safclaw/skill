package main

import (
	"fmt"
	"os"

	"github.com/Safclaw/skill/pkg/config"
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

	// TODO: 实现完整的下载和安装逻辑
	// 3. 检查缓存
	// 4. 下载（如果缓存未命中）
	// 5. 安装
	// 6. 执行 hooks

	fmt.Println("Note: Full implementation pending")
	fmt.Printf("Module: %s\n", modPath.String())
	fmt.Printf("Host: %s\n", modPath.Host)
	fmt.Printf("Namespace: %s\n", modPath.Namespace)
	fmt.Printf("Name: %s\n", modPath.Name)
	fmt.Printf("Version: %s\n", modPath.Version)
	fmt.Printf("Install Dir: %s\n", installDir)

	return nil
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
