package api

import (
	"bytes"
	"net/http"
	"time"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// Realizar solicitud GET
func (c *APIClient) Get(endpoint string) (*http.Response, int, time.Duration, error) {
	start := time.Now()
	resp, err := c.Client.Get(endpoint)
	duration := time.Since(start)

	if err != nil {
		return nil, 0, duration, err
	}

	return resp, resp.StatusCode, duration, nil
}

// Realizar solicitud POST
func (c *APIClient) Post(endpoint string, data string) (*http.Response, int, time.Duration, error) {
	start := time.Now()
	resp, err := c.Client.Post(endpoint, "application/json", bytes.NewBufferString(data))
	duration := time.Since(start)

	if err != nil {
		return nil, 0, duration, err
	}

	return resp, resp.StatusCode, duration, nil
}
