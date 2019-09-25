package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"books-list/db"

	"github.com/gorilla/mux"
)

// Book properties
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book

// BooksAPI database
var bdb *db.DB

func main() {
	CreateConfig()
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	addr := os.Getenv("ADDRESS")
	dbport := os.Getenv("DBPORT")

	router := mux.NewRouter()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=booksapi sslmode=disable",
		addr, dbport, user, pass)
	bdb = db.InitDB(psqlInfo)
	defer bdb.DB.Close()
	bdb.NewDatabase("booksapi")
	bdb.NewBooksTable()

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
	var book Book
	books = make([]Book, 0)

	rows, err := bdb.DB.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}

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
