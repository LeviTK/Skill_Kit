package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateSymlink 创建软链接
func CreateSymlink(source, target string, isProject bool) error {
	// 展开路径
	source = ResolvePath(source)
	target = ResolvePath(target)

	// 确保源存在
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source not found: %s", source)
	}

	// 确保目标目录存在
	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 检查目标是否存在
	if info, err := os.Lstat(target); err == nil {
		// 目标已存在
		if info.Mode()&os.ModeSymlink != 0 {
			// 是软链接，删除后重建
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove existing symlink: %v", err)
			}
		} else if info.IsDir() {
			// 是实体目录，报错保护
			return fmt.Errorf("target is a real directory (not symlink): %s", target)
		} else {
			// 是实体文件，报错保护
			return fmt.Errorf("target is a real file (not symlink): %s", target)
		}
	}

	// 创建软链接
	return os.Symlink(source, target)
}

// RemoveSymlink 移除软链接
func RemoveSymlink(target string) error {
	target = ResolvePath(target)

	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil // 不存在视为成功
	}
	if err != nil {
		return err
	}

	// 只删除软链接
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("target is not a symlink: %s", target)
	}

	return os.Remove(target)
}

// IsSymlink 检查是否为软链接
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// ReadSymlink 读取软链接目标
func ReadSymlink(path string) (string, error) {
	return os.Readlink(path)
}

// ResolvePath 展开路径中的 ~ 和相对路径
func ResolvePath(parts ...string) string {
	path := filepath.Join(parts...)

	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
