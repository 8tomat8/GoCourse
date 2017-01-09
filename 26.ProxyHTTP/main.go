package main

import (
	"flag"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
)

var listen = flag.String("l", "127.0.0.1:11112", "host:port")
var remote = flag.String("r", "google.com", "Host name")

type ProxyHandler struct{
	remote string
}

func main() {
	flag.Parse()

	h := ProxyHandler{*remote}

	log.Printf("Starting local HTTP server on http://%v ...", *listen)
	log.Fatal(http.ListenAndServe(*listen, &h))
}

func  (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := h.copyRequest(r)

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	res.Body.Close()

	for k, v := range res.Header {
		for _, elem := range v {
			w.Header().Add(k, elem)
		}
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

}

func (h *ProxyHandler) copyRequest(r *http.Request) *http.Request {
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}

	return &http.Request{
		Method: r.Method,
		URL: &url.URL{
			Scheme:   r.URL.Scheme,
			Host:     h.remote,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
		},
		Host:       *remote,
		Header:     r.Header,
		Body:       r.Body,
		ContentLength: r.ContentLength,
	}
}
