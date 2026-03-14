package hook

import (
	"runtime"
)

// Platform 平台工具
type Platform struct{}

// NewPlatform 创建平台工具
func NewPlatform() *Platform {
	return &Platform{}
}

// GetCurrentOS 获取当前操作系统
func (p *Platform) GetCurrentOS() string {
	return runtime.GOOS
}

// SelectScript 根据当前平台选择合适的脚本
func (p *Platform) SelectScript(scripts []HookScript) *HookScript {
	currentOS := p.GetCurrentOS()

	// 第一次遍历：查找完全匹配的脚本
	for i := range scripts {
		if contains(scripts[i].Platforms, currentOS) {
			return &scripts[i]
		}
	}

	// 第二次遍历：查找没有指定平台的脚本（通用脚本）
	for i := range scripts {
		if len(scripts[i].Platforms) == 0 {
			return &scripts[i]
		}
	}

	// 没有找到合适的脚本
	return nil
}

// IsSupported 检查脚本是否支持当前平台
func (p *Platform) IsSupported(script HookScript) bool {
	if len(script.Platforms) == 0 {
		return true // 没有指定平台，支持所有平台
	}

	return contains(script.Platforms, runtime.GOOS)
}

// GetCommandForPlatform 根据平台获取命令
func (p *Platform) GetCommandForPlatform(script HookScript) string {
	if script.Command != "" {
		return script.Command
	}

	// 根据平台返回默认命令
	switch runtime.GOOS {
	case "windows":
		return "powershell"
	case "darwin", "linux":
		return "bash"
	default:
		return "sh"
	}
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
