package cmd

import (
	"fmt"

	"github.com/kweonminsung/ascii-player/pkg/tui"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play [file]",
	Short: "Play ASCII animations from an MP4 video file",
	Long:  `Play ASCII animations from a specified MP4 video file. The video will be converted to ASCII art in real-time and displayed in the terminal. Supports options for Mode, FPS, looping, and resolution.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filename string
		if len(args) > 0 {
			filename = args[0]
		} else {
			fmt.Println("Error: Please specify an MP4 file to play")
			fmt.Println("Usage: ascii-player play <video.mp4>")
			return
		}

		fps, _ := cmd.Flags().GetInt("fps")
		loop, _ := cmd.Flags().GetBool("loop")
		resolution, _ := cmd.Flags().GetString("resolution")
		color, _ := cmd.Flags().GetBool("color")
		mode, _ := cmd.Flags().GetString("mode")

		fmt.Printf("Starting %s player for file: %s\n", mode, filename)
		fmt.Printf("Settings - FPS: %d, Loop: %t, Resolution: %s, Color: %t, Mode: %s\n", fps, loop, resolution, color, mode)

		// Create and start TUI player
		player := tui.NewPlayer(filename, fps, loop, resolution, color, mode)

		err := player.Play()
		if err != nil {
			fmt.Printf("Error during playback: %v\n", err)
			return
		}

	},
}

func init() {
	playCmd.Flags().BoolP("color", "c", false, "Enable colored output")
	playCmd.Flags().IntP("fps", "f", 30, "Frames per second for playback")
	playCmd.Flags().BoolP("loop", "l", false, "Loop the animation")
	playCmd.Flags().StringP("resolution", "r", "high", "Resolution quality (low, medium, high, ultra)")
	playCmd.Flags().StringP("mode", "m", "ascii", "Player mode (ascii, pixel)")

	rootCmd.AddCommand(playCmd)
}
