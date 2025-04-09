package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"stockgame/internal/database"
	"stockgame/internal/util"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var maxWorkers = runtime.NumCPU() * 2 // Double the CPU cores

func createTables(db *sql.DB) {
	startTime := time.Now()
	// Delete existing records before inserting
	_, err := db.Exec("DROP TABLE stocks;")
	if err != nil {
		log.Fatal(err)
	}

	// Delete existing records before inserting
	_, err = db.Exec("DROP TABLE stocks_info;")
	if err != nil {
		log.Fatal(err)
	}

	// Create the stocks table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stocks (
		id SERIAL PRIMARY KEY,
    symbol VARCHAR NULL,
    date DATE NOT NULL,
    open FLOAT NOT NULL,
    high FLOAT NOT NULL,
    low FLOAT NOT NULL,
    "close" FLOAT NOT NULL,
    adj_close FLOAT NOT NULL,
    volume BIGINT NOT NULL
);`)
	if err != nil {
		println("Cannot create stocks table")
		panic(err)
	}

	// Create the stocks table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stocks_info (
		id SERIAL PRIMARY KEY,
    symbol VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
		symbol_uuid VARCHAR NOT NULL
);`)
	if err != nil {
		println("Cannot create stocks_info table")
		panic(err)
	}
	fmt.Println("Data deletion completed.")
	fmt.Printf("Time taken drop table + create table: %v\n", time.Since(startTime))
}

func insertCompanyInfo(db *sql.DB) {
	relativePath := "./data/raw/symbols_valid_meta.csv"
	absolutePath := filepath.Join(util.GetProjectRoot(), relativePath)
	startTime := time.Now()

	// Delete existing records before inserting
	_, err := db.Exec("TRUNCATE TABLE stocks_info;")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted existing records")

	// Open the file and read the data line by line
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 12
	reader.Comma = ','
	reader.LazyQuotes = true

	// Insert the data into the database
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO stocks_info (symbol, name, symbol_uuid) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Skip the header
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Generate a UUID for the symbol
		uuid := uuid.New()

		_, err = stmt.Exec(row[1], row[2], uuid)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
	}
	fmt.Println("Data insertion completed.")
	fmt.Printf("Time taken stock_info: %v\n", time.Since(startTime))
}
func insertStocksParallel(db *sql.DB) {
	dirPath := "./data/raw/stocks/"
	startTime := time.Now()

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}
	fmt.Println("Found", len(files), "files in directory")

	fileChan := make(chan string, len(files)) // Channel for file paths
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				processStockFile(db, filePath)
			}
		}()
	}

	// Send file paths to the channel
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			filePath := filepath.Join(dirPath, file.Name())
			fileChan <- filePath
		}
	}
	close(fileChan) // Close channel after sending all file paths

	// Wait for workers to finish
	wg.Wait()

	fmt.Println("Parallel data insertion completed.")
	fmt.Printf("Time stock taken: %v\n", time.Since(startTime))
}

// Process a single stock CSV file
func processStockFile(db *sql.DB, filePath string) {
	symbol := strings.TrimSuffix(filepath.Base(filePath), ".csv")
	absolutePath := filepath.Join(util.GetProjectRoot(), filePath)

	// Preprocess the CSV file
	cleanedFilePath, err := preprocessCSV(absolutePath, symbol)
	if err != nil {
		fmt.Printf("Error preprocessing CSV file %s: %v\n", filePath, err)
		return
	}
	defer os.Remove(cleanedFilePath) // Remove temp file after processing

	// Use COPY command for fast bulk loading
	query := fmt.Sprintf(`
		COPY stocks (date, open, high, low, close, adj_close, volume, symbol)
		FROM '%s'
		WITH (FORMAT csv, HEADER TRUE, DELIMITER ',', QUOTE '"', ESCAPE '\', NULL '');
	`, cleanedFilePath)

	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("Error copying CSV file %s: %v\n", filePath, err)
	}
}

func addIndex(db *sql.DB) {
	startTime := time.Now()
	// Create indexes for the stocks table
	_, err := db.Exec(`CREATE INDEX idx_stocks_symbol_date ON stocks (symbol, date DESC);`)
	if err != nil {
		println("Cannot create index for stocks table")
		panic(err)
	}
	fmt.Println("Adding index completed")
	fmt.Printf("Time adding index: %v\n", time.Since(startTime))
}

func preprocessCSV(filePath string, symbol string) (string, error) {
	tempFilePath := fmt.Sprintf("%s_cleaned.csv", filePath)

	inputFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Read the header and write it to the new file
	header, err := reader.Read()
	if err != nil {
		return "", err
	}
	writer.Write(header)

	// Read and process each row
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Ensure the row has all expected columns
		if len(row) != 7 {
			continue // Skip rows with incorrect number of columns
		}

		// Check if any required field is empty (such as "open" column)
		if row[1] == "" || row[6] == "" {
			// Skip rows where essential columns are missing (e.g., "open" or "volume")
			continue
		}

		// Clean volume column (index 6) if necessary (convert to integer format)
		volume := row[6]
		if strings.Contains(volume, ".") {
			volume = strings.Split(volume, ".")[0] // Trim decimal part
		}
		row[6] = volume
		row = append(row, symbol) // Add an empty string to the end of the row

		// Write the cleaned row to the output file
		writer.Write(row)
	}

	return tempFilePath, nil
}

func main() {
	// Create the SQL Lite database if it doesn't exist
	// Create a connection to the SQL Lite database
	println("Max workers: ", maxWorkers)
	database.ConnectDB()
	db := database.GetDB()
	startTime := time.Now()
	createTables(db)
	insertStocksParallel(db)
	insertCompanyInfo(db)
	addIndex(db)
	fmt.Printf("Total time taken: %v\n", time.Since(startTime))
	defer db.Close()
}
