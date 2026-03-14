package downloader

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// GitHubDownloader GitHub 下载器
type GitHubDownloader struct {
	client *http.Client
}

// NewGitHubDownloader 创建 GitHub 下载器
func NewGitHubDownloader() *GitHubDownloader {
	return &GitHubDownloader{
		client: &http.Client{},
	}
}

// Download 从 GitHub 下载 skill
func (d *GitHubDownloader) Download(ctx context.Context, host, namespace, name, version string) (*DownloadResult, error) {
	if host != "github.com" {
		return nil, fmt.Errorf("not a github host: %s", host)
	}

	// 构建下载 URL
	downloadURL := fmt.Sprintf("https://codeload.github.com/%s/%s/zip/%s", namespace, name, version)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "skill-download-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// 计算 SHA256 并保存文件
	hash := sha256.New()
	writer := io.MultiWriter(tmpFile, hash)

	size, err := io.Copy(writer, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	return &DownloadResult{
		Path:        tmpFile.Name(),
		Version:     version,
		Checksum:    checksum,
		Size:        size,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}

// ListVersions 列出 GitHub 仓库的可用版本（通过 tags）
func (d *GitHubDownloader) ListVersions(ctx context.Context, host, namespace, name string) ([]string, error) {
	if host != "github.com" {
		return nil, fmt.Errorf("not a github host: %s", host)
	}

	// 使用 GitHub API 获取 tags
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", namespace, name)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list versions: %d", resp.StatusCode)
	}

	// TODO: 解析 JSON 响应
	// 这里简化处理，实际需要使用 encoding/json 解析
	return []string{"latest"}, nil
}

// ExtractZip 解压 zip 文件到目标目录
func ExtractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// 跳过目录
		if f.FileInfo().IsDir() {
			continue
		}

		// 构建目标路径
		targetPath := filepath.Join(destDir, f.Name)

		// 确保目标目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// 提取文件
		if err := extractFile(f, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(f *zip.File, targetPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// NormalizeZipPath 规范化 zip 解压后的路径（移除顶层目录）
func NormalizeZipPath(sourceDir, destDir string) error {
	// 查找顶层目录
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	if len(entries) != 1 {
		return fmt.Errorf("expected single top-level directory, got %d entries", len(entries))
	}

	topLevel := entries[0]
	if !topLevel.IsDir() {
		return fmt.Errorf("expected top-level entry to be a directory")
	}

	sourcePath := filepath.Join(sourceDir, topLevel.Name())

	// 移动所有文件到目标目录
	return moveDirectory(sourcePath, destDir)
}

func moveDirectory(source, dest string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// 移动文件
		return os.Rename(path, targetPath)
	})
}

// CleanupEmptyDirs 清理空目录
func CleanupEmptyDirs(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		// 检查目录是否为空
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			return os.Remove(path)
		}

		return nil
	})
}
