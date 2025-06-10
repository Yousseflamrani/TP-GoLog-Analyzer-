package analyser

import (
	"bufio"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/config"
	"os"
	"sort"
	"time"
)

// LogAnalyser structure pour analyser une source de log
type LogAnalyser struct {
	source      config.LogSource
	filterLevel string
	parser      LogParser
}

// SourceAnalysisResult résultat d'analyse pour une source
type SourceAnalysisResult struct {
	SourceID    string         `json:"source_id"`
	SourceType  string         `json:"source_type"`
	SourcePath  string         `json:"source_path"`
	TotalLines  int            `json:"total_lines"`
	ParsedLines int            `json:"parsed_lines"`
	ErrorLines  int            `json:"error_lines"`
	LevelStats  map[string]int `json:"level_stats"`
	HourlyStats map[int]int    `json:"hourly_stats"`
	Errors      []ErrorStat    `json:"errors"`
	Duration    time.Duration  `json:"duration"`
	Error       error          `json:"error,omitempty"`
}

// GlobalAnalysisResult résultat d'analyse globale
type GlobalAnalysisResult struct {
	TotalSources      int                     `json:"total_sources"`
	SuccessfulSources int                     `json:"successful_sources"`
	ErrorSources      int                     `json:"error_sources"`
	TotalLines        int                     `json:"total_lines"`
	LevelStats        map[string]int          `json:"level_stats"`
	TypeStats         map[string]int          `json:"type_stats"`
	TopErrors         []ErrorStat             `json:"top_errors"`
	SourceResults     []*SourceAnalysisResult `json:"source_results"`
	AnalysisDuration  time.Duration           `json:"analysis_duration"`
	StartTime         time.Time               `json:"start_time"`
}

// ErrorStat représente une statistique d'erreur
type ErrorStat struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
	Level   string `json:"level"`
	Source  string `json:"source"`
}

// NewLogAnalyser crée un nouvel analyseur pour une source
func NewLogAnalyser(source config.LogSource, filterLevel string) *LogAnalyser {
	analyzer := &LogAnalyser{
		source:      source,
		filterLevel: filterLevel,
	}

	// Initialiser le parser selon le type
	analyzer.parser = GetParserForType(source.Type)
	return analyzer
}

// NewGlobalAnalysisResult crée un nouveau résultat global
func NewGlobalAnalysisResult() *GlobalAnalysisResult {
	return &GlobalAnalysisResult{
		LevelStats:    make(map[string]int),
		TypeStats:     make(map[string]int),
		TopErrors:     make([]ErrorStat, 0),
		SourceResults: make([]*SourceAnalysisResult, 0),
		StartTime:     time.Now(),
	}
}

// AnalyseSource analyse une source de log
func (la *LogAnalyser) AnalyseSource() *SourceAnalysisResult {
	startTime := time.Now()

	result := &SourceAnalysisResult{
		SourceID:    la.source.ID,
		SourceType:  la.source.Type,
		SourcePath:  la.source.Path,
		LevelStats:  make(map[string]int),
		HourlyStats: make(map[int]int),
		Errors:      make([]ErrorStat, 0),
	}

	// Vérifier si le fichier existe
	if _, err := os.Stat(la.source.Path); os.IsNotExist(err) {
		result.Error = fmt.Errorf("fichier non trouvé: %s", la.source.Path)
		result.Duration = time.Since(startTime)
		return result
	}

	// Ouvrir le fichier
	file, err := os.Open(la.source.Path)
	if err != nil {
		result.Error = fmt.Errorf("impossible d'ouvrir le fichier: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	defer file.Close()

	// Analyser ligne par ligne
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	errorCounts := make(map[string]int)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		result.TotalLines++

		// Parser la ligne
		entry, err := la.parser.ParseLine(line, lineNumber)
		if err != nil {
			result.ErrorLines++
			continue
		}

		// Appliquer le filtre de niveau si spécifié
		if la.filterLevel != "" && entry.Level != la.filterLevel {
			continue
		}

		result.ParsedLines++

		// Enrichir l'entrée avec les infos de la source
		entry.SourceID = la.source.ID
		entry.SourceType = la.source.Type

		// Statistiques par niveau
		result.LevelStats[entry.Level]++

		// Statistiques horaires
		hour := entry.Timestamp.Hour()
		result.HourlyStats[hour]++

		// Compter les erreurs
		if entry.Level == "ERROR" || entry.Level == "FATAL" {
			errorCounts[entry.Message]++
		}
	}

	// Créer les statistiques d'erreurs
	la.createErrorStats(errorCounts, result)

	result.Duration = time.Since(startTime)

	if err := scanner.Err(); err != nil {
		result.Error = fmt.Errorf("erreur lors de la lecture: %w", err)
	}

	return result
}

// createErrorStats crée les statistiques d'erreurs
func (la *LogAnalyser) createErrorStats(errorCounts map[string]int, result *SourceAnalysisResult) {
	type errorPair struct {
		message string
		count   int
	}

	var pairs []errorPair
	for msg, count := range errorCounts {
		pairs = append(pairs, errorPair{msg, count})
	}

	// Trier par nombre d'occurrences (décroissant)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Convertir en ErrorStat
	for _, pair := range pairs {
		result.Errors = append(result.Errors, ErrorStat{
			Message: pair.message,
			Count:   pair.count,
			Level:   "ERROR",
			Source:  la.source.ID,
		})
	}
}

// AddSourceResult ajoute un résultat de source au résultat global
func (gar *GlobalAnalysisResult) AddSourceResult(result *SourceAnalysisResult) {
	gar.TotalSources++
	gar.SourceResults = append(gar.SourceResults, result)

	if result.Error != nil {
		gar.ErrorSources++
		return
	}

	gar.SuccessfulSources++
	gar.TotalLines += result.TotalLines

	// Agréger les statistiques de niveau
	for level, count := range result.LevelStats {
		gar.LevelStats[level] += count
	}

	// Agréger les statistiques de type
	gar.TypeStats[result.SourceType] += result.TotalLines

	// Agréger les erreurs
	gar.TopErrors = append(gar.TopErrors, result.Errors...)
}

// Finalize finalise le résultat global
func (gar *GlobalAnalysisResult) Finalize() {
	gar.AnalysisDuration = time.Since(gar.StartTime)

	// Trier les erreurs par nombre d'occurrences
	sort.Slice(gar.TopErrors, func(i, j int) bool {
		return gar.TopErrors[i].Count > gar.TopErrors[j].Count
	})

	// Limiter aux 10 premières erreurs
	if len(gar.TopErrors) > 10 {
		gar.TopErrors = gar.TopErrors[:10]
	}
}
