package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/analyser"
	"os"
	"path/filepath"
)

// ExportAnalysisToJSON exporte les résultats d'analyse vers un fichier JSON
func ExportAnalysisToJSON(filePath string, results *analyser.GlobalAnalysisResult) error {
	// Créer le répertoire parent si nécessaire
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", dir, err)
	}

	// Sérialiser en JSON avec indentation
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser les résultats: %w", err)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", filePath, err)
	}

	return nil
}

// ExportAnalysisToCSV exporte un résumé vers CSV
func ExportAnalysisToCSV(filePath string, results *analyser.GlobalAnalysisResult) error {
	// Créer le répertoire parent si nécessaire
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", dir, err)
	}

	// Construire le contenu CSV
	csvContent := "SourceID,SourceType,TotalLines,ParsedLines,ErrorLines,Duration\n"

	for _, source := range results.SourceResults {
		if source.Error == nil {
			csvContent += fmt.Sprintf("%s,%s,%d,%d,%d,%v\n",
				source.SourceID, source.SourceType, source.TotalLines,
				source.ParsedLines, source.ErrorLines, source.Duration)
		} else {
			csvContent += fmt.Sprintf("%s,%s,0,0,0,ERROR: %v\n",
				source.SourceID, source.SourceType, source.Error)
		}
	}

	// Écrire dans le fichier
	if err := os.WriteFile(filePath, []byte(csvContent), 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", filePath, err)
	}

	return nil
}

// ExportSourceAnalysis exporte l'analyse d'une source spécifique
func ExportSourceAnalysis(filePath string, result *analyser.SourceAnalysisResult) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser le résultat: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible d'écrire dans le fichier %s: %w", filePath, err)
	}

	return nil
}
