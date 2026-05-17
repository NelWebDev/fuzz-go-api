package fuzz

import (
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
	fuzzEndpointWithBody(f, http.MethodPut, putEndpointSeeds, putBodySeeds)
}

func FuzzPatchEndpoint(f *testing.F) {
	fuzzEndpointWithBody(f, http.MethodPatch, patchEndpointSeeds, patchBodySeeds)
}

func FuzzDeleteEndpoint(f *testing.F) {
	fuzzEndpointWithoutBody(f, http.MethodDelete, deleteEndpointSeeds)
}

func fuzzEndpointWithBody(f *testing.F, method string, endpointSeeds func(*utils.Config) []string, bodySeeds func(*utils.Config) []string) {
	if os.Getenv("FUZZ_API_EXTERNAL") != "1" {
		f.Skip("establece FUZZ_API_EXTERNAL=1 para ejecutar fuzzing contra la API configurada")
	}

	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuracion: %v", err)
	}

	addEndpointBodySeeds(f.Add, endpointSeeds(config), bodySeeds(config))

	client := api.NewAPIClient(config.BaseURL)
	f.Fuzz(func(t *testing.T, seed string, requestBody string) {
		requestURL, err := client.ResolveEndpoint(seed)
		if err != nil {
			t.Skipf("Semilla con endpoint invalido %q: %v", seed, err)
		}

		resp, statusCode, duration, err := client.Request(method, seed, requestBody)
		if err != nil {
			logger.LogRequest(method, requestURL, seed, 0, duration, requestBody, fmt.Sprintf("Error: %v", err))
			if logErr := logger.LogFinding(method, requestURL, seed, 0, duration, requestBody, "", err.Error()); logErr != nil {
				t.Logf("Error al registrar el hallazgo %s: %v", method, logErr)
			}
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
			if logErr := logger.LogFinding(method, requestURL, seed, statusCode, duration, requestBody, string(responseBody), ""); logErr != nil {
				t.Logf("Error al registrar el hallazgo %s: %v", method, logErr)
			}
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}

func fuzzEndpointWithoutBody(f *testing.F, method string, endpointSeeds func(*utils.Config) []string) {
	if os.Getenv("FUZZ_API_EXTERNAL") != "1" {
		f.Skip("establece FUZZ_API_EXTERNAL=1 para ejecutar fuzzing contra la API configurada")
	}

	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuracion: %v", err)
	}

	for _, seed := range endpointSeeds(config) {
		f.Add(seed)
	}

	client := api.NewAPIClient(config.BaseURL)
	f.Fuzz(func(t *testing.T, seed string) {
		requestURL, err := client.ResolveEndpoint(seed)
		if err != nil {
			t.Skipf("Semilla con endpoint invalido %q: %v", seed, err)
		}

		resp, statusCode, duration, err := client.Request(method, seed, "")
		if err != nil {
			logger.LogRequest(method, requestURL, seed, 0, duration, "", fmt.Sprintf("Error: %v", err))
			if logErr := logger.LogFinding(method, requestURL, seed, 0, duration, "", "", err.Error()); logErr != nil {
				t.Logf("Error al registrar el hallazgo %s: %v", method, logErr)
			}
			t.Errorf("Error en la solicitud %s: %v", method, err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error al leer la respuesta %s: %v", method, err)
			return
		}
		logger.LogRequest(method, requestURL, seed, statusCode, duration, "", string(responseBody))

		if statusCode >= 500 {
			if logErr := logger.LogFinding(method, requestURL, seed, statusCode, duration, "", string(responseBody), ""); logErr != nil {
				t.Logf("Error al registrar el hallazgo %s: %v", method, logErr)
			}
			t.Errorf("Error del servidor: %d para la semilla: %s", statusCode, seed)
		}
	})
}
