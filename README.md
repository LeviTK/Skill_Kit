# Skill Kit

Install, manage and distribute AI agent skills across multiple coding agents from a single repository.

Supports **Claude Code**, **Cursor**, **Amp**, **OpenCode**, **Codex**, and [9 more](#supported-platforms).

## Quick Start

```bash
# Install
make build && make install

# Download skills from GitHub
sk add vercel-labs/agent-skills

# Distribute to all platforms
sk use my-skill
```

## What is Skill Kit?

Skill Kit is a cross-platform AI skill distribution hub that:

- **Downloads** skills from any git repository (GitHub, GitLab, or local)
- **Manages** skills in a unified repository (`~/.config/agent/`)
- **Distributes** via symlinks to 14+ AI coding tools
- **Syncs** changes instantly - edit once, update everywhere

## Usage

### Source Formats

The `<source>` argument accepts multiple formats:

```bash
# GitHub shorthand
sk add vercel-labs/agent-skills

# Full GitHub URL
sk add https://github.com/vercel-labs/agent-skills

# Direct path to a skill in a repo
sk add https://github.com/owner/repo/tree/main/skills/my-skill

# GitLab URL
sk add https://gitlab.com/org/repo

# Local directory
sk add ./local/skills
```

### Commands

| Command | Description |
|---------|-------------|
| `sk` | Interactive menu (recommended) |
| `sk add <source>` | Download skills from git repo or local path |
| `sk use <module> [platform]` | Distribute skill to platform(s) via symlink |
| `sk list` | List all modules and their link status |
| `sk platforms` | Show registered platforms |
| `sk info <module>` | Show module details and aliases |
| `sk remove <module> [platform]` | Remove symlinks for a module |
| `sk status` | Health check: detect broken symlinks |
| `sk sync` | Sync all modules to all platforms |
| `sk init` | Initialize the agent repository |

### Examples

```bash
# Interactive mode
sk

# Download and select skills to install
sk add vercel-labs/agent-skills

# Distribute a skill to all default platforms
sk use my-skill

# Distribute to specific platform
sk use my-skill claude

# Use custom link name
sk use my-skill cursor --as py-coder

# List all modules with status
sk list

# Remove skill from all platforms
sk remove my-skill
```

## Supported Platforms

Skills can be distributed to any of these supported agents:

| Platform | Project Path | Global Path |
|----------|--------------|-------------|
| OpenCode | `.opencode/skills/` | `~/.config/opencode/skills/` |
| Claude Code | `.claude/skills/` | `~/.claude/skills/` |
| OpenAI Codex | `.codex/skills/` | `~/.codex/skills/` |
| Cursor | `.cursor/skills/` | `~/.cursor/skills/` |
| Amp | `.agents/skills/` | `~/.config/agents/skills/` |
| Kilo Code | `.kilocode/skills/` | `~/.kilocode/skills/` |
| Roo Code | `.roo/skills/` | `~/.roo/skills/` |
| Goose | `.goose/skills/` | `~/.config/goose/skills/` |
| Gemini CLI | `.gemini/skills/` | `~/.gemini/skills/` |
| Antigravity | `.agent/skills/` | `~/.gemini/antigravity/skills/` |
| GitHub Copilot | `.github/skills/` | `~/.copilot/skills/` |
| Clawdbot | `skills/` | `~/.clawdbot/skills/` |
| Droid | `.factory/skills/` | `~/.factory/skills/` |
| Windsurf | `.windsurf/skills/` | `~/.codeium/windsurf/skills/` |

## Directory Structure

```
~/.config/agent/
├── platforms.toml        # Platform registry
├── skill/                # Skill pool
│   └── my-skill/
│       ├── SKILL.md      # Skill documentation
│       └── skillkit.toml # Optional: custom config
└── agent/                # Agent pool
    └── my-agent/
```

## Platform Configuration

Edit `~/.config/agent/platforms.toml` to add or modify platforms:

```toml
[platforms.claude]
name = "Claude Code"
project = ".claude/"
global = "~/.claude/"
skill_dir = "skills"
agent_dir = "agents"

# Default platforms for sync
default_platforms = ["claude", "cursor", "amp"]
```

## Module Aliases

Create `skillkit.toml` in module directory to customize link names:

```toml
[link]
default = "python-coder"

[link.overrides]
claude = "py-expert"
cursor = "py-coder"
```

## Creating Skills

Skills are directories containing a `SKILL.md` file with YAML frontmatter:

```markdown
---
name: my-skill
description: What this skill does and when to use it
---

# My Skill

Instructions for the agent...
```

## Related Links

- [Vercel Agent Skills Repository](https://github.com/vercel-labs/agent-skills)
- [Agent Skills Specification](https://agentskills.io)
- [add-skill CLI](https://github.com/vercel-labs/add-skill)

## License

MIT
