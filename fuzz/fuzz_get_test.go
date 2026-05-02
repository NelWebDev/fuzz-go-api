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

func FuzzGetEndpoint(f *testing.F) {
	if os.Getenv("FUZZ_API_EXTERNAL") != "1" {
		f.Skip("establece FUZZ_API_EXTERNAL=1 para ejecutar fuzzing contra la API configurada")
	}

	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuración: %v", err)
	}

	client := api.NewAPIClient(config.BaseURL)
	for _, seed := range []string{
		config.Endpoints.Get,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities?completed=true&completed=false",
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

		resp, statusCode, duration, err := client.Get(seed)
		if err != nil {
			logger.LogRequest("GET", requestURL, seed, 0, duration, "", fmt.Sprintf("Error: %v", err))
			t.Errorf("Error en la solicitud GET: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error al leer la respuesta GET: %v", err)
			return
		}
		logger.LogRequest("GET", requestURL, seed, statusCode, duration, "", string(body))

		// Manejo de códigos HTTP.
		if statusCode >= 500 {
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
