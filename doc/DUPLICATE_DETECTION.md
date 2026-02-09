# Skill 重复检测机制

本文档记录 Skill Kit 中重复 skill 检测的当前实现逻辑、覆盖场景及已知局限。

## 检测场景

### 1. `sk add` 发现阶段（`download.go` - `DiscoverSkills`）

- **检测方式**：`skillKey = category:name:hash`，用 `seenSkills` map 去重
- **范围**：仅在单次 `sk add` 操作内，防止同一仓库中发现多个相同技能
- **Hash 来源**：SKILL.md/AGENT.md 文件内容的 SHA256
- **局限**：不检查本地仓库，只防同源重复

### 2. `sk add` 安装阶段（`download.go` - `InstallSkill`）

- **检测方式**：
  1. 检查目标目录 `~/.config/agent/{category}/{name}` 是否存在
  2. 存在则调用 `CheckInstallDuplicate` 比较 SKILL.md 的 hash

- **判定结果**：

| 结果 | 含义 | 处理 |
|------|------|------|
| `new` | 本地不存在 | 正常安装 |
| `identical` | 同名且内容相同 | 跳过，提示已存在 |
| `different` | 同名但内容不同 | 跳过，警告 |

- **局限**：
  - 只按目录名匹配，不同名但内容完全相同的 skill 不会被检测
  - `different` 时只跳过，没有提供覆盖/重命名选项

### 3. `sk init --collect` 收集阶段（`collect.go` - `RunCollect`）

- **检测方式**：
  1. 扫描所有平台 global 目录，收集实体目录和软链接
  2. 按 `category:name` 分组（`GroupByName`）
  3. `HasConflict()` 检查同组内实体 skill 的 hash 是否有多个不同值

- **判定逻辑**：

| 情况 | 处理 |
|------|------|
| 组内只有 1 个实体 | 直接移动到中央仓库 |
| 多个实体，hash 全相同 | 取第一个，其余删除替换为链接 |
| 多个实体，hash 不同 | 弹出 `ConflictResolveMenu` 让用户选择 |
| 全是软链接 | 检查指向是否正确，不正确则修复 |
| 中央仓库已有同名 | 比较 hash，相同跳过，不同也跳过并警告 |

- **局限**：
  - Hash 只算 SKILL.md/AGENT.md 单个文件，其他文件差异不参与
  - 软链接的 hash 不参与冲突判定
  - 中央仓库已有且内容不同时只跳过，没有覆盖选项

## 已知弱点

1. **Hash 粒度太粗**：只算 SKILL.md/AGENT.md 一个文件的 SHA256，忽略目录中其他文件（prompt 模板、配置文件等）的差异。两个 md 相同但辅助文件不同的 skill 会被判定为 identical。

2. **不同名相同内容无法检测**：三个场景都是先按名称匹配，两个功能完全相同但目录名不同的 skill 不会被识别为重复。

3. **`different` 场景缺少用户选择**：`sk add` 遇到同名不同内容时只是跳过；`init --collect` 遇到仓库已有不同内容也只是跳过。均没有提供覆盖或重命名选项。

## 改进方向

- 将 hash 范围扩展到整个目录内容（遍历所有文件计算联合 hash）
- 增加跨名称的内容相似度检测（基于目录 hash 或文件指纹）
- `sk add` 遇到 `different` 时提供覆盖/跳过/重命名三选一交互
