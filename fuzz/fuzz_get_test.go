package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/utils"
	"net/url"
	"testing"
)

func FuzzGetEndpoint(f *testing.F) {
	// Cargar la configuración
	config, err := utils.LoadConfig("../config/config.json")
	if err != nil {
		f.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Crear el cliente API con la baseURL desde la configuración
	client := api.NewAPIClient(config.BaseURL)

	// Semilla para fuzzing
	f.Add(config.Endpoints.Get)

	f.Fuzz(func(t *testing.T, endpoint string) {
		// Escapar la URL para evitar problemas con caracteres especiales
		escapedEndpoint := url.QueryEscape(endpoint)
		url := fmt.Sprintf("%s%s", config.BaseURL, escapedEndpoint)

		// Realizar la solicitud GET
		resp, status, err := client.Get(url)
		if err != nil {
			t.Errorf("Error en la solicitud GET: %v", err)
			return
		}
		defer resp.Body.Close()

		// Comprobar el código de estado
		if resp.StatusCode != status {
			t.Errorf("Código de estado inesperado: %d, esperado: %d", resp.StatusCode, status)
		}
	})
}
