package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

// Crear un nuevo cliente API
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Método para realizar una petición GET
func (a *APIClient) Get(endpoint string) ([]byte, int, error) {
	url := a.baseURL + endpoint
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// Método para realizar una petición POST
func (a *APIClient) Post(endpoint, data string) ([]byte, int, error) {
	url := a.baseURL + endpoint
	resp, err := a.client.Post(url, "application/json", bytes.NewBufferString(data))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
