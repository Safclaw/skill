package downloader

import (
	"context"
	"io"
)

// DownloadResult 下载结果
type DownloadResult struct {
	Path        string // 下载的临时文件路径
	Version     string // 实际版本
	Checksum    string // SHA256 校验和
	Size        int64  // 文件大小（字节）
	ContentType string // 内容类型
}

// Downloader 下载器接口
type Downloader interface {
	// Download 下载指定版本的 skill
	Download(ctx context.Context, host, namespace, name, version string) (*DownloadResult, error)

	// ListVersions 列出可用版本
	ListVersions(ctx context.Context, host, namespace, name string) ([]string, error)
}

// ProgressWriter 进度写入器
type ProgressWriter struct {
	Total      int64
	Written    int64
	OnProgress func(written, total int64)
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pw.Written += int64(n)
	if pw.OnProgress != nil {
		pw.OnProgress(pw.Written, pw.Total)
	}
	return n, nil
}

// NopCloser 创建一个不执行任何操作的 Close 方法
func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{r}
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
