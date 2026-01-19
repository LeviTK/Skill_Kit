package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "home directory expansion",
			input:    []string{"~/.config"},
			expected: filepath.Join(home, ".config"),
		},
		{
			name:     "multiple parts",
			input:    []string{"~", ".config", "agent"},
			expected: filepath.Join(home, ".config", "agent"),
		},
		{
			name:     "absolute path unchanged",
			input:    []string{"/usr/local/bin"},
			expected: "/usr/local/bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.input...)
			if result != tt.expected {
				t.Errorf("ResolvePath(%v) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateAndRemoveSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	targetLink := filepath.Join(tmpDir, "target", "link")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	if err := CreateSymlink(sourceDir, targetLink, false); err != nil {
		t.Fatalf("CreateSymlink failed: %v", err)
	}

	if !IsSymlink(targetLink) {
		t.Error("expected target to be a symlink")
	}

	realPath, err := ReadSymlink(targetLink)
	if err != nil {
		t.Fatalf("ReadSymlink failed: %v", err)
	}
	if realPath != sourceDir {
		t.Errorf("symlink points to %s, expected %s", realPath, sourceDir)
	}

	if err := RemoveSymlink(targetLink); err != nil {
		t.Fatalf("RemoveSymlink failed: %v", err)
	}

	if IsSymlink(targetLink) {
		t.Error("symlink should have been removed")
	}
}

func TestCreateSymlinkProtectsRealFiles(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	targetFile := filepath.Join(tmpDir, "realfile")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	if err := os.WriteFile(targetFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create real file: %v", err)
	}

	err := CreateSymlink(sourceDir, targetFile, false)
	if err == nil {
		t.Error("expected error when target is a real file")
	}
}

func TestRemoveSymlinkProtectsRealFiles(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "realfile")

	if err := os.WriteFile(targetFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create real file: %v", err)
	}

	err := RemoveSymlink(targetFile)
	if err == nil {
		t.Error("expected error when target is not a symlink")
	}
}

func TestIsSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	if IsSymlink(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("nonexistent path should not be a symlink")
	}

	realFile := filepath.Join(tmpDir, "realfile")
	if err := os.WriteFile(realFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if IsSymlink(realFile) {
		t.Error("real file should not be reported as symlink")
	}
}
