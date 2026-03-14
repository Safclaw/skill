package cache

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Manager 缓存管理器
type Manager struct {
	cacheDir string
}

// NewManager 创建缓存管理器
func NewManager(cacheDir string) *Manager {
	return &Manager{
		cacheDir: cacheDir,
	}
}

// CacheInfo 缓存信息
type CacheInfo struct {
	Path       string
	Version    string
	Checksum   string
	Downloaded time.Time
	Size       int64
}

// Get 从缓存获取 skill
func (m *Manager) Get(host, namespace, name, version string) (*CacheInfo, error) {
	cachePath := m.getCachePath(host, namespace, name, version)

	// 检查缓存是否存在
	info, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("cache path is not a directory")
	}

	// 读取校验和文件
	checksumFile := filepath.Join(cachePath, "checksum.sha256")
	checksumData, err := os.ReadFile(checksumFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read checksum file")
	}

	return &CacheInfo{
		Path:       cachePath,
		Version:    version,
		Checksum:   string(checksumData),
		Downloaded: info.ModTime(),
		Size:       0, // TODO: 计算目录大小
	}, nil
}

// Put 将 skill 放入缓存
func (m *Manager) Put(host, namespace, name, version string, zipPath string, checksum string) (string, error) {
	cachePath := m.getCachePath(host, namespace, name, version)

	// 确保缓存目录存在
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return "", err
	}

	// 创建临时解压目录
	tmpDir, err := os.MkdirTemp("", "skill-cache-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// 解压 zip 文件
	if err := extractZip(zipPath, tmpDir); err != nil {
		return "", fmt.Errorf("failed to extract zip: %w", err)
	}

	// 规范化路径（移除顶层目录）
	if err := normalizeZipPath(tmpDir, cachePath); err != nil {
		return "", fmt.Errorf("failed to normalize path: %w", err)
	}

	// 写入校验和文件
	checksumFile := filepath.Join(cachePath, "checksum.sha256")
	if err := os.WriteFile(checksumFile, []byte(checksum), 0644); err != nil {
		return "", err
	}

	return cachePath, nil
}

// Verify 验证缓存完整性
func (m *Manager) Verify(host, namespace, name, version string) error {
	cachePath := m.getCachePath(host, namespace, name, version)

	// 读取存储的校验和
	checksumFile := filepath.Join(cachePath, "checksum.sha256")
	storedChecksum, err := os.ReadFile(checksumFile)
	if err != nil {
		return fmt.Errorf("failed to read stored checksum")
	}

	// 重新计算校验和
	computedChecksum, err := m.computeChecksum(cachePath)
	if err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	if string(storedChecksum) != computedChecksum {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

// Clean 清理缓存
func (m *Manager) Clean(all bool) error {
	if all {
		// 清理所有缓存
		return os.RemoveAll(m.cacheDir)
	}

	// TODO: 实现 LRU 清理策略
	// 目前只清理超过 30 天的缓存
	cutoff := time.Now().AddDate(0, 0, -30)

	return filepath.Walk(m.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			return os.RemoveAll(path)
		}

		return nil
	})
}

// getCachePath 获取缓存路径
func (m *Manager) getCachePath(host, namespace, name, version string) string {
	// Windows 兼容性：使用 ! 前缀
	escapedHost := "!" + host
	escapedNamespace := namespace
	escapedName := name

	return filepath.Join(m.cacheDir, escapedHost, escapedNamespace, escapedName+"@"+version)
}

// computeChecksum 计算目录的校验和
func (m *Manager) computeChecksum(dirPath string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过 checksum 文件本身
		if info.Name() == "checksum.sha256" {
			return nil
		}

		// 写入文件信息
		fmt.Fprintf(hash, "%s %d\n", info.Name(), info.Size())

		// 如果是文件，读取内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(hash, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// extractZip 解压 zip 文件
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		targetPath := filepath.Join(destDir, f.Name)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

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

// normalizeZipPath 规范化 zip 解压后的路径
func normalizeZipPath(sourceDir, destDir string) error {
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

	return moveDirectory(sourcePath, destDir)
}

func moveDirectory(source, dest string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return os.Rename(path, targetPath)
	})
}
