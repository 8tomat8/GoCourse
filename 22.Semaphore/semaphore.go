package main

import (
	"encoding/json"
	"flag"
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"github.com/golang/glog"
)

type empty struct{}
type semaphore chan empty

func main() {
	var limit = flag.Uint("limit", 2, "Limit of concurrent workers.")
	var data = flag.String("data", "[]", "JSON list of urls to download.")
	flag.Parse()

	var links = []string{}

	json.Unmarshal([]byte(data), &links)

	sem := make(chan empty)
	for _, link := range links {
		sem <- empty{}
		go download(link, &sem)
	}

}

func download(link string, sem <-chan empty) {
	tokens := strings.Split(link, "/")
	fileName := tokens[len(tokens)-1]
	glog.Info("Downloading", link, "to", fileName)

	output, err := os.Create(fileName)
	if err != nil {
		if os.IsExist(err) {
			glog.Error(fmt.Sprintf("File %v already exist. Skipping it.", fileName))
			return
		}
		glog.Error("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(link)
	if err != nil {
		glog.Error(fmt.Sprintf("Error while downloading %v - %v. Skipping it.", link, err))
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		glog.Error(fmt.Sprintf("Error while downloading body %v - %v. Skipping it.", link, err))
		return
	}

	<-sem
}