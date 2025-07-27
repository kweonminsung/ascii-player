package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ascii-player",
	Short: "A real-time ASCII/Pixel art video player for the command line",
	Long: `ASCII Player - Convert and play videos as ASCII/Pixel art in real-time

Commands:
  play     Play local video files (MP4, AVI, etc.)
  youtube  Play YouTube videos by URL
  config   Manage configuration settings

Examples:
  ascii-player play video.mp4 --mode ascii --fps 30
  ascii-player youtube https://youtube.com/watch?v=... --mode pixel
  ascii-player config show`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ASCII Player - Real-time ASCII/Pixel art video player")
		fmt.Println("Use 'ascii-player --help' for more information.")
		fmt.Println("")
		fmt.Println("Quick start:")
		fmt.Println("  ascii-player play video.mp4")
		fmt.Println("  ascii-player youtube <youtube_url>")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add all subcommands to root
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(youtubeCmd)

	// Global flags
	// rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	// rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode")
}
