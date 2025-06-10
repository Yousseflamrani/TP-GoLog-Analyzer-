package analyser

import (
	"errors"
	"fmt"
)

// FileNotFoundError erreur personnalisée pour fichier introuvable
type FileNotFoundError struct {
	Path string
	Err  error
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("fichier introuvable: %s (%v)", e.Path, e.Err)
}

func (e *FileNotFoundError) Unwrap() error {
	return e.Err
}

// FileAccessError erreur personnalisée pour fichier inaccessible
type FileAccessError struct {
	Path string
	Err  error
}

func (e *FileAccessError) Error() string {
	return fmt.Sprintf("fichier inaccessible: %s (%v)", e.Path, e.Err)
}

func (e *FileAccessError) Unwrap() error {
	return e.Err
}

// ParseError erreur personnalisée pour erreur de parsing
type ParseError struct {
	LogID   string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("erreur de parsing pour %s: %s", e.LogID, e.Message)
}

// IsFileNotFound vérifie si l'erreur est de type FileNotFoundError
func IsFileNotFound(err error) bool {
	var fileNotFoundErr *FileNotFoundError
	return errors.As(err, &fileNotFoundErr)
}

// IsFileAccess vérifie si l'erreur est de type FileAccessError
func IsFileAccess(err error) bool {
	var fileAccessErr *FileAccessError
	return errors.As(err, &fileAccessErr)
}

// IsParseError vérifie si l'erreur est de type ParseError
func IsParseError(err error) bool {
	var parseErr *ParseError
	return errors.As(err, &parseErr)
}
