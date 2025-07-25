package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Data represents the structure of the configuration file
type Data struct {
	FPS        int    `json:"fps"`
	Loop       bool   `json:"loop"`
	Resolution string `json:"resolution"`
}

// Manager handles configuration operations
type Manager struct {
	configPath string
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".ascii-player-config.json")
	return &Manager{
		configPath: configPath,
	}, nil
}

// GetConfigPath returns the path to the config file
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// CreateDefault creates a default configuration
func (m *Manager) CreateDefault() *Data {
	return &Data{
		FPS:        30,
		Loop:       false,
		Resolution: "high", // ultra, high, medium, low
	}
}

// Load loads configuration from file
func (m *Manager) Load() (*Data, error) {
	// If config file doesn't exist, create default
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		defaultConfig := m.CreateDefault()

		err = m.Save(defaultConfig)
		if err != nil {
			return nil, err
		}

		return defaultConfig, nil
	}

	// Read existing config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, err
	}

	var config Data
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves configuration to file
func (m *Manager) Save(config *Data) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// OpenWithEditor opens the config file with the system's default editor
func (m *Manager) OpenWithEditor() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad.exe", m.configPath)
	case "darwin":
		cmd = exec.Command("open", "-t", m.configPath)
	case "linux":
		// Try different editors in order of preference
		editors := []string{"code", "gedit", "nano", "vim", "vi"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				cmd = exec.Command(editor, m.configPath)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable editor found")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// GetFileModTime returns the modification time of the config file
func (m *Manager) GetFileModTime() (time.Time, error) {
	info, err := os.Stat(m.configPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// Display formats and returns the configuration as a string
func (config *Data) Display() string {
	return fmt.Sprintf(`FPS:           %d
Loop:          %t
Resolution:    %s`, config.FPS, config.Loop, config.Resolution,
	)
}

// Edit opens the config file in an editor and returns whether it was modified
func (m *Manager) Edit() (bool, *Data, error) {
	// Load or create config to ensure file exists
	_, err := m.Load()
	if err != nil {
		return false, nil, fmt.Errorf("error loading config: %v", err)
	}

	// Get initial modification time
	initialModTime, err := m.GetFileModTime()
	if err != nil {
		return false, nil, fmt.Errorf("error getting file modification time: %v", err)
	}

	// Open with default editor
	err = m.OpenWithEditor()
	if err != nil {
		return false, nil, fmt.Errorf("error opening editor: %v", err)
	}

	// Check if file was modified
	finalModTime, err := m.GetFileModTime()
	if err != nil {
		return false, nil, fmt.Errorf("error checking file modification: %v", err)
	}

	wasModified := finalModTime.After(initialModTime)

	// Load and return the updated config
	updatedConfig, err := m.Load()
	if err != nil {
		return wasModified, nil, fmt.Errorf("error reloading config: %v", err)
	}

	return wasModified, updatedConfig, nil
}

// Reset resets the configuration to default values
func (m *Manager) Reset() (*Data, error) {
	defaultConfig := m.CreateDefault()
	err := m.Save(defaultConfig)
	if err != nil {
		return nil, err
	}
	return defaultConfig, nil
}
