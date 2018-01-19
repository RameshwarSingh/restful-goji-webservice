package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/gorilla/mux"
)

type book struct {
	ISBN    string "json:isbn"
	Title   string "json:name"
	Authors string "json:author"
	Price   string "json:price"
}

var bookStore = []book{
	book{
		ISBN:    "01231231",
		Title:   "Programming with Go Lang",
		Authors: "Mark Summerfield",
		Price:   "$34.12",
	},
	book{
		ISBN:    "01231233",
		Title:   "Go Lang: RESTful with Goji",
		Authors: "Anthony Brian",
		Price:   "$14.32",
	},
}

func main() {
	mux := gojiMux()
	http.ListenAndServe("localhost:8080", mux)
}

func gojiMux() *goji.Mux {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/books"), allBooks)
	mux.HandleFunc(pat.Get("/books/:isbn"), bookByISBN)
	mux.Use(logging)
	return mux
}

func gorillaMux() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/books", allBooks)
	mux.HandleFunc("/books/{isbn}", bookByISBN)
	mux.Use(logging)
	http.Handle("/", mux)
	return mux
}

func allBooks(w http.ResponseWriter, r *http.Request) {
	jsonOut, _ := json.Marshal(bookStore)
	fmt.Fprintf(w, string(jsonOut))
}

// Search book by ISBN.
func bookByISBN(w http.ResponseWriter, r *http.Request) {
	isbn := pat.Param(r, "isbn")
	for _, b := range bookStore {
		if b.ISBN == isbn {
			jsonOut, _ := json.Marshal(b)
			fmt.Fprintf(w, string(jsonOut))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func logging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %v\n", r.URL)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
