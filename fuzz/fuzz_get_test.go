package fuzz

import (
	"net/http"
	"testing"
)

func FuzzGetEndpoint(f *testing.F) {
	fuzzEndpointWithoutBody(f, http.MethodGet, getEndpointSeeds)
}
