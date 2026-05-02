package fuzz

import (
	"encoding/json"
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
	for _, seed := range []string{
		config.Endpoints.Post,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities?validate=true&validate=false",
		"/Activities?page=-1&pageSize=999999",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, seed string) {
		requestURL, err := client.ResolveEndpoint(seed)
		if err != nil {
			t.Skipf("Semilla con endpoint inválido %q: %v", seed, err)
		}

		bodyJSON, err := json.Marshal(config.RequestBody)
		if err != nil {
			t.Fatalf("Error al serializar el cuerpo POST: %v", err)
		}

		resp, statusCode, duration, err := client.Post(seed, string(bodyJSON))
		if err != nil {
			logger.LogRequest("POST", requestURL, seed, 0, duration, string(bodyJSON), fmt.Sprintf("Error: %v", err))
			t.Errorf("Error en la solicitud POST: %v", err)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error al leer la respuesta POST: %v", err)
			return
		}
		logger.LogRequest("POST", requestURL, seed, statusCode, duration, string(bodyJSON), string(respBody))

		// Manejo de códigos HTTP.
		if statusCode >= 500 {
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
