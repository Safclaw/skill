package signature

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ComputeFile 计算文件的 SHA256 签名
func ComputeFile(filePath string) (string, error) {
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

// ComputeDirectory 计算目录的 SHA256 签名
func ComputeDirectory(dirPath string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过 checksum 文件本身
		if info.Name() == "checksum.sha256" {
			return nil
		}

		// 写入文件信息（用于保证文件列表的一致性）
		fmt.Fprintf(hash, "%s %d %d\n", info.Name(), info.Size(), info.ModTime().Unix())

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

// VerifyFile 验证文件签名
func VerifyFile(filePath, expectedSig string) error {
	actualSig, err := ComputeFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to compute signature: %w", err)
	}

	if actualSig != expectedSig {
		return fmt.Errorf("signature mismatch: expected %s, got %s", expectedSig, actualSig)
	}

	return nil
}

// VerifyDirectory 验证目录签名
func VerifyDirectory(dirPath, expectedSig string) error {
	actualSig, err := ComputeDirectory(dirPath)
	if err != nil {
		return fmt.Errorf("failed to compute signature: %w", err)
	}

	if actualSig != expectedSig {
		return fmt.Errorf("signature mismatch: expected %s, got %s", expectedSig, actualSig)
	}

	return nil
}

// WriteSignature 将签名写入文件
func WriteSignature(filePath, signature string) error {
	return os.WriteFile(filePath, []byte(signature), 0644)
}

// ReadSignature 从文件读取签名
func ReadSignature(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
