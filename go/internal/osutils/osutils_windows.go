//go:build windows

package osutils

import (
	"os"
	"strings"
	"syscall"
)

func IsExecutable(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil {
		return false
	}
	mode := stat.Mode()
	if !mode.IsRegular() {
		return false
	}
	if len(path) < 4 {
		return false
	}
	return strings.EqualFold(path[len(path)-4:], ".exe")
}

func GetRealExecutable() (string, error) {
	return os.Executable()
}

// On Windows, only the 0o200 bit (owner writable) of mode is used;
// it controls whether the file's read-only attribute is set or cleared.
// The other bits are currently unused.
// https://pkg.go.dev/os#Chmod
//
// Let's use this as an optimization.
func AddPermissions(path string, perms os.FileMode) error {
	if perms&syscall.S_IWRITE == 0 {
		return nil
	}
	return addPermissions(path, perms)
}

func RemovePermissions(path string, perms os.FileMode) error {
	if perms&syscall.S_IWRITE == 0 {
		return nil
	}
	return removePermissions(path, perms)
}
