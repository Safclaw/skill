package manifest

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// SkillsManifest 技能清单
type SkillsManifest struct {
	Skills []SkillEntry `yaml:"skills"`
}

// SkillEntry 技能条目
type SkillEntry struct {
	Name    string `yaml:"name"`    // 可读名称
	Dir     string `yaml:"dir"`     // 相对路径
	Version string `yaml:"version"` // 版本号
	Sig     string `yaml:"sig"`     // SHA256 签名
}

// ReadManifest 读取清单文件
func ReadManifest(path string) (*SkillsManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回空清单
			return &SkillsManifest{Skills: []SkillEntry{}}, nil
		}
		return nil, err
	}

	var manifest SkillsManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// WriteManifest 写入清单文件
func WriteManifest(path string, manifest *SkillsManifest) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 序列化
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	// 添加 YAML 头部注释
	content := "# skill reps\n---\n" + string(data)

	// 写入文件
	return os.WriteFile(path, []byte(content), 0644)
}

// AddSkill 添加技能到清单
func (m *SkillsManifest) AddSkill(entry SkillEntry) {
	// 检查是否已存在
	for i, e := range m.Skills {
		if e.Dir == entry.Dir {
			// 更新现有条目
			m.Skills[i] = entry
			return
		}
	}

	// 添加新条目
	m.Skills = append(m.Skills, entry)
}

// RemoveSkill 从清单删除技能
func (m *SkillsManifest) RemoveSkill(dir string) bool {
	for i, e := range m.Skills {
		if e.Dir == dir {
			m.Skills = append(m.Skills[:i], m.Skills[i+1:]...)
			return true
		}
	}
	return false
}

// GetSkill 获取技能信息
func (m *SkillsManifest) GetSkill(dir string) *SkillEntry {
	for _, e := range m.Skills {
		if e.Dir == dir {
			return &e
		}
	}
	return nil
}

// HasSkill 检查技能是否存在
func (m *SkillsManifest) HasSkill(dir string) bool {
	return m.GetSkill(dir) != nil
}

// ListSkills 列出所有技能
func (m *SkillsManifest) ListSkills() []SkillEntry {
	return m.Skills
}
