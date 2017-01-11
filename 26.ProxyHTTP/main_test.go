package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"bytes"
	"net"
)

var testRequest = []byte("Test Request. 123123123")
var testResponse = []byte("Test Response. 321321321")

func TestRegexpHandler_ServeHTTP(t *testing.T) {

	host := "127.123.123.112:12345"

	listener, err := net.Listen("tcp", host)
	if err != nil {
		t.Fatal(err)
	}
	go http.Serve(listener, &Handler{})

	req, err := http.NewRequest("POST", "http://"+host+"/", bytes.NewReader(testRequest))

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	h := ProxyHandler{host}
	h.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	
	if bytes.Compare(w.Body.Bytes(), testResponse) != 0 {
		t.Errorf("Response is corrupted! Received: %v", w.Body.String())
	}
}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if bytes.Compare(data, testRequest) != 0 {
		w.Write([]byte("Request is corrupted! Received: "))
		w.Write(data)
		return
	}

	w.Write(testResponse)
}
