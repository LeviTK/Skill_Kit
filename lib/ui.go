package lib

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"golang.org/x/term"
)

// UI 交互式界面组件

const (
	Version = "0.2.0"
	Tagline = "Cross-Platform AI Skill Distribution Hub"
)

// ShowBanner 显示品牌横幅
func ShowBanner() {
	banner := `
` + ColorCyan + `  _     _       _   _                  _    ` + ColorReset + `
` + ColorCyan + ` | |   (_)_ __ | | |_ _ __ __ _  ___| | __` + ColorReset + `
` + ColorCyan + ` | |   | | '_ \| |/ / '__/ _` + "`" + ` |/ __| |/ /` + ColorReset + `
` + ColorCyan + ` | |___| | | | |   <| | | (_| | (__|   < ` + ColorReset + `
` + ColorCyan + ` |_____|_|_| |_|_|\_\_|  \__,_|\___|_|\_\` + ColorReset + `
`
	fmt.Print(banner)
	fmt.Printf("\n %s%s%s\n\n", ColorGray, Tagline, ColorReset)
}

// ShowHelp 显示帮助信息
func ShowHelp() {
	ShowBanner()

	fmt.Printf("%sUSAGE%s\n", ColorBlue, ColorReset)
	fmt.Printf("  %slt%s                          Interactive menu (recommended)\n", ColorGreen, ColorReset)
	fmt.Printf("  %slt <command> [options]%s      Run a specific command\n\n", ColorGreen, ColorReset)

	fmt.Printf("%sCOMMANDS%s\n", ColorBlue, ColorReset)
	for _, cmd := range Commands {
		display := cmd.Name
		if cmd.Name == "help" {
			display = "-h, --help"
		} else if cmd.Name == "version" {
			display = "-v, --version"
		}
		fmt.Printf("  %s%-20s%s %s\n", ColorGreen, display, ColorReset, cmd.Description)
	}

	fmt.Println()
	fmt.Printf("%sOPTIONS%s\n", ColorBlue, ColorReset)
	fmt.Printf("  %s%-20s%s %s\n", ColorGreen, "--dry-run", ColorReset, "Preview without making changes")
	fmt.Printf("  %s%-20s%s %s\n", ColorGreen, "--global", ColorReset, "Use global scope (default)")
	fmt.Printf("  %s%-20s%s %s\n", ColorGreen, "--project", ColorReset, "Use project scope")
	fmt.Printf("  %s%-20s%s %s\n", ColorGreen, "--as <name>", ColorReset, "Override link name")

	fmt.Println()
	fmt.Printf("%sEXAMPLES%s\n", ColorBlue, ColorReset)
	fmt.Printf("  %ssk%s                          Start interactive menu\n", ColorGreen, ColorReset)
	fmt.Printf("  %ssk use my-skill%s             Sync skill to all platforms\n", ColorGreen, ColorReset)
	fmt.Printf("  %ssk use my-skill amp%s         Sync skill to specific platform\n", ColorGreen, ColorReset)
	fmt.Printf("  %ssk list%s                     Show all modules and status\n", ColorGreen, ColorReset)
	fmt.Printf("  %ssk remove my-skill%s          Remove skill from all platforms\n", ColorGreen, ColorReset)

	fmt.Println()
	fmt.Printf("%sUNINSTALL%s\n", ColorBlue, ColorReset)
	fmt.Printf("  %ssudo rm /usr/local/bin/sk /usr/local/bin/skillkit%s\n", ColorGray, ColorReset)
	fmt.Printf("  %srm -rf ~/.config/agent%s      (optional: remove config)\n", ColorGray, ColorReset)

	fmt.Println()
}

// ShowVersion 显示版本信息
func ShowVersion() {
	fmt.Printf("\nSkill Kit version %s\n", Version)
	fmt.Printf("OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go: %s\n\n", runtime.Version())
}

// MenuOption 菜单选项
type MenuOption struct {
	Number      int
	Name        string
	Description string
	Command     string
}

// MainMenuOptions 主菜单选项
var MainMenuOptions = []MenuOption{
	{1, "Use", "Distribute skill to platforms", "use"},
	{2, "List", "Show all modules and status", "list"},
	{3, "Platforms", "View registered platforms", "platforms"},
	{4, "Defaults", "Set default platforms for sync", "defaults"},
	{5, "Init", "Initialize repository", "init"},
}

// ShowMainMenu 显示主菜单
func ShowMainMenu(selected int, cfg *Config) {
	ClearScreen()
	ShowBanner()

	for _, opt := range MainMenuOptions {
		if opt.Number == selected {
			fmt.Printf("  %s%s %d. %-12s%s %s\n",
				ColorCyan, IconArrow, opt.Number, opt.Name, ColorReset, opt.Description)
		} else {
			fmt.Printf("    %d. %-12s %s\n", opt.Number, opt.Name, opt.Description)
		}
	}

	fmt.Println()
	fmt.Printf("  %s↑↓ Navigate  |  →/Enter Select  |  Q Quit  |  H Help%s\n", ColorGray, ColorReset)
	fmt.Println()
}

// ClearScreen 清屏
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// HideCursor 隐藏光标
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor 显示光标
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// ReadKey 读取单个按键
func ReadKey() string {
	// 设置终端为 raw 模式
	fd := int(os.Stdin.Fd())
	if runtime.GOOS != "windows" {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			return "QUIT"
		}
		defer term.Restore(fd, oldState)
	}

	buf := make([]byte, 6) // 增加缓冲区以支持 Shift+方向键
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return "QUIT"
	}

	switch buf[0] {
	case 'q', 'Q':
		return "QUIT"
	case 'h', 'H':
		return "HELP"
	case 'j':
		return "DOWN"
	case 'J':
		return "MOVEDOWN"
	case 'k':
		return "UP"
	case 'K':
		return "MOVEUP"
	case 'l', 'L':
		return "RIGHT"
	case ' ':
		return "SPACE"
	case 'r', 'R':
		return "RESET"
	case 'y', 'Y':
		return "YES"
	case 'n', 'N':
		return "NO"
	case 'v', 'V':
		return "VMODE"
	case 'a', 'A':
		return "SELECTALL"
	case '\t':
		return "TAB"
	case '\r', '\n':
		return "ENTER"
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return fmt.Sprintf("NUM:%c", buf[0])
	case 27: // ESC 序列
		if n >= 3 && buf[1] == '[' {
			// 检查 Shift+方向键: ESC [1;2A/B/C/D
			if n >= 6 && buf[2] == '1' && buf[3] == ';' && buf[4] == '2' {
				switch buf[5] {
				case 'A':
					return "MOVEUP"
				case 'B':
					return "MOVEDOWN"
				}
			}
			// 普通方向键
			switch buf[2] {
			case 'A':
				return "UP"
			case 'B':
				return "DOWN"
			case 'C':
				return "RIGHT"
			case 'D':
				return "LEFT"
			}
		}
		return "ESC"
	}

	return "OTHER"
}

// WaitForKey 等待任意按键
func WaitForKey() {
	fmt.Print("Press any key to continue...")
	ReadKey()
}

// InteractiveMenu 交互式菜单主循环
func InteractiveMenu() string {
	selected := 1
	maxOptions := len(MainMenuOptions)

	HideCursor()
	defer ShowCursor()

	for {
		// 加载配置以显示默认平台
		cfg, _ := LoadConfig()
		if cfg == nil {
			cfg = &Config{Platforms: make(map[string]Platform)}
		}

		ShowMainMenu(selected, cfg)

		key := ReadKey()

		switch key {
		case "UP":
			if selected > 1 {
				selected--
			}
		case "DOWN":
			if selected < maxOptions {
				selected++
			}
		case "RIGHT", "ENTER":
			return MainMenuOptions[selected-1].Command
		case "NUM:1":
			return MainMenuOptions[0].Command
		case "NUM:2":
			return MainMenuOptions[1].Command
		case "NUM:3":
			return MainMenuOptions[2].Command
		case "NUM:4":
			return MainMenuOptions[3].Command
		case "NUM:5":
			return MainMenuOptions[4].Command
		case "HELP":
			ClearScreen()
			ShowHelp()
			fmt.Print("Press any key to continue...")
			ReadKey()
		case "QUIT":
			return "quit"
		}
	}
}

// SelectOption 通用选择项
type SelectOption struct {
	Key   string
	Label string
}

// SelectMenuResult 选择菜单结果
type SelectMenuResult struct {
	Key    string
	Back   bool // 用户选择返回上一级
	Cancel bool // 用户取消（退出）
}

// SelectMenu 通用选择菜单，返回选中项的 Key，支持返回上一级
func SelectMenu(title string, options []SelectOption) SelectMenuResult {
	if len(options) == 0 {
		return SelectMenuResult{Cancel: true}
	}

	selected := 0
	maxOptions := len(options)

	HideCursor()
	defer ShowCursor()

	for {
		ClearScreen()
		fmt.Printf("\n  %s %s\n\n", Blue(IconArrow), title)

		for i, opt := range options {
			if i == selected {
				fmt.Printf("  %s %s\n", Cyan(IconArrow), White(opt.Label))
			} else {
				fmt.Printf("    %s\n", opt.Label)
			}
		}

		fmt.Println()
		fmt.Printf("  %s↑↓ Navigate  |  Enter Select  |  ← Back  |  Q Quit%s\n", ColorGray, ColorReset)

		key := ReadKey()

		switch key {
		case "UP":
			if selected > 0 {
				selected--
			}
		case "DOWN":
			if selected < maxOptions-1 {
				selected++
			}
		case "ENTER", "RIGHT":
			return SelectMenuResult{Key: options[selected].Key}
		case "LEFT":
			return SelectMenuResult{Back: true}
		case "QUIT":
			return SelectMenuResult{Cancel: true}
		}
	}
}

// ModuleListResult 模块列表操作结果
type ModuleListResult struct {
	Module  *Module
	Modules []*Module // 用于批量操作
	Action  string    // "sync", "sync_default", "sync_all_default", "detail", "back", "quit"
}

// SelectModuleMenu 模块选择菜单（显示同步平台信息）
func SelectModuleMenu(cfg *Config) SelectMenuResult {
	modules, err := ListModules(cfg)
	if err != nil || len(modules) == 0 {
		return SelectMenuResult{Cancel: true}
	}

	options := make([]SelectOption, len(modules))
	for i, m := range modules {
		options[i] = SelectOption{
			Key:   m.Name,
			Label: fmt.Sprintf("%s (%s)", m.Name, m.Category),
		}
	}

	return SelectMenu("Select Module", options)
}

// ModuleListMenu 增强版模块列表菜单，显示已同步平台，支持 v 模式快速同步到默认平台
func ModuleListMenu(cfg *Config) ModuleListResult {
	modules, err := ListModules(cfg)
	if err != nil || len(modules) == 0 {
		return ModuleListResult{Action: "quit"}
	}

	selected := 0
	maxOptions := len(modules)
	selectAll := false // 全选状态

	HideCursor()
	defer ShowCursor()

	for {
		ClearScreen()

		// 标题
		if selectAll {
			fmt.Printf("\n  %s %s %s\n", Blue(IconArrow), "Select Module", Magenta("[ALL]"))
		} else {
			fmt.Printf("\n  %s %s\n", Blue(IconArrow), "Select Module")
		}

		// 显示默认平台（在 v 模式或始终显示）
		defaultPlatformNames := getDefaultPlatformNames(cfg)
		if len(defaultPlatformNames) > 0 {
			fmt.Printf("  %s %s\n", Gray("Default:"), Yellow(joinStrings(defaultPlatformNames, ", ")))
		} else {
			fmt.Printf("  %s %s\n", Gray("Default:"), Gray("(all platforms)"))
		}

		// 显示当前选中模块的已同步平台
		currentMod := modules[selected]
		syncedPlatforms := getSyncedPlatformNames(cfg, currentMod)
		if len(syncedPlatforms) > 0 {
			fmt.Printf("  %s %s\n", Gray("Synced:"), Cyan(joinStrings(syncedPlatforms, ", ")))
		} else {
			fmt.Printf("  %s %s\n", Gray("Synced:"), Gray("(none)"))
		}

		// 显示当前选中模块的描述（带换行对齐）
		if currentMod.Description != "" {
			printWrappedDesc(currentMod.Description)
		}

		fmt.Println()

		// 模块列表
		for i, m := range modules {
			prefix := "  "
			if selectAll {
				prefix = Green("✓ ")
			}
			if i == selected {
				fmt.Printf("%s%s %s %s\n", prefix, Cyan(IconArrow), White(m.Name), Gray("("+m.Category+")"))
			} else {
				fmt.Printf("%s  %s %s\n", prefix, m.Name, Gray("("+m.Category+")"))
			}
		}

		fmt.Println()
		fmt.Printf("  %s↑↓ Navigate  |  Enter Sync to Default  |  A Select All  |  → Details  |  ← Back%s\n", ColorGray, ColorReset)

		key := ReadKey()

		switch key {
		case "UP":
			if selected > 0 {
				selected--
			}
			selectAll = false
		case "DOWN":
			if selected < maxOptions-1 {
				selected++
			}
			selectAll = false
		case "SELECTALL":
			selectAll = !selectAll
		case "ESC":
			selectAll = false
		case "ENTER":
			if selectAll {
				return ModuleListResult{Modules: modules, Action: "sync_all_default"}
			}
			return ModuleListResult{Module: modules[selected], Action: "sync_default"}
		case "RIGHT":
			return ModuleListResult{Module: modules[selected], Action: "detail"}
		case "LEFT":
			return ModuleListResult{Action: "back"}
		case "QUIT":
			return ModuleListResult{Action: "quit"}
		}
	}
}

// getDefaultPlatformNames 获取默认平台的名称列表
func getDefaultPlatformNames(cfg *Config) []string {
	var names []string
	for _, key := range cfg.DefaultPlatforms {
		if p, ok := cfg.Platforms[key]; ok {
			names = append(names, p.Name)
		}
	}
	return names
}

// getSyncedPlatformNames 获取已同步的平台名称列表
func getSyncedPlatformNames(cfg *Config, mod *Module) []string {
	var names []string
	for key, p := range cfg.Platforms {
		targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := targetDir + "/" + mod.GetLinkName(key)
		if IsSymlink(targetPath) {
			names = append(names, p.Name)
		}
	}
	return names
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// printWrappedDesc 打印描述文本，支持换行对齐
func printWrappedDesc(desc string) {
	prefix := "  " + Gray("Desc:") + " "
	indent := "        " // 8 空格，与 "  Desc: " 对齐

	// 获取终端宽度，默认 80，使用 2/3
	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}
	maxWidth := width * 2 / 3
	if maxWidth < 40 {
		maxWidth = 40
	}

	// 按 rune 处理，避免中文截断
	runes := []rune(desc)
	lineWidth := 0
	firstLine := true

	for _, r := range runes {
		charWidth := 1
		if r > 127 {
			charWidth = 2 // 中文等宽字符占 2 格
		}

		if lineWidth+charWidth > maxWidth && lineWidth > 0 {
			fmt.Println()
			fmt.Print(indent)
			lineWidth = 0
			firstLine = false
		}

		if lineWidth == 0 && firstLine {
			fmt.Print(prefix)
		}

		fmt.Print(string(r))
		lineWidth += charWidth
	}
	fmt.Println()
}

// ListMenu 交互式列表菜单，支持模式切换和平台排序
func ListMenu(cfg *Config) {
	modules, err := ListModules(cfg)
	if err != nil {
		return
	}

	// 模式: 0 = 按模块分组, 1 = 按平台分组
	mode := 0
	selected := 0                              // 平台视图中的选中项
	platformKeys := cfg.GetOrderedPlatformKeys() // 使用有序的平台列表

	HideCursor()
	defer ShowCursor()

	for {
		ClearScreen()

		if mode == 0 {
			// 模式1: 按模块分组
			fmt.Printf("\n  %s %s %s\n\n", Blue(IconFolder), "Modules", Gray("[Tab: Platform view]"))

			if len(modules) == 0 {
				fmt.Printf("  %s No modules found\n", Yellow(IconWarning))
			} else {
				for _, mod := range modules {
					fmt.Printf("  %s %s %s\n", Cyan(IconArrow), White(mod.Name), Gray("("+mod.Category+")"))
					// 获取已同步的平台（按顺序）
					for _, key := range platformKeys {
						p := cfg.Platforms[key]
						targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
						targetPath := targetDir + "/" + mod.GetLinkName(key)
						if IsSymlink(targetPath) {
							fmt.Printf("      %s %s\n", Green(IconSuccess), p.Name)
						}
					}
				}
			}

			fmt.Println()
			fmt.Printf("  %sTab Switch View  |  ← Back  |  Q Quit%s\n", ColorGray, ColorReset)
		} else {
			// 模式2: 按平台分组（支持排序）
			fmt.Printf("\n  %s %s %s\n\n", Blue(IconFolder), "Platforms", Gray("[Tab: Module view]"))

			for i, key := range platformKeys {
				p := cfg.Platforms[key]
				prefix := "  "
				if i == selected {
					prefix = Cyan(IconArrow) + " "
				}
				fmt.Printf("%s%s %s\n", prefix, White(p.Name), Gray("("+key+")"))
				// 获取该平台下的已同步模块
				hasModule := false
				for _, mod := range modules {
					targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
					targetPath := targetDir + "/" + mod.GetLinkName(key)
					if IsSymlink(targetPath) {
						fmt.Printf("      %s %s %s\n", Green(IconSuccess), mod.Name, Gray("("+mod.Category+")"))
						hasModule = true
					}
				}
				if !hasModule {
					fmt.Printf("      %s\n", Gray("(no modules)"))
				}
			}

			fmt.Println()
			fmt.Printf("  %sTab Switch  |  ↑↓ Select  |  Shift+↑↓ Move  |  ← Back%s\n", ColorGray, ColorReset)
		}

		key := ReadKey()

		switch key {
		case "TAB":
			mode = 1 - mode // 切换模式
			selected = 0
		case "UP":
			if mode == 1 && selected > 0 {
				selected--
			}
		case "DOWN":
			if mode == 1 && selected < len(platformKeys)-1 {
				selected++
			}
		case "MOVEUP": // J 键上移
			if mode == 1 && selected > 0 {
				platformKeys[selected], platformKeys[selected-1] = platformKeys[selected-1], platformKeys[selected]
				selected--
				cfg.PlatformOrder = platformKeys
				SaveConfig(cfg)
			}
		case "MOVEDOWN": // K 键下移
			if mode == 1 && selected < len(platformKeys)-1 {
				platformKeys[selected], platformKeys[selected+1] = platformKeys[selected+1], platformKeys[selected]
				selected++
				cfg.PlatformOrder = platformKeys
				SaveConfig(cfg)
			}
		case "LEFT", "QUIT", "ESC":
			return
		}
	}
}

// SelectPlatformMenu 平台选择菜单，包含 "All Platforms" 选项
func SelectPlatformMenu(cfg *Config) SelectMenuResult {
	options := []SelectOption{
		{Key: "", Label: "All Platforms"},
	}

	for key, p := range cfg.Platforms {
		options = append(options, SelectOption{
			Key:   key,
			Label: fmt.Sprintf("%s (%s)", p.Name, key),
		})
	}

	return SelectMenu("Select Platform", options)
}

// ConfirmDialog 确认对话框（在当前位置弹出），返回 true 确认，false 取消
func ConfirmDialog(message string) bool {
	HideCursor()
	defer ShowCursor()

	// 在当前位置显示确认框，不清屏
	fmt.Println()
	fmt.Printf("  %s %s\n", Yellow(IconWarning), message)
	fmt.Printf("  %s[Enter] Confirm  |  [ESC] Cancel%s", ColorGray, ColorReset)

	for {
		key := ReadKey()
		switch key {
		case "ENTER", "YES":
			fmt.Println()
			return true
		case "ESC", "NO", "QUIT", "LEFT":
			fmt.Println()
			return false
		}
	}
}

// ModuleDetailResult 模块详情操作结果
type ModuleDetailResult struct {
	Action   string   // "apply", "back", "quit"
	ToSync   []string // 需要同步的平台
	ToRemove []string // 需要删除的平台
}

// ModuleDetailMenu 模块详情菜单（选中=同步，取消=删除）
func ModuleDetailMenu(cfg *Config, mod *Module) ModuleDetailResult {
	// 获取所有平台，并检查当前同步状态
	type platformState struct {
		key      string
		name     string
		selected bool
		synced   bool // 当前是否已同步
	}

	platforms := make([]platformState, 0)
	for key, p := range cfg.Platforms {
		targetDir := ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := targetDir + "/" + mod.GetLinkName(key)
		synced := IsSymlink(targetPath)
		platforms = append(platforms, platformState{
			key:      key,
			name:     p.Name,
			selected: synced, // 默认选中已同步的
			synced:   synced,
		})
	}

	selected := 0

	HideCursor()
	defer ShowCursor()

	for {
		ClearScreen()

		// 模块信息
		fmt.Println()
		fmt.Printf("  %s %s\n", Blue("Module:"), White(mod.Name))
		fmt.Printf("  %s %s\n", Blue("Category:"), mod.Category)
		if mod.Description != "" {
			fmt.Printf("  %s %s\n", Blue("Desc:"), Gray(mod.Description))
		}
		fmt.Printf("  %s %s\n", Blue("Path:"), Gray(mod.Path))
		fmt.Println()

		// 平台列表
		fmt.Printf("  %s\n\n", Blue("Platforms (✓=sync, ✗=remove):"))

		for i, p := range platforms {
			checkbox := "[ ]"
			if p.selected {
				checkbox = Green("[✓]")
			}

			status := ""
			if p.synced && !p.selected {
				status = Yellow(" (will remove)")
			} else if !p.synced && p.selected {
				status = Green(" (will sync)")
			} else if p.synced && p.selected {
				status = Gray(" (synced)")
			}

			if i == selected {
				fmt.Printf("  %s %s %s%s\n", Cyan(IconArrow), checkbox, White(p.name), status)
			} else {
				fmt.Printf("    %s %s%s\n", checkbox, p.name, status)
			}
		}

		fmt.Println()
		fmt.Printf("  %s↑↓ Navigate  |  Space/Enter Toggle  |  ← Back & Apply  |  Q Quit%s\n", ColorGray, ColorReset)

		key := ReadKey()

		switch key {
		case "UP":
			if selected > 0 {
				selected--
			}
		case "DOWN":
			if selected < len(platforms)-1 {
				selected++
			}
		case "SPACE", "ENTER":
			platforms[selected].selected = !platforms[selected].selected
		case "LEFT":
			// 返回时自动应用变更
			var toSync, toRemove []string
			for _, p := range platforms {
				if p.selected && !p.synced {
					toSync = append(toSync, p.key)
				} else if !p.selected && p.synced {
					toRemove = append(toRemove, p.key)
				}
			}
			return ModuleDetailResult{Action: "apply", ToSync: toSync, ToRemove: toRemove}
		case "QUIT":
			return ModuleDetailResult{Action: "quit"}
		}
	}
}

// Spinner 加载动画
type Spinner struct {
	chars   []string
	index   int
	message string
	running bool
	done    chan bool
}

// NewSpinner 创建新的 Spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		chars:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:   0,
		message: message,
		running: false,
		done:    make(chan bool),
	}
}

// Start 启动 Spinner
func (s *Spinner) Start() {
	s.running = true
	go func() {
		for s.running {
			fmt.Printf("\r  %s%s%s %s", ColorBlue, s.chars[s.index], ColorReset, s.message)
			s.index = (s.index + 1) % len(s.chars)
			select {
			case <-s.done:
				return
			case <-time.After(80 * time.Millisecond):
				// 继续动画
			}
		}
	}()
}

// Stop 停止 Spinner
func (s *Spinner) Stop() {
	s.running = false
	s.done <- true
	fmt.Print("\r\033[2K") // 清除行
}

// PrintTable 打印表格（用于 dry-run）
func PrintTable(headers []string, rows [][]string) {
	// 计算列宽
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// 打印表头
	fmt.Print("  ")
	for i, h := range headers {
		fmt.Printf("%s%-*s%s  ", ColorBlue, widths[i], h, ColorReset)
	}
	fmt.Println()

	// 打印分隔线
	fmt.Print("  ")
	for _, w := range widths {
		for j := 0; j < w; j++ {
			fmt.Print("─")
		}
		fmt.Print("  ")
	}
	fmt.Println()

	// 打印数据行
	for _, row := range rows {
		fmt.Print("  ")
		for i, cell := range row {
			if i < len(widths) {
				// 根据内容着色
				colored := cell
				if cell == "CREATE" {
					colored = Green(cell)
				} else if cell == "UPDATE" {
					colored = Yellow(cell)
				} else if cell == "SKIP" {
					colored = Gray(cell)
				} else if cell == "ERROR" {
					colored = Red(cell)
				}
				fmt.Printf("%-*s  ", widths[i]+len(colored)-len(cell), colored)
			}
		}
		fmt.Println()
	}
}

// DefaultsPlatformsMenu 默认平台选择菜单
func DefaultsPlatformsMenu(cfg *Config) {
	// 获取当前选中的默认平台
	selectedDefaults := make(map[string]bool)
	for _, key := range cfg.DefaultPlatforms {
		selectedDefaults[key] = true
	}

	type platformState struct {
		key      string
		name     string
		selected bool
	}

	platforms := make([]platformState, 0)
	for key, p := range cfg.Platforms {
		platforms = append(platforms, platformState{
			key:      key,
			name:     p.Name,
			selected: selectedDefaults[key],
		})
	}

	selected := 0

	HideCursor()
	defer ShowCursor()

	for {
		ClearScreen()
		fmt.Println()
		fmt.Printf("  %s %s\n\n", Blue("Set Default Platforms for Sync"), White("(used when pressing Enter on a module)"))

		// 显示已选数量
		selectedCount := 0
		for _, p := range platforms {
			if p.selected {
				selectedCount++
			}
		}
		if selectedCount > 0 {
			fmt.Printf("  %s Selected: %d platform(s)\n\n", Gray(IconInfo), selectedCount)
		} else {
			fmt.Printf("  %s Selected: %s\n\n", Gray(IconInfo), Gray("(press Enter to sync all)"))
		}

		// 平台列表
		for i, p := range platforms {
			checkbox := "[ ]"
			if p.selected {
				checkbox = Green("[✓]")
			}

			if i == selected {
				fmt.Printf("  %s %s %s\n", Cyan(IconArrow), checkbox, White(p.name))
			} else {
				fmt.Printf("    %s %s\n", checkbox, p.name)
			}
		}

		fmt.Println()
		fmt.Printf("  %s↑↓ Navigate  |  Space/Enter Toggle  |  ← Save & Exit%s\n", ColorGray, ColorReset)

		key := ReadKey()

		switch key {
		case "UP":
			if selected > 0 {
				selected--
			}
		case "DOWN":
			if selected < len(platforms)-1 {
				selected++
			}
		case "SPACE", "ENTER":
			platforms[selected].selected = !platforms[selected].selected
		case "LEFT":
			// 保存设置
			var defaults []string
			for _, p := range platforms {
				if p.selected {
					defaults = append(defaults, p.key)
				}
			}
			cfg.DefaultPlatforms = defaults
			SaveConfig(cfg)
			return
		case "QUIT":
			return
		}
	}
}
