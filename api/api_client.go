package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIClient estructura para interactuar con la API
type APIClient struct {
	BaseURL string
	client  *http.Client
}

// NewAPIClient crea una nueva instancia de APIClient
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		client:  &http.Client{},
	}
}

// Get realiza una solicitud GET
func (api *APIClient) Get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", api.BaseURL, endpoint)
	resp, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Post realiza una solicitud POST
func (api *APIClient) Post(endpoint string, data string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", api.BaseURL, endpoint)
	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
