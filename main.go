package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Book properties
type Book struct {
	ID     string `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {
	CreateConfig()
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	books = append(books,
		Book{ID: "1", Title: "Golang pointers", Author: "Mr. Golang", Year: "2010"},
		Book{ID: "2", Title: "Goroutines", Author: "Mr. Goroutine", Year: "2011"},
		Book{ID: "3", Title: "Golang routers", Author: "Mr. Router", Year: "2012"},
		Book{ID: "4", Title: "Golang concurrency", Author: "Mr. Currency", Year: "2013"},
		Book{ID: "5", Title: "Golang good parts", Author: "Mr. Good", Year: "2014"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Printf("Listening on port :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe(): ", err)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	books = append(books, book)

	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, book := range books {
		if book.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
