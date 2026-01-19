package lib

import "fmt"

// 自定义错误类型

// ModuleNotFoundError 模块未找到错误
type ModuleNotFoundError struct {
	Name string
}

func (e *ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module not found: %s", e.Name)
}

// PlatformNotFoundError 平台未找到错误
type PlatformNotFoundError struct {
	Name string
}

func (e *PlatformNotFoundError) Error() string {
	return fmt.Sprintf("unknown platform: %s", e.Name)
}

// ConfigNotFoundError 配置文件未找到错误
type ConfigNotFoundError struct {
	Path string
}

func (e *ConfigNotFoundError) Error() string {
	return fmt.Sprintf("config file not found: %s", e.Path)
}

// SymlinkError 软链接操作错误
type SymlinkError struct {
	Op     string // "create", "remove", "read"
	Path   string
	Reason string
}

func (e *SymlinkError) Error() string {
	return fmt.Sprintf("symlink %s failed for %s: %s", e.Op, e.Path, e.Reason)
}

// IsModuleNotFound 检查是否为模块未找到错误
func IsModuleNotFound(err error) bool {
	_, ok := err.(*ModuleNotFoundError)
	return ok
}

// IsPlatformNotFound 检查是否为平台未找到错误
func IsPlatformNotFound(err error) bool {
	_, ok := err.(*PlatformNotFoundError)
	return ok
}
