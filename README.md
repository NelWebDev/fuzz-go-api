# Fuzzing API Testing

This project is a fuzz testing framework for API endpoints. It uses Golang to perform fuzz tests on GET and POST requests, with configurations stored in a JSON file for flexibility and ease of reuse.

## Features

- **GET and POST Testing**: Test API endpoints for GET and POST requests using fuzzing techniques.
- **Configuration**: Easily manage base URL, endpoints, and request bodies via a JSON configuration file.
- **Error Handling**: Logs errors effectively to identify issues.
- **Fuzzing Reports**: Generate test reports automatically.

## Project Structure

```plaintext
fuzzing-api/
|-- config/
|   `-- config.json           # Configuration file (base URL, endpoints, and request bodies)
|-- fuzz/
|   |-- fuzz_get_test.go      # Fuzz tests for GET requests
|   `-- fuzz_post_test.go     # Fuzz tests for POST requests
|-- utils/
|   `-- utils.go              # Utility functions (configuration loading, helpers)
|-- api/
    `-- client.go             # API client implementation
```

## Configuration

The `config.json` file should be placed in the `config/` directory. Example:

```json
{
  "base_url": "https://example.com/api/v1/",
  "endpoints": {
    "get": "/resources",
    "post": "/resources"
  },
  "bodies": {
    "post": {
      "id": 1,
      "name": "example",
      "completed": false
    }
  }
}
```

## Usage

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/fuzzing-api.git
   cd fuzzing-api
   ```

2. **Run Tests:**
   ```bash
   go test -fuzz=FuzzGetEndpoint -fuzztime=30s ./fuzz
   go test -fuzz=FuzzPostEndpoint -fuzztime=30s ./fuzz
   ```

3. **Modify Configuration:**
   Edit the `config/config.json` file to update base URLs, endpoints, or request bodies for your API.

## Requirements

- **Go version**: 1.20+

