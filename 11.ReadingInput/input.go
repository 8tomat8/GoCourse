package main

import (
	"bufio"
	"io"
	"log"
	"strings"
	"net/http"
	"fmt"
	"flag"
)

// ReadAll reads all the lines of text from r and returns
// all the data read as a string
func ReadAll(r io.Reader) string {
	sc := bufio.NewScanner(r)
	var lines []string

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	return strings.Join(lines, "\n")
}

func main() {
	var url = flag.String("url", "https://httpbin.org/get?a=1&b=2", "Link to read from.")
	flag.Parse()

	resp, err := http.Get(*url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	rv := ReadAll(resp.Body)
	fmt.Println(rv)
}
