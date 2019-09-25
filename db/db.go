package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL
)

// DB is just a *sql.DB, made to create receivers
type DB struct {
	DB *sql.DB
}

// InitDB creates database connection
func InitDB(dataSourceName string) *DB {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	return &DB{db}
}

// NewDatabase creates database with chosen name
func (d DB) NewDatabase(dbName string) {
	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, dbName, dbName)
	_, err := d.DB.Exec(statement)
	if err != nil {
		log.Panic(err)
	}
}

// NewBooksTable creates Table books
func (d DB) NewBooksTable() {
	stmt, err := d.DB.Prepare(`CREATE TABLE IF NOT EXISTS books(
		id serial PRIMARY KEY,
		title VARCHAR (50) NOT NULL,
		author VARCHAR (50) NOT NULL,
		year VARCHAR (50) NOT NULL
	 );`)
	if err != nil {
		log.Panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Panic(err)
	}
}
