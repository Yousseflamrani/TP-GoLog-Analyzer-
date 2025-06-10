package cmd

import (
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/analyser"
	"github.com/axellelanca/go_loganizer/internal/config"
	"github.com/axellelanca/go_loganizer/internal/reporter"
	"github.com/spf13/cobra"
	"sync"
)

var (
	configFile  string
	outputFile  string
	maxWorkers  int
	verbose     bool
	filterLevel string
	filterType  string
	specificID  string
)

var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "Analyse les fichiers de logs définis dans la configuration.",
	Long:  `La commande 'analyse' parse les fichiers de logs configurés et extrait des statistiques, erreurs, et patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile == "" {
			fmt.Println("Erreur: le fichier de configuration (--config) est obligatoire.")
			return
		}

		// Charger la configuration des logs
		logSources, err := config.LoadLogSourcesFromFile(configFile)
		if err != nil {
			fmt.Printf("❌ Erreur lors du chargement de la configuration: %v\n", err)
			return
		}

		if len(logSources) == 0 {
			fmt.Println("⚠️  Aucune source de log trouvée dans la configuration.")
			return
		}

		// Filtrer les sources si nécessaire
		filteredSources := filterLogSources(logSources, filterType, specificID)

		if verbose {
			fmt.Printf("🔍 Analyse de %d sources de logs\n", len(filteredSources))
			fmt.Printf("⚡ Nombre de workers: %d\n", maxWorkers)
		}

		// Lancer l'analyse de toutes les sources
		results, err := analyseMultipleLogSources(filteredSources, maxWorkers, verbose, filterLevel)
		if err != nil {
			fmt.Printf("❌ Erreur lors de l'analyse: %v\n", err)
			return
		}

		// Afficher les résultats
		displayResults(results, verbose)

		// Exporter si demandé
		if outputFile != "" {
			err := reporter.ExportAnalysisToJSON(outputFile, results)
			if err != nil {
				fmt.Printf("❌ Erreur lors de l'exportation: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", outputFile)
			}
		}
	},
}

func filterLogSources(sources []config.LogSource, typeFilter, idFilter string) []config.LogSource {
	var filtered []config.LogSource

	for _, source := range sources {
		// Filtrer par ID si spécifié
		if idFilter != "" && source.ID != idFilter {
			continue
		}

		// Filtrer par type si spécifié
		if typeFilter != "" && source.Type != typeFilter {
			continue
		}

		filtered = append(filtered, source)
	}

	return filtered
}

func analyseMultipleLogSources(sources []config.LogSource, maxWorkers int, verbose bool, filterLevel string) (*analyser.GlobalAnalysisResult, error) {
	globalResult := analyser.NewGlobalAnalysisResult()

	// Utiliser un worker pool pour traiter les sources en parallèle
	sourceChan := make(chan config.LogSource, len(sources))
	resultChan := make(chan *analyser.SourceAnalysisResult, len(sources))

	var wg sync.WaitGroup

	// Lancer les workers
	numWorkers := maxWorkers
	if numWorkers > len(sources) {
		numWorkers = len(sources)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for source := range sourceChan {
				result := analyseLogSource(source, filterLevel, verbose)
				resultChan <- result
			}
		}()
	}

	// Envoyer les sources aux workers
	go func() {
		for _, source := range sources {
			sourceChan <- source
		}
		close(sourceChan)
	}()

	// Collecter les résultats
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Agréger tous les résultats
	for result := range resultChan {
		globalResult.AddSourceResult(result)
	}

	globalResult.Finalize()
	return globalResult, nil
}

func analyseLogSource(source config.LogSource, filterLevel string, verbose bool) *analyser.SourceAnalysisResult {
	if verbose {
		fmt.Printf("📂 Analyse de %s (%s)...\n", source.ID, source.Path)
	}

	analyzer := analyser.NewLogAnalyser(source, filterLevel)
	result := analyzer.AnalyseSource()

	if verbose {
		if result.Error != nil {
			fmt.Printf("❌ %s: %v\n", source.ID, result.Error)
		} else {
			fmt.Printf("✅ %s: %d lignes analysées\n", source.ID, result.TotalLines)
		}
	}

	return result
}

func displayResults(results *analyser.GlobalAnalysisResult, verbose bool) {
	fmt.Println("\n📈 === RÉSULTATS D'ANALYSE GLOBALE ===")
	fmt.Printf("📁 Sources analysées: %d\n", results.TotalSources)
	fmt.Printf("✅ Sources réussies: %d\n", results.SuccessfulSources)
	fmt.Printf("❌ Sources en erreur: %d\n", results.ErrorSources)
	fmt.Printf("📄 Lignes totales: %d\n", results.TotalLines)
	fmt.Printf("⏱️  Temps d'analyse: %v\n", results.AnalysisDuration)

	if len(results.LevelStats) > 0 {
		fmt.Println("\n🏷️  === RÉPARTITION PAR NIVEAU (GLOBAL) ===")
		for level, count := range results.LevelStats {
			fmt.Printf("%s: %d\n", level, count)
		}
	}

	if len(results.TypeStats) > 0 {
		fmt.Println("\n📊 === RÉPARTITION PAR TYPE DE LOG ===")
		for logType, count := range results.TypeStats {
			fmt.Printf("%s: %d lignes\n", logType, count)
		}
	}

	if len(results.TopErrors) > 0 {
		fmt.Println("\n🚨 === TOP ERREURS ===")
		for i, err := range results.TopErrors {
			if i >= 5 { // Limiter à 5
				break
			}
			fmt.Printf("%d. %s (occurrences: %d, source: %s)\n", i+1, err.Message, err.Count, err.Source)
		}
	}

	if verbose && len(results.SourceResults) > 0 {
		fmt.Println("\n📋 === DÉTAILS PAR SOURCE ===")
		for _, sourceResult := range results.SourceResults {
			if sourceResult.Error != nil {
				fmt.Printf("❌ %s (%s): %v\n", sourceResult.SourceID, sourceResult.SourceType, sourceResult.Error)
			} else {
				fmt.Printf("✅ %s (%s): %d lignes, %d erreurs\n",
					sourceResult.SourceID, sourceResult.SourceType, sourceResult.TotalLines, sourceResult.ErrorLines)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(analyseCmd)

	analyseCmd.Flags().StringVarP(&configFile, "config", "c", "", "Fichier JSON de configuration des sources de logs")
	analyseCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Fichier de sortie pour les résultats JSON")
	analyseCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 4, "Nombre de workers pour le traitement parallèle")
	analyseCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Mode verbeux")
	analyseCmd.Flags().StringVarP(&filterLevel, "level", "l", "", "Filtrer par niveau de log (DEBUG, INFO, WARN, ERROR)")
	analyseCmd.Flags().StringVarP(&filterType, "type", "t", "", "Filtrer par type de log (nginx-access, mysql-error, etc.)")
	analyseCmd.Flags().StringVarP(&specificID, "id", "i", "", "Analyser uniquement une source spécifique par son ID")

	analyseCmd.MarkFlagRequired("config")
}
