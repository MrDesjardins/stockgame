package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stockgame/internal/database"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) {

	// Create the stocks table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS stocks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		date TEXT NOT NULL,
		open REAL NOT NULL,
		high REAL NOT NULL,
		low REAL NOT NULL,
		close REAL NOT NULL,
		adj_close REAL NOT NULL,
		volume INTEGER NOT NULL
	);`)
	if err != nil {
		println("Cannot create table")
		panic(err)
	}

}

func insertStocks(db *sql.DB) {
	dirPath := "./data/raw/stocks/"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		println("Error reading directory")
		log.Fatal(err)
	}

	totalFiles := len(files)
	fmt.Printf("Total files: %d\n", totalFiles)

	startTime := time.Now()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Delete all stocks from the table
	_, err = tx.Exec("DELETE FROM stocks;")
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the statement once
	stmt, err := tx.Prepare(`INSERT INTO stocks (symbol, date, open, high, low, close, adj_close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Ensure statement is closed

	for i, file := range files {
		fileName := file.Name()
		stockName := fileName[:len(fileName)-4] // Remove the .csv extension

		fmt.Printf("%d/%d - %s\n", i+1, totalFiles, stockName)

		filePath := filepath.Join(dirPath, fileName)

		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		// Skip the first line (header)
		lines := strings.Split(string(data), "\n")[1:]

		// Insert each line using the prepared statement
		for _, line := range lines {
			// Skip empty lines
			if strings.TrimSpace(line) == "" {
				continue
			}

			// Split the line into columns
			columns := strings.Split(line, ",")
			if len(columns) < 7 {
				log.Printf("Skipping invalid row in %s: %v\n", filePath, columns)
				continue
			}

			_, err := stmt.Exec(stockName, columns[0], columns[1], columns[2], columns[3], columns[4], columns[5], columns[6])
			if err != nil {
				log.Printf("Error inserting row in %s: %v\n", filePath, err)
				continue
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	endTime := time.Now()
	fmt.Printf("Time taken: %v\n", endTime.Sub(startTime))

	fmt.Println("Data insertion completed successfully.")
}

func main() {
	// Create the SQL Lite database if it doesn't exist
	// Create a connection to the SQL Lite database
	database.ConnectDB()
	db := database.GetDB()

	createTables(db)
	insertStocks(db)

	defer db.Close()
}
