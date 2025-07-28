package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

const cacheDirName = "console-cinema"

// GetCacheDir returns the path to the application's cache directory.
func GetCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to temp dir if user cache dir is not available
		userCacheDir = os.TempDir()
	}
	cachePath := filepath.Join(userCacheDir, cacheDirName)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}
	return cachePath, nil
}

// CreateTemp creates a new temporary file in the cache directory.
func CreateTemp(pattern string) (*os.File, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}
	return os.CreateTemp(cacheDir, pattern)
}

// ClearCache removes all files from the cache directory.
func ClearCache() error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}
	dir, err := os.Open(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, nothing to clear
		}
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(cacheDir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
