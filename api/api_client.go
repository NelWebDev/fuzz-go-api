package api

import (
	"net/http"
	"strings"
)

// NewAPIClient crea un cliente con la URL base
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

type APIClient struct {
	BaseURL string
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
