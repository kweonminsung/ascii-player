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
