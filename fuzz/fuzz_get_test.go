package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"net/url"
	"testing"
)

func FuzzGetEndpoint(f *testing.F) {
	client := api.NewAPIClient("https://fakerestapi.azurewebsites.net/api/v1")
	log := logger.NewLogger("generatedSeeds/get_seeds.json")

	// Semillas iniciales
	f.Add("1")
	f.Add("2")
	f.Add("invalid_id")

	f.Fuzz(func(t *testing.T, id string) {
		escapedID := url.QueryEscape(id)
		endpoint := fmt.Sprintf("/Activities/%s", escapedID)

		body, statusCode, err := client.Get(endpoint)
		if err != nil {
			t.Logf("Error en GET: %v, ID: %s", err, id)
			return
		}

		if statusCode >= 400 {
			t.Logf("GET falló con estado %d, cuerpo: %s", statusCode, string(body))
			return
		}

		// Registrar solo semillas nuevas
		log.LogSeed(endpoint)
	})

	// Generar el reporte HTML en la carpeta /fuzz/reports
	defer log.GenerateHTMLReport("reports/fuzz_get_report.html")
}
