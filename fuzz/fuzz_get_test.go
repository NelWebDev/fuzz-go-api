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
	log := &logger.Logger{}

	// Semillas iniciales
	f.Add("1")
	f.Add("2")
	f.Add("3")

	// Función de fuzzing
	f.Fuzz(func(t *testing.T, id string) {
		// Escapar el id para asegurar que sea una URL válida
		escapedID := url.QueryEscape(id)

		// Crear la URL para la petición GET usando el ID escapado
		endpoint := fmt.Sprintf("/Activities/%s", escapedID)

		// Realizar la petición GET
		_, err := client.Get(endpoint)
		if err != nil {
			t.Errorf("Error en GET: %v", err)
		}

		// Registrar la semilla
		log.LogSeed(endpoint)
	})

	// Al terminar el fuzzing, generar el reporte HTML
	f.Cleanup(func() {
		err := log.GenerateHTMLReport("fuzz_get_report.html")
		if err != nil {
			// En lugar de t.Errorf, imprimimos el error en la consola
			fmt.Printf("Error al generar el reporte HTML: %v\n", err)
		}
	})
}
