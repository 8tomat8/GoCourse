package main

import (
	"os"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
)

func main() {
	url := os.Args[1]

	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	rv, err := ioutil.ReadAll(io.LimitReader(r.Body, 2 ^ 1024))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rv))
}
