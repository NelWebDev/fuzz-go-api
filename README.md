# Fuzzing API Testing

This project is a small Go framework for fuzz testing HTTP API endpoints. It currently supports GET, POST, PUT, PATCH, and DELETE requests, reads endpoint settings from `config/config.json`, and logs each request outcome.

## Features

- **HTTP verb fuzzing**: Exercise GET, POST, PUT, PATCH, and DELETE endpoints with Go fuzz tests.
- **JSON configuration**: Manage the base URL, endpoints, and POST body from one file.
- **Request logging**: Log method, endpoint, seed, response status, duration, request body, and response body.
- **Reproducible findings**: Store request errors and 5xx responses as JSON Lines artifacts for later triage.
- **Local client tests**: Validate the API client without calling external services.

## Project Structure

```plaintext
fuzzing-api/
|-- api/
|   |-- api_client.go        # API client implementation
|   `-- api_client_test.go   # Local client tests with httptest
|-- config/
|   `-- config.json          # Base URL, endpoints, and request body
|-- fuzz/
|   |-- fuzz_get_test.go     # Fuzz tests for GET requests
|   |-- fuzz_post_test.go    # Fuzz tests for POST requests
|   `-- fuzz_write_methods_test.go # Fuzz tests for PUT, PATCH, and DELETE requests
|-- logger/
|   |-- logger.go            # Request logging and findings helpers
|   `-- logger_test.go       # Findings artifact tests
|-- utils/
|   |-- utils.go             # Configuration loading
|   `-- validator.go         # HTTP status validator helper
|-- go.mod
`-- README.md
```

## Configuration

The `config/config.json` file uses this shape:

```json
{
  "baseURL": "https://example.com/api/v1",
  "endpoints": {
    "get": "/resources",
    "post": "/resources",
    "put": "/resources/1",
    "patch": "/resources/1",
    "delete": "/resources/1"
  },
  "requestBody": {
    "id": 1,
    "title": "Test Activity",
    "dueDate": "2026-01-01T00:00:00Z",
    "completed": false
  }
}
```

## Usage

Run local package tests:

```bash
go test ./api ./logger ./utils
```

Run all tests:

```bash
go test ./...
```

Run fuzz tests:

```bash
FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzPostEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzPutEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzPatchEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzDeleteEndpoint -fuzztime=30s ./fuzz
```

On PowerShell:

```powershell
$env:FUZZ_API_EXTERNAL = "1"
go test -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
go test -fuzz=FuzzPostEndpoint -fuzztime=30s ./fuzz
go test -fuzz=FuzzPutEndpoint -fuzztime=30s ./fuzz
go test -fuzz=FuzzPatchEndpoint -fuzztime=30s ./fuzz
go test -fuzz=FuzzDeleteEndpoint -fuzztime=30s ./fuzz
```

The fuzz tests call the API configured in `config/config.json`, so they are opt-in and require network access plus a reachable target service.

When a fuzz test hits a request error or an HTTP 5xx response, it appends a reproducible finding to `artifacts/fuzz-findings.jsonl`. Override the output path with `FUZZ_API_FINDINGS`:

```bash
FUZZ_API_FINDINGS=./tmp/findings.jsonl FUZZ_API_EXTERNAL=1 go test -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
```

On PowerShell:

```powershell
$env:FUZZ_API_FINDINGS = ".\tmp\findings.jsonl"
$env:FUZZ_API_EXTERNAL = "1"
go test -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
```

## Requirements

- Go 1.20+
