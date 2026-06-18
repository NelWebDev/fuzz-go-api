# Fuzzing API Testing

This project is a small Go framework for fuzz testing HTTP API endpoints. It currently supports GET, POST, PUT, PATCH, and DELETE requests, reads endpoint settings from `config/config.json`, and logs each request outcome.

## Features

- **HTTP verb fuzzing**: Exercise GET, POST, PUT, PATCH, and DELETE endpoints with separate Go fuzz targets.
- **Verb-specific seeds**: Keep endpoint and request-body seeds tuned to each HTTP method.
- **Shared fuzz execution**: Run all HTTP methods through common helpers for consistent logging, finding capture, and 5xx handling.
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
|   |-- fuzz_get_test.go     # GET fuzz target wrapper
|   |-- fuzz_post_test.go    # POST fuzz target wrapper
|   |-- fuzz_write_methods_test.go # Shared fuzz helpers plus PUT, PATCH, and DELETE wrappers
|   `-- seeds_test.go        # Verb-specific endpoint and body seeds
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
FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzPostEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzPutEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzPatchEndpoint -fuzztime=30s ./fuzz
FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzDeleteEndpoint -fuzztime=30s ./fuzz
```

On PowerShell:

```powershell
$env:FUZZ_API_EXTERNAL = "1"
go test -run=^$ -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
go test -run=^$ -fuzz=FuzzPostEndpoint -fuzztime=30s ./fuzz
go test -run=^$ -fuzz=FuzzPutEndpoint -fuzztime=30s ./fuzz
go test -run=^$ -fuzz=FuzzPatchEndpoint -fuzztime=30s ./fuzz
go test -run=^$ -fuzz=FuzzDeleteEndpoint -fuzztime=30s ./fuzz
```

The fuzz tests call the API configured in `config/config.json`, so they are opt-in and require network access plus a reachable target service.
The `-run=^$` flag keeps the command focused on the selected fuzz target instead of running the seed corpus for every fuzz test in the package first.

Each fuzz target is intentionally run separately because each HTTP verb has different inputs worth mutating:

| Target | Initial seeds | What it mutates |
| --- | ---: | --- |
| `FuzzGetEndpoint` | 11 | Endpoint paths, IDs, query strings, and encoded path segments |
| `FuzzPostEndpoint` | 17 | POST endpoints plus request bodies for creation scenarios |
| `FuzzPutEndpoint` | 17 | PUT endpoints plus full replacement request bodies |
| `FuzzPatchEndpoint` | 19 | PATCH endpoints plus partial, null, unknown-field, and malformed bodies |
| `FuzzDeleteEndpoint` | 10 | Endpoint paths, IDs, query strings, and encoded path segments |

`POST`, `PUT`, and `PATCH` use two fuzz arguments: the endpoint seed and the request-body seed. The seed corpus keeps the campaign focused by testing all endpoint variants with the configured body, then testing body variants against the configured endpoint.

When a fuzz test hits a request error or an HTTP 5xx response, it appends a reproducible finding to `artifacts/fuzz-findings.jsonl`. Override the output path with `FUZZ_API_FINDINGS`:

```bash
FUZZ_API_FINDINGS=./tmp/findings.jsonl FUZZ_API_EXTERNAL=1 go test -run=^$ -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
```

On PowerShell:

```powershell
$env:FUZZ_API_FINDINGS = ".\tmp\findings.jsonl"
$env:FUZZ_API_EXTERNAL = "1"
go test -run=^$ -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
```

## Requirements

- Go 1.20+
