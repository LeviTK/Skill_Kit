package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CollectedSkill 从平台目录收集到的技能
type CollectedSkill struct {
	Name     string
	Category string // skill / agent
	Path     string // 原始完整路径
	Platform string // 来源平台 key
	PlatName string // 来源平台显示名
	Hash     string // SKILL.md/AGENT.md 的 SHA256
	Desc     string // 描述
	IsLink   bool   // 是否已是软链接
	LinkDest string // 软链接目标（仅 IsLink=true 时有效）
}

// CollectGroup 按 category+name 分组后的结果
type CollectGroup struct {
	Name     string
	Category string
	Skills   []*CollectedSkill
}

// CollectResult 收集流程的最终结果
type CollectResult struct {
	Moved   int
	Linked  int
	Skipped int
	Fixed   int
	Errors  int
}

// ScanPlatformSkills 扫描所有平台的 global skill/agent 目录
func ScanPlatformSkills(cfg *Config) []*CollectedSkill {
	var collected []*CollectedSkill

	for key, p := range cfg.Platforms {
		for _, category := range []string{"skill", "agent"} {
			catDir := p.GetCategoryDir(category)
			dir := ResolvePath(p.Global, catDir)

			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if !entry.IsDir() && entry.Type()&os.ModeSymlink == 0 {
					continue
				}

				fullPath := filepath.Join(dir, entry.Name())
				info, err := os.Lstat(fullPath)
				if err != nil {
					continue
				}

				cs := &CollectedSkill{
					Name:     entry.Name(),
					Category: category,
					Path:     fullPath,
					Platform: key,
					PlatName: p.Name,
				}

				if info.Mode()&os.ModeSymlink != 0 {
					cs.IsLink = true
					cs.LinkDest, _ = os.Readlink(fullPath)
				}

				// 读取 SKILL.md/AGENT.md 计算 hash 和描述
				mdName := "SKILL.md"
				if category == "agent" {
					mdName = "AGENT.md"
				}

				// 对软链接，解析实际路径后读取
				readPath := fullPath
				if cs.IsLink {
					readPath = cs.LinkDest
				}

				mdPath := filepath.Join(readPath, mdName)
				if data, err := os.ReadFile(mdPath); err == nil {
					sum := sha256.Sum256(data)
					cs.Hash = hex.EncodeToString(sum[:])
					_, cs.Desc = parseFrontmatter(string(data))
				}

				collected = append(collected, cs)
			}
		}
	}

	return collected
}

// GroupByName 按 category+name 分组
func GroupByName(skills []*CollectedSkill) []*CollectGroup {
	groupMap := make(map[string]*CollectGroup)

	for _, s := range skills {
		key := s.Category + ":" + s.Name
		if g, ok := groupMap[key]; ok {
			g.Skills = append(g.Skills, s)
		} else {
			groupMap[key] = &CollectGroup{
				Name:     s.Name,
				Category: s.Category,
				Skills:   []*CollectedSkill{s},
			}
		}
	}

	var groups []*CollectGroup
	for _, g := range groupMap {
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Category != groups[j].Category {
			return groups[i].Category < groups[j].Category
		}
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// HasConflict 判断分组是否存在内容冲突（多个不同 hash）
func (g *CollectGroup) HasConflict() bool {
	hashes := make(map[string]bool)
	for _, s := range g.Skills {
		if !s.IsLink && s.Hash != "" {
			hashes[s.Hash] = true
		}
	}
	return len(hashes) > 1
}

// UniqueHashes 返回去重后的 hash 列表
func (g *CollectGroup) UniqueHashes() []string {
	seen := make(map[string]bool)
	var hashes []string
	for _, s := range g.Skills {
		if s.Hash != "" && !seen[s.Hash] {
			hashes = append(hashes, s.Hash)
			seen[s.Hash] = true
		}
	}
	return hashes
}

// RealSkills 返回非软链接的实体 skill
func (g *CollectGroup) RealSkills() []*CollectedSkill {
	var real []*CollectedSkill
	for _, s := range g.Skills {
		if !s.IsLink {
			real = append(real, s)
		}
	}
	return real
}

// MoveToRepo 将选中的 skill 移动到中央仓库
func MoveToRepo(skill *CollectedSkill, repoPath string) (string, error) {
	targetDir := filepath.Join(repoPath, skill.Category, skill.Name)

	// 如果目标已存在，检查是否相同
	if info, err := os.Stat(targetDir); err == nil {
		if info.IsDir() {
			return targetDir, fmt.Errorf("already exists in repo")
		}
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent dir: %w", err)
	}

	// 移动目录
	if err := os.Rename(skill.Path, targetDir); err != nil {
		// 跨文件系统时 Rename 会失败，回退到复制+删除
		if err2 := copyDir(skill.Path, targetDir); err2 != nil {
			return "", fmt.Errorf("move failed: %w", err2)
		}
		// 复制成功后删除源
		os.RemoveAll(skill.Path)
	}

	return targetDir, nil
}

// ReplaceWithSymlink 将平台目录中的实体/错误链接替换为指向中央仓库的软链接
func ReplaceWithSymlink(platformPath string, repoTargetDir string) error {
	info, err := os.Lstat(platformPath)
	if err != nil {
		// 路径不存在，直接创建软链接
		if err := os.MkdirAll(filepath.Dir(platformPath), 0755); err != nil {
			return err
		}
		return os.Symlink(repoTargetDir, platformPath)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		// 已是软链接，检查指向
		dest, _ := os.Readlink(platformPath)
		if dest == repoTargetDir {
			return nil // 已正确
		}
		// 指向错误，删除重建
		os.Remove(platformPath)
		return os.Symlink(repoTargetDir, platformPath)
	}

	if info.IsDir() {
		// 实体目录（已被移动走的情况，或残留），删除后创建链接
		os.RemoveAll(platformPath)
		return os.Symlink(repoTargetDir, platformPath)
	}

	return fmt.Errorf("unexpected file at %s", platformPath)
}

// RunCollect 执行完整的收集流程
func RunCollect(cfg *Config) *CollectResult {
	result := &CollectResult{}

	fmt.Printf("\n%s Scanning platform directories...\n\n", Blue(IconInfo))

	collected := ScanPlatformSkills(cfg)
	if len(collected) == 0 {
		fmt.Printf("  %s No skills found in platform directories.\n\n", Yellow(IconWarning))
		return result
	}

	// 统计
	realCount := 0
	linkCount := 0
	for _, s := range collected {
		if s.IsLink {
			linkCount++
		} else {
			realCount++
		}
	}
	fmt.Printf("  Found %d skill(s): %d real, %d symlinks\n\n", len(collected), realCount, linkCount)

	groups := GroupByName(collected)

	// 预览摘要
	conflictCount := 0
	for _, g := range groups {
		if g.HasConflict() {
			conflictCount++
		}
	}

	fmt.Printf("  %s %d unique skill(s), %d with conflicts\n", Blue(IconInfo), len(groups), conflictCount)
	if conflictCount > 0 {
		fmt.Printf("  %s Conflicts will be resolved interactively.\n", Yellow(IconWarning))
	}
	fmt.Println()

	if !ConfirmDialog(fmt.Sprintf("Proceed with collecting %d skill(s)?", len(groups))) {
		fmt.Printf("\n  %s Cancelled.\n\n", Yellow(IconWarning))
		return result
	}

	fmt.Println()

	for _, g := range groups {
		realSkills := g.RealSkills()

		if len(realSkills) == 0 {
			// 全是软链接，检查指向是否正确
			repoTarget := filepath.Join(cfg.RepoPath, g.Category, g.Name)
			for _, s := range g.Skills {
				if s.LinkDest != repoTarget {
					err := ReplaceWithSymlink(s.Path, repoTarget)
					if err != nil {
						fmt.Printf("  %s Fix link %s → %s: %v\n", Red(IconError), s.PlatName, g.Name, err)
						result.Errors++
					} else {
						fmt.Printf("  %s Fixed link: %s/%s\n", Green(IconSuccess), s.PlatName, g.Name)
						result.Fixed++
					}
				} else {
					result.Skipped++
				}
			}
			continue
		}

		// 确定要保留的 skill
		var chosen *CollectedSkill
		if len(realSkills) == 1 {
			chosen = realSkills[0]
		} else if !g.HasConflict() {
			// 多个来源但内容相同，取第一个
			chosen = realSkills[0]
			fmt.Printf("  %s %s/%s: %d copies (identical), using from %s\n",
				Blue(IconInfo), g.Category, g.Name, len(realSkills), chosen.PlatName)
		} else {
			// 内容冲突，交互选择
			chosen = ConflictResolveMenu(g)
			if chosen == nil {
				fmt.Printf("  %s Skipped %s/%s\n", Yellow(IconWarning), g.Category, g.Name)
				result.Skipped++
				continue
			}
		}

		// 检查中央仓库是否已存在
		repoTarget := filepath.Join(cfg.RepoPath, g.Category, g.Name)
		if _, err := os.Stat(repoTarget); err == nil {
			// 已存在于仓库，检查内容
			existingHash := hashSkillMd(repoTarget, g.Category)
			if existingHash == chosen.Hash {
				fmt.Printf("  %s %s/%s already in repo (identical)\n", Gray("○"), g.Category, g.Name)
			} else if existingHash != "" && chosen.Hash != "" {
				fmt.Printf("  %s %s/%s exists in repo with different content, skipping\n",
					Yellow(IconWarning), g.Category, g.Name)
				result.Skipped++
				// 仍然修复软链接
				repoTarget = filepath.Join(cfg.RepoPath, g.Category, g.Name)
			}
		} else {
			// 移动到仓库
			dest, err := MoveToRepo(chosen, cfg.RepoPath)
			if err != nil {
				fmt.Printf("  %s Move %s/%s: %v\n", Red(IconError), g.Category, g.Name, err)
				result.Errors++
				continue
			}
			repoTarget = dest
			fmt.Printf("  %s %s/%s → repo (from %s)\n", Green(IconSuccess), g.Category, g.Name, chosen.PlatName)
			result.Moved++
		}

		// 替换所有出现位置为软链接
		for _, s := range g.Skills {
			if s.Path == chosen.Path && !s.IsLink {
				// 已被移动的那个，路径已不存在，创建链接
				err := ReplaceWithSymlink(s.Path, repoTarget)
				if err != nil {
					fmt.Printf("  %s Link %s/%s: %v\n", Red(IconError), s.PlatName, g.Name, err)
					result.Errors++
				} else {
					result.Linked++
				}
			} else if !s.IsLink {
				// 其他实体目录，删除后创建链接
				os.RemoveAll(s.Path)
				err := os.Symlink(repoTarget, s.Path)
				if err != nil {
					fmt.Printf("  %s Link %s/%s: %v\n", Red(IconError), s.PlatName, g.Name, err)
					result.Errors++
				} else {
					result.Linked++
				}
			} else {
				// 已有软链接，修复指向
				err := ReplaceWithSymlink(s.Path, repoTarget)
				if err != nil {
					result.Errors++
				} else if s.LinkDest != repoTarget {
					result.Fixed++
				}
			}
		}
	}

	return result
}

// hashSkillMd 计算仓库中已有 skill 的 md 文件 hash
func hashSkillMd(dir string, category string) string {
	mdName := "SKILL.md"
	if category == "agent" {
		mdName = "AGENT.md"
	}
	data, err := os.ReadFile(filepath.Join(dir, mdName))
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// HashExistingSkill 计算本地仓库中已安装 skill 的 hash（用于 sk add 去重）
func HashExistingSkill(cfg *Config, name, category string) string {
	dir := filepath.Join(cfg.RepoPath, category, name)
	return hashSkillMd(dir, category)
}

// CheckInstallDuplicate 安装前检查是否与本地已有 skill 重复
// 返回: "new"（不存在）, "identical"（完全相同）, "different"（同名但内容不同）
func CheckInstallDuplicate(skill *DiscoveredSkill, cfg *Config) string {
	targetDir := filepath.Join(cfg.RepoPath, skill.Category, skill.Name)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return "new"
	}
	existingHash := hashSkillMd(targetDir, skill.Category)
	if existingHash != "" && skill.Hash != "" && existingHash == skill.Hash {
		return "identical"
	}
	return "different"
}

// shortHash 缩短 hash 显示
func shortHash(h string) string {
	if len(h) > 8 {
		return h[:8]
	}
	if h == "" {
		return "(no hash)"
	}
	return h
}

// truncateDesc 截断描述用于显示
func truncateDesc(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-3]) + "..."
}

// FormatCollectSummary 格式化收集结果摘要
func FormatCollectSummary(r *CollectResult) string {
	parts := []string{}
	if r.Moved > 0 {
		parts = append(parts, fmt.Sprintf("%s Moved: %d", Green(IconSuccess), r.Moved))
	}
	if r.Linked > 0 {
		parts = append(parts, fmt.Sprintf("%s Linked: %d", Cyan(IconLink), r.Linked))
	}
	if r.Fixed > 0 {
		parts = append(parts, fmt.Sprintf("%s Fixed: %d", Blue(IconInfo), r.Fixed))
	}
	if r.Skipped > 0 {
		parts = append(parts, fmt.Sprintf("%s Skipped: %d", Gray("○"), r.Skipped))
	}
	if r.Errors > 0 {
		parts = append(parts, fmt.Sprintf("%s Errors: %d", Red(IconError), r.Errors))
	}
	return strings.Join(parts, "  ")
}
