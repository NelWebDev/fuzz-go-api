package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	seedSet map[string]struct{} // Almacena semillas únicas
	file    string              // Nombre del archivo para persistencia
}

// Crear un nuevo Logger con persistencia
func NewLogger(file string) *Logger {
	logger := &Logger{
		seedSet: make(map[string]struct{}),
		file:    file,
	}
	logger.loadSeedsFromFile() // Cargar semillas existentes del archivo
	return logger
}

// Registrar una nueva semilla (si no existe)
func (l *Logger) LogSeed(seed string) {
	if _, exists := l.seedSet[seed]; !exists {
		l.seedSet[seed] = struct{}{}
		l.saveSeedsToFile() // Guardar en el archivo después de registrar
	}
}

// Obtener todas las semillas como un slice
func (l *Logger) GetSeeds() []string {
	seeds := make([]string, 0, len(l.seedSet))
	for seed := range l.seedSet {
		seeds = append(seeds, seed)
	}
	return seeds
}

// Generar un reporte HTML
func (l *Logger) GenerateHTMLReport(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	htmlContent := "<!DOCTYPE html><html><head><title>Reporte de Fuzzing</title></head><body>"
	htmlContent += fmt.Sprintf("<h1>Reporte de Fuzzing</h1><p>Fecha: %s</p>", time.Now().Format(time.RFC1123))
	htmlContent += "<h2>Semillas Generadas</h2><ul>"

	for _, seed := range l.GetSeeds() {
		htmlContent += fmt.Sprintf("<li>%s</li>", seed)
	}
	htmlContent += "</ul></body></html>"

	_, err = file.WriteString(htmlContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	fmt.Printf("Reporte HTML generado en: %s\n", filename)
	return nil
}

// Guardar semillas en un archivo JSON
func (l *Logger) saveSeedsToFile() error {
	data, err := json.MarshalIndent(l.GetSeeds(), "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar semillas: %v", err)
	}

	err = os.WriteFile(l.file, data, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar semillas en el archivo: %v", err)
	}

	return nil
}

// Cargar semillas desde un archivo JSON
func (l *Logger) loadSeedsFromFile() error {
	if _, err := os.Stat(l.file); os.IsNotExist(err) {
		return nil // Si el archivo no existe, no hacemos nada
	}

	data, err := os.ReadFile(l.file)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de semillas: %v", err)
	}

	var seeds []string
	err = json.Unmarshal(data, &seeds)
	if err != nil {
		return fmt.Errorf("error al deserializar semillas: %v", err)
	}

	for _, seed := range seeds {
		l.seedSet[seed] = struct{}{}
	}

	return nil
}
