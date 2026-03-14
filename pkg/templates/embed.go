package templates

import (
	"embed"
)

//go:embed empty/*
var EmptyTemplate embed.FS

// GetEmptyTemplateFS 返回嵌入的空模板文件系统
func GetEmptyTemplateFS() embed.FS {
	return EmptyTemplate
}
