package analyzer

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/Yousseflamrani/TP-GoLog-Analyzer-/internal/config"
)

var (
	ErrFileNotFound = errors.New("fichier introuvable")
	ErrParseFailed  = errors.New("échec du parsing")
)

type Result struct {
	LogID  string
	Status string
	Error  error
}

func ProcessLog(log config.LogConfig, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, err := os.Stat(log.Path); os.IsNotExist(err) {
		results <- Result{LogID: log.ID, Status: "FAILED", Error: ErrFileNotFound}
		return
	}

	time.Sleep(3 * time.Second)

	results <- Result{LogID: log.ID, Status: "OK"}
}
