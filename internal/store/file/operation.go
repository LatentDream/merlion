package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
)

func validatePath(path string) (string, error) {
	// Handle tilde expansion
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	if strings.Contains(absPath, "..") {
		return "", fmt.Errorf("path contains directory traversal: %s", absPath)
	}
	return absPath, nil
}

func ensureDirectoryExists(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error checking directory: %w", err)
	}

	return os.MkdirAll(path, 0o755)
}

func getFileTimes(path string) (createdAt, updatedAt time.Time, err error) {
	t, err := times.Stat(path)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	updatedAt = t.ModTime()

	// BirthTime() returns creation time on systems that support it
	// Falls back to ModTime() if not available
	if t.HasBirthTime() {
		createdAt = t.BirthTime()
	} else {
		createdAt = updatedAt
	}

	return createdAt, updatedAt, nil
}

func getBoolOrDefault(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
