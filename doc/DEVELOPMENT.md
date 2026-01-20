# Skill Kit 开发文档

## 1. 项目概览

Skill Kit 是一个用 Go 编写的跨平台 AI 技能分发工具。它作为一个中心枢纽，通过符号链接管理并将“技能”（提示词/指令集合）同步到各种 AI 编程助手（Claude Code, Cursor 等）。

## 2. 架构

该项目遵循标准的 Go CLI 结构：

```
skillkit/
├── cmd/
│   └── main.go           # 入口点，命令路由，交互循环
├── lib/
│   ├── config.go         # 配置加载 (platforms.toml)
│   ├── module.go         # 技能/Agent 发现与元数据解析
│   ├── commands.go       # 命令注册表
│   ├── link.go           # 软链接管理 (包含安全检查)
│   ├── ui.go             # TUI 组件 (菜单，选择器)
│   └── ...
├── doc/                  # 文档
└── platforms.toml        # 默认平台定义
```

### 核心概念

*   **仓库 (Repository) (`~/.config/agent/`)**：所有下载的技能和 Agent 的中心存储位置。
*   **模块 (Module)**：功能单元（“技能”或“Agent”）。
    *   元数据从 `SKILL.md` 或 `AGENT.md` (YAML frontmatter) 读取。
    *   配置从 `skillkit.toml` (链接名称，别名) 读取。
*   **平台 (Platform)**：使用技能的外部工具（如 Claude, Cursor）。在 `platforms.toml` 中定义。
*   **分发 (Distribution)**：从仓库创建符号链接到平台配置目录的过程。

## 3. 关键组件

### 配置 (`lib/config.go`)
加载 `platforms.toml`。支持环境变量覆盖：
- `SKILLKIT_REPO`: 覆盖默认仓库路径。
- `SKILLKIT_CONFIG`: 覆盖特定配置文件路径。

### 模块管理 (`lib/module.go`)
- **发现**：扫描 `skill/` 和 `agent/` 目录。
- **解析**：从 Markdown 文件提取描述，从 TOML 文件提取别名。
- **多态性**：类似地处理“技能”（被动提示词）和“Agent”（可执行单元）。

### 命令处理 (`cmd/main.go`)
- **交互模式**：如果未提供参数，进入 TUI 循环 (`lib.InteractiveMenu`)。
- **CLI 模式**：直接执行命令，如 `add`, `use`, `sync`。

## 4. 开发工作流

### 先决条件
- Go 1.24+
- Make

### 构建
```bash
make build
# 输出二进制文件在 bin/skillkit
```

### 测试
```bash
go test -v ./...
```

### 添加新命令
1.  在 `lib/Commands` (`lib/commands.go`) 中定义命令元数据。
2.  在 `cmd/main.go` 中实现处理函数 (例如 `handleNewCommand`)。
3.  在 `main()` 和 `handleInteractiveCommand` 的 `switch` 语句中添加该 case。

### 添加新平台
更新 `platforms.toml`。通常不需要更改代码，除非平台具有标准配置无法处理的独特目录结构。

## 5. 风格指南
- **日志**：使用 `lib` 包的颜色函数 (`lib.Red`, `lib.Blue`) 和图标 (`lib.IconSuccess`) 以保持输出一致。
- **错误处理**：将错误返回给调用者；在适用的情况下使用 `lib.Errors` 定义。
- **UI**：保持 TUI 交互简单且支持键盘导航。
