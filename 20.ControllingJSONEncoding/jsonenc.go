package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
)

type response struct {
	Origin string `json:"origin"`
}

func main() {
	r, err := http.Get("http://httpbin.org/ip")
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	rv := response{}
	json.Unmarshal(body, &rv)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("My IP adress is:", rv.Origin)
}
