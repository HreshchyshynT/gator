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

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	return write(c)
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to create config path: %v", err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return config, fmt.Errorf("Failed to read config file: %v", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Can't find user home dir: %v", err)
	}

	return path.Join(homeDir, configName), nil
}

func write(config *Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("Failed to Unmarshal config: %v", err)
	}

	return os.WriteFile(configPath, data, os.ModeExclusive)
}
