package lib

// Command 命令定义
type Command struct {
	Name        string
	Description string
	Usage       string
}

// Commands 命令注册表
var Commands = []Command{
	{"use", "Distribute a skill/agent to platform(s)", "sk use <module> [platform]"},
	{"list", "List all modules and their link status", "sk list"},
	{"platforms", "Show registered platforms", "sk platforms"},
	{"info", "Show module details and aliases", "sk info <module>"},
	{"remove", "Remove symlinks for a module", "sk remove <module> [platform]"},
	{"status", "Health check: detect broken symlinks", "sk status"},
	{"init", "Initialize the agent repository", "sk init"},
	{"help", "Show help message", "sk -h"},
	{"version", "Show version", "sk -v"},
}

// GetCommandByName 根据名称获取命令
func GetCommandByName(name string) *Command {
	for _, cmd := range Commands {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

// GetCommandNames 获取所有命令名称（用于补全）
func GetCommandNames() []string {
	names := make([]string, len(Commands))
	for i, cmd := range Commands {
		names[i] = cmd.Name
	}
	return names
}
