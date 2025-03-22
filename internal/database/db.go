package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Connect to DB
func ConnectDB() {
	var err error

	db, err = sql.Open("sqlite", "./data/db/stockgame.db")
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
