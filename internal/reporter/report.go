package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Report struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

func ExportJSON(results []Report, path string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
