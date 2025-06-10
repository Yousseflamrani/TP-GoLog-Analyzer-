package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func LoadConfig(path string) ([]LogConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture: %w", err)
	}
	var configs []LogConfig
	if err := json.Unmarshal(file, &configs); err != nil {
		return nil, fmt.Errorf("erreur de parsing JSON: %w", err)
	}
	return configs, nil
}
