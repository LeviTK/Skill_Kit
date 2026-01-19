package main

import (
	"fmt"
	"os"

	"skillkit/lib"
)

func main() {
	if len(os.Args) < 2 {
		// 无参数时显示交互式菜单（循环）
		for {
			cmd := lib.InteractiveMenu()
			if cmd == "quit" {
				lib.ClearScreen()
				os.Exit(0)
			}
			// 执行选中的命令（交互模式）
			handleInteractiveCommand(cmd)
		}
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "--help", "-h", "help":
		lib.ShowHelp()
	case "--version", "-v", "-V", "version":
		lib.ShowVersion()
	default:
		handleCommand(cmd, args)
	}
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "use":
		handleUse(args)
	case "list":
		handleList(args)
	case "platforms":
		handlePlatforms(args)
	case "info":
		handleInfo(args)
	case "remove":
		handleRemove(args)
	case "sync":
		handleSync(args)
	case "status":
		handleStatus(args)
	case "init":
		handleInit(args)
	default:
		fmt.Printf("%s Unknown command: %s\n", lib.Red(lib.IconError), cmd)
		fmt.Println("Run 'sk --help' for usage information.")
		os.Exit(1)
	}
}

// handleInteractiveCommand 处理交互模式下的命令，返回 true 继续循环，false 退出
func handleInteractiveCommand(cmd string) bool {
	cfg, err := lib.LoadConfig()
	if err != nil {
		lib.ClearScreen()
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		return true
	}

	switch cmd {
	case "use":
		return handleInteractiveUse(cfg)

	case "list":
		handleInteractiveList(cfg)

	case "platforms":
		lib.ClearScreen()
		handlePlatforms([]string{})
		lib.WaitForKey()

	case "info":
		return handleInteractiveInfo(cfg)

	case "remove":
		return handleInteractiveRemove(cfg)

	case "defaults":
		handleInteractiveDefaults(cfg)

	case "status":
		lib.ClearScreen()
		handleStatus([]string{})
		lib.WaitForKey()

	case "init":
		lib.ClearScreen()
		handleInit([]string{})
		lib.WaitForKey()

	default:
		lib.ClearScreen()
		fmt.Printf("%s Unknown command: %s\n", lib.Red(lib.IconError), cmd)
	}

	return true
}

// handleInteractiveUse 交互式 use 命令，支持回退和平台多选
func handleInteractiveUse(cfg *lib.Config) bool {
	for {
		// 使用增强版模块列表
		result := lib.ModuleListMenu(cfg)

		switch result.Action {
		case "back":
			return true // 返回主菜单
		case "quit":
			return false
		case "sync_default":
			// 同步单个模块到默认平台
			mod := result.Module
			syncModuleToDefaultPlatforms(cfg, mod)
			// 不返回主页面，继续留在 Use 子页面
		case "sync_all_default":
			// 全选：同步所有模块到默认平台
			syncAllModulesToDefaultPlatforms(cfg, result.Modules)
			// 不返回主页面，继续留在 Use 子页面
		case "detail":
			// 右键进入详情页手动管理
			if !handleModuleDetail(cfg, result.Module) {
				return false
			}
		}
	}
}

// syncModuleToDefaultPlatforms 同步单个模块到默认平台
func syncModuleToDefaultPlatforms(cfg *lib.Config, mod *lib.Module) bool {
	targetPlatforms := getTargetPlatforms(cfg)
	if len(targetPlatforms) == 0 {
		return false
	}

	msg := fmt.Sprintf("Sync '%s' to %d default platform(s)?", mod.Name, len(targetPlatforms))
	if !lib.ConfirmDialog(msg) {
		return false
	}

	lib.ClearScreen()
	fmt.Println()
	for platKey, p := range targetPlatforms {
		ln := mod.GetLinkName(platKey)
		targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := targetDir + "/" + ln

		err := lib.CreateSymlink(mod.Path, targetPath, false)
		if err != nil {
			fmt.Printf("  %s %s → %s: %v\n", lib.Red(lib.IconError), mod.Name, platKey, err)
		} else {
			fmt.Printf("  %s %s %s %s\n", lib.Green(lib.IconSuccess), mod.Name, lib.Cyan(lib.IconLink), p.Name)
		}
	}
	return true
}

