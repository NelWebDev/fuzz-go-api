package logger

import (
	"fmt"
	"log"
	"time"
)

// LogRequest registra las solicitudes HTTP en un archivo de log.
func LogRequest(method, endpoint, seed string, statusCode int, duration time.Duration, requestBody, responseBody string) {
	logEntry := fmt.Sprintf(
		"Hora: %s\nMétodo: %s\nEndpoint: %s\nSemilla: %s\nCuerpo de la solicitud: %s\nCódigo HTTP: %d\nDuración: %v\nRespuesta: %s\n\n",
		time.Now().Format(time.RFC3339),
		method,
		endpoint,
		seed,
		requestBody,
		statusCode,
		duration,
		responseBody,
	)
	log.Println(logEntry)
}
