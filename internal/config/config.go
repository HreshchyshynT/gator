package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const configName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c Config) String() string {
	return fmt.Sprintf("Config(\nDbUrl=%v\nCurrentUserName=%v\n)\n", c.DbUrl, c.CurrentUserName)
}

func Read() (Config, error) {
	var config Config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("Can't find user home dir: %v", err)
	}

	configPath := path.Join(homeDir, configName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("Failed to create config path: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("Failed to read config file: %v", err)
	}

	return config, nil
}
