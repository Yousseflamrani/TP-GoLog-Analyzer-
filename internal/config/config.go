package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LogSource représente une configuration de log du fichier JSON (selon le TP)
type LogSource struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// LoadLogSourcesFromFile charge les configurations depuis un fichier JSON
func LoadLogSourcesFromFile(configPath string) ([]LogSource, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier de configuration %s: %w", configPath, err)
	}

	var configs []LogSource
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("impossible de parser le JSON dans %s: %w", configPath, err)
	}

	return configs, nil
}

// SaveLogSourcesToFile sauvegarde les configurations dans un fichier JSON
func SaveLogSourcesToFile(configPath string, configs []LogSource) error {
	// BONUS: Créer les répertoires parent si nécessaire
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer les répertoires %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser les configurations: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", configPath, err)
	}

	return nil
}
