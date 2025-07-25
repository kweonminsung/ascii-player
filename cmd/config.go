package cmd

import (
	"fmt"

	"github.com/kweonminsung/ascii-player/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Printf("Error creating config manager: %v\n", err)
			return
		}

		configData, err := manager.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("Current configuration:")
		fmt.Println("======================")
		fmt.Println(configData.Display())
		fmt.Printf("\nConfig file location: %s\n", manager.GetConfigPath())
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit ascii-player configuration",
	Long:  `Opens the configuration file in the default editor. After saving and closing the editor, displays the updated configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Printf("Error creating config manager: %v\n", err)
			return
		}

		fmt.Printf("Opening config file: %s\n", manager.GetConfigPath())
		fmt.Println("Please edit the configuration and save the file...")

		wasModified, updatedConfig, err := manager.Edit()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if wasModified {
			fmt.Println("\nConfiguration file was updated!")
		} else {
			fmt.Println("\nNo changes detected.")
		}

		fmt.Println("\nUpdated configuration:")
		fmt.Println("======================")
		fmt.Println(updatedConfig.Display())
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Printf("Error creating config manager: %v\n", err)
			return
		}

		configData, err := manager.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("Current configuration:")
		fmt.Println("======================")
		fmt.Println(configData.Display())
		fmt.Printf("\nConfig file location: %s\n", manager.GetConfigPath())
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset the configuration file to default values. This will overwrite any existing configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := config.NewManager()
		if err != nil {
			fmt.Printf("Error creating config manager: %v\n", err)
			return
		}

		fmt.Printf("Resetting configuration to defaults...\n")

		defaultConfig, err := manager.Reset()
		if err != nil {
			fmt.Printf("Error resetting config: %v\n", err)
			return
		}

		fmt.Println("Configuration has been reset to defaults!")

		fmt.Println("\nDefault configuration:")
		fmt.Println("======================")
		fmt.Println(defaultConfig.Display())
		fmt.Printf("\nConfig file location: %s\n", manager.GetConfigPath())
	},
}

func init() {
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configResetCmd)
}
