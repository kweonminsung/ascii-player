package cmd

import (
	"regexp"

	"github.com/spf13/cobra"
)

var youtubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "YouTube video processing commands",
	Long:  `Download and convert YouTube videos to ASCII art animations.`,
}

func isValidYouTubeURL(url string) bool {
	// Simple regex to validate YouTube URLs
	patterns := []string{
		`^https?://(www\.)?youtube\.com/watch\?v=[\w-]+`,
		`^https?://youtu\.be/[\w-]+`,
		`^https?://(www\.)?youtube\.com/embed/[\w-]+`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}
	return false
}

func init() {

}
