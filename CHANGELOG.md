# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2025-01-20

### Added
- **`sk add` command**: Download skills from GitHub, GitLab, or local directories
  - Support GitHub shorthand (`owner/repo`)
  - Support full URLs with branch/path (`https://github.com/owner/repo/tree/main/path`)
  - Support GitLab URLs
  - Support local paths (`./local/path`)
  - Interactive skill selection with multi-select
  - Auto-discovery of SKILL.md/AGENT.md files
- **14 platform support**: OpenCode, Claude Code, Codex, Cursor, Amp, Kilo Code, Roo Code, Goose, Gemini CLI, Antigravity, GitHub Copilot, Clawdbot, Droid, Windsurf
- **CJK input handling**: Improved keyboard input for Chinese/Japanese/Korean input methods
- **Ctrl+C support**: Exit with Ctrl+C in interactive mode

### Changed
- Renamed project from Linktrack to Skill Kit
- Binary renamed from `lt` to `sk`
- Config file renamed from `linktrack.toml` to `skillkit.toml`
- Environment variables renamed from `LINKTRACK_*` to `SKILLKIT_*`
- Updated ASCII banner

### Fixed
- Fixed usage messages showing old command name
- Fixed Makefile build path issues

## [0.0.1] - 2025-01-20

### Added
- Initial release as Skill Kit
- Symlink-based skill distribution
- Interactive TUI menu
- Platform configuration via TOML
- Module aliases support
