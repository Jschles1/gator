package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}
	filePath := home + "/" + configFileName
	return filePath, nil
}

func Read() (*Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	configData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg *Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	newConfig, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, newConfig, 0777)
	return nil
}
