package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogFindingWritesJSONL(t *testing.T) {
	findingsPath := filepath.Join(t.TempDir(), "findings.jsonl")
	t.Setenv("FUZZ_API_FINDINGS", findingsPath)

	err := LogFinding("GET", "https://example.com/api", "/api", 500, 1500*time.Millisecond, "", "boom", "")
	if err != nil {
		t.Fatalf("LogFinding returned error: %v", err)
	}

	data, err := os.ReadFile(findingsPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	var finding Finding
	if err := json.Unmarshal(data, &finding); err != nil {
		t.Fatalf("finding is not valid JSON: %v", err)
	}

	if finding.Method != "GET" || finding.StatusCode != 500 || finding.Seed != "/api" {
		t.Fatalf("finding = %+v", finding)
	}
	if finding.DurationMS != 1500 {
		t.Fatalf("DurationMS = %d, want 1500", finding.DurationMS)
	}
}
