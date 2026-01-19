# AGENTS.md

## Overview

**Skill Kit** - 跨平台 AI 技能分发枢纽，通过软链接将统一仓库中的 AI Skills/Agents 分发到不同 AI 工具的配置目录。

## Tech Stack

- **Language**: Go 1.24+
- **Dependencies**: 
  - `github.com/pelletier/go-toml/v2` (TOML 解析)
  - `golang.org/x/term` (终端 Raw 模式)
- **Build Tool**: Make

## Project Structure

```
skill-kit/
├── cmd/main.go           # 入口，命令路由与参数解析
├── lib/                  # 核心库
│   ├── color.go          # 终端颜色与图标
│   ├── commands.go       # 命令定义注册表
│   ├── config.go         # 配置加载 (platforms.toml)
│   ├── errors.go         # 自定义错误类型
│   ├── errors_test.go    # 错误类型测试
│   ├── link.go           # 软链接操作 (创建/删除/检查)
│   ├── link_test.go      # 软链接测试
│   ├── module.go         # 模块发现与管理 (含描述读取)
│   ├── module_test.go    # 模块测试
│   └── ui.go             # 交互式菜单、选择器、确认弹窗
├── bin/                  # 构建产物
├── scripts/install.sh    # 安装脚本
├── platforms.toml        # 默认平台配置
├── Makefile              # 构建脚本
└── README.md
```

## Commands

```bash
# 构建
make build

# 运行测试
make test
# 或
go test -v ./...

# 类型检查/编译检查
go build ./...

# 清理构建产物
make clean

# 安装到系统
make install

# 初始化用户配置目录
make init
```

## CLI Usage

```bash
# 交互式菜单（推荐）
skillkit

# 命令行模式
sk use <module> [platform] [--global|--project] [--as <name>]
sk list
sk platforms
sk info <module>
sk remove <module> [platform]
sk sync [--dry-run]
sk status
sk init
```

## Interactive Mode

交互式菜单支持键盘导航：

| 按键 | 功能 |
|------|------|
| `↑/↓` 或 `k/j` | 上下导航 |
| `←` | 返回上一级 / 应用变更 |
| `→` 或 `Enter` | 进入下一级 / 确认 |
| `Space` | 切换选中状态 |
| `Q` | 退出程序 |
| `Y/N` | 确认/取消弹窗 |

### Use 命令交互流程

1. **模块列表页**: 显示所有模块，顶部显示当前选中模块的已同步平台
   - `Enter`: 重新同步到已同步平台
   - `→`: 进入详情页手动管理平台

2. **模块详情页**: 显示模块信息和平台列表
   - 从 SKILL.md/AGENT.md 读取 `description` 字段显示描述
   - `[✓]` 选中 = 同步, `[ ]` 取消 = 删除
   - `←` 返回时自动应用变更（需二次确认）

## Environment Variables

- `SKILLKIT_REPO` - 覆盖模块仓库路径 (默认 `~/.config/agent`)
- `SKILLKIT_CONFIG` - 覆盖配置文件路径 (默认 `~/.config/agent/platforms.toml`)

## Key Concepts

- **Module**: `~/.config/agent/skill/` 或 `~/.config/agent/agent/` 下的目录
- **Platform**: 定义在 `platforms.toml` 中的 AI 工具配置
- **Symlink**: 从源模块到目标平台配置目录的软链接
- **Description**: 从模块的 SKILL.md/AGENT.md frontmatter 中读取

## Guidelines

- 遵循现有代码风格，所有库代码放在 `lib/` 包
- 命令处理逻辑在 `cmd/main.go` 的 `handle*` 函数中
- 交互式命令处理在 `handleInteractive*` 函数中
- 使用 `lib.ResolvePath()` 处理路径展开 (`~` 和相对路径)
- 终端输出使用 `lib/color.go` 中的颜色函数和图标常量
- 新增命令需在 `lib/commands.go` 的 `Commands` 切片中注册
- UI 组件在 `lib/ui.go` 中：
  - `InteractiveMenu()` - 主菜单
  - `ModuleListMenu()` - 模块列表（显示同步状态）
  - `ModuleDetailMenu()` - 模块详情（平台多选）
  - `SelectMenu()` - 通用选择菜单
  - `ConfirmDialog()` - 确认弹窗

## Notes for Agents

- 修改代码后运行 `go build ./...` 确保编译通过
- 软链接操作有保护机制，不会覆盖实体文件/目录
- 配置文件位于 `~/.config/agent/platforms.toml`
- 使用 `golang.org/x/term` 处理终端 Raw 模式（跨平台兼容）
- 模块描述从 SKILL.md/AGENT.md 的 frontmatter `description:` 字段读取
