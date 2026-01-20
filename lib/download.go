package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ParsedSource 解析后的源地址
type ParsedSource struct {
	Type      string // github, gitlab, git, local
	URL       string // Git URL
	Ref       string // 分支/标签
	Subpath   string // 仓库内子路径
	LocalPath string // 本地路径
}

// DiscoveredSkill 发现的技能
type DiscoveredSkill struct {
	Name        string
	Description string
	Path        string
	Category    string // skill 或 agent
	Hash        string
}

// ParseSource 解析源地址字符串
// 支持格式:
//   - 本地路径: ./path, ../path, /absolute/path
//   - GitHub URL: https://github.com/owner/repo
//   - GitHub URL with branch: https://github.com/owner/repo/tree/branch
//   - GitHub URL with path: https://github.com/owner/repo/tree/branch/path/to/skill
//   - GitHub shorthand: owner/repo, owner/repo/path/to/skill
//   - GitLab URL: https://gitlab.com/owner/repo
//   - Direct git URL: git@github.com:owner/repo.git
func ParseSource(input string) *ParsedSource {
	// 本地路径
	if isLocalPath(input) {
		absPath, _ := filepath.Abs(input)
		return &ParsedSource{
			Type:      "local",
			URL:       absPath,
			LocalPath: absPath,
		}
	}

	// GitHub URL with path: https://github.com/owner/repo/tree/branch/path
	githubTreeWithPath := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/tree/([^/]+)/(.+)`)
	if m := githubTreeWithPath.FindStringSubmatch(input); m != nil {
		return &ParsedSource{
			Type:    "github",
			URL:     fmt.Sprintf("https://github.com/%s/%s.git", m[1], m[2]),
			Ref:     m[3],
			Subpath: m[4],
		}
	}

	// GitHub URL with branch only: https://github.com/owner/repo/tree/branch
	githubTree := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/tree/([^/]+)$`)
	if m := githubTree.FindStringSubmatch(input); m != nil {
		return &ParsedSource{
			Type: "github",
			URL:  fmt.Sprintf("https://github.com/%s/%s.git", m[1], m[2]),
			Ref:  m[3],
		}
	}

	// GitHub URL: https://github.com/owner/repo
	githubRepo := regexp.MustCompile(`github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)
	if m := githubRepo.FindStringSubmatch(input); m != nil {
		repo := strings.TrimSuffix(m[2], ".git")
		return &ParsedSource{
			Type: "github",
			URL:  fmt.Sprintf("https://github.com/%s/%s.git", m[1], repo),
		}
	}

	// GitLab URL with path: https://gitlab.com/owner/repo/-/tree/branch/path
	gitlabTreeWithPath := regexp.MustCompile(`gitlab\.com/([^/]+)/([^/]+)/-/tree/([^/]+)/(.+)`)
	if m := gitlabTreeWithPath.FindStringSubmatch(input); m != nil {
		return &ParsedSource{
			Type:    "gitlab",
			URL:     fmt.Sprintf("https://gitlab.com/%s/%s.git", m[1], m[2]),
			Ref:     m[3],
			Subpath: m[4],
		}
	}

	// GitLab URL: https://gitlab.com/owner/repo
	gitlabRepo := regexp.MustCompile(`gitlab\.com/([^/]+)/([^/]+?)(?:\.git)?$`)
	if m := gitlabRepo.FindStringSubmatch(input); m != nil {
		repo := strings.TrimSuffix(m[2], ".git")
		return &ParsedSource{
			Type: "gitlab",
			URL:  fmt.Sprintf("https://gitlab.com/%s/%s.git", m[1], repo),
		}
	}

	// GitHub shorthand: owner/repo 或 owner/repo/path/to/skill
	shorthand := regexp.MustCompile(`^([^/:]+)/([^/:]+)(?:/(.+))?$`)
	if m := shorthand.FindStringSubmatch(input); m != nil && !strings.Contains(input, ":") {
		parsed := &ParsedSource{
			Type: "github",
			URL:  fmt.Sprintf("https://github.com/%s/%s.git", m[1], m[2]),
		}
		if m[3] != "" {
			parsed.Subpath = m[3]
		}
		return parsed
	}

	// 直接 git URL
	return &ParsedSource{
		Type: "git",
		URL:  input,
	}
}

func isLocalPath(input string) bool {
	if filepath.IsAbs(input) {
		return true
	}
	if strings.HasPrefix(input, "./") || strings.HasPrefix(input, "../") {
		return true
	}
	if input == "." || input == ".." {
		return true
	}
	return false
}

// CloneRepo 克隆仓库到临时目录
func CloneRepo(url string, ref string) (string, error) {
	tempDir, err := os.MkdirTemp("", "skillkit-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	args := []string{"clone", "--depth", "1"}
	if ref != "" {
		args = append(args, "--branch", ref)
	}
	args = append(args, url, tempDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil
}

// CleanupTempDir 清理临时目录
func CleanupTempDir(dir string) error {
	// 安全检查：确保是临时目录
	tempDir := os.TempDir()
	if !strings.HasPrefix(dir, tempDir) {
		return fmt.Errorf("refusing to delete non-temp directory: %s", dir)
	}
	return os.RemoveAll(dir)
}

// DiscoverSkills 在目录中发现技能
func DiscoverSkills(basePath string, subpath string) ([]*DiscoveredSkill, error) {
	searchPath := basePath
	if subpath != "" {
		searchPath = filepath.Join(basePath, subpath)
	}

	var skills []*DiscoveredSkill
	seenSkills := make(map[string]bool)

	// 检查是否直接指向一个 skill
	if hasSkillFile(searchPath) {
		skill := parseSkillFile(searchPath)
		if skill != nil {
			skills = append(skills, skill)
			return skills, nil
		}
	}

	// 搜索常见位置
	priorityDirs := []string{
		searchPath,
		filepath.Join(searchPath, "skills"),
		filepath.Join(searchPath, "skill"),
		filepath.Join(searchPath, "agent"),
		filepath.Join(searchPath, "agents"),
		filepath.Join(searchPath, ".claude/skills"),
		filepath.Join(searchPath, ".cursor/skills"),
		filepath.Join(searchPath, ".codex/skills"),
	}

	for _, dir := range priorityDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillDir := filepath.Join(dir, entry.Name())
			if hasSkillFile(skillDir) {
				skill := parseSkillFile(skillDir)
				if skill != nil {
					key := skillKey(skill)
					if !seenSkills[key] {
						skills = append(skills, skill)
						seenSkills[key] = true
					}
				}
			}
		}
	}

	// 如果没找到，递归搜索
	if len(skills) == 0 {
		skills = findSkillsRecursive(searchPath, seenSkills, 0, 5)
	}

	return skills, nil
}

func hasSkillFile(dir string) bool {
	// 检查 SKILL.md 或 AGENT.md
	for _, name := range []string{"SKILL.md", "AGENT.md"} {
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

func parseSkillFile(dir string) *DiscoveredSkill {
	var filePath string
	var category string

	if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err == nil {
		filePath = filepath.Join(dir, "SKILL.md")
		category = "skill"
	} else if _, err := os.Stat(filepath.Join(dir, "AGENT.md")); err == nil {
		filePath = filepath.Join(dir, "AGENT.md")
		category = "agent"
	} else {
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	hash := hashContent(content)

	// 解析 frontmatter
	name, description := parseFrontmatter(string(content))
	if name == "" {
		// 使用目录名作为名称
		name = filepath.Base(dir)
	}

	return &DiscoveredSkill{
		Name:        name,
		Description: description,
		Path:        dir,
		Category:    category,
		Hash:        hash,
	}
}

func parseFrontmatter(content string) (name, description string) {
	lines := strings.Split(content, "\n")
	inFrontmatter := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter {
			if strings.HasPrefix(line, "name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				name = strings.Trim(name, `"'`)
			} else if strings.HasPrefix(line, "description:") {
				description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
				description = strings.Trim(description, `"'`)
			}
		}
	}
	return
}

