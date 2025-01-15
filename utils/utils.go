package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config representa la estructura del archivo de configuración
type Config struct {
	BaseURL   string `json:"baseURL"`
	Endpoints struct {
		Get  string `json:"get"`
		Post string `json:"post"`
	} `json:"endpoints"`
	Bodies struct {
		Post interface{} `json:"post"`
	} `json:"bodies"`
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de configuración: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error al decodificar el archivo de configuración: %w", err)
	}

	return &config, nil
}
