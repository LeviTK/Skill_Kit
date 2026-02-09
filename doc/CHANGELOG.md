# 更新日志

本文件记录 Skill Kit 项目的所有重要变更。

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
