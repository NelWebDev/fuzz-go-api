package logger

import (
	"os"
	"text/template"
	"time"
)

// Logger estructura para manejar las semillas
type Logger struct {
	Seeds []string
}

// LogSeed registra una nueva semilla
func (l *Logger) LogSeed(seed string) {
	l.Seeds = append(l.Seeds, seed)
}

// GenerateHTMLReport genera un archivo HTML con las semillas registradas
func (l *Logger) GenerateHTMLReport(filename string) error {
	reportTemplate := `<!DOCTYPE html>
<html>
<head>
	<title>Reporte de Fuzzing</title>
</head>
<body>
	<h1>Reporte de Fuzzing</h1>
	<p>Fecha: {{.Date}}</p>
	<h2>Semillas Generadas</h2>
	<ul>
	{{range .Seeds}}
		<li>{{.}}</li>
	{{end}}
	</ul>
</body>
</html>`

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return err
	}

	// Crear el archivo HTML
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el reporte en el archivo
	data := struct {
		Date  string
		Seeds []string
	}{
		Date:  time.Now().Format(time.RFC1123),
		Seeds: l.Seeds,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
