package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_DefaultPath(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "opencode-orc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
events:
  include:
    - step_start
    - text
  rules:
    step_start:
      keep:
        - sessionID
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Load config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	// Verify
	if len(cfg.Events.Include) != 2 {
		t.Errorf("expected 2 include events, got %d", len(cfg.Events.Include))
	}
}
