package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zalando/go-keyring"
)

const (
	appName    = "AI Terminal Pro"
	configFile = "config.json"
)

// Settings represents user configuration
type Settings struct {
	LiteLLMEndpoint string `json:"litellm_endpoint"`
	VirtualKey      string `json:"-"` // Not stored in JSON, use keyring
	Model           string `json:"model"`
	Theme           string `json:"theme"`
	FontSize        int    `json:"font_size"`
	FontFamily      string `json:"font_family"`
	CursorStyle     string `json:"cursor_style"`
	AIShortcut      string `json:"ai_shortcut"`
	SafetyMode      string `json:"safety_mode"` // strict, normal, off
}

// DefaultSettings returns default configuration
func DefaultSettings() *Settings {
	return &Settings{
		LiteLLMEndpoint: "",
		Model:           "qwen3-terminal",
		Theme:           "dark",
		FontSize:        14,
		FontFamily:      "JetBrains Mono",
		CursorStyle:     "block",
		AIShortcut:      "ctrl+k",
		SafetyMode:      "normal",
	}
}

// Load reads settings from disk
func Load() (*Settings, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	settings := DefaultSettings()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return defaults if config doesn't exist
			return settings, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, settings); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Load virtual key from keyring
	key, err := keyring.Get(appName, "virtual_key")
	if err == nil {
		settings.VirtualKey = key
	}

	return settings, nil
}

// Save writes settings to disk
func (s *Settings) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save virtual key to keyring
	if s.VirtualKey != "" {
		if err := keyring.Set(appName, "virtual_key", s.VirtualKey); err != nil {
			return fmt.Errorf("failed to store key: %w", err)
		}
	}

	// Save other settings to JSON (without virtual key)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getConfigPath returns the OS-appropriate config path
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(home, "Library", "Application Support", "ai-terminal")
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "ai-terminal")
	default: // Linux and others
		configDir = filepath.Join(home, ".config", "ai-terminal")
	}

	return filepath.Join(configDir, configFile), nil
}

// GetConfigDir returns the configuration directory
func GetConfigDir() (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(configPath), nil
}

// GetOSType returns the current operating system
func (s *Settings) GetOSType() string {
	return runtime.GOOS
}

// GetShell returns the default shell for the OS
func (s *Settings) GetShell() string {
	switch runtime.GOOS {
	case "windows":
		return "powershell"
	default:
		return "bash"
	}
}
