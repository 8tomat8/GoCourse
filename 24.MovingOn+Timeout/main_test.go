package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
)

func TestMakeHTTPCalls(t *testing.T) {
	workers := 500
	testValue := "fw4vwsv3wv3f3a8nfyow4gw4gw4gwgwgwefe fewf wef4w5huohf7937h f9o78vbw8bcieubeveiff"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testValue)
	}))
	if rv := makeHTTPCalls(&workers, server); rv != testValue {
		t.Error("Input != Output")
	}
}