func findSkillsRecursive(dir string, seen map[string]bool, depth, maxDepth int) []*DiscoveredSkill {
	if depth > maxDepth {
		return nil
	}

	var skills []*DiscoveredSkill
	skipDirs := map[string]bool{
		"node_modules": true, ".git": true, "dist": true, "build": true,
	}

	if hasSkillFile(dir) {
		skill := parseSkillFile(dir)
		if skill != nil {
			key := skillKey(skill)
			if !seen[key] {
				skills = append(skills, skill)
				seen[key] = true
			}
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return skills
	}

	for _, entry := range entries {
		if !entry.IsDir() || skipDirs[entry.Name()] {
			continue
		}
		subSkills := findSkillsRecursive(filepath.Join(dir, entry.Name()), seen, depth+1, maxDepth)
		skills = append(skills, subSkills...)
	}

	return skills
}

func skillKey(skill *DiscoveredSkill) string {
	hash := skill.Hash
	if hash == "" {
		hash = "nohash"
	}
	// TODO: Add directory-content hashing for dedupe (backlog).
	return skill.Category + ":" + skill.Name + ":" + hash
}

func hashContent(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// InstallSkill 安装技能到本地仓库
func InstallSkill(skill *DiscoveredSkill, cfg *Config) error {
	// 目标目录
	targetBase := filepath.Join(cfg.RepoPath, skill.Category)
	targetDir := filepath.Join(targetBase, skill.Name)

	// 检查是否已存在
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("skill '%s' already exists at %s", skill.Name, targetDir)
	}

	// 创建目标目录
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 复制文件
	if err := copyDir(skill.Path, targetDir); err != nil {
		os.RemoveAll(targetDir)
		return fmt.Errorf("failed to copy skill: %w", err)
	}

	return nil
}

func copyDir(src, dst string) error {
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
			if err := copyDir(srcPath, dstPath); err != nil {
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

// SelectSkillsInteractive 交互式选择技能
func SelectSkillsInteractive(skills []*DiscoveredSkill) []*DiscoveredSkill {
	if len(skills) == 0 {
		return nil
	}

	if len(skills) == 1 {
		fmt.Printf("\n%s Found 1 skill: %s\n", Blue(IconInfo), skills[0].Name)
		if skills[0].Description != "" {
			fmt.Printf("  %s\n", Gray(skills[0].Description))
		}
		fmt.Print("\nInstall? [Y/n] ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "" || input == "y" || input == "yes" {
			return skills
		}
		return nil
	}

	// 多个技能，使用多选菜单
	fmt.Printf("\n%s Found %d skills:\n\n", Blue(IconInfo), len(skills))

	selected := make([]bool, len(skills))
	for i := range selected {
		selected[i] = true // 默认全选
	}

	cursor := 0
	for {
		ClearScreen()
		fmt.Printf("\n%s Select skills to install (Space=toggle, Enter=confirm, A=all, Q=quit):\n\n", Blue(IconInfo))

		for i, skill := range skills {
			marker := "[ ]"
			if selected[i] {
				marker = Green("[✓]")
			}
			prefix := "  "
			if i == cursor {
				prefix = Cyan("> ")
			}
			fmt.Printf("%s%s %s", prefix, marker, skill.Name)
			if skill.Description != "" {
				fmt.Printf(" %s", Gray("- "+truncate(skill.Description, 50)))
			}
			fmt.Println()
		}

		key := ReadKey()
		switch key {
		case "UP":
			if cursor > 0 {
				cursor--
			}
		case "DOWN":
			if cursor < len(skills)-1 {
				cursor++
			}
		case "SPACE":
			selected[cursor] = !selected[cursor]
		case "SELECTALL":
			allSelected := true
			for _, s := range selected {
				if !s {
					allSelected = false
					break
				}
			}
			for i := range selected {
				selected[i] = !allSelected
			}
		case "ENTER":
			var result []*DiscoveredSkill
			for i, skill := range skills {
				if selected[i] {
					result = append(result, skill)
				}
			}
			return result
		case "QUIT", "ESC":
			return nil
		}
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
