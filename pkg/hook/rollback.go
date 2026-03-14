package hook

import (
	"fmt"
	"os"
	"path/filepath"
)

// RollbackManager 回滚管理器
type RollbackManager struct {
	executedScripts []ExecuteResult
	skillPath       string
}

// NewRollbackManager 创建回滚管理器
func NewRollbackManager(skillPath string) *RollbackManager {
	return &RollbackManager{
		executedScripts: make([]ExecuteResult, 0),
		skillPath:       skillPath,
	}
}

// RecordExecution 记录执行的脚本
func (rm *RollbackManager) RecordExecution(result ExecuteResult) {
	rm.executedScripts = append(rm.executedScripts, result)
}

// Rollback 执行回滚
func (rm *RollbackManager) Rollback() error {
	if len(rm.executedScripts) == 0 {
		return nil // 没有执行任何脚本，无需回滚
	}

	// 反向遍历已执行的脚本
	for i := len(rm.executedScripts) - 1; i >= 0; i-- {
		result := rm.executedScripts[i]

		// TODO: 根据脚本类型执行相应的回滚操作
		// 目前仅记录日志，实际回滚逻辑需要根据具体脚本定制
		fmt.Printf("Rollback: script %s was executed (duration: %v)\n",
			result.Script, result.Duration)
	}

	return nil
}

// Cleanup 清理安装目录（用于 post_add 失败时的回滚）
func (rm *RollbackManager) Cleanup() error {
	if rm.skillPath == "" {
		return fmt.Errorf("skill path is empty")
	}

	// 删除整个技能目录
	return os.RemoveAll(rm.skillPath)
}

// Backup 备份技能目录（用于 pre_remove 失败时的恢复）
func (rm *RollbackManager) Backup(backupPath string) error {
	if rm.skillPath == "" {
		return fmt.Errorf("skill path is empty")
	}

	// 确保备份目录存在
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return err
	}

	// 复制整个目录
	return copyDirectory(rm.skillPath, backupPath)
}

// Restore 从备份恢复
func (rm *RollbackManager) Restore(backupPath string) error {
	if rm.skillPath == "" {
		return fmt.Errorf("skill path is empty")
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(rm.skillPath), 0755); err != nil {
		return err
	}

	// 从备份恢复
	return copyDirectory(backupPath, rm.skillPath)
}

// copyDirectory 复制目录
func copyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
