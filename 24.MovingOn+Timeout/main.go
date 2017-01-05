package main

import (
	"time"
	"net/http/httptest"
	"net/http"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"flag"
	randLib "math/rand"
)

func main() {
	var workers = flag.Int("workers", 5, "Number of workers.")
	flag.Parse()

	rand := randLib.New(randLib.NewSource(time.Now().UnixNano()))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
		fmt.Fprintf(w, "Усё чотко")
	}))
	defer server.Close()

	fmt.Println(makeHTTPCalls(workers, server))

}

func makeHTTPCalls(workers *int, server *httptest.Server) string {
	rv := make(chan string)
	for i := 0; i <= *workers; i++ {
		go makeHTTPCall(server, &rv)
	}

	return <-rv
}

func makeHTTPCall(server *httptest.Server, ch *chan string) error {
	client := http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(server.URL)

	if err != nil {
		glog.Error(err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return err
	}
	resp.Body.Close()

	*ch <- string(body)

	return nil
}
