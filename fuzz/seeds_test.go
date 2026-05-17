package fuzz

import (
	"encoding/json"
	"fuzzing-api/utils"
)

func getEndpointSeeds(config *utils.Config) []string {
	return []string{
		config.Endpoints.Get,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities?completed=true&completed=false",
		"/Activities?page=-1&pageSize=999999",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	}
}

func postEndpointSeeds(config *utils.Config) []string {
	return []string{
		config.Endpoints.Post,
		"/Activities",
		"/Activities/",
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/abc",
		"/Activities?validate=true&validate=false",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	}
}

func putEndpointSeeds(config *utils.Config) []string {
	return []string{
		config.Endpoints.Put,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities/1?overwrite=true&overwrite=false",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	}
}

func patchEndpointSeeds(config *utils.Config) []string {
	return []string{
		config.Endpoints.Patch,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities/1?partial=true&partial=false",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	}
}

func deleteEndpointSeeds(config *utils.Config) []string {
	return []string{
		config.Endpoints.Delete,
		"/Activities/0",
		"/Activities/-1",
		"/Activities/2147483647",
		"/Activities/000001",
		"/Activities/abc",
		"/Activities/1.5",
		"/Activities/1?force=true&force=false",
		"/Activities/%2e%2e/%2e%2e",
		"/Activities/%20",
	}
}

func postBodySeeds(config *utils.Config) []string {
	return []string{
		configuredRequestBody(config),
		`{}`,
		`{"id":0,"title":"","dueDate":"0001-01-01T00:00:00Z","completed":false}`,
		`{"id":-1,"title":"negative id","dueDate":"2026-01-01T00:00:00Z","completed":false}`,
		`{"id":2147483647,"title":"max int id","dueDate":"9999-12-31T23:59:59Z","completed":true}`,
		`{"id":"1","title":123,"dueDate":false,"completed":"yes"}`,
		`{"unknown":"field","nested":{"value":[1,true,null]}}`,
		`not-json`,
	}
}

func putBodySeeds(config *utils.Config) []string {
	return []string{
		configuredRequestBody(config),
		`{"id":1,"title":"","dueDate":"2026-01-01T00:00:00Z","completed":false}`,
		`{"id":1,"title":"replacement","dueDate":"0001-01-01T00:00:00Z","completed":true}`,
		`{"id":0,"title":"id mismatch","dueDate":"2026-01-01T00:00:00Z","completed":false}`,
		`{"id":2147483647,"title":"large id","dueDate":"9999-12-31T23:59:59Z","completed":true}`,
		`{"id":"1","title":null,"dueDate":123,"completed":null}`,
		`[]`,
		`not-json`,
	}
}

func patchBodySeeds(config *utils.Config) []string {
	return []string{
		configuredRequestBody(config),
		`{}`,
		`{"title":""}`,
		`{"title":null}`,
		`{"completed":true}`,
		`{"completed":null}`,
		`{"dueDate":"not-a-date"}`,
		`{"unknown":"field"}`,
		`[]`,
		`not-json`,
	}
}

func addEndpointBodySeeds(add func(...any), endpoints []string, bodies []string) {
	if len(endpoints) == 0 || len(bodies) == 0 {
		return
	}

	for _, endpoint := range endpoints {
		add(endpoint, bodies[0])
	}
	for _, body := range bodies[1:] {
		add(endpoints[0], body)
	}
}

func configuredRequestBody(config *utils.Config) string {
	bodyJSON, err := json.Marshal(config.RequestBody)
	if err != nil {
		return `{}`
	}
	return string(bodyJSON)
}
