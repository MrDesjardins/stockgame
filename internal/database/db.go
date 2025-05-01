package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

var CONTEXT_TIMEOUT = 30 * time.Second

type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func configureDB() {
	maxConnection, err := getMaxConnections()
	if err != nil {
		println("Error getting max connections:", err)
		return
	}
	maxConnectionReduced := int(float64(maxConnection) * 0.75)
	// More conservative connection pool settings
	db.SetMaxOpenConns(maxConnectionReduced)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 1) // Add this line to clean up idle connections

	// Test the connection
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	fmt.Printf("Connected to the database with connection pool configuration using %d connections\n", maxConnectionReduced)
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
	configureDB()
}

func ConnectDBFullPath(dsn string) {
	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		println("Error connecting to the database")
		panic(err)
	}

	configureDB()
}
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// GetRawDB returns the raw sql.DB instance.
// This is useful for cases where you need to use the raw database connection like when creating tables
func GetRawDB() *sql.DB {
	if db == nil {
		println("Database connection is not initialized. Call ConnectDB() first.")
		panic("Database connection is nil")
	}
	return db
}

// Get DB connection
// This is the way to get fetch data form the database
func GetDB() DBInterface {
	if db == nil {
		println("Database connection is not initialized. Call ConnectDB() first.")
		panic("Database connection is nil")
	}
	return db
}

func getMaxConnections() (int, error) {
	if db != nil {
		var maxConnections int
		err := db.QueryRow("SHOW max_connections;").Scan(&maxConnections)
		if err != nil {
			return 0, err
		}

		// Reserve connections for PostgreSQL system processes and other applications
		// Typically use 75% of max_connections as a safe limit
		safeConnections := int(float64(maxConnections) * 0.75)

		return safeConnections, nil
	}
	return 0, fmt.Errorf("db is nil")
}

func GetDBStats() sql.DBStats {
	if db == nil {
		panic("Database connection is nil")
	}
	return db.Stats()
}

// Use this function to log connection usage periodically
func LogDBStats() {
	if db == nil {
		return
	}
	stats := db.Stats()
	fmt.Printf("DB Stats: Open=%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v\n",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration)
}
