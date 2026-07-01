package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DefaultPath returns the absolute path of the user's config file:
// $UserConfigDir/stex/config.json. The directory is not created, Save does it on demand.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "stex", "config.json"), nil
}

// ErrNotFound is returned by Load when the config file does not exist. Callers should treat
// this as "use defaults" without warning the user.
var ErrNotFound = fmt.Errorf("config: file not found")

// Load reads the config file into a Config. A missing file returns ErrNotFound so the caller
// can distinguish first run from a corrupt file.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrNotFound
		}
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return c, nil
}

// Save writes c to the user's config file ($UserConfigDir/stex/config.json). The parent
// directory is created if it does not exist. The Filter field is not exported to JSON, so it
// is always nil after a round trip.
func (c Config) Save() error {
	path, err := DefaultPath()
	if err != nil {
		return fmt.Errorf("config path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
