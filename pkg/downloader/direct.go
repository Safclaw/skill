package downloader

import (
	"context"
	"fmt"
)

// DirectDownloader 直连下载器（不使用代理）
type DirectDownloader struct {
	github *GitHubDownloader
	// TODO: 添加 GitLab, Gitee 下载器
}

// NewDirectDownloader 创建直连下载器
func NewDirectDownloader() *DirectDownloader {
	return &DirectDownloader{
		github: NewGitHubDownloader(),
	}
}

// Download 下载指定版本的 skill
func (d *DirectDownloader) Download(ctx context.Context, host, namespace, name, version string) (*DownloadResult, error) {
	switch host {
	case "github.com":
		return d.github.Download(ctx, host, namespace, name, version)
	case "gitlab.com":
		// TODO: 实现 GitLab 下载器
		return nil, fmt.Errorf("GitLab download not implemented yet")
	case "gitee.com":
		// TODO: 实现 Gitee 下载器
		return nil, fmt.Errorf("Gitee download not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported host: %s", host)
	}
}

// ListVersions 列出可用版本
func (d *DirectDownloader) ListVersions(ctx context.Context, host, namespace, name string) ([]string, error) {
	switch host {
	case "github.com":
		return d.github.ListVersions(ctx, host, namespace, name)
	default:
		return []string{"latest"}, nil
	}
}
