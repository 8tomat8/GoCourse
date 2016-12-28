package handlers

import (
	"net/http"
	"github.com/8tomat8/GoCourse/21.LineCountingHTTPService/library"
	"github.com/8tomat8/GoCourse/21.LineCountingHTTPService/common"
	"encoding/json"
	"strings"
)

func Book(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	name := path[len(path) - 1]
	encoder := json.NewEncoder(w)
	l := library.Get()
	book, err := l.GetBook(name)
	if err != nil {
		w.WriteHeader(404)
		encoder.Encode(map[string]string{"error": "Hmmm... It seems like there is no book with name " + name})
		return
	}
	encoder.Encode(book)

}

func Books(w http.ResponseWriter, r *http.Request) {
	l := library.Get()
	encoder := json.NewEncoder(w)
	files, err := common.GetFiles(l.Path)
	if err != nil {
		encoder.Encode(err)
	}
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	encoder.Encode(map[string][]string{"books": result})

}
