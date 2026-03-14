package lister

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Safclaw/skill/pkg/manifest"
)

// Lister 列表查询器
type Lister struct {
	installDir string
}

// NewLister 创建列表查询器
func NewLister(installDir string) *Lister {
	return &Lister{
		installDir: installDir,
	}
}

// SkillInfo 技能信息
type SkillInfo struct {
	Name      string
	ModuleDir string
	Version   string
	Signature string
	Path      string
}

// ListSkills 列出所有已安装的 skill
func (l *Lister) ListSkills() ([]SkillInfo, error) {
	manifestPath := filepath.Join(l.installDir, ".skills.yaml")

	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	infos := make([]SkillInfo, 0, len(m.Skills))
	for _, entry := range m.Skills {
		info := SkillInfo{
			Name:      entry.Name,
			ModuleDir: entry.Dir,
			Version:   entry.Version,
			Signature: entry.Sig,
			Path:      filepath.Join(l.installDir, "reps", entry.Dir),
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// GetSkillInfo 获取单个 skill 的详细信息
func (l *Lister) GetSkillInfo(moduleDir string) (*SkillInfo, error) {
	manifestPath := filepath.Join(l.installDir, ".skills.yaml")

	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	entry := m.GetSkill(moduleDir)
	if entry == nil {
		return nil, fmt.Errorf("skill not found: %s", moduleDir)
	}

	return &SkillInfo{
		Name:      entry.Name,
		ModuleDir: entry.Dir,
		Version:   entry.Version,
		Signature: entry.Sig,
		Path:      filepath.Join(l.installDir, "reps", entry.Dir),
	}, nil
}

// IsInstalled 检查 skill 是否已安装
func (l *Lister) IsInstalled(moduleDir string) (bool, error) {
	manifestPath := filepath.Join(l.installDir, ".skills.yaml")

	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return false, err
	}

	return m.HasSkill(moduleDir), nil
}

// Count 返回已安装的 skill 数量
func (l *Lister) Count() (int, error) {
	manifestPath := filepath.Join(l.installDir, ".skills.yaml")

	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return 0, err
	}

	return len(m.Skills), nil
}

// PrintList 打印技能列表（用于 CLI）
func (l *Lister) PrintList() error {
	infos, err := l.ListSkills()
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		fmt.Println("No skills installed")
		return nil
	}

	fmt.Printf("Installed skills in %s:\n\n", l.installDir)
	for _, info := range infos {
		fmt.Printf("  %s@%s\n", info.ModuleDir, info.Version)
		fmt.Printf("    Path: %s\n", info.Path)
		if info.Signature != "" {
			fmt.Printf("    Signature: %s\n", info.Signature[:16]+"...")
		}
		fmt.Println()
	}

	return nil
}

// VerifyInstallations 验证所有安装的完整性
func (l *Lister) VerifyInstallations() error {
	infos, err := l.ListSkills()
	if err != nil {
		return err
	}

	errors := make([]error, 0)
	for _, info := range infos {
		// 检查目录是否存在
		if _, err := os.Stat(info.Path); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("%s: directory not found", info.ModuleDir))
			continue
		}

		// TODO: 验证签名
		// if info.Signature != "" {
		//     if err := signature.VerifyDirectory(info.Path, info.Signature); err != nil {
		//         errors = append(errors, fmt.Errorf("%s: %w", info.ModuleDir, err))
		//     }
		// }
	}

	if len(errors) > 0 {
		fmt.Println("Verification failed:")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("verification failed with %d errors", len(errors))
	}

	fmt.Println("All installations verified successfully")
	return nil
}
