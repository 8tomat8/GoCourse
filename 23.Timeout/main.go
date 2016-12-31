package main

import (
	randLib "math/rand"
	"time"
	"net/http/httptest"
	"net/http"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"flag"
)

func main() {
	flag.Parse()
	rand := randLib.New(randLib.NewSource(time.Now().UnixNano()))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * time.Duration(rand.Intn(3)))
		fmt.Fprintf(w, "Усё чотко")
	}))
	defer server.Close()

	body, err := makeHTTPCall(server)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}

func makeHTTPCall(server *httptest.Server) (string, error){
	client := http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(server.URL)
	defer resp.Body.Close()

	if err != nil {
		glog.Error(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	return string(body), nil
}