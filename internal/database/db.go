package database

import (
	"database/sql"
	"path/filepath"
	"stockgame/internal/util"

	_ "github.com/marcboeker/go-duckdb" // DuckDB driver
	_ "modernc.org/sqlite"
)

var db *sql.DB

// Connect to DB
func ConnectDB() {
	var err error

	// Get the absolute path to the database file
	dbPath := filepath.Join(util.GetProjectRoot(), "data", "db", "stockgame.duckdb")
	println("Database path: ", dbPath)
	db, err = sql.Open("duckdb", dbPath)
	if err != nil {
		println("Error connecting to the database")
		panic(err)
	}

	println("Connected to the database")
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Get DB connection
func GetDB() *sql.DB {
	if db == nil {
		println("Database connection is not initialized. Call ConnectDB() first.")
		panic("Database connection is nil")
	}
	// return the DB connection
	return db
}
