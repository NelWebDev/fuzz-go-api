package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveEndpoint(t *testing.T) {
	client := NewAPIClient("https://example.com/api/v1")

	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "leading slash",
			endpoint: "/Activities",
			want:     "https://example.com/api/v1/Activities",
		},
		{
			name:     "without leading slash",
			endpoint: "Activities",
			want:     "https://example.com/api/v1/Activities",
		},
		{
			name:     "query string",
			endpoint: "/Activities?id=1",
			want:     "https://example.com/api/v1/Activities?id=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ResolveEndpoint(tt.endpoint)
			if err != nil {
				t.Fatalf("ResolveEndpoint returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ResolveEndpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetAndPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/Activities" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/api/v1/Activities")
		}

		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("ReadAll returned error: %v", err)
			}
			if string(body) != `{"id":1}` {
				t.Fatalf("body = %q, want %q", body, `{"id":1}`)
			}
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("method = %q", r.Method)
		}
	}))
	defer server.Close()

	client := NewAPIClient(server.URL + "/api/v1")

	resp, statusCode, _, err := client.Get("/Activities")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	resp.Body.Close()
	if statusCode != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", statusCode, http.StatusOK)
	}

	resp, statusCode, _, err = client.Post("/Activities", `{"id":1}`)
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	resp.Body.Close()
	if statusCode != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d", statusCode, http.StatusCreated)
	}

	resp, statusCode, _, err = client.Put("/Activities", `{"id":1}`)
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	resp.Body.Close()
	if statusCode != http.StatusOK {
		t.Fatalf("PUT status = %d, want %d", statusCode, http.StatusOK)
	}

	resp, statusCode, _, err = client.Patch("/Activities", `{"id":1}`)
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	resp.Body.Close()
	if statusCode != http.StatusOK {
		t.Fatalf("PATCH status = %d, want %d", statusCode, http.StatusOK)
	}

	resp, statusCode, _, err = client.Delete("/Activities")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	resp.Body.Close()
	if statusCode != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want %d", statusCode, http.StatusNoContent)
	}
}
