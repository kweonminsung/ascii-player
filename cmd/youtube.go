package cmd

import (
	"fmt"
	"regexp"

	"github.com/kweonminsung/ascii-player/pkg/tui"
	"github.com/spf13/cobra"
)

var youtubeCmd = &cobra.Command{
	Use:   "youtube [url]",
	Short: "Play ASCII/Pixel animations from a YouTube video",
	Long:  `Play ASCII/Pixel animations from a specified YouTube video URL. The video will be streamed and converted to ASCII art or pixel art in real-time and displayed in the terminal. Supports options for mode, FPS, looping, and resolution.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var youtubeURL string
		if len(args) > 0 {
			youtubeURL = args[0]
		} else {
			fmt.Println("Error: Please specify a YouTube URL to play")
			fmt.Println("Usage: ascii-player youtube <youtube_url>")
			fmt.Println("Example: ascii-player youtube https://www.youtube.com/watch?v=dQw4w9WgXcQ")
			return
		}

		// Check if it's a valid YouTube URL
		if !isValidYouTubeURL(youtubeURL) {
			fmt.Println("Error: Invalid YouTube URL provided")
			fmt.Println("Supported formats:")
			fmt.Println("  - https://www.youtube.com/watch?v=VIDEO_ID")
			fmt.Println("  - https://youtu.be/VIDEO_ID")
			fmt.Println("  - https://www.youtube.com/embed/VIDEO_ID")
			return
		}

		fps, _ := cmd.Flags().GetInt("fps")
		loop, _ := cmd.Flags().GetBool("loop")
		resolution, _ := cmd.Flags().GetString("resolution")
		color, _ := cmd.Flags().GetBool("color")
		mode, _ := cmd.Flags().GetString("mode")

		fmt.Printf("Starting %s player for YouTube video: %s\n", mode, youtubeURL)
		fmt.Printf("Settings - FPS: %d, Loop: %t, Resolution: %s, Color: %t, Mode: %s\n", fps, loop, resolution, color, mode)

		// Create and start TUI player
		player := tui.NewPlayer(youtubeURL, fps, loop, resolution, color, mode)

		err := player.Play()
		if err != nil {
			fmt.Printf("Error during playback: %v\n", err)
			return
		}

	},
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
	rootCmd.AddCommand(youtubeCmd)
	youtubeCmd.Flags().BoolP("color", "c", false, "Enable colored output")
	youtubeCmd.Flags().IntP("fps", "f", 30, "Frames per second for playback")
	youtubeCmd.Flags().BoolP("loop", "l", false, "Loop the animation")
	youtubeCmd.Flags().StringP("resolution", "r", "high", "Resolution quality (low, medium, high, ultra)")
	youtubeCmd.Flags().StringP("mode", "m", "ascii", "Player mode (ascii, pixel)")
}
