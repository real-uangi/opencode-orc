package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/real-uangi/opencode-orc/types"
	"gopkg.in/yaml.v3"
)

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(home, ".config", "opencode-orc", "config.yaml")
}

// DefaultConfig returns a default configuration
func DefaultConfig() *types.Config {
	return &types.Config{
		Events: types.EventsConfig{
			Include: []string{
				types.EventTypeStepStart,
				types.EventTypeToolUse,
				types.EventTypeText,
				types.EventTypeStepFinish,
				types.EventTypeError,
			},
			Rules: map[string]types.EventRule{
				types.EventTypeStepStart: {
					Keep: []string{"sessionID"},
				},
				types.EventTypeToolUse: {
					Keep: []string{"part.tool", "part.state.status", "part.state.input", "part.state.error", "part.state.metadata.exit", "part.title"},
				},
				types.EventTypeText: {
					Keep: []string{"part.text"},
				},
				types.EventTypeStepFinish: {
					Keep: []string{"part.reason"},
				},
				types.EventTypeError: {
					Keep: []string{"error.name", "error.data.message"},
				},
			},
		},
		Output: types.OutputConfig{
			Format: "text",
			Pretty: false,
		},
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg types.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// WriteDefaultConfig writes a default config file if it doesn't exist
func WriteDefaultConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // File exists
	}

	if err := EnsureConfigDir(path); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	cfg := DefaultConfig()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
