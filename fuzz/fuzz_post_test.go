package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"fuzzing-api/utils"
	"io"
	"os"
	"testing"
)

func FuzzPostEndpoint(f *testing.F) {
	if os.Getenv("FUZZ_API_EXTERNAL") != "1" {
		f.Skip("establece FUZZ_API_EXTERNAL=1 para ejecutar fuzzing contra la API configurada")
	}

	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuración: %v", err)
	}

	client := api.NewAPIClient(config.BaseURL)
	addEndpointBodySeeds(f.Add, postEndpointSeeds(config), postBodySeeds(config))

	f.Fuzz(func(t *testing.T, seed string, requestBody string) {
		requestURL, err := client.ResolveEndpoint(seed)
		if err != nil {
			t.Skipf("Semilla con endpoint inválido %q: %v", seed, err)
		}

		resp, statusCode, duration, err := client.Post(seed, requestBody)
		if err != nil {
			logger.LogRequest("POST", requestURL, seed, 0, duration, requestBody, fmt.Sprintf("Error: %v", err))
			if logErr := logger.LogFinding("POST", requestURL, seed, 0, duration, requestBody, "", err.Error()); logErr != nil {
				t.Logf("Error al registrar el hallazgo POST: %v", logErr)
			}
			t.Errorf("Error en la solicitud POST: %v", err)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error al leer la respuesta POST: %v", err)
			return
		}
		logger.LogRequest("POST", requestURL, seed, statusCode, duration, requestBody, string(respBody))

		// Manejo de códigos HTTP.
		if statusCode >= 500 {
			if logErr := logger.LogFinding("POST", requestURL, seed, statusCode, duration, requestBody, string(respBody), ""); logErr != nil {
				t.Logf("Error al registrar el hallazgo POST: %v", logErr)
			}
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
