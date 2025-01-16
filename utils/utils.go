package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseURL   string `json:"baseURL"`
	Endpoints struct {
		Get  string `json:"get"`
		Post string `json:"post"`
	} `json:"endpoints"`
	RequestBody map[string]interface{} `json:"requestBody"`
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
