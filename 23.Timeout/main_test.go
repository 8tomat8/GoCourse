package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"time"
	"fmt"
)

func TestMakeHTTPCall(t *testing.T) {
	testValue := "fw4vwsv3wv3f3a8nfyow4gw4gw4gwgwgwefe fewf wef4w5huohf7937h f9o78vbw8bcieubeveiff"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second / 2)
		fmt.Fprintf(w, testValue)
	}))

	if rv, err := makeHTTPCall(server); err != nil {
		t.Error(err)
	} else if rv != testValue {
		t.Error("Input != Output")
	}
}
