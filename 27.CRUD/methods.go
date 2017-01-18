package main

import (
	"database/sql"
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"strings"
	"io/ioutil"
)

func getId(r *http.Request) string {
	pathParams := strings.Split(r.URL.Path, "/")
	return pathParams[len(pathParams) - 1]
}

func processError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err)

	log.Println(err)
}

func get(w http.ResponseWriter, r *http.Request) (err error) {
	var book Book

	id := getId(r)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(string(http.StatusBadRequest))
	}

	err = dbInst.QueryRow("SELECT id, name, author, pages FROM books WHERE id=?", id).Scan(&book.Id, &book.Name, &book.Author, &book.Pages)
	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case err != nil:
		processError(err, w)
	default:
		data, err := json.Marshal(book)
		if err != nil {
			processError(err, w)
		}
		fmt.Fprint(w, string(data))
	}
	return
}

func post(w http.ResponseWriter, r *http.Request) (err error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processError(err, w)
	}

	var booksList books
	err = json.Unmarshal(data, &booksList)
	if err != nil {
		processError(err, w)
		return
	}

	s, err := dbInst.Prepare("INSERT INTO books(name, author, pages) VALUES(?,?,?)")
	if err != nil {
		processError(err, w)
		return
	}
	for _, book := range booksList {
		_, err = s.Exec(book.Name, book.Author, book.Pages)
		if err != nil {
			processError(err, w)
		}
	}
	return
}

func delete(w http.ResponseWriter, r *http.Request) (err error) {
	id := getId(r)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(string(http.StatusBadRequest))
	}

	s, err := dbInst.Prepare("DELETE FROM books WHERE id=?")
	result, err := s.Exec(id)
	if err != nil {
		processError(err, w)
		return
	}
	changed, err := result.RowsAffected()
	if err != nil {
		processError(err, w)
		return
	}
	if changed == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return
}

func put(w http.ResponseWriter, r *http.Request) (err error) {
	id := getId(r)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New(string(http.StatusBadRequest))
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		processError(err, w)
	}

	var book Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		processError(err, w)
		return
	}

	s, err := dbInst.Prepare("UPDATE books SET name=?, author=?, pages=? WHERE id=?")
	result, err := s.Exec(book.Name, book.Author, book.Pages, id)
	if err != nil {
		processError(err, w)
		return
	}
	changed, err := result.RowsAffected()
	if err != nil {
		processError(err, w)
		return
	}
	if changed == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return
}