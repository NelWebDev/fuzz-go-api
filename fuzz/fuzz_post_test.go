package fuzz

import (
	"fmt"
	"fuzzing-api/api"
	"fuzzing-api/logger"
	"testing"
)

func FuzzPostEndpoint(f *testing.F) {
	client := api.NewAPIClient("https://fakerestapi.azurewebsites.net/api/v1")
	log := &logger.Logger{}

	// Semillas iniciales
	f.Add(1, "Sample Title", "2025-01-13T12:24:47.807Z", true)
	f.Add(2, "Another Title", "2025-01-14T12:24:47.807Z", false)

	// Función de fuzzing
	f.Fuzz(func(t *testing.T, id int, title, dueDate string, completed bool) {
		// Crear el cuerpo de la solicitud POST
		data := fmt.Sprintf(`{"id": %d, "title": "%s", "dueDate": "%s", "completed": %v}`, id, title, dueDate, completed)

		// Realizar la petición POST
		_, err := client.Post("/Activities", data)
		if err != nil {
			t.Errorf("Error en POST: %v", err)
		}

		// Registrar la semilla
		log.LogSeed(data)
	})

	// Al terminar el fuzzing, generar el reporte HTML
	f.Cleanup(func() {
		err := log.GenerateHTMLReport("fuzz_post_report.html")
		if err != nil {
			// En lugar de t.Errorf, imprimimos el error en la consola
			fmt.Printf("Error al generar el reporte HTML: %v\n", err)
		}
	})
}