// syncAllModulesToDefaultPlatforms 同步所有模块到默认平台
func syncAllModulesToDefaultPlatforms(cfg *lib.Config, modules []*lib.Module) bool {
	targetPlatforms := getTargetPlatforms(cfg)
	if len(targetPlatforms) == 0 {
		return false
	}

	msg := fmt.Sprintf("Sync ALL %d modules to %d default platform(s)?", len(modules), len(targetPlatforms))
	if !lib.ConfirmDialog(msg) {
		return false
	}

	lib.ClearScreen()
	fmt.Println()
	success := 0
	failed := 0
	for _, mod := range modules {
		for platKey, p := range targetPlatforms {
			ln := mod.GetLinkName(platKey)
			targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
			targetPath := targetDir + "/" + ln

			err := lib.CreateSymlink(mod.Path, targetPath, false)
			if err != nil {
				fmt.Printf("  %s %s → %s: %v\n", lib.Red(lib.IconError), mod.Name, platKey, err)
				failed++
			} else {
				fmt.Printf("  %s %s %s %s\n", lib.Green(lib.IconSuccess), mod.Name, lib.Cyan(lib.IconLink), p.Name)
				success++
			}
		}
	}
	fmt.Println()
	fmt.Printf("  %s: %d  %s: %d\n", lib.Green("Success"), success, lib.Red("Failed"), failed)
	lib.WaitForKey()
	return true
}

// getTargetPlatforms 获取目标平台（默认平台或全部平台）
func getTargetPlatforms(cfg *lib.Config) map[string]lib.Platform {
	if len(cfg.DefaultPlatforms) > 0 {
		platforms := make(map[string]lib.Platform)
		for _, key := range cfg.DefaultPlatforms {
			if p, ok := cfg.Platforms[key]; ok {
				platforms[key] = p
			}
		}
		return platforms
	}
	return cfg.Platforms
}

// handleModuleDetail 处理模块详情页
func handleModuleDetail(cfg *lib.Config, mod *lib.Module) bool {
	for {
		detailResult := lib.ModuleDetailMenu(cfg, mod)

		switch detailResult.Action {
		case "quit":
			return false
		case "apply":
			// 有变更时执行
			if len(detailResult.ToSync) > 0 || len(detailResult.ToRemove) > 0 {
				// 二次确认
				msg := fmt.Sprintf("Apply changes? (+%d sync, -%d remove)", len(detailResult.ToSync), len(detailResult.ToRemove))
				if lib.ConfirmDialog(msg) {
					lib.ClearScreen()
					fmt.Println()

					// 执行同步
					for _, platKey := range detailResult.ToSync {
						p := cfg.Platforms[platKey]
						ln := mod.GetLinkName(platKey)
						targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
						targetPath := targetDir + "/" + ln

						err := lib.CreateSymlink(mod.Path, targetPath, false)
						if err != nil {
							fmt.Printf("  %s %s → %s: %v\n", lib.Red(lib.IconError), mod.Name, platKey, err)
						} else {
							fmt.Printf("  %s %s %s %s\n", lib.Green(lib.IconSuccess), mod.Name, lib.Cyan(lib.IconLink), p.Name)
						}
					}

					// 执行删除
					for _, platKey := range detailResult.ToRemove {
						p := cfg.Platforms[platKey]
						ln := mod.GetLinkName(platKey)
						targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
						targetPath := targetDir + "/" + ln

						err := lib.RemoveSymlink(targetPath)
						if err != nil {
							fmt.Printf("  %s Remove %s from %s: %v\n", lib.Red(lib.IconError), mod.Name, platKey, err)
						} else {
							fmt.Printf("  %s Removed from %s\n", lib.Yellow(lib.IconWarning), p.Name)
						}
					}
					return true
				}
			}
			// 无变更或取消确认，返回模块列表
			return true
		}
	}
}

// handleInteractiveInfo 交互式 info 命令
func handleInteractiveInfo(cfg *lib.Config) bool {
	modResult := lib.SelectModuleMenu(cfg)
	if modResult.Cancel {
		return false
	}
	if modResult.Back {
		return true
	}
	lib.ClearScreen()
	handleInfo([]string{modResult.Key})
	return true
}

