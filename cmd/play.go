package cmd

import (
	"fmt"
	"regexp"

	"github.com/kweonminsung/console-cinema/pkg/tui/player"
	"github.com/spf13/cobra"
)

// isYouTubeURL checks if the given string is a YouTube URL
func isYouTubeURL(url string) bool {
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

var playCmd = &cobra.Command{
	Use:   "play [file]",
	Short: "Play ASCII/Pixel animations from a local video file",
	Long:  `Play ASCII/Pixel animations from a specified local video file (MP4, AVI, etc.). The video will be converted to ASCII art or pixel art in real-time and displayed in the terminal. Supports options for mode, FPS and looping.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filename string
		if len(args) > 0 {
			filename = args[0]
		} else {
			fmt.Println("Error: Please specify a local video file to play")
			fmt.Println("Usage: console-cinema play <video.mp4>")
			fmt.Println("For YouTube videos, use: console-cinema youtube <url>")
			return
		}

		// Check if it's a YouTube URL and reject it
		if isYouTubeURL(filename) {
			fmt.Println("Error: YouTube URLs are not supported in 'play' command")
			fmt.Println("For YouTube videos, use: console-cinema youtube <url>")
			return
		}

		fps, _ := cmd.Flags().GetInt("fps")
		loop, _ := cmd.Flags().GetBool("loop")
		color, _ := cmd.Flags().GetBool("color")
		mode, _ := cmd.Flags().GetString("mode")

		fmt.Printf("Starting %s player for local file: %s\n", mode, filename)
		fmt.Printf("Settings - FPS: %d, Loop: %t, Color: %t, Mode: %s\n", fps, loop, color, mode)

		// Create and start TUI player
		player := player.NewPlayer(filename, fps, loop, color, mode)

		err := player.Play()
		if err != nil {
			fmt.Printf("Error during playback: %v\n", err)
			return
		}

	},
}

func init() {
	playCmd.Flags().BoolP("color", "c", true, "Enable colored output")
	playCmd.Flags().IntP("fps", "f", 30, "Frames per second for playback")
	playCmd.Flags().BoolP("loop", "l", false, "Loop the animation")
	playCmd.Flags().StringP("mode", "m", "pixel", "Player mode (ascii, pixel)")

	rootCmd.AddCommand(playCmd)
}
