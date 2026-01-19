package lib

import (
	"errors"
	"testing"
)

func TestModuleNotFoundError(t *testing.T) {
	err := &ModuleNotFoundError{Name: "test-module"}

	expected := "module not found: test-module"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}

	if !IsModuleNotFound(err) {
		t.Error("IsModuleNotFound should return true")
	}

	if IsModuleNotFound(errors.New("other error")) {
		t.Error("IsModuleNotFound should return false for other errors")
	}
}

func TestPlatformNotFoundError(t *testing.T) {
	err := &PlatformNotFoundError{Name: "unknown-platform"}

	expected := "unknown platform: unknown-platform"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}

	if !IsPlatformNotFound(err) {
		t.Error("IsPlatformNotFound should return true")
	}

	if IsPlatformNotFound(errors.New("other error")) {
		t.Error("IsPlatformNotFound should return false for other errors")
	}
}

func TestConfigNotFoundError(t *testing.T) {
	err := &ConfigNotFoundError{Path: "/path/to/config"}

	expected := "config file not found: /path/to/config"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestSymlinkError(t *testing.T) {
	err := &SymlinkError{
		Op:     "create",
		Path:   "/path/to/link",
		Reason: "permission denied",
	}

	expected := "symlink create failed for /path/to/link: permission denied"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
