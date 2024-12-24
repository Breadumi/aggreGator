package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = "/.gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {

	filePath, err := getFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading json config file")
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding json config file")
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + configFileName, nil
}

func write(c *Config) error {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
