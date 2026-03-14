package main

import (
	"github.com/Safclaw/skill/pkg/config"
	"github.com/Safclaw/skill/pkg/lister"
	"github.com/spf13/cobra"
)

var (
	listGlobalFlag    bool
	listWorkspaceFlag bool
	listWorkspacePath string
)

func initListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed skills",
		Long: `List installed skills.

Examples:
  # List global skills
  skill list -g
  
  # List workspace skills
  skill list -w
  
  # List skills in specified directory
  skill list --workspace /path/to/dir`,
		Aliases: []string{"ls"},
		RunE:    runList,
	}

	cmd.Flags().BoolVarP(&listGlobalFlag, "global", "g", false, "List global skills")
	cmd.Flags().BoolVarP(&listWorkspaceFlag, "workspace", "w", false, "List workspace skills")
	cmd.Flags().StringVar(&listWorkspacePath, "workspace-path", "", "List skills in specified directory")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// 确定安装目录
	cfg := config.LoadConfig()
	installDir, err := getInstallDir(cfg, listGlobalFlag, listWorkspaceFlag, listWorkspacePath)
	if err != nil {
		return err
	}

	// 创建列表器
	l := lister.NewLister(installDir)

	// 打印列表
	return l.PrintList()
}
