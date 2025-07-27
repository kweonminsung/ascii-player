package cmd

import (
	"fmt"
	"regexp"

	"github.com/kweonminsung/console-cinema/pkg/tui/player"
	"github.com/spf13/cobra"
)

var youtubeCmd = &cobra.Command{
	Use:   "youtube [url]",
	Short: "Play ASCII/Pixel animations from a YouTube video",
	Long:  `Play ASCII/Pixel animations from a specified YouTube video URL. The video will be streamed and converted to ASCII art or pixel art in real-time and displayed in the terminal. Supports options for mode, FPS, looping, and resolution.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		youtubeURL := args[0]
		// The 'explore' command is a subcommand, so if the arg is 'explore',
		// cobra will handle it. We just need to avoid treating 'explore' as a URL.
		if youtubeURL == "explore" {
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
		color, _ := cmd.Flags().GetBool("color")
		mode, _ := cmd.Flags().GetString("mode")

		fmt.Printf("Starting %s player for YouTube video: %s\n", mode, youtubeURL)
		fmt.Printf("Settings - FPS: %d, Loop: %t, Color: %t, Mode: %s\n", fps, loop, color, mode)

		// Create and start TUI player
		player := player.NewPlayer(youtubeURL, fps, loop, color, mode)

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
	youtubeCmd.Flags().BoolP("color", "c", true, "Enable colored output")
	youtubeCmd.Flags().IntP("fps", "f", 30, "Frames per second for playback")
	youtubeCmd.Flags().BoolP("loop", "l", false, "Loop the animation")
	youtubeCmd.Flags().StringP("mode", "m", "pixel", "Player mode (ascii, pixel)")
}
