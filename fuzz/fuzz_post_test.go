package fuzz

import (
	"encoding/json"
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"fuzzing-api/utils"
	"io/ioutil"
	"net/url"
	"testing"
)

func FuzzPostEndpoint(f *testing.F) {
	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuración: %v", err)
	}

	client := api.NewAPIClient(config.BaseURL)
	f.Add(config.Endpoints.Post)

	f.Fuzz(func(t *testing.T, seed string) {
		escapedSeed := url.QueryEscape(seed)
		url := fmt.Sprintf("%s%s", config.BaseURL, escapedSeed)

		body := config.RequestBody
		bodyJSON, _ := json.Marshal(body)

		resp, statusCode, duration, err := client.Post(url, string(bodyJSON))
		if err != nil {
			logger.LogRequest("POST", url, seed, 0, duration, string(bodyJSON), fmt.Sprintf("Error: %v", err))
			t.Errorf("Error en la solicitud POST: %v", err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)
		logger.LogRequest("POST", url, seed, statusCode, duration, string(bodyJSON), string(respBody))

		// Manejo de códigos HTTP
		if statusCode >= 500 {
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