// handleInteractiveRemove 交互式 remove 命令，支持回退
func handleInteractiveRemove(cfg *lib.Config) bool {
	for {
		modResult := lib.SelectModuleMenu(cfg)
		if modResult.Cancel {
			return false
		}
		if modResult.Back {
			return true
		}

		platResult := lib.SelectPlatformMenu(cfg)
		if platResult.Cancel {
			return false
		}
		if platResult.Back {
			continue
		}

		lib.ClearScreen()
		args := []string{modResult.Key}
		if platResult.Key != "" {
			args = append(args, platResult.Key)
		}
		handleRemove(args)
		return true
	}
}

func handleUse(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sk use <module> [platform] [--global|--project] [--as <name>]")
		os.Exit(1)
	}

	module := args[0]
	platform := ""
	scope := "global"
	linkName := ""
	dryRun := false

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--global":
			scope = "global"
		case "--project":
			scope = "project"
		case "--as":
			if i+1 < len(args) {
				linkName = args[i+1]
				i++
			}
		case "--dry-run":
			dryRun = true
		default:
			if platform == "" && !hasPrefix(args[i], "--") {
				platform = args[i]
			}
		}
	}

	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	mod, err := lib.FindModule(cfg, module)
	if err != nil {
		fmt.Printf("%s %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	platforms := cfg.Platforms
	if platform != "" {
		if p, ok := cfg.Platforms[platform]; ok {
			platforms = map[string]lib.Platform{platform: p}
		} else {
			fmt.Printf("%s Unknown platform: %s\n", lib.Red(lib.IconError), platform)
			os.Exit(1)
		}
	}

	if dryRun {
		// 表格化输出
		fmt.Printf("\n%s Preview: %s → %d platform(s)\n\n", lib.Blue(lib.IconInfo), module, len(platforms))
		headers := []string{"Module", "Platform", "Target Path", "Action"}
		var rows [][]string

		for name, p := range platforms {
			ln := linkName
			if ln == "" {
				ln = mod.GetLinkName(name)
			}
			targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
			targetPath := targetDir + "/" + ln

			action := "CREATE"
			if lib.IsSymlink(targetPath) {
				action = "UPDATE"
			}
			rows = append(rows, []string{module, name, targetPath, action})
		}

		lib.PrintTable(headers, rows)
		fmt.Println()
	} else {
		fmt.Println()
		for name, p := range platforms {
			ln := linkName
			if ln == "" {
				ln = mod.GetLinkName(name)
			}

			targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
			targetPath := targetDir + "/" + ln

			err := lib.CreateSymlink(mod.Path, targetPath, scope == "project")
			if err != nil {
				fmt.Printf("  %s %s → %s: %v\n", lib.Red(lib.IconError), module, name, err)
			} else {
				fmt.Printf("  %s %s %s %s (%s)\n", lib.Green(lib.IconSuccess), module, lib.Cyan(lib.IconLink), targetPath, lib.Gray(name))
			}
		}
		fmt.Println()
	}
}

func handleList(args []string) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	modules, err := lib.ListModules(cfg)
	if err != nil {
		fmt.Printf("%s Error listing modules: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	if len(modules) == 0 {
		fmt.Printf("\n%s No modules found in ~/.config/agent/\n", lib.Yellow(lib.IconWarning))
		fmt.Println("  Run 'lt init' to initialize the repository.")
		fmt.Println()
		return
	}

	fmt.Printf("\n%s Modules:\n\n", lib.Blue(lib.IconFolder))
	for _, mod := range modules {
		status := lib.GetLinkStatus(cfg, mod)
		fmt.Printf("  %s %s %s\n", lib.Cyan(lib.IconArrow), lib.White(mod.Name), lib.Gray("("+mod.Category+")"))
		if len(status) > 0 {
			for i, s := range status {
				prefix := "  │   ├──"
				if i == len(status)-1 {
					prefix = "  │   └──"
				}
				fmt.Printf("%s %s\n", lib.Gray(prefix), s)
			}
		} else {
			fmt.Printf("  %s %s\n", lib.Gray("│   └──"), lib.Gray("(not linked)"))
		}
	}
	fmt.Println()
}

func handlePlatforms(args []string) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	fmt.Printf("\n%s Registered Platforms (%d):\n\n", lib.Blue(lib.IconInfo), len(cfg.Platforms))
	for key, p := range cfg.Platforms {
		fmt.Printf("  %s %s %s\n", lib.Cyan(lib.IconArrow), lib.White(p.Name), lib.Gray("("+key+")"))
		fmt.Printf("      Project: %s\n", lib.Gray(p.Project+p.SkillDir+"/"))
		fmt.Printf("      Global:  %s\n", lib.Gray(p.Global+p.SkillDir+"/"))
		fmt.Println()
	}
}

