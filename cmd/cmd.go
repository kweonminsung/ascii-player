package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ascii-player",
	Short: "A real-time ASCII art video streamer for the command line",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ASCII Player - Real-time ASCII art video streamer")
		fmt.Println("Use 'ascii-player --help' for more information.")
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
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(youtubeCmd)

	// Global flags
	// rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	// rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode")
}
