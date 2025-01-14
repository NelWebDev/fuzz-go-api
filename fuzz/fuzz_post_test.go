package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"testing"
)

func FuzzPostEndpoint(f *testing.F) {
	client := api.NewAPIClient("https://fakerestapi.azurewebsites.net/api/v1")
	log := logger.NewLogger("generatedSeeds/post_seeds.json")

	// Semillas iniciales
	f.Add(1, "Title 1", "2025-01-15T12:00:00Z", true)
	f.Add(2, "Title 2", "invalid_date", false)

	f.Fuzz(func(t *testing.T, id int, title, dueDate string, completed bool) {
		body := fmt.Sprintf(`{"id": %d, "title": "%s", "dueDate": "%s", "completed": %v}`, id, title, dueDate, completed)

		respBody, statusCode, err := client.Post("/Activities", body)
		if err != nil {
			t.Logf("Error en POST: %v, Body: %s", err, body)
			return
		}

		if statusCode >= 400 {
			t.Logf("POST falló con estado %d, cuerpo: %s", statusCode, string(respBody))
			return
		}

		// Registrar solo semillas nuevas
		log.LogSeed(body)
	})

	// Generar el reporte HTML en la carpeta /fuzz/reports
	defer log.GenerateHTMLReport("reports/fuzz_post_report.html")
}