func handleInfo(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sk info <module>")
		os.Exit(1)
	}

	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	mod, err := lib.FindModule(cfg, args[0])
	if err != nil {
		fmt.Printf("%s %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("  %s %s\n", lib.Blue("Module:"), lib.White(mod.Name))
	fmt.Printf("  %s %s\n", lib.Blue("Category:"), mod.Category)
	fmt.Printf("  %s %s\n", lib.Blue("Path:"), mod.Path)

	if len(mod.Aliases) > 0 {
		fmt.Printf("  %s\n", lib.Blue("Aliases:"))
		for platform, alias := range mod.Aliases {
			fmt.Printf("    %s %s %s\n", platform, lib.Cyan(lib.IconLink), alias)
		}
	}
	fmt.Println()
}

func handleRemove(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sk remove <module> [platform]")
		os.Exit(1)
	}

	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	module := args[0]
	platform := ""
	if len(args) > 1 {
		platform = args[1]
	}

	mod, err := lib.FindModule(cfg, module)
	if err != nil {
		fmt.Printf("%s %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	platforms := cfg.Platforms
	if platform != "" {
		if p, ok := cfg.Platforms[platform]; ok {
			platforms = map[string]lib.Platform{platform: p}
		} else {
			fmt.Printf("%s Unknown platform: %s\n", lib.Red(lib.IconError), platform)
			os.Exit(1)
		}
	}

	fmt.Println()
	for name, p := range platforms {
		ln := mod.GetLinkName(name)
		targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
		targetPath := targetDir + "/" + ln

		err := lib.RemoveSymlink(targetPath)
		if err != nil {
			fmt.Printf("  %s %s from %s: %v\n", lib.Red(lib.IconError), module, name, err)
		} else {
			fmt.Printf("  %s Removed %s from %s\n", lib.Green(lib.IconSuccess), module, name)
		}
	}
	fmt.Println()
}

func handleSync(args []string) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	modules, err := lib.ListModules(cfg)
	if err != nil {
		fmt.Printf("%s Error listing modules: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	if len(modules) == 0 {
		fmt.Printf("\n%s No modules to sync.\n\n", lib.Yellow(lib.IconWarning))
		return
	}

	dryRun := false
	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
		}
	}

	totalLinks := len(modules) * len(cfg.Platforms)

	if dryRun {
		fmt.Printf("\n%s Preview: %d modules → %d platforms = %d symlinks\n\n",
			lib.Blue(lib.IconInfo), len(modules), len(cfg.Platforms), totalLinks)

		headers := []string{"Module", "Platform", "Target Path", "Action"}
		var rows [][]string

		for _, mod := range modules {
			for name, p := range cfg.Platforms {
				ln := mod.GetLinkName(name)
				targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
				targetPath := targetDir + "/" + ln

				action := "CREATE"
				if lib.IsSymlink(targetPath) {
					action = "UPDATE"
				}
				rows = append(rows, []string{mod.Name, name, targetPath, action})
			}
		}

		lib.PrintTable(headers, rows)
		fmt.Println()
	} else {
		fmt.Printf("\n%s Syncing %d modules to %d platforms...\n\n",
			lib.Blue(lib.IconInfo), len(modules), len(cfg.Platforms))

		success := 0
		failed := 0

		for _, mod := range modules {
			for name, p := range cfg.Platforms {
				ln := mod.GetLinkName(name)
				targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
				targetPath := targetDir + "/" + ln

				err := lib.CreateSymlink(mod.Path, targetPath, false)
				if err != nil {
					fmt.Printf("  %s %s → %s: %v\n", lib.Red(lib.IconError), mod.Name, name, err)
					failed++
				} else {
					fmt.Printf("  %s %s %s %s\n", lib.Green(lib.IconSuccess), mod.Name, lib.Cyan(lib.IconLink), name)
					success++
				}
			}
		}

		fmt.Println()
		fmt.Printf("  %s: %d  %s: %d\n\n", lib.Green("Success"), success, lib.Red("Failed"), failed)
	}
}

