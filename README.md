# Skill Kit

**跨平台 AI 技能分发枢纽** - 将统一仓库中的 AI Skills/Agents 通过软链接分发到不同 AI 工具的配置目录。

## 功能

- **统一管理**：在 `~/.config/agent/` 维护所有 Skills 和 Agents
- **多平台支持**：Claude、Cursor、Copilot、OpenCode、Gemini、Windsurf 等 8 个平台
- **软链接分发**：一处修改，多处生效
- **自定义别名**：支持为不同平台指定不同的链接名称

## 快速开始

```bash
# 构建
make build

# 初始化配置目录
make init

# 安装到系统
make install
```

## 使用

```bash
# 分发模块到指定平台
sk use python-coder claude

# 分发到所有平台
sk use python-coder

# 使用自定义链接名
sk use python-coder cursor --as py-coder

# 列出所有模块及状态
sk list

# 查看已注册平台
sk platforms

# 查看模块详情
sk info python-coder

# 移除软链接
sk remove python-coder claude

# 同步所有模块
sk sync
```

## 目录结构

```
~/.config/agent/
├── platforms.toml        # 平台注册表
├── skill/                # 技能池
│   └── python-coder/
│       ├── AGENT.md      # 技能文档
│       └── skillkit.toml # 可选：自定义配置
└── agent/                # 代理池
    └── writer-bot/
```

## 平台配置

编辑 `~/.config/agent/platforms.toml` 添加或修改平台：

```toml
[platforms.claude]
name = "Claude Code"
project = ".claude/"
global = "~/.claude/"
skill_dir = "skills"
agent_dir = "agents"
```

## 模块别名

在模块目录下创建 `skillkit.toml` 自定义链接名：

```toml
[link]
default = "python-coder"

[link.overrides]
claude = "py-expert"
cursor = "py-coder"
```

## 支持的平台

| 平台               | 全局路径                       |
| ------------------ | ------------------------------ |
| Claude Code        | `~/.claude/skills/`              |
| GitHub Copilot     | `~/.copilot/skills/`             |
| Google Antigravity | `~/.gemini/antigravity/skills/`  |
| Cursor             | `~/.cursor/skills/`              |
| OpenCode           | `~/.config/opencode/skill/`      |
| OpenAI Codex       | `~/.codex/skills/`               |
| Gemini CLI         | `~/.gemini/skills/`              |
| Windsurf           | `~/.codeium/windsurf/skills/`    |

## License

MIT
