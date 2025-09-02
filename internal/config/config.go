package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/IPGeolocation/cli/internal/utils"
)

type Config struct {
	ApiKey string `json:"apikey"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ipgeolocation", "config.json")
}

func Save(cfg Config) error {
	path := configPath()
	os.MkdirAll(filepath.Dir(path), 0755)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, data, 0600)
}

func Load() (Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.ApiKey != "" {
		decrypted, err := utils.DecryptString(cfg.ApiKey)
		if err == nil {
			cfg.ApiKey = decrypted
		}
	}

	return cfg, nil
}
