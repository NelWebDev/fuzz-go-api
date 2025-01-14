package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

// SaveSeeds guarda las semillas generadas en un archivo JSON
func SaveSeeds(filename string, seeds []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error al crear el archivo de semillas: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(seeds); err != nil {
		return fmt.Errorf("error al guardar las semillas: %w", err)
	}
	return nil
}

// GenerateHTMLReport genera un reporte HTML a partir de las semillas
func GenerateHTMLReport(filename string, seeds []string) error {
	tmpl := `<html>
		<head><title>Fuzzing Report</title></head>
		<body>
		<h1>Reporte de Fuzzing</h1>
		<ul>
			{{range .}}
				<li>{{.}}</li>
			{{end}}
		</ul>
		</body>
		</html>`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error al parsear el template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error al crear el archivo HTML: %w", err)
	}
	defer file.Close()

	if err := t.Execute(file, seeds); err != nil {
		return fmt.Errorf("error al generar el reporte HTML: %w", err)
	}

	return nil
}
