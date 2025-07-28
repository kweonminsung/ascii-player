package ytdlp

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

//go:embed yt-dlp
var ytDlpBinary []byte

//go:embed yt-dlp.exe
var ytDlpBinaryWin []byte

// GetExecutablePath creates a temporary executable file for yt-dlp and returns its path.
// The caller is responsible for removing the file when done.
func GetExecutablePath() (string, error) {
	var binary []byte
	var fileName string
	if runtime.GOOS == "windows" {
		binary = ytDlpBinaryWin
		fileName = "yt-dlp.exe"
	} else {
		binary = ytDlpBinary
		fileName = "yt-dlp"
	}

	tmpFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for yt-dlp: %w", err)
	}

	perm := fs.FileMode(0755)
	if err := os.WriteFile(tmpFile.Name(), binary, perm); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write yt-dlp binary to temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to close temp file for yt-dlp: %w", err)
	}

	return tmpFile.Name(), nil
}
