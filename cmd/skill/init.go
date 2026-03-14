package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Safclaw/skill/pkg/config"
	"github.com/Safclaw/skill/pkg/template"
	"github.com/spf13/cobra"
)

var (
	templateFlag string
)

func initInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <module-name>",
		Short: "Initialize a new skill",
		Long: `Initialize a new skill project.

Examples:
  skill init github.com/myorg/my-skill
  skill init github.com/myorg/my-skill --template github.com/Safclaw/skill/empty
  skill init github.com/myorg/my-skill --template ./my-template`,
		Args: cobra.ExactArgs(1),
		RunE: runInit,
	}

	cmd.Flags().StringVar(&templateFlag, "template", "", "Template to use (local path or remote repository)")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 确定模板路径
	tmplPath := templateFlag
	if tmplPath == "" {
		// 使用默认模板
		tmplPath = template.GetDefaultTemplate()
	}

	// 创建模板管理器
	cfg := config.DefaultConfig()
	tmplManager := template.NewTemplateManager(cfg.CacheDir)

	// 从模板复制文件
	ctx := context.Background()
	if err := tmplManager.CopyTemplate(ctx, tmplPath, workDir, moduleName); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	fmt.Printf("\nSuccessfully initialized skill in current directory\n")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit skill.yaml with your skill metadata")
	fmt.Println("  2. Implement your skill in skill.md")
	fmt.Println("  3. Add installation scripts in scripts/")
	fmt.Println("  4. Test your skill locally")
	fmt.Println("  5. Publish to a Git repository")

	return nil
}

