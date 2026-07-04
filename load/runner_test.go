package load

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fuzzing-api/api"
)

func TestRunExecutesConfiguredRequests(t *testing.T) {
	var seen int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := api.NewAPIClient(server.URL)
	report, err := Run(context.Background(), client, []Scenario{
		{Name: "health", Method: http.MethodGet, Endpoint: "/health"},
	}, RunnerOptions{
		Users:    2,
		Requests: 5,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if seen != 5 {
		t.Fatalf("seen = %d, want 5", seen)
	}
	if report.TotalRequests != 5 || report.Successful != 5 || report.Failed != 0 {
		t.Fatalf("report = %+v", report)
	}
	if len(report.Samples) != 5 {
		t.Fatalf("samples = %d, want 5", len(report.Samples))
	}
}

func TestRunUsesDefaultRequestCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := api.NewAPIClient(server.URL)
	report, err := Run(context.Background(), client, []Scenario{
		{Name: "one", Method: http.MethodPost, Endpoint: "/one", Body: []byte(`{"id":1}`)},
		{Name: "two", Method: http.MethodGet, Endpoint: "/two"},
	}, RunnerOptions{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if report.TotalRequests != 2 {
		t.Fatalf("TotalRequests = %d, want 2", report.TotalRequests)
	}
	if report.Users != 1 {
		t.Fatalf("Users = %d, want 1", report.Users)
	}
}

func TestRunStopsOnDuration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := api.NewAPIClient(server.URL)
	report, err := Run(context.Background(), client, []Scenario{
		{Name: "loop", Method: http.MethodGet, Endpoint: "/loop"},
	}, RunnerOptions{
		Users:    1,
		Duration: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if report.TotalRequests == 0 {
		t.Fatal("TotalRequests = 0, want at least one request")
	}
}
