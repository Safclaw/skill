package hook

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Executor Hook 执行器
type Executor struct {
	platform    *Platform
	checksummer *ChecksumVerifier
}

// NewExecutor 创建执行器
func NewExecutor() *Executor {
	return &Executor{
		platform:    NewPlatform(),
		checksummer: NewChecksumVerifier(),
	}
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	Script   string
	Output   string
	Duration time.Duration
	Error    error
}

// Execute 执行 Hook
func (e *Executor) Execute(stage HookStage, skillPath string, hooks []HookConfig) ([]ExecuteResult, error) {
	results := make([]ExecuteResult, 0)

	for _, hook := range hooks {
		if hook.Stage != stage {
			continue
		}

		// 选择适合当前平台的脚本
		script := e.platform.SelectScript(hook.Scripts)
		if script == nil {
			// 没有适合当前平台的脚本，跳过
			continue
		}

		// 验证脚本路径
		if len(script.Args) == 0 {
			return results, fmt.Errorf("script args is empty")
		}

		scriptPath := script.Args[0]
		if err := validateScriptPath(scriptPath); err != nil {
			return results, fmt.Errorf("invalid script path: %w", err)
		}

		// 验证 checksum
		if script.Checksum != "" {
			fullPath := filepath.Join(skillPath, scriptPath)
			if err := e.checksummer.Verify(fullPath, script.Checksum); err != nil {
				return results, fmt.Errorf("checksum verification failed: %w", err)
			}
		}

		// 执行脚本
		result := e.executeScript(script, skillPath)
		results = append(results, result)

		// 如果执行失败，返回错误
		if result.Error != nil {
			return results, result.Error
		}
	}

	return results, nil
}

// executeScript 执行单个脚本
func (e *Executor) executeScript(script *HookScript, skillPath string) ExecuteResult {
	result := ExecuteResult{
		Script: script.Args[0],
	}

	startTime := time.Now()
	defer func() {
		result.Duration = time.Since(startTime)
	}()

	// 确定工作目录
	workDir := script.WorkingDir
	if workDir == "" {
		workDir = "."
	}
	fullWorkDir := filepath.Join(skillPath, workDir)

	// 准备命令
	cmd := exec.Command(script.Command, script.Args...)
	cmd.Dir = fullWorkDir

	// 设置环境变量
	env := os.Environ()
	for k, v := range script.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// 执行并捕获输出
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Error = err

	if err != nil {
		// 检查是否是超时错误
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ProcessState != nil && exitErr.ProcessState.Exited() {
				result.Error = fmt.Errorf("script exited with code %d: %s",
					exitErr.ProcessState.ExitCode(), result.Output)
			}
		}
	}

	return result
}

// ExecuteWithTimeout 带超时控制的执行
func (e *Executor) ExecuteWithTimeout(stage HookStage, skillPath string, hooks []HookConfig, timeout time.Duration) ([]ExecuteResult, error) {
	// 为每个 hook 设置独立的超时
	for _, hook := range hooks {
		if hook.Stage != stage {
			continue
		}

		hookTimeout := time.Duration(hook.Timeout) * time.Second
		if hookTimeout == 0 {
			hookTimeout = 30 * time.Second // 默认 30 秒
		}

		// 使用较小的超时时间
		if hookTimeout > timeout {
			hookTimeout = timeout
		}

		// 创建带超时的 context
		ctx, cancel := context.WithTimeout(context.Background(), hookTimeout)
		defer cancel()

		// 选择脚本
		script := e.platform.SelectScript(hook.Scripts)
		if script == nil {
			continue
		}

		// 执行脚本（带 context）
		result := e.executeScriptWithContext(ctx, script, skillPath)
		if result.Error != nil {
			return []ExecuteResult{result}, result.Error
		}
	}

	return nil, nil
}

// executeScriptWithContext 带 context 执行脚本
func (e *Executor) executeScriptWithContext(ctx context.Context, script *HookScript, skillPath string) ExecuteResult {
	result := ExecuteResult{
		Script: script.Args[0],
	}

	startTime := time.Now()
	defer func() {
		result.Duration = time.Since(startTime)
	}()

	// 确定工作目录
	workDir := script.WorkingDir
	if workDir == "" {
		workDir = "."
	}
	fullWorkDir := filepath.Join(skillPath, workDir)

	// 准备命令
	cmd := exec.CommandContext(ctx, script.Command, script.Args...)
	cmd.Dir = fullWorkDir

	// 设置环境变量
	env := os.Environ()
	for k, v := range script.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// 执行并捕获输出
	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("script timeout")
		} else {
			result.Error = err
		}
	}

	return result
}
