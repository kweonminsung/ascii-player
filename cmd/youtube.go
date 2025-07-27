package cmd

import (
	"fmt"
	"regexp"

	"github.com/kweonminsung/console-cinema/pkg/tui/player"
	"github.com/spf13/cobra"
)

var youtubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "Play or explore YouTube videos",
	Long:  `A container for YouTube related commands like play and explore.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var youtubePlayCmd = &cobra.Command{
	Use:   "play [url]",
	Short: "Play a YouTube video from a URL",
	Long:  `Play ASCII/Pixel animations from a specified YouTube video URL.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		youtubeURL := args[0]

		// Check if it's a valid YouTube URL
		if !isValidYouTubeURL(youtubeURL) {
			fmt.Println("Error: Invalid YouTube URL provided")
			fmt.Println("Supported formats:")
			fmt.Println("  - https://www.youtube.com/watch?v=VIDEO_ID")
			fmt.Println("  - https://youtu.be/VIDEO_ID")
			fmt.Println("  - https://www.youtube.com/embed/VIDEO_ID")
			return
		}

		// Note: Flags are inherited from the parent youtubeCmd
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
	youtubeCmd.AddCommand(youtubePlayCmd)

	// Flags for both play and explore (via playVideo)
	youtubeCmd.PersistentFlags().BoolP("color", "c", true, "Enable colored output")
	youtubeCmd.PersistentFlags().IntP("fps", "f", 30, "Frames per second for playback")
	youtubeCmd.PersistentFlags().BoolP("loop", "l", false, "Loop the animation")
	youtubeCmd.PersistentFlags().StringP("mode", "m", "pixel", "Player mode (ascii, pixel)")
}
