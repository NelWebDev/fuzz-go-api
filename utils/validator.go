package utils

import (
	"fmt"
	"log"
	"os"
)

type Validator struct {
	logger *log.Logger
}

// NewValidator crea un validador con un archivo de log para registrar códigos HTTP y semillas.
func NewValidator(logFile string) (*Validator, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de log: %v", err)
	}
	logger := log.New(file, "HTTP_LOG: ", log.LstdFlags)
	return &Validator{logger: logger}, nil
}

// ValidateAndLog registra el código de estado junto con la semilla utilizada y valida el rango esperado.
func (v *Validator) ValidateAndLog(status int, seed string, expectedRange string) bool {
	v.logger.Printf("Seed: %s, HTTP Status: %d", seed, status)

	switch expectedRange {
	case "2xx":
		return status >= 200 && status < 300
	case "4xx":
		return status >= 400 && status < 500
	case "5xx":
		return status >= 500 && status < 600
	default:
		return false
	}
}

// Close cierra el archivo de log.
func (v *Validator) Close() error {
	if logFile, ok := v.logger.Writer().(*os.File); ok {
		return logFile.Close()
	}
	return nil
}
