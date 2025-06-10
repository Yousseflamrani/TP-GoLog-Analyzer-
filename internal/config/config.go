package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogSource représente une source de log configurée
type LogSource struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// LogEntry représente une entrée de log parsée
type LogEntry struct {
	Timestamp  time.Time         `json:"timestamp"`
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	Source     string            `json:"source"`
	SourceID   string            `json:"source_id"`
	SourceType string            `json:"source_type"`
	Fields     map[string]string `json:"fields"`
	LineNumber int               `json:"line_number"`
}

// LoadLogSourcesFromFile charge les sources de logs depuis un fichier JSON
func LoadLogSourcesFromFile(filePath string) ([]LogSource, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	var sources []LogSource
	if err := json.Unmarshal(data, &sources); err != nil {
		return nil, fmt.Errorf("impossible de parser le JSON dans %s: %w", filePath, err)
	}
	return sources, nil
}

// SaveLogSourcesToFile sauvegarde les sources de logs dans un fichier JSON
func SaveLogSourcesToFile(filePath string, sources []LogSource) error {
	// Créer le répertoire parent si nécessaire
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser en JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", filePath, err)
	}
	return nil
}

// ValidateLogSource valide une source de log
func ValidateLogSource(source LogSource) error {
	if source.ID == "" {
		return fmt.Errorf("ID manquant")
	}
	if source.Path == "" {
		return fmt.Errorf("chemin manquant")
	}
	if source.Type == "" {
		return fmt.Errorf("type manquant")
	}
	return nil
}
