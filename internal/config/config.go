// Package config contains the user config stored in ~/.config/merlion.json
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

type UserConfig struct {
	Theme          string `json:"theme"`
	InfoHidden     bool   `json:"infoHidden"`
	InfoBottom     bool   `json:"infoBottom"`
	CompactView    bool   `json:"compactView"`
	DefaultToCloud bool   `json:"defaultToCloud"`
}

var (
	instance UserConfig
	once     sync.Once
)

func getConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}
	configDir := filepath.Join(userConfigDir, "merlion")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	return configDir, nil
}

// Load loads the user config from the config file
func Load() *UserConfig {
	once.Do(func() {
		configDir, err := getConfigDir()
		if err != nil {
			log.Fatalf("Failed to get config dir: %v", err)
		}

		var config UserConfig
		configPath := filepath.Join(configDir, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			config = UserConfig{
				Theme:       "neotokyo",
				InfoHidden:  false,
				InfoBottom:  true,
				CompactView: false,
			}
			config.Save()
		} else {
			data, err := os.ReadFile(configPath)
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}
			if err := json.Unmarshal(data, &config); err != nil {
				log.Fatalf("Failed to parse config: %v", err)
			}
		}

		instance = config
	})

	return &instance
}

// Save saves the user config to the config file
func (config *UserConfig) Save() error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		filepath.Join(configDir, "config.json"),
		data,
		0o600,
	); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
