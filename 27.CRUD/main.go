package main

import (
	"flag"
	"net"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

var dbInst *sql.DB

type handler struct{}

type Book struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Pages  int `json:"pages"`
}

type books []Book

func main() {
	host := flag.String("host", "127.0.0.1:8181", "Host:port to listen.")
	flag.Parse()

	db, err := sql.Open("sqlite3", "./27.CRUD/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	dbInst = db
	defer db.Close()

	listener, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.Serve(listener, &handler{}))
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	case http.MethodPut:
		put(w, r)
	case http.MethodDelete:
		delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(" Use POST, GET, PUT or DELETE"))
	}
}
