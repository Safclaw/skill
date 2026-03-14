package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

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
  skill init github.com/myorg/my-skill --template github.com/Safclaw/skill/empty`,
		Args: cobra.ExactArgs(1),
		RunE: runInit,
	}

	cmd.Flags().StringVar(&templateFlag, "template", "", "Template to use")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	// 创建目录
	if err := os.MkdirAll(moduleName, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建基础文件结构
	files := map[string]string{
		"skill.yaml":         skillYamlTemplate,
		"skill.md":           skillMdTemplate,
		"README.md":          readmeTemplate,
		"scripts/setup.sh":   setupShTemplate,
		"scripts/unsetup.sh": unsetupShTemplate,
	}

	for file, content := range files {
		filePath := filepath.Join(moduleName, file)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		// 创建文件
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		// 渲染模板
		tmpl, err := template.New(file).Parse(content)
		if err != nil {
			f.Close()
			return err
		}

		data := map[string]interface{}{
			"ModuleName": moduleName,
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return err
		}

		f.Close()
		fmt.Printf("Created: %s\n", filePath)
	}

	fmt.Printf("\nSuccessfully initialized skill in %s/\n", moduleName)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit skill.yaml with your skill metadata")
	fmt.Println("  2. Implement your skill in skill.md")
	fmt.Println("  3. Add installation scripts in scripts/")
	fmt.Println("  4. Test your skill locally")
	fmt.Println("  5. Publish to a Git repository")

	return nil
}

const skillYamlTemplate = `# ==============================
# Skill Metadata
# ==============================
name: {{.ModuleName}}
description: "A new skill"
# entrypoint: "skill.md"

# Authors
authors:
  - name: "Your Name"
    email: "your.email@example.com"

license: "MIT"
tags: []

# Dependencies
# dependencies:
#   - name: github.com/Safclaw/skills/xlsx
#     version: v1.0.0

# Permissions (for reference only, enforced by runtime)
# permissions:
#   storage: []
#   network: []
#   execution: []

# Hooks
# hooks:
#   - stage: "post_add"
#     reason: "Setup script"
#     timeout: 30
#     scripts:
#       - command: "bash"
#         platforms: ["linux", "darwin"]
#         args: ["./scripts/setup.sh"]
`

const skillMdTemplate = `# {{.ModuleName}}

## Usage

Describe how to use this skill here.

## Implementation

Implement your skill logic here.
`

const readmeTemplate = `# {{.ModuleName}}

A new skill for SafeClaw.

## Installation

` + "```bash" + `
skill add {{.ModuleName}}@latest
` + "```" + `

## Usage

TODO: Add usage instructions

## Development

` + "```bash" + `
# Initialize the skill
skill init {{.ModuleName}}

# Test locally
# TODO: Add testing instructions
` + "```" + `

## License

MIT
`

const setupShTemplate = `#!/bin/bash
set -e

echo "Setting up {{.ModuleName}}..."

# Add your setup logic here

echo "Setup complete!"
`

const unsetupShTemplate = `#!/bin/bash
set -e

echo "Cleaning up {{.ModuleName}}..."

# Add your cleanup logic here

echo "Cleanup complete!"
`
