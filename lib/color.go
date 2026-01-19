package lib

// ANSI 颜色代码
const (
	// 基础颜色
	ColorReset  = "\033[0m"
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
	ColorPurple  = "\033[0;35m"
	ColorMagenta = "\033[1;35m"
	ColorCyan    = "\033[0;36m"
	ColorGray   = "\033[0;90m"
	ColorWhite  = "\033[1;37m"

	// 图标
	IconSuccess = "✓"
	IconError   = "✗"
	IconArrow   = "▶"
	IconDot     = "●"
	IconLink    = "⟶"
	IconFolder  = "📁"
	IconFile    = "📄"
	IconWarning = "⚠"
	IconInfo    = "ℹ"
)

// 快捷函数
func Green(s string) string  { return ColorGreen + s + ColorReset }
func Red(s string) string    { return ColorRed + s + ColorReset }
func Yellow(s string) string { return ColorYellow + s + ColorReset }
func Blue(s string) string   { return ColorBlue + s + ColorReset }
func Purple(s string) string  { return ColorPurple + s + ColorReset }
func Magenta(s string) string { return ColorMagenta + s + ColorReset }
func Cyan(s string) string    { return ColorCyan + s + ColorReset }
func Gray(s string) string   { return ColorGray + s + ColorReset }
func White(s string) string  { return ColorWhite + s + ColorReset }

// 状态输出
func Success(msg string) string { return Green(IconSuccess) + " " + msg }
func Error(msg string) string   { return Red(IconError) + " " + msg }
func Warning(msg string) string { return Yellow(IconWarning) + " " + msg }
func Info(msg string) string    { return Blue(IconInfo) + " " + msg }
