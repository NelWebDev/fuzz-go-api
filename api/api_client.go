package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// ResolveEndpoint builds a request URL from the configured base URL and a relative endpoint.
func (c *APIClient) ResolveEndpoint(endpoint string) (string, error) {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("base URL must include scheme and host")
	}

	ref, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %w", err)
	}

	resolved := *base
	resolved.Path = joinPath(base.Path, ref.Path)
	resolved.RawQuery = ref.RawQuery
	resolved.Fragment = ""

	return resolved.String(), nil
}

func joinPath(basePath, endpointPath string) string {
	joined := path.Join(basePath, endpointPath)
	if joined == "." {
		return ""
	}
	if strings.HasSuffix(endpointPath, "/") && !strings.HasSuffix(joined, "/") {
		joined += "/"
	}
	if strings.HasPrefix(basePath, "/") && !strings.HasPrefix(joined, "/") {
		joined = "/" + joined
	}
	return joined
}

// Get performs a GET request against an endpoint relative to BaseURL.
func (c *APIClient) Get(endpoint string) (*http.Response, int, time.Duration, error) {
	return c.Request(http.MethodGet, endpoint, "")
}

// Post performs a POST request against an endpoint relative to BaseURL.
func (c *APIClient) Post(endpoint string, data string) (*http.Response, int, time.Duration, error) {
	return c.Request(http.MethodPost, endpoint, data)
}

// Put performs a PUT request against an endpoint relative to BaseURL.
func (c *APIClient) Put(endpoint string, data string) (*http.Response, int, time.Duration, error) {
	return c.Request(http.MethodPut, endpoint, data)
}

// Patch performs a PATCH request against an endpoint relative to BaseURL.
func (c *APIClient) Patch(endpoint string, data string) (*http.Response, int, time.Duration, error) {
	return c.Request(http.MethodPatch, endpoint, data)
}

// Delete performs a DELETE request against an endpoint relative to BaseURL.
func (c *APIClient) Delete(endpoint string) (*http.Response, int, time.Duration, error) {
	return c.Request(http.MethodDelete, endpoint, "")
}

// Request performs an HTTP request against an endpoint relative to BaseURL.
func (c *APIClient) Request(method, endpoint string, data string) (*http.Response, int, time.Duration, error) {
	requestURL, err := c.ResolveEndpoint(endpoint)
	if err != nil {
		return nil, 0, 0, err
	}

	var body io.Reader
	if data != "" {
		body = bytes.NewBufferString(data)
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, 0, 0, err
	}
	if data != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	resp, err := c.Client.Do(req)
	duration := time.Since(start)
	if err != nil {
		return nil, 0, duration, err
	}

	return resp, resp.StatusCode, duration, nil
}
