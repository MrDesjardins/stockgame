package database

import (
	"database/sql"

	_ "github.com/marcboeker/go-duckdb" // DuckDB driver
	_ "modernc.org/sqlite"
)

var db *sql.DB

// Connect to DB
func ConnectDB() {
	var err error

	db, err = sql.Open("duckdb", "./data/db/stockgame.duckdb")
	if err != nil {
		println("Error connecting to the database")
		panic(err)
	}

	println("Connected to the database")
}

// Get DB connection
func GetDB() *sql.DB {
	// return the DB connection
	return db
}
