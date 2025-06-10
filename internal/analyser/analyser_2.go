package analyser

import (
	"encoding/json"
	"fmt"
	"github.com/axellelanca/go_loganizer/internal/config"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LogParser interface pour différents types de parsers
type LogParser interface {
	ParseLine(line string, lineNumber int) (*config.LogEntry, error)
}

// GetParserForType retourne le parser approprié selon le type
func GetParserForType(logType string) LogParser {
	switch logType {
	case "nginx-access":
		return &NginxAccessParser{}
	case "nginx-error":
		return &NginxErrorParser{}
	case "apache-access":
		return &ApacheAccessParser{}
	case "apache-error":
		return &ApacheErrorParser{}
	case "mysql-error":
		return &MySQLErrorParser{}
	case "custom-app":
		return &CustomAppParser{}
	case "json":
		return &JSONLogParser{}
	default:
		return &GenericLogParser{}
	}
}

// GenericLogParser parser générique
type GenericLogParser struct{}

func (p *GenericLogParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	if strings.TrimSpace(line) == "" {
		return nil, fmt.Errorf("ligne vide")
	}

	// Pattern générique: timestamp level message
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s+\[?(\w+)\]?\s+(.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 4 {
		// Fallback: essayer de détecter au moins le niveau
		level := extractLogLevel(line)
		return &config.LogEntry{
			Timestamp:  time.Now(),
			Level:      level,
			Message:    line,
			Source:     "generic",
			Fields:     make(map[string]string),
			LineNumber: lineNumber,
		}, nil
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		timestamp = time.Now()
	}

	return &config.LogEntry{
		Timestamp:  timestamp,
		Level:      strings.ToUpper(matches[2]),
		Message:    matches[3],
		Source:     "generic",
		Fields:     make(map[string]string),
		LineNumber: lineNumber,
	}, nil
}

// NginxAccessParser parser pour les logs d'accès Nginx
type NginxAccessParser struct{}

func (p *NginxAccessParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format Nginx: IP - - [timestamp] "method path protocol" status size "referer" "user-agent"
	re := regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "([^"]*)" (\d+) (\S+) "([^"]*)" "([^"]*)"`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 8 {
		return nil, fmt.Errorf("format Nginx access invalide")
	}

	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		timestamp = time.Now()
	}

	status, _ := strconv.Atoi(matches[4])
	level := "INFO"
	if status >= 500 {
		level = "ERROR"
	} else if status >= 400 {
		level = "WARN"
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   fmt.Sprintf("%s %s -> %s", matches[1], matches[3], matches[4]),
		Source:    "nginx-access",
		Fields: map[string]string{
			"ip":         matches[1],
			"request":    matches[3],
			"status":     matches[4],
			"size":       matches[5],
			"referer":    matches[6],
			"user_agent": matches[7],
		},
		LineNumber: lineNumber,
	}, nil
}

// NginxErrorParser parser pour les logs d'erreur Nginx
type NginxErrorParser struct{}

func (p *NginxErrorParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format Nginx error: timestamp [level] pid#tid: message
	re := regexp.MustCompile(`^(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (\d+)#(\d+): (.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 6 {
		return nil, fmt.Errorf("format Nginx error invalide")
	}

	timestamp, err := time.Parse("2006/01/02 15:04:05", matches[1])
	if err != nil {
		timestamp = time.Now()
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[5],
		Source:    "nginx-error",
		Fields: map[string]string{
			"pid": matches[3],
			"tid": matches[4],
		},
		LineNumber: lineNumber,
	}, nil
}

// ApacheAccessParser parser pour les logs d'accès Apache
type ApacheAccessParser struct{}

func (p *ApacheAccessParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format Apache commun: IP - - [timestamp] "method path protocol" status size
	re := regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "([^"]*)" (\d+) (\S+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 6 {
		return nil, fmt.Errorf("format Apache access invalide")
	}

	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		timestamp = time.Now()
	}

	status, _ := strconv.Atoi(matches[4])
	level := "INFO"
	if status >= 500 {
		level = "ERROR"
	} else if status >= 400 {
		level = "WARN"
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   fmt.Sprintf("%s %s -> %s", matches[1], matches[3], matches[4]),
		Source:    "apache-access",
		Fields: map[string]string{
			"ip":      matches[1],
			"request": matches[3],
			"status":  matches[4],
			"size":    matches[5],
		},
		LineNumber: lineNumber,
	}, nil
}

// ApacheErrorParser parser pour les logs d'erreur Apache
type ApacheErrorParser struct{}

func (p *ApacheErrorParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format Apache error: [timestamp] [level] [pid tid] message
	re := regexp.MustCompile(`^\[([^\]]+)\] \[([^\]]+)\] \[([^\]]+)\] (.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 5 {
		return nil, fmt.Errorf("format Apache error invalide")
	}

	timestamp, err := time.Parse("Mon Jan 02 15:04:05.000000 2006", matches[1])
	if err != nil {
		timestamp = time.Now()
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[4],
		Source:    "apache-error",
		Fields: map[string]string{
			"process": matches[3],
		},
		LineNumber: lineNumber,
	}, nil
}

