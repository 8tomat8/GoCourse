package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
)

func main() {
	r, err := http.Get("http://httpbin.org/ip")
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	rv := make(map[string]string)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&rv)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("My IP adress is: ", rv["origin"])
}
