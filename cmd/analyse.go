package cmd

import (
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/analyser"
	"github.com/axellelanca/go_loganizer/internal/config"
	"github.com/axellelanca/go_loganizer/internal/reporter"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var (
	configPath   string
	outputPath   string
	statusFilter string
)

var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "Analyse les fichiers de logs définis dans la configuration.",
	Long:  `La commande 'analyse' parse les fichiers de logs configurés et simule l'analyse avec gestion d'erreurs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if configPath == "" {
			fmt.Println("❌ Erreur: le fichier de configuration (--config) est obligatoire.")
			return
		}

		// Charger la configuration des logs
		logConfigs, err := config.LoadLogSourcesFromFile(configPath)
		if err != nil {
			fmt.Printf("❌ Erreur lors du chargement de la configuration: %v\n", err)
			return
		}

		if len(logConfigs) == 0 {
			fmt.Println("⚠️  Aucune configuration de log trouvée.")
			return
		}

		fmt.Printf("🚀 Analyse de %d fichiers de logs en parallèle...\n\n", len(logConfigs))

		// Analyser les logs en parallèle
		results := analyzeLogsParallel(logConfigs)

		// Filtrer les résultats si nécessaire
		if statusFilter != "" {
			results = filterResultsByStatus(results, statusFilter)
		}

		// Afficher les résultats
		displayResults(results)

		// Exporter si demandé
		if outputPath != "" {
			// BONUS: Horodatage des exports
			finalOutputPath := addTimestampToFilename(outputPath)

			err := reporter.ExportAnalysisToJSON(finalOutputPath, results)
			if err != nil {
				fmt.Printf("❌ Erreur lors de l'exportation: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", finalOutputPath)
			}
		}
	},
}

// analyzeLogsParallel traite tous les logs en parallèle avec goroutines
func analyzeLogsParallel(logConfigs []config.LogSource) []analyser.LogResult {
	var wg sync.WaitGroup
	resultsChan := make(chan analyser.LogResult, len(logConfigs))

	// Lancer une goroutine pour chaque log
	for _, logConfig := range logConfigs {
		wg.Add(1)
		go func(cfg config.LogSource) {
			defer wg.Done()

			// Analyser le log avec simulation d'erreurs
			result := analyser.AnalyzeLog(cfg)
			resultsChan <- result
		}(logConfig)
	}

	// Attendre que toutes les goroutines finissent
	wg.Wait()
	close(resultsChan)

	// Collecter tous les résultats
	var results []analyser.LogResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// filterResultsByStatus filtre les résultats par statut
func filterResultsByStatus(results []analyser.LogResult, status string) []analyser.LogResult {
	var filtered []analyser.LogResult
	for _, result := range results {
		if result.Status == status {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// displayResults affiche les résultats sur la console
func displayResults(results []analyser.LogResult) {
	fmt.Println("📊 === RÉSULTATS D'ANALYSE ===")

	successCount := 0
	failedCount := 0

	for _, result := range results {
		status := "✅"
		if result.Status == "FAILED" {
			status = "❌"
			failedCount++
		} else {
			successCount++
		}

		fmt.Printf("%s %s (%s) - %s: %s\n",
			status, result.LogID, result.FilePath, result.Status, result.Message)

		if result.ErrorDetails != "" {
			fmt.Printf("   Détails de l'erreur: %s\n", result.ErrorDetails)
		}
	}

	fmt.Printf("\n📈 Résumé: %d succès, %d échecs sur %d fichiers\n",
		successCount, failedCount, len(results))
}

// addTimestampToFilename ajoute un timestamp au nom de fichier (BONUS)
func addTimestampToFilename(filename string) string {
	now := time.Now()
	timestamp := now.Format("060102") // Format AAMMJJ

	// Insérer le timestamp avant l'extension
	if len(filename) > 5 && filename[len(filename)-5:] == ".json" {
		base := filename[:len(filename)-5]
		return fmt.Sprintf("%s_%s.json", timestamp, base)
	}

	return fmt.Sprintf("%s_%s", timestamp, filename)
}

func init() {
	rootCmd.AddCommand(analyseCmd)

	analyseCmd.Flags().StringVarP(&configPath, "config", "c", "", "Chemin vers le fichier de configuration JSON (obligatoire)")
	analyseCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Chemin vers le fichier de sortie JSON (optionnel)")
	analyseCmd.Flags().StringVar(&statusFilter, "status", "", "Filtrer par statut (OK, FAILED)")

	analyseCmd.MarkFlagRequired("config")
}
