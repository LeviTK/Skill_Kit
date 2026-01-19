package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModuleGetLinkName(t *testing.T) {
	mod := &Module{
		Name:     "python-coder",
		Category: "skill",
		Path:     "/path/to/module",
		Aliases: map[string]string{
			"claude": "py-expert",
			"cursor": "py-dev",
		},
	}

	tests := []struct {
		platform string
		expected string
	}{
		{"claude", "py-expert"},
		{"cursor", "py-dev"},
		{"copilot", "python-coder"},
	}

	for _, tt := range tests {
		result := mod.GetLinkName(tt.platform)
		if result != tt.expected {
			t.Errorf("GetLinkName(%s) = %s, expected %s", tt.platform, result, tt.expected)
		}
	}
}

func TestFindModule(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "skill", "test-skill")
	agentDir := filepath.Join(tmpDir, "agent", "test-agent")

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatalf("failed to create agent dir: %v", err)
	}

	cfg := &Config{RepoPath: tmpDir}

	skill, err := FindModule(cfg, "test-skill")
	if err != nil {
		t.Fatalf("FindModule failed for skill: %v", err)
	}
	if skill.Category != "skill" {
		t.Errorf("expected category 'skill', got '%s'", skill.Category)
	}

	agent, err := FindModule(cfg, "test-agent")
	if err != nil {
		t.Fatalf("FindModule failed for agent: %v", err)
	}
	if agent.Category != "agent" {
		t.Errorf("expected category 'agent', got '%s'", agent.Category)
	}

	_, err = FindModule(cfg, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent module")
	}
	if !IsModuleNotFound(err) {
		t.Errorf("expected ModuleNotFoundError, got %T", err)
	}
}

func TestListModules(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "skill")
	agentDir := filepath.Join(tmpDir, "agent")

	os.MkdirAll(filepath.Join(skillDir, "skill-a"), 0755)
	os.MkdirAll(filepath.Join(skillDir, "skill-b"), 0755)
	os.MkdirAll(filepath.Join(agentDir, "agent-a"), 0755)

	cfg := &Config{RepoPath: tmpDir}

	modules, err := ListModules(cfg)
	if err != nil {
		t.Fatalf("ListModules failed: %v", err)
	}

	if len(modules) != 3 {
		t.Errorf("expected 3 modules, got %d", len(modules))
	}

	skillCount := 0
	agentCount := 0
	for _, m := range modules {
		if m.Category == "skill" {
			skillCount++
		} else if m.Category == "agent" {
			agentCount++
		}
	}

	if skillCount != 2 {
		t.Errorf("expected 2 skills, got %d", skillCount)
	}
	if agentCount != 1 {
		t.Errorf("expected 1 agent, got %d", agentCount)
	}
}

func TestLoadModuleWithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	modDir := filepath.Join(tmpDir, "skill", "custom-mod")
	os.MkdirAll(modDir, 0755)

	configContent := `[link]
default = "custom-name"

[link.overrides]
claude = "claude-alias"
cursor = "cursor-alias"
`
	configPath := filepath.Join(modDir, "skillkit.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg := &Config{RepoPath: tmpDir}
	mod, err := FindModule(cfg, "custom-mod")
	if err != nil {
		t.Fatalf("FindModule failed: %v", err)
	}

	if mod.Name != "custom-name" {
		t.Errorf("expected name 'custom-name', got '%s'", mod.Name)
	}

	if mod.GetLinkName("claude") != "claude-alias" {
		t.Errorf("expected claude alias 'claude-alias', got '%s'", mod.GetLinkName("claude"))
	}

	if mod.GetLinkName("copilot") != "custom-name" {
		t.Errorf("expected default name for copilot, got '%s'", mod.GetLinkName("copilot"))
	}
}
