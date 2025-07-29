package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	appDirName = "console-cinema"
	logDirName = "logs"
	tmpDirName = "tmp"
)

// getAppDir returns the path to the application's base directory.
func getAppDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to temp dir if user cache dir is not available
		userCacheDir = os.TempDir()
	}
	appPath := filepath.Join(userCacheDir, appDirName)
	if err := os.MkdirAll(appPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}
	return appPath, nil
}

// GetLogDir returns the path to the application's log directory.
func GetLogDir() (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	logPath := filepath.Join(appDir, logDirName)
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}
	return logPath, nil
}

// GetTmpDir returns the path to the application's temporary directory.
func GetTmpDir() (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	tmpPath := filepath.Join(appDir, tmpDirName)
	if err := os.MkdirAll(tmpPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create tmp directory: %w", err)
	}
	return tmpPath, nil
}

// CreateTemp creates a new temporary file in the tmp directory.
func CreateTemp(pattern string) (*os.File, error) {
	tmpDir, err := GetTmpDir()
	if err != nil {
		return nil, err
	}
	return os.CreateTemp(tmpDir, pattern)
}

// ClearTmpDir removes all files from the tmp directory.
func ClearTmpDir() error {
	tmpDir, err := GetTmpDir()
	if err != nil {
		return err
	}
	dir, err := os.Open(tmpDir)
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
		err = os.RemoveAll(filepath.Join(tmpDir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
