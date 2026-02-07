// Package config contains the user config stored in ~/.config/merlion.json
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"merlion/internal/vault/cloud"
	"merlion/internal/vault/files"
	"merlion/internal/vault/sqlite"

	"github.com/charmbracelet/log"
)

type Vault struct {
	Provider string `json:"provider"`
	Path     string `json:"path"`
	Name     string `json:"name"`
}

type UserConfig struct {
	Theme          string  `json:"theme"`
	InfoHidden     bool    `json:"infoHidden"`
	InfoBottom     bool    `json:"infoBottom"`
	CompactView    bool    `json:"compactView"`
	DefaultToCloud bool    `json:"defaultToCloud"`
	Vaults         []Vault `json:"vaults"`
}

var (
	instance UserConfig
	once     sync.Once
)

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}
	configDir := filepath.Join(homeDir, ".config")
	merlionDir := filepath.Join(configDir, "merlion")
	if err := os.MkdirAll(merlionDir, 0o700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	log.Infof("merlion Dir: %s", configDir)
	return merlionDir, nil
}

func (c *UserConfig) validate() error {
	validProviders := map[string]bool{
		cloud.Type:  true,
		sqlite.Type: true,
		files.Type:  true,
	}

	for _, vault := range c.Vaults {
		if !validProviders[vault.Provider] {
			return fmt.Errorf("provider must be one of: cloud, sqlite, files")
		}
	}
	return nil
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
		log.Info("Reading :s", configPath)
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

		log.Info("Config", config)

		if err := config.validate(); err != nil {
			log.Fatalf("Invalid config: %v", err)
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

	path := filepath.Join(configDir, "config.json")
	log.Info("Saving config: ", path)

	if err := os.WriteFile(
		path,
		data,
		0o600,
	); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
