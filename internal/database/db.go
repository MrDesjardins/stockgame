package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Connect to DB
func ConnectDB() {
	var err error

	dsn := "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable connect_timeout=5"
	db, err = sql.Open("postgres", dsn)
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
