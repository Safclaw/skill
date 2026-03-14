package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser Hook 配置解析器
type Parser struct{}

// NewParser 创建解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 从 YAML 文件解析 Hook 配置
func (p *Parser) Parse(filePath string) ([]HookConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return p.ParseBytes(data)
}

// ParseBytes 从字节解析 Hook 配置
func (p *Parser) ParseBytes(data []byte) ([]HookConfig, error) {
	var config struct {
		Hooks []HookConfig `yaml:"hooks"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Hooks, nil
}

// Validate 验证 Hook 配置
func (p *Parser) Validate(hooks []HookConfig, skillPath string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	for i, hook := range hooks {
		// 验证 stage
		if hook.Stage != PostAdd && hook.Stage != PreRemove {
			result.Errors = append(result.Errors,
				fmt.Sprintf("hook[%d]: invalid stage '%s', must be 'post_add' or 'pre_remove'", i, hook.Stage))
			result.Valid = false
		}

		// 验证 scripts
		if len(hook.Scripts) == 0 {
			result.Errors = append(result.Errors,
				fmt.Sprintf("hook[%d]: no scripts defined", i))
			result.Valid = false
			continue
		}

		for j, script := range hook.Scripts {
			// 验证 command
			if script.Command == "" {
				result.Errors = append(result.Errors,
					fmt.Sprintf("hook[%d].script[%d]: command is required", i, j))
				result.Valid = false
			}

			// 验证 platforms
			if len(script.Platforms) == 0 {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("hook[%d].script[%d]: no platforms specified, will run on all platforms", i, j))
			}

			// 验证 args
			if len(script.Args) == 0 {
				result.Errors = append(result.Errors,
					fmt.Sprintf("hook[%d].script[%d]: args is required", i, j))
				result.Valid = false
				continue
			}

			// 验证脚本路径安全性
			scriptPath := script.Args[0]
			if err := validateScriptPath(scriptPath); err != nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf("hook[%d].script[%d]: %v", i, j, err))
				result.Valid = false
			}

			// 验证 checksum 格式
			if script.Checksum != "" && !strings.HasPrefix(script.Checksum, "sha256:") {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("hook[%d].script[%d]: checksum should start with 'sha256:'", i, j))
			}

			// 验证脚本文件是否存在（如果 skillPath 已提供）
			if skillPath != "" {
				fullPath := filepath.Join(skillPath, scriptPath)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					result.Errors = append(result.Errors,
						fmt.Sprintf("hook[%d].script[%d]: script file not found: %s", i, j, fullPath))
					result.Valid = false
				}
			}
		}
	}

	return result
}

// validateScriptPath 验证脚本路径安全性
func validateScriptPath(path string) error {
	// 必须是相对路径
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute path not allowed: %s", path)
	}

	// 不能包含 ..
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	// 必须在 scripts/ 目录下
	if !strings.HasPrefix(path, "scripts/") {
		return fmt.Errorf("script must be in scripts/ directory: %s", path)
	}

	return nil
}
