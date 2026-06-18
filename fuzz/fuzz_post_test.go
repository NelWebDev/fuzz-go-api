package fuzz

import (
	"net/http"
	"testing"
)

func FuzzPostEndpoint(f *testing.F) {
	fuzzEndpointWithBody(f, http.MethodPost, postEndpointSeeds, postBodySeeds)
}
