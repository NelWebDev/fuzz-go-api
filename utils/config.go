package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

// Config es la estructura que almacenará la configuración desde el archivo JSON.
type Config struct {
	BaseURL   string `json:"baseURL"`
	Endpoints struct {
		Get  string `json:"get"`
		Post string `json:"post"`
	} `json:"endpoints"`
}

// LoadConfig carga la configuración desde un archivo JSON.
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de configuración: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error al cargar la configuración: %v", err)
	}
	return &config, nil
}

// EscaparURL asegura que cualquier URL esté correctamente codificada.
func EscaparURL(endpoint string) string {
	return url.PathEscape(endpoint)
}
