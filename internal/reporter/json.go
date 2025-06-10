package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/analyser"
	"os"
	"path/filepath"
)

// ExportAnalysisToJSON exporte les résultats vers un fichier JSON (selon le TP)
func ExportAnalysisToJSON(outputPath string, results []analyser.LogResult) error {
	// BONUS: Créer les répertoires parent si nécessaire
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer les répertoires %s: %w", dir, err)
	}

	// Sérialiser en JSON avec indentation
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser les résultats: %w", err)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", outputPath, err)
	}

	return nil
}

// ExportResultsToCSV exporte un résumé vers CSV (bonus)
func ExportResultsToCSV(outputPath string, results []analyser.LogResult) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer les répertoires %s: %w", dir, err)
	}

	// Construire le contenu CSV
	csvContent := "LogID,FilePath,Status,Message,ErrorDetails\n"

	for _, result := range results {
		csvContent += fmt.Sprintf("%s,%s,%s,%s,%s\n",
			result.LogID, result.FilePath, result.Status,
			result.Message, result.ErrorDetails)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(outputPath, []byte(csvContent), 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", outputPath, err)
	}

	return nil
}