func handleStatus(args []string) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		fmt.Printf("%s Error loading config: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	modules, err := lib.ListModules(cfg)
	if err != nil {
		fmt.Printf("%s Error listing modules: %v\n", lib.Red(lib.IconError), err)
		os.Exit(1)
	}

	if len(modules) == 0 {
		fmt.Printf("\n%s No modules found.\n\n", lib.Yellow(lib.IconWarning))
		return
	}

	fmt.Printf("\n%s Health Check\n\n", lib.Blue(lib.IconInfo))

	healthy := 0
	broken := 0
	missing := 0

	for _, mod := range modules {
		for name, p := range cfg.Platforms {
			ln := mod.GetLinkName(name)
			targetDir := lib.ResolvePath(p.Global, p.GetCategoryDir(mod.Category))
			targetPath := targetDir + "/" + ln

			if lib.IsSymlink(targetPath) {
				realPath, _ := lib.ReadSymlink(targetPath)
				if realPath == mod.Path {
					healthy++
				} else {
					fmt.Printf("  %s %s → %s: broken (points to %s)\n",
						lib.Red(lib.IconError), mod.Name, name, realPath)
					broken++
				}
			} else if _, err := os.Stat(targetPath); err == nil {
				fmt.Printf("  %s %s → %s: blocked by real file/dir\n",
					lib.Yellow(lib.IconWarning), mod.Name, name)
				broken++
			} else {
				missing++
			}
		}
	}

	fmt.Println()
	fmt.Printf("  %s Healthy: %d  %s Broken: %d  %s Not linked: %d\n\n",
		lib.Green(lib.IconSuccess), healthy,
		lib.Red(lib.IconError), broken,
		lib.Gray("○"), missing)

	if broken > 0 {
		fmt.Printf("  %s Run 'sk sync' to fix broken links.\n\n", lib.Blue(lib.IconInfo))
	}
}

func handleInit(args []string) {
	home, _ := os.UserHomeDir()
	repoPath := home + "/.config/agent"

	// 创建目录结构
	dirs := []string{
		repoPath,
		repoPath + "/skill",
		repoPath + "/agent",
	}

	fmt.Println()
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("  %s Failed to create %s: %v\n", lib.Red(lib.IconError), dir, err)
		} else {
			fmt.Printf("  %s Created %s\n", lib.Green(lib.IconSuccess), dir)
		}
	}

	// 复制 platforms.toml 如果不存在
	configPath := repoPath + "/platforms.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 读取当前目录的 platforms.toml
		execPath, _ := os.Executable()
		srcPath := execPath[:len(execPath)-len("bin/skillkit")] + "platforms.toml"

		if data, err := os.ReadFile(srcPath); err == nil {
			if err := os.WriteFile(configPath, data, 0644); err == nil {
				fmt.Printf("  %s Created %s\n", lib.Green(lib.IconSuccess), configPath)
			}
		}
	} else {
		fmt.Printf("  %s %s already exists\n", lib.Yellow(lib.IconWarning), configPath)
	}

	fmt.Println()
	fmt.Printf("  %s Repository initialized at %s\n\n", lib.Green(lib.IconSuccess), repoPath)
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// handleInteractiveDefaults 交互式设置默认平台
func handleInteractiveDefaults(cfg *lib.Config) {
	lib.DefaultsPlatformsMenu(cfg)
}

// handleInteractiveList 交互式列表，支持模式切换
func handleInteractiveList(cfg *lib.Config) {
	lib.ListMenu(cfg)
}
