package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Connect to DB
func ConnectDB(host, port, user, password, dbname string) {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5", host, port, user, password, dbname)
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
