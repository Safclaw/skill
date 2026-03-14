package hook

// HookStage Hook 阶段
type HookStage string

const (
	PostAdd   HookStage = "post_add"   // 安装后执行
	PreRemove HookStage = "pre_remove" // 卸载前执行
)

// HookConfig Hook 配置
type HookConfig struct {
	Stage   HookStage    `yaml:"stage"`
	Reason  string       `yaml:"reason,omitempty"`
	Timeout int          `yaml:"timeout,omitempty"` // 秒，默认 30 秒
	Scripts []HookScript `yaml:"scripts"`
}

// HookScript Hook 脚本配置
type HookScript struct {
	Command    string            `yaml:"command"`            // bash, powershell, python 等
	Platforms  []string          `yaml:"platforms"`          // linux, darwin, windows
	Args       []string          `yaml:"args"`               // 脚本参数
	WorkingDir string            `yaml:"working_dir"`        // 工作目录
	Checksum   string            `yaml:"checksum,omitempty"` // SHA256 校验和
	Env        map[string]string `yaml:"env,omitempty"`      // 环境变量
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}
