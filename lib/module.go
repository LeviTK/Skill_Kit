package lib

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Module 模块信息
type Module struct {
	Name        string
	Category    string // "skill" or "agent"
	Path        string
	Aliases     map[string]string // platform -> link_name
	Description string            // 从 SKILL.md/AGENT.md 读取的描述
}

// ModuleConfig 模块配置文件 (skillkit.toml)
type ModuleConfig struct {
	Link struct {
		Default   string            `toml:"default"`
		Overrides map[string]string `toml:"overrides"`
	} `toml:"link"`
}

// GetLinkName 获取模块在指定平台的链接名
func (m Module) GetLinkName(platform string) string {
	if alias, ok := m.Aliases[platform]; ok {
		return alias
	}
	return m.Name
}

// FindModule 查找模块
func FindModule(cfg *Config, name string) (*Module, error) {
	// 先在 skill 目录查找
	skillPath := filepath.Join(cfg.RepoPath, "skill", name)
	if info, err := os.Stat(skillPath); err == nil && info.IsDir() {
		return loadModule(name, "skill", skillPath)
	}

	// 再在 agent 目录查找
	agentPath := filepath.Join(cfg.RepoPath, "agent", name)
	if info, err := os.Stat(agentPath); err == nil && info.IsDir() {
		return loadModule(name, "agent", agentPath)
	}

	return nil, &ModuleNotFoundError{Name: name}
}

// ListModules 列出所有模块
func ListModules(cfg *Config) ([]*Module, error) {
	var modules []*Module

	// 遍历 skill 目录
	skillDir := filepath.Join(cfg.RepoPath, "skill")
	if entries, err := os.ReadDir(skillDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				mod, err := loadModule(entry.Name(), "skill", filepath.Join(skillDir, entry.Name()))
				if err == nil {
					modules = append(modules, mod)
				}
			}
		}
	}

	// 遍历 agent 目录
	agentDir := filepath.Join(cfg.RepoPath, "agent")
	if entries, err := os.ReadDir(agentDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				mod, err := loadModule(entry.Name(), "agent", filepath.Join(agentDir, entry.Name()))
				if err == nil {
					modules = append(modules, mod)
				}
			}
		}
	}

	return modules, nil
}

// loadModule 加载模块信息
func loadModule(name, category, path string) (*Module, error) {
	mod := &Module{
		Name:     name,
		Category: category,
		Path:     path,
		Aliases:  make(map[string]string),
	}

	// 尝试读取 skillkit.toml
	configPath := filepath.Join(path, "skillkit.toml")
	if data, err := os.ReadFile(configPath); err == nil {
		var modCfg ModuleConfig
		if err := toml.Unmarshal(data, &modCfg); err == nil {
			if modCfg.Link.Default != "" {
				mod.Name = modCfg.Link.Default
			}
			for platform, alias := range modCfg.Link.Overrides {
				mod.Aliases[platform] = alias
			}
		}
	}

	// 读取描述信息（从 SKILL.md 或 AGENT.md）
	mod.Description = loadModuleDescription(path, category)

	return mod, nil
}

// loadModuleDescription 从 SKILL.md/AGENT.md 读取描述
func loadModuleDescription(path, category string) string {
	// 尝试读取 SKILL.md 或 AGENT.md
	var mdPath string
	if category == "skill" {
		mdPath = filepath.Join(path, "SKILL.md")
	} else {
		mdPath = filepath.Join(path, "AGENT.md")
	}

	data, err := os.ReadFile(mdPath)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := splitLines(content)

	// 查找 description 字段（在 frontmatter 中）或第一段描述
	inFrontmatter := false
	for _, line := range lines {
		if line == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			// 查找 description: 字段
			if len(line) > 12 && line[:12] == "description:" {
				desc := line[12:]
				// 去除引号和空格
				desc = trimQuotes(desc)
				return desc // 返回完整描述，不截断
			}
		}
	}

	// 如果没有 frontmatter，返回第一行非空非标题内容
	for _, line := range lines {
		if line == "" || line[0] == '#' || line == "---" {
			continue
		}
		return line // 返回完整描述，不截断
	}

	return ""
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimQuotes(s string) string {
	s = trimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// GetSyncedPlatformKeys 获取已同步的平台 key 列表
func GetSyncedPlatformKeys(cfg *Config, mod *Module) []string {
	var keys []string
	for key, p := range cfg.Platforms {
		targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := targetDir + "/" + mod.GetLinkName(key)
		if IsSymlink(targetPath) {
			keys = append(keys, key)
		}
	}
	return keys
}

// GetLinkStatus 获取模块的链接状态
func GetLinkStatus(cfg *Config, mod *Module) []string {
	var status []string

	for name, p := range cfg.Platforms {
		ln := mod.GetLinkName(name)
		targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := filepath.Join(targetDir, ln)

		if IsSymlink(targetPath) {
			realPath, _ := ReadSymlink(targetPath)
			if realPath == mod.Path {
				status = append(status, fmt.Sprintf("%s ✓", name))
			} else {
				status = append(status, fmt.Sprintf("%s ✗ (broken)", name))
			}
		}
	}

	return status
}
