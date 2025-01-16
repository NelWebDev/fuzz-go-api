package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"fuzzing-api/utils"
	"io/ioutil"
	"net/url"
	"testing"
)

func FuzzGetEndpoint(f *testing.F) {
	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuración: %v", err)
	}

	client := api.NewAPIClient(config.BaseURL)
	f.Add(config.Endpoints.Get)

	f.Fuzz(func(t *testing.T, seed string) {
		escapedSeed := url.QueryEscape(seed)
		url := fmt.Sprintf("%s%s", config.BaseURL, escapedSeed)

		resp, statusCode, duration, err := client.Get(url)
		if err != nil {
			logger.LogRequest("GET", url, seed, 0, duration, "", fmt.Sprintf("Error: %v", err))
			t.Errorf("Error en la solicitud GET: %v", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		logger.LogRequest("GET", url, seed, statusCode, duration, "", string(body))

		// Manejo de códigos HTTP
		if statusCode >= 500 {
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
