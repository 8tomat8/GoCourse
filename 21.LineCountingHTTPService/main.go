package main

import (
	"net/http"
	"log"
	"flag"
	"fmt"
	"regexp"
	"github.com/8tomat8/GoCourse/LineCountingHTTPService/handlers"
	"github.com/8tomat8/GoCourse/LineCountingHTTPService/library"
)

type handlerRoute struct {
	pattern *regexp.Regexp
	handler http.Handler
}
type route struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	routes []*handlerRoute
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &handlerRoute{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func registerEndpoints(routes []route) (*RegexpHandler, error) {
	log.Print("Registering endpoints...")

	handler := RegexpHandler{}
	for _, r := range routes {
		pattern, err := regexp.Compile(r.pattern)
		if err != nil {
			return nil, err
		}
		handler.HandleFunc(pattern, http.HandlerFunc(r.handler))
	}
	log.Print("DONE!")
	return &handler, nil
}

func main() {
	var host = flag.String("host", "127.0.0.1", "Host;")
	var port = flag.String("port", "8181", "Port;")
	var dir = flag.String("dir", "LineCountingHTTPService/books/", "Directory with book files.")
	flag.Parse()

	log.Print("Adding new library...")
	library.New(*dir)
	log.Print("DONE!")

	routes := []route{
		{"/books/.{1,}\\.txt$", handlers.Book},
		{"/books/$", handlers.Books},
	}
	handler, err := registerEndpoints(routes)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting HTTP server on http://%v:%v ...", *host, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", *host, *port), handler))
}
