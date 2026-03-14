package hook

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ChecksumVerifier 校验和验证器
type ChecksumVerifier struct{}

// NewChecksumVerifier 创建验证器
func NewChecksumVerifier() *ChecksumVerifier {
	return &ChecksumVerifier{}
}

// Verify 验证文件校验和
func (v *ChecksumVerifier) Verify(filePath string, expectedChecksum string) error {
	if expectedChecksum == "" {
		return nil // 没有指定 checksum，跳过验证
	}

	// 计算实际 checksum
	actualChecksum, err := v.computeFileChecksum(filePath)
	if err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	// 移除 "sha256:" 前缀进行比较
	expectedClean := strings.TrimPrefix(expectedChecksum, "sha256:")

	if actualChecksum != expectedClean {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedClean, actualChecksum)
	}

	return nil
}

// computeFileChecksum 计算文件的 SHA256 checksum
func (v *ChecksumVerifier) computeFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ComputeScriptChecksum 计算脚本的 checksum（用于生成配置）
func (v *ChecksumVerifier) ComputeScriptChecksum(skillPath string, scriptPath string) (string, error) {
	fullPath := filepath.Join(skillPath, scriptPath)
	checksum, err := v.computeFileChecksum(fullPath)
	if err != nil {
		return "", err
	}

	return "sha256:" + checksum, nil
}

// VerifyAllScripts 验证所有脚本的 checksum
func (v *ChecksumVerifier) VerifyAllScripts(skillPath string, scripts []HookScript) []error {
	errors := make([]error, 0)

	for i, script := range scripts {
		if script.Checksum == "" {
			continue // 跳过没有 checksum 的脚本
		}

		if len(script.Args) == 0 {
			errors = append(errors, fmt.Errorf("script[%d]: no args specified", i))
			continue
		}

		scriptPath := script.Args[0]
		fullPath := filepath.Join(skillPath, scriptPath)

		if err := v.Verify(fullPath, script.Checksum); err != nil {
			errors = append(errors, fmt.Errorf("script[%d] (%s): %w", i, scriptPath, err))
		}
	}

	return errors
}
