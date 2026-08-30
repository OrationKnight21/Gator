package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func GetConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, configFileName), nil
}
func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
