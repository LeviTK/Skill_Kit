package lib

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Config 全局配置
type Config struct {
	RepoPath         string              `toml:"-"`
	Platforms        map[string]Platform `toml:"platforms"`
	DefaultPlatforms []string            `toml:"default_platforms"` // 默认同步的平台列表
	PlatformOrder    []string            `toml:"platform_order"`    // 平台显示顺序
}

// Platform 平台配置
type Platform struct {
	Name     string `toml:"name"`
	Project  string `toml:"project"`
	Global   string `toml:"global"`
	SkillDir string `toml:"skill_dir"`
	AgentDir string `toml:"agent_dir"`
}

// GetCategoryDir 根据类别返回目录名
func (p Platform) GetCategoryDir(category string) string {
	if category == "agent" {
		return p.AgentDir
	}
	return p.SkillDir
}

// LoadConfig 加载配置
// 配置路径优先级: SKILLKIT_CONFIG 环境变量 > ~/.config/agent/platforms.toml > 可执行文件目录
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	repoPath := filepath.Join(home, ".config", "agent")

	// 支持环境变量覆盖配置路径
	if envRepo := os.Getenv("SKILLKIT_REPO"); envRepo != "" {
		repoPath = ResolvePath(envRepo)
	}

	configPath := filepath.Join(repoPath, "platforms.toml")

	// 支持环境变量直接指定配置文件
	if envConfig := os.Getenv("SKILLKIT_CONFIG"); envConfig != "" {
		configPath = ResolvePath(envConfig)
	} else if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果配置不存在，尝试从项目目录读取
		execPath, _ := os.Executable()
		execDir := filepath.Dir(execPath)
		configPath = filepath.Join(execDir, "..", "platforms.toml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.RepoPath = repoPath
	return &cfg, nil
}

// SaveConfig 保存配置
func SaveConfig(cfg *Config) error {
	configPath := filepath.Join(cfg.RepoPath, "platforms.toml")
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// GetOrderedPlatformKeys 获取有序的平台 key 列表
func (cfg *Config) GetOrderedPlatformKeys() []string {
	// 如果有保存的顺序，使用保存的顺序
	if len(cfg.PlatformOrder) > 0 {
		// 验证并过滤有效的 key
		var ordered []string
		seen := make(map[string]bool)
		for _, key := range cfg.PlatformOrder {
			if _, ok := cfg.Platforms[key]; ok && !seen[key] {
				ordered = append(ordered, key)
				seen[key] = true
			}
		}
		// 添加新增的平台（不在 order 中的）
		for key := range cfg.Platforms {
			if !seen[key] {
				ordered = append(ordered, key)
			}
		}
		return ordered
	}

	// 默认顺序
	var keys []string
	for key := range cfg.Platforms {
		keys = append(keys, key)
	}
	return keys
}
