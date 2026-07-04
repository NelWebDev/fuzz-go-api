# Fuzzing API Testing

This project is a small Go framework for testing HTTP API endpoints. It currently has two separated testing paths:

- **Fuzz testing**: Mutates endpoints and request bodies for GET, POST, PUT, PATCH, and DELETE.
- **Load testing**: Runs a small set of configured API calls concurrently and reports basic performance metrics.

## Features

- **HTTP verb fuzzing**: Exercise GET, POST, PUT, PATCH, and DELETE endpoints with separate Go fuzz targets.
- **Verb-specific seeds**: Keep endpoint and request-body seeds tuned to each HTTP method.
- **Shared fuzz execution**: Run all HTTP methods through common helpers for consistent logging, finding capture, and 5xx handling.
- **JSON configuration**: Manage the base URL, endpoints, and POST body from one file.
- **Request logging**: Log method, endpoint, seed, response status, duration, request body, and response body.
- **Reproducible findings**: Store request errors and 5xx responses as JSON Lines artifacts for later triage.
- **Basic load testing**: Run configured scenarios with concurrent users, fixed request counts, or duration-based execution.
- **Load reports**: Store per-request JSON Lines logs plus summary metrics such as success count, failures, RPS, average latency, p95, and p99.
- **Separated or combined runs**: Execute fuzzing, load testing, or both while keeping artifacts organized by test type.
- **Local client tests**: Validate the API client without calling external services.

## Project Structure

```plaintext
fuzzing-api/
|-- api/
|   |-- api_client.go        # API client implementation
|   `-- api_client_test.go   # Local client tests with httptest
|-- config/
|   |-- config.json          # Base URL, endpoints, and request body
|   `-- load.json            # Load testing scenarios
|-- cmd/
|   `-- loadtest/
|       `-- main.go          # Load test CLI entrypoint
|-- fuzz/
|   |-- fuzz_get_test.go     # GET fuzz target wrapper
|   |-- fuzz_post_test.go    # POST fuzz target wrapper
|   |-- fuzz_write_methods_test.go # Shared fuzz helpers plus PUT, PATCH, and DELETE wrappers
|   `-- seeds_test.go        # Verb-specific endpoint and body seeds
|-- load/
|   |-- config.go            # Load test config loading and validation
|   |-- runner.go            # Concurrent load runner and metrics aggregation
|   `-- runner_test.go       # Local load runner tests with httptest
|-- logger/
|   |-- logger.go            # Request logging and findings helpers
|   `-- logger_test.go       # Findings artifact tests
|-- scripts/
|   |-- run-fuzz.ps1         # PowerShell runner that stores artifacts per execution
|   |-- run-load.ps1         # PowerShell load runner with separated artifacts
|   `-- run-api-tests.ps1    # Combined fuzz + load runner
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

The `config/load.json` file keeps load testing separate from fuzzing:

```json
{
  "baseURL": "https://example.com/api/v1",
  "scenarios": [
    {
      "name": "list-resources",
      "method": "GET",
      "endpoint": "/resources"
    },
    {
      "name": "create-resource",
      "method": "POST",
      "endpoint": "/resources",
      "body": {
        "id": 1,
        "title": "Load Test Activity",
        "dueDate": "2026-01-01T00:00:00Z",
        "completed": false
      }
    }
  ]
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

GitHub Actions writes a lightweight CI diagnostics summary for every run and uploads `reports/ci` as an artifact for 14 days. The summary includes the latest lines from module verification, formatting, vet, and test logs so failures are easier to triage from the workflow page.

Run every fuzz target on PowerShell and save artifacts per execution:

```powershell
.\scripts\run-fuzz.ps1 -FuzzTime 30s
```

The runner creates a timestamped directory under `artifacts/` for each execution:

```plaintext
artifacts/
`-- 2026-06-28_10-30-15/
    |-- FuzzGetEndpoint.log
    |-- FuzzPostEndpoint.log
    |-- FuzzPutEndpoint.log
    |-- FuzzPatchEndpoint.log
    |-- FuzzDeleteEndpoint.log
    |-- fuzz-findings.jsonl
    `-- summary.txt
```

By default the script keeps the latest 10 execution directories and removes older ones. Override that with `-KeepRuns`, or use `-KeepRuns 0` to disable cleanup:

```powershell
.\scripts\run-fuzz.ps1 -FuzzTime 5m -KeepRuns 20
```

Request bodies, response bodies, and error text are truncated to 8192 bytes before being stored in logs or findings. Override the limit with `-MaxLogBytes`; use `0` to disable truncation:

```powershell
.\scripts\run-fuzz.ps1 -FuzzTime 30s -MaxLogBytes 4096
```

For longer campaigns, keep `fuzz-findings.jsonl` and summaries but suppress per-request logging noise with `-QuietRequests`:

```powershell
.\scripts\run-fuzz.ps1 -FuzzTime 10m -QuietRequests
```

Run a small load test and save artifacts separately:

```powershell
.\scripts\run-load.ps1 -Users 5 -Requests 25
```

The runner creates a timestamped directory under `artifacts/load/`:

```plaintext
artifacts/
`-- load/
    `-- 2026-07-04_12-00-00/
        |-- load-console.log
        |-- load-requests.jsonl
        `-- summary.json
```

Use a duration-based run instead of a fixed request count:

```powershell
.\scripts\run-load.ps1 -Users 10 -Requests 0 -Duration 30s
```

Run only selected load scenarios:

```powershell
.\scripts\run-load.ps1 -Scenario list-activities,get-activity -Users 5 -Requests 20
```

### Reading API health from load results

Load test results are written to `summary.json` and `load-requests.jsonl`. A small run is useful as a smoke test: it confirms that the configured API is reachable, the selected calls respond, and the runner can collect latency and status data. It is not enough by itself to prove that the API will behave well under sustained traffic.

Start by checking these fields in `summary.json`:

| Metric | What to check |
| --- | --- |
| `successful` / `totalRequests` | Availability for the tested calls. |
| `failed` | Connection errors, timeouts, unexpected status codes, or failed requests. |
| `statusCodes` | HTTP response distribution. Watch for `5xx` and unexpected `4xx`. |
| `avgMs` | Normal response time across the run. |
| `p95Ms` and `p99Ms` | Slowest user-facing responses. These matter more than the average. |
| `requestsPerSec` | Throughput reached by the test. |

An initial health guideline for small API smoke/load runs:

```plaintext
Healthy:
- failed == 0
- success rate >= 99%
- no HTTP 5xx responses
- p95Ms < 500
- p99Ms < 1000
```

Example from a small run:

```plaintext
Requests: 10
Success: 10
Failed: 0
Status codes: 200 => 10
Avg: 95.60ms
P95: 296ms
P99: 296ms
```

That result is healthy for a smoke test, but confidence is still limited because it only uses a few requests. Increase users, duration, and scenario coverage gradually to evaluate sustained API health.

Run fuzzing and load testing together while keeping their logs separated:

```powershell
.\scripts\run-api-tests.ps1 -FuzzTime 30s -LoadUsers 5 -LoadRequests 25 -QuietFuzzRequests
```

The combined runner writes its own summary and stores fuzz/load artifacts below the same timestamped run:

```plaintext
artifacts/
`-- api-tests/
    `-- 2026-07-04_12-05-00/
        |-- summary.txt
        |-- fuzz/
        `-- load/
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
