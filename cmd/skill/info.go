package main

import (
	"fmt"

	"github.com/Safclaw/skill/pkg/config"
	"github.com/Safclaw/skill/pkg/lister"
	"github.com/spf13/cobra"
)

func initInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <module>",
		Short: "Show skill information",
		Long: `Show detailed information about a skill.

Examples:
  skill info github.com/Safclaw/skills/read-json`,
		Args: cobra.ExactArgs(1),
		RunE: runInfo,
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	moduleDir := args[0]

	// 默认使用全局安装目录
	cfg := config.LoadConfig()
	installDir := cfg.GlobalDir

	// 创建列表器
	l := lister.NewLister(installDir)

	// 获取技能信息
	info, err := l.GetSkillInfo(moduleDir)
	if err != nil {
		return err
	}

	// 打印信息
	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Module: %s\n", info.ModuleDir)
	fmt.Printf("Version: %s\n", info.Version)
	fmt.Printf("Path: %s\n", info.Path)
	if info.Signature != "" {
		fmt.Printf("Signature: %s\n", info.Signature)
	}

	return nil
}
