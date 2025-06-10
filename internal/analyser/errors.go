package analyser

import "fmt"

// ParseError erreur de parsing personnalisée
type ParseError struct {
	Line       int
	Content    string
	Reason     string
	ParserType string
	SourceID   string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("erreur de parsing ligne %d dans %s (%s): %s",
		e.Line, e.SourceID, e.ParserType, e.Reason)
}

// NewParseError crée une nouvelle erreur de parsing
func NewParseError(line int, content, reason, parserType, sourceID string) *ParseError {
	return &ParseError{
		Line:       line,
		Content:    content,
		Reason:     reason,
		ParserType: parserType,
		SourceID:   sourceID,
	}
}

// SourceError erreur liée à une source de log
type SourceError struct {
	SourceID   string
	SourcePath string
	Reason     string
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("erreur source %s (%s): %s", e.SourceID, e.SourcePath, e.Reason)
}
