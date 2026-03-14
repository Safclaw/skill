package main

import (
	"fmt"

	"github.com/Safclaw/skill/pkg/config"
	"github.com/spf13/cobra"
)

var (
	removeGlobalFlag    bool
	removeWorkspaceFlag bool
	removeWorkspacePath string
	removeNoHooksFlag   bool
)

func initRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <module>",
		Short: "Remove a skill",
		Long: `Remove a skill from the installation.

Examples:
  # Remove from global installation
  skill remove github.com/Safclaw/skills/read-json -g
  
  # Remove from workspace
  skill remove github.com/Safclaw/skills/read-json -w
  
  # Remove from specified directory
  skill remove github.com/Safclaw/skills/read-json --workspace /path/to/dir`,
		Args: cobra.ExactArgs(1),
		RunE: runRemove,
	}

	cmd.Flags().BoolVarP(&removeGlobalFlag, "global", "g", false, "Remove from global installation")
	cmd.Flags().BoolVarP(&removeWorkspaceFlag, "workspace", "w", false, "Remove from workspace")
	cmd.Flags().StringVar(&removeWorkspacePath, "workspace-path", "", "Remove from specified directory")
	cmd.Flags().BoolVar(&removeNoHooksFlag, "no-hooks", false, "Skip running hooks")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	moduleDir := args[0]

	// 确定安装目录
	cfg := config.LoadConfig()
	installDir, err := getInstallDir(cfg, removeGlobalFlag, removeWorkspaceFlag, removeWorkspacePath)
	if err != nil {
		return err
	}

	fmt.Printf("Removing %s from %s...\n", moduleDir, installDir)

	// TODO: 实现完整的卸载逻辑
	// 1. 执行 pre_remove hooks
	// 2. 删除技能目录
	// 3. 更新清单文件

	fmt.Println("Note: Full implementation pending")
	fmt.Printf("Module: %s\n", moduleDir)
	fmt.Printf("Install Dir: %s\n", installDir)

	return nil
}