// MySQLErrorParser parser pour les logs d'erreur MySQL
type MySQLErrorParser struct{}

func (p *MySQLErrorParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format MySQL: timestamp thread [level] message
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z)\s+(\d+)\s+\[(\w+)\]\s+(.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 5 {
		// Essayer un format plus simple
		re2 := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s+\[(\w+)\]\s+(.+)$`)
		matches2 := re2.FindStringSubmatch(line)
		if len(matches2) < 4 {
			return nil, fmt.Errorf("format MySQL error invalide")
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", matches2[1])
		if err != nil {
			timestamp = time.Now()
		}

		return &config.LogEntry{
			Timestamp:  timestamp,
			Level:      strings.ToUpper(matches2[2]),
			Message:    matches2[3],
			Source:     "mysql-error",
			Fields:     make(map[string]string),
			LineNumber: lineNumber,
		}, nil
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05.000000Z", matches[1])
	if err != nil {
		timestamp = time.Now()
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(matches[3]),
		Message:   matches[4],
		Source:    "mysql-error",
		Fields: map[string]string{
			"thread": matches[2],
		},
		LineNumber: lineNumber,
	}, nil
}

// CustomAppParser parser pour les logs d'application custom
type CustomAppParser struct{}

func (p *CustomAppParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	// Format custom: [timestamp] LEVEL component: message
	re := regexp.MustCompile(`^\[([^\]]+)\]\s+(\w+)\s+([^:]+):\s+(.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 5 {
		// Format alternatif plus simple
		re2 := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s+(\w+)\s+(.+)$`)
		matches2 := re2.FindStringSubmatch(line)
		if len(matches2) < 4 {
			return &config.LogEntry{
				Timestamp:  time.Now(),
				Level:      extractLogLevel(line),
				Message:    line,
				Source:     "custom-app",
				Fields:     make(map[string]string),
				LineNumber: lineNumber,
			}, nil
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", matches2[1])
		if err != nil {
			timestamp = time.Now()
		}

		return &config.LogEntry{
			Timestamp:  timestamp,
			Level:      strings.ToUpper(matches2[2]),
			Message:    matches2[3],
			Source:     "custom-app",
			Fields:     make(map[string]string),
			LineNumber: lineNumber,
		}, nil
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		timestamp = time.Now()
	}

	return &config.LogEntry{
		Timestamp: timestamp,
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[4],
		Source:    "custom-app",
		Fields: map[string]string{
			"component": matches[3],
		},
		LineNumber: lineNumber,
	}, nil
}

// JSONLogParser parser pour les logs JSON
type JSONLogParser struct{}

func (p *JSONLogParser) ParseLine(line string, lineNumber int) (*config.LogEntry, error) {
	var jsonLog map[string]interface{}

	if err := json.Unmarshal([]byte(line), &jsonLog); err != nil {
		return nil, fmt.Errorf("JSON invalide: %w", err)
	}

	entry := &config.LogEntry{
		Fields:     make(map[string]string),
		LineNumber: lineNumber,
		Source:     "json",
	}

	// Extraire les champs communs
	if ts, ok := jsonLog["timestamp"].(string); ok {
		if timestamp, err := time.Parse(time.RFC3339, ts); err == nil {
			entry.Timestamp = timestamp
		}
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	if level, ok := jsonLog["level"].(string); ok {
		entry.Level = strings.ToUpper(level)
	} else {
		entry.Level = "INFO"
	}

	if msg, ok := jsonLog["message"].(string); ok {
		entry.Message = msg
	} else if msg, ok := jsonLog["msg"].(string); ok {
		entry.Message = msg
	}

	// Convertir tous les autres champs
	for key, value := range jsonLog {
		if key != "timestamp" && key != "level" && key != "message" && key != "msg" {
			entry.Fields[key] = fmt.Sprintf("%v", value)
		}
	}

	return entry, nil
}

// extractLogLevel tente d'extraire le niveau de log d'une ligne
func extractLogLevel(line string) string {
	line = strings.ToUpper(line)
	levels := []string{"FATAL", "ERROR", "WARN", "WARNING", "INFO", "DEBUG", "TRACE"}

	for _, level := range levels {
		if strings.Contains(line, level) {
			return level
		}
	}
	return "INFO"
}
