# 更新日志

本文件记录 Skill Kit 项目的所有重要变更。

## [0.3.0] - 2026-02-09

### 新增

- **`sk init --collect` 命令**：扫描所有平台目录，将实体 skill 收集到中央仓库，用软链接替换
  - 自动检测跨平台重复 skill（基于目录名 + SKILL.md 内容 hash）
  - 内容相同的重复自动合并，内容不同的弹出交互菜单让用户选择保留版本
  - 已有软链接检查指向是否正确，不正确则自动修复
- **`sk add` 安装时重复检测**：安装前基于 hash 检查本地仓库是否已有相同 skill
  - 内容完全相同则跳过（不再报错）
  - 同名但内容不同则警告并跳过
- **冲突解决交互菜单**（`ConflictResolveMenu`）：显示来源平台、hash、描述、路径，用户逐个选择保留版本
- **重复检测机制文档**：`doc/DUPLICATE_DETECTION.md`

### 改进

- 更新 AGENTS.md：补充 collect.go、download.go 等新文件，新增收集流程和重复检测章节

## [0.2.1] - 2026-02-09

### 改进

- **模块详情页取消二次确认**：在模块详情页按 `←` 返回时，有变更（同步/删除）直接执行，不再弹出确认弹窗，操作更流畅
- **修复描述区域跳动问题**：模块列表页预计算所有模块描述的最大行数，固定描述区域高度，切换模块时列表不再上下跳动
- **重构描述渲染逻辑**：抽取 `getDescMaxWidth`、`countWrappedLines`、`printWrappedDescFixed` 函数，提升代码复用性

## [0.2.0] - 2025-02-08

### 新增

- 新增 `CLAUDE.MD` 软链接指向 `AGENTS.md`，统一 AI Agent 指导文件

## [0.1.0] - 2025-01-20

### 新增

- **`sk add` 命令**：从 GitHub、GitLab 或本地路径下载技能
  - 支持 GitHub 简写（`owner/repo`）
  - 支持带分支/路径的完整 URL（`https://github.com/owner/repo/tree/main/path`）
  - 支持 GitLab URL
  - 支持本地路径（`./local/path`）
  - 交互式技能多选
  - 自动发现 SKILL.md/AGENT.md 文件
- **14 个平台支持**：OpenCode、Claude Code、Codex、Cursor、Amp、Kilo Code、Roo Code、Goose、Gemini CLI、Antigravity、GitHub Copilot、Clawdbot、Droid、Windsurf
- **CJK 输入处理**：改进中文/日文/韩文输入法环境下的键盘输入

### 变更

- 项目从 Linktrack 重命名为 Skill Kit
- 二进制文件从 `lt` 重命名为 `sk`
- 配置文件从 `linktrack.toml` 重命名为 `skillkit.toml`
- 环境变量从 `LINKTRACK_*` 重命名为 `SKILLKIT_*`

### 修复

- 修复使用说明中显示旧命令名的问题
- 修复 Makefile 构建路径问题

## [0.0.1] - 2025-01-20

### 新增

- 初始版本发布
- 基于软链接的技能分发机制
- 交互式 TUI 菜单
- TOML 平台配置
- 模块别名支持
