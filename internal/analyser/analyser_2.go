package analyser

import (
	"errors"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/config"
	"os"
)

// AddLogConfig ajoute une nouvelle configuration de log (BONUS add-log)
func AddLogConfig(configFile, id, path, logType string) error {
	// Charger les configurations existantes
	existingConfigs, err := config.LoadLogSourcesFromFile(configFile)
	if err != nil {
		// Si le fichier n'existe pas, créer une liste vide
		if os.IsNotExist(err) {
			existingConfigs = []config.LogSource{}
		} else {
			return fmt.Errorf("erreur lors du chargement: %w", err)
		}
	}

	// Vérifier si l'ID existe déjà
	for _, cfg := range existingConfigs {
		if cfg.ID == id {
			return fmt.Errorf("un log avec l'ID '%s' existe déjà", id)
		}
	}

	// Créer la nouvelle configuration
	newConfig := config.LogSource{
		ID:   id,
		Path: path,
		Type: logType,
	}

	// Ajouter à la liste
	existingConfigs = append(existingConfigs, newConfig)

	// Sauvegarder
	return config.SaveLogSourcesToFile(configFile, existingConfigs)
}

// ValidateLogConfig valide une configuration de log
func ValidateLogConfig(cfg config.LogSource) error {
	if cfg.ID == "" {
		return errors.New("ID ne peut pas être vide")
	}
	if cfg.Path == "" {
		return errors.New("path ne peut pas être vide")
	}
	if cfg.Type == "" {
		return errors.New("type ne peut pas être vide")
	}
	return nil
}
