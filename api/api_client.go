package api

import (
	"net/http"
	"strings"
)

// APIClient representa el cliente de la API
type APIClient struct {
	BaseURL string
}

// NewAPIClient crea un cliente con la URL base
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

// GET request
func (c *APIClient) Get(endpoint string) (*http.Response, int, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, 0, err
	}
	return resp, resp.StatusCode, nil
}

// POST request
func (c *APIClient) Post(endpoint string, data string) (*http.Response, int, error) {
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(data))
	if err != nil {
		return nil, 0, err
	}
	return resp, resp.StatusCode, nil
}
