package ytdlp

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

//go:embed yt-dlp.exe
var ytDlpBinaryWin []byte

//go:embed yt-dlp_macos
var ytDlpBinaryMac []byte

//go:embed yt-dlp
var ytDlpBinary []byte

// GetExecutablePath creates a temporary executable file for yt-dlp and returns its path.
// The caller is responsible for removing the file when done.
func GetExecutablePath() (string, error) {
	var binary []byte
	var fileName string

	// Determine the correct binary and filename based on the OS
	switch runtime.GOOS {
	case "windows":
		binary = ytDlpBinaryWin
		fileName = "yt-dlp.exe"
	case "darwin":
		binary = ytDlpBinaryMac
		fileName = "yt-dlp_macos"
	case "linux":
		binary = ytDlpBinary
		fileName = "yt-dlp"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create a temporary file. It's created with default permissions (e.g., 0600).
	tmpFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for yt-dlp: %w", err)
	}
	defer tmpFile.Close() // Ensure the file handle is closed

	// Write the binary data to the temporary file.
	// Note: os.WriteFile is simpler but we already have a file handle.
	if _, err := tmpFile.Write(binary); err != nil {
		os.Remove(tmpFile.Name()) // Clean up on error
		return "", fmt.Errorf("failed to write yt-dlp binary to temp file: %w", err)
	}

	// IMPORTANT: Explicitly set the executable permission on the file.
	// This is the key fix.
	perm := fs.FileMode(0755) // rwxr-xr-x
	if err := os.Chmod(tmpFile.Name(), perm); err != nil {
		os.Remove(tmpFile.Name()) // Clean up on error
		return "", fmt.Errorf("failed to set executable permission on temp file: %w", err)
	}
	
	// The file handle must be closed before the file can be executed by another process.
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to close temp file for yt-dlp: %w", err)
	}


	// Return the path of the now-executable file
	return tmpFile.Name(), nil
}
