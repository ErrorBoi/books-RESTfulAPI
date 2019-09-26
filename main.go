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
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)

	row := bdb.DB.QueryRow(`SELECT * FROM books
	WHERE id=$1`, params["id"])

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := bdb.DB.QueryRow(`INSERT INTO books (title, author, year)
	VALUES($1, $2, $3)
	RETURNING id;`, book.Title, book.Author, book.Year).Scan(&bookID)

	logFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := bdb.DB.Exec(`UPDATE books
	SET title=$1, author=$2, year=$3
	WHERE id=$4
	RETURNING id;`, book.ID, book.Author, book.Year, book.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := bdb.DB.Exec(`DELETE FROM books
	WHERE id=$1`, params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
