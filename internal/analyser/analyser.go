package analyser

import (
	"errors"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/config"
	"math/rand"
	"os"
	"time"
)

// LogResult représente le résultat d'analyse d'un log (selon le TP)
type LogResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

// AnalyzeLog analyse un fichier de log selon les spécifications du TP
func AnalyzeLog(cfg config.LogSource) LogResult {
	result := LogResult{
		LogID:    cfg.ID,
		FilePath: cfg.Path,
	}

	// 1. Vérifier si le fichier existe et est lisible
	if err := checkFileAccess(cfg.Path); err != nil {
		result.Status = "FAILED"

		// Utiliser les erreurs personnalisées avec errors.As
		var fileNotFoundErr *FileNotFoundError
		var fileAccessErr *FileAccessError

		if errors.As(err, &fileNotFoundErr) {
			result.Message = "Fichier introuvable."
			result.ErrorDetails = err.Error()
		} else if errors.As(err, &fileAccessErr) {
			result.Message = "Fichier inaccessible."
			result.ErrorDetails = err.Error()
		} else {
			result.Message = "Erreur d'accès au fichier."
			result.ErrorDetails = err.Error()
		}

		return result
	}

	// 2. Simuler l'analyse avec un délai aléatoire (50 à 200 ms)
	analysisDelay := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(analysisDelay)

	// 3. Simuler une erreur de parsing aléatoire (10% de chance)
	if rand.Float32() < 0.1 {
		parseErr := &ParseError{
			LogID:   cfg.ID,
			Message: "Erreur lors du parsing du fichier de log",
		}

		result.Status = "FAILED"
		result.Message = "Erreur de parsing."
		result.ErrorDetails = parseErr.Error()
		return result
	}

	// 4. Succès
	result.Status = "OK"
	result.Message = "Analyse terminée avec succès."
	result.ErrorDetails = ""

	return result
}

// checkFileAccess vérifie l'existence et l'accessibilité du fichier
func checkFileAccess(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileNotFoundError{
				Path: filePath,
				Err:  err,
			}
		}
		return &FileAccessError{
			Path: filePath,
			Err:  err,
		}
	}

	// Vérifier si c'est un fichier (pas un répertoire)
	if info.IsDir() {
		return &FileAccessError{
			Path: filePath,
			Err:  fmt.Errorf("le chemin pointe vers un répertoire, pas un fichier"),
		}
	}

	// Vérifier les permissions de lecture
	file, err := os.Open(filePath)
	if err != nil {
		return &FileAccessError{
			Path: filePath,
			Err:  err,
		}
	}
	file.Close()

	return nil
}
