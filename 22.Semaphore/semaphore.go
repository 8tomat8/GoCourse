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
	"sync"
)

type empty struct{}
type semaphore chan empty

func main() {
	var limit = flag.Uint("limit", 2, "Limit of concurrent workers.")
	var data = flag.String("data", "[]", "JSON list of urls to download.")
	flag.Parse()
	var links = []string{}

	err := json.Unmarshal([]byte(*data), &links)
	if err != nil {
		panic(err)
	}

	sem := make(chan empty, int(*limit))
	wg := &sync.WaitGroup{}
	for _, link := range links {
		sem <- empty{}
		wg.Add(1)
		go download(link, &sem, wg)
	}
	wg.Wait()
}

func download(link string, sem *chan empty, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-*sem }()
	tokens := strings.Split(link, "/")
	fileName := tokens[len(tokens) - 1]

	if _, err := os.Stat("./" + fileName); err == nil {
		glog.Error(fmt.Sprintf("File %v already exist. Skipping it.", fileName))
		return
	}

	glog.Info("Downloading ", link, " to ", fileName)
	output, err := os.Create(fileName)
	if err != nil {
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

	_, err = io.Copy(output, response.Body)
	if err != nil {
		glog.Error(fmt.Sprintf("Error while downloading body %v - %v. Skipping it.", link, err))
		return
	}
}
