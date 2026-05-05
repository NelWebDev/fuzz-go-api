package fuzz

import (
	"encoding/json"
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"fuzzing-api/utils"
	"io"
	"net/http"
	"os"
	"testing"
)

func FuzzPutEndpoint(f *testing.F) {
	fuzzEndpoint(f, http.MethodPut, configuredPutEndpoint, true)
}

func FuzzPatchEndpoint(f *testing.F) {
	fuzzEndpoint(f, http.MethodPatch, configuredPatchEndpoint, true)
}

func FuzzDeleteEndpoint(f *testing.F) {
	fuzzEndpoint(f, http.MethodDelete, configuredDeleteEndpoint, false)
}

func fuzzEndpoint(f *testing.F, method string, configuredEndpoint func(*utils.Config) string, includeBody bool) {
	if os.Getenv("FUZZ_API_EXTERNAL") != "1" {
		f.Skip("establece FUZZ_API_EXTERNAL=1 para ejecutar fuzzing contra la API configurada")
	}

	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuracion: %v", err)
	}

	for _, seed := range writeMethodSeeds(configuredEndpoint(config)) {
		f.Add(seed)
	}

	client := api.NewAPIClient(config.BaseURL)
	f.Fuzz(func(t *testing.T, seed string) {
		requestURL, err := client.ResolveEndpoint(seed)
		if err != nil {
			t.Skipf("Semilla con endpoint invalido %q: %v", seed, err)
		}

		var requestBody string
		if includeBody {
			bodyJSON, err := json.Marshal(config.RequestBody)
			if err != nil {
				t.Fatalf("Error al serializar el cuerpo %s: %v", method, err)
			}
			requestBody = string(bodyJSON)
		}

		resp, statusCode, duration, err := client.Request(method, seed, requestBody)
		if err != nil {
			logger.LogRequest(method, requestURL, seed, 0, duration, requestBody, fmt.Sprintf("Error: %v", err))
			t.Errorf("Error en la solicitud %s: %v", method, err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error al leer la respuesta %s: %v", method, err)
			return
		}
		logger.LogRequest(method, requestURL, seed, statusCode, duration, requestBody, string(responseBody))

		if statusCode >= 500 {
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}

func writeMethodSeeds(configuredEndpoint string) []string {
	return []string{
		configuredEndpoint,
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
	}
}

func configuredPutEndpoint(config *utils.Config) string {
	return config.Endpoints.Put
}

func configuredPatchEndpoint(config *utils.Config) string {
	return config.Endpoints.Patch
}

func configuredDeleteEndpoint(config *utils.Config) string {
	return config.Endpoints.Delete
}
