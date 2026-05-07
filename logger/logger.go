package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const defaultFindingsPath = "../artifacts/fuzz-findings.jsonl"

var findingsMu sync.Mutex

type Finding struct {
	Timestamp    string `json:"timestamp"`
	Method       string `json:"method"`
	Endpoint     string `json:"endpoint"`
	Seed         string `json:"seed"`
	StatusCode   int    `json:"statusCode"`
	DurationMS   int64  `json:"durationMs"`
	RequestBody  string `json:"requestBody,omitempty"`
	ResponseBody string `json:"responseBody,omitempty"`
	Error        string `json:"error,omitempty"`
}

// LogRequest registra las solicitudes HTTP.
func LogRequest(method, endpoint, seed string, statusCode int, duration time.Duration, requestBody, responseBody string) {
	logEntry := fmt.Sprintf(
		"Time: %s\nMethod: %s\nEndpoint: %s\nSeed: %s\nRequest body: %s\nHTTP status: %d\nDuration: %v\nResponse: %s\n\n",
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

// LogFinding appends a reproducible fuzz finding as JSON Lines.
func LogFinding(method, endpoint, seed string, statusCode int, duration time.Duration, requestBody, responseBody, errText string) error {
	path := os.Getenv("FUZZ_API_FINDINGS")
	if path == "" {
		path = defaultFindingsPath
	}

	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating findings directory: %w", err)
		}
	}

	findingsMu.Lock()
	defer findingsMu.Unlock()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening findings file: %w", err)
	}
	defer file.Close()

	finding := Finding{
		Timestamp:    time.Now().Format(time.RFC3339),
		Method:       method,
		Endpoint:     endpoint,
		Seed:         seed,
		StatusCode:   statusCode,
		DurationMS:   duration.Milliseconds(),
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		Error:        errText,
	}

	if err := json.NewEncoder(file).Encode(finding); err != nil {
		return fmt.Errorf("error writing finding: %w", err)
	}

	return nil
}
