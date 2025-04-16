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
const tableNameStocks = "stocks"
const tableNameStocksInfo = "stocks_info"
const maxFieldCVS = 12

func createTables(db *sql.DB) {
	startTime := time.Now()

	// Use fmt.Sprintf for better clarity and safety with constants
	dropStocksQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableNameStocks)
	_, err := db.Exec(dropStocksQuery)
	if err != nil {
		log.Fatal(err)
	}

	dropStocksInfoQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableNameStocksInfo)
	_, err = db.Exec(dropStocksInfoQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Use fmt.Sprintf for the CREATE TABLE statements as well
	createStocksQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        symbol VARCHAR NULL,
        date DATE NOT NULL,
        open FLOAT NOT NULL,
        high FLOAT NOT NULL,
        low FLOAT NOT NULL,
        "close" FLOAT NOT NULL,
        adj_close FLOAT NOT NULL,
        volume BIGINT NOT NULL
    );`, tableNameStocks)

	_, err = db.Exec(createStocksQuery)
	if err != nil {
		log.Printf("Cannot create %s table: %v", tableNameStocks, err)
		panic(err)
	}

	createStocksInfoQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        symbol VARCHAR NOT NULL,
        name VARCHAR NOT NULL,
        symbol_uuid VARCHAR NOT NULL
    );`, tableNameStocksInfo)

	_, err = db.Exec(createStocksInfoQuery)
	if err != nil {
		log.Printf("Cannot create %s table: %v", tableNameStocksInfo, err)
		panic(err)
	}

	fmt.Println("Data deletion completed.")
	fmt.Printf("Time taken drop table + create table: %v\n", time.Since(startTime))
}

func insertCompanyInfo(db *sql.DB) {
	relativePath := "./data/raw/symbols_valid_meta.csv"
	absolutePath := filepath.Join(util.GetProjectRoot(), relativePath)
	startTime := time.Now()

	truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s;", tableNameStocksInfo)
	_, err := db.Exec(truncateQuery)
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
	reader.FieldsPerRecord = maxFieldCVS
	reader.Comma = ','
	reader.LazyQuotes = true

	// Insert the data into the database
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Use fmt.Sprintf for the prepared statement
	insertQuery := fmt.Sprintf("INSERT INTO %s (symbol, name, symbol_uuid) VALUES ($1, $2, $3)", tableNameStocksInfo)
	stmt, err := tx.Prepare(insertQuery)
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
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
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

	// Transaction for bulk insert
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// Start worker goroutines
	for range maxWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				processStockFile(filePath, tx)
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
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	fmt.Println("Parallel data insertion completed.")
	fmt.Printf("Time stock taken: %v\n", time.Since(startTime))
}

// Process a single stock CSV file
func processStockFile(filePath string, tx *sql.Tx) {
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
		COPY %s (date, open, high, low, close, adj_close, volume, symbol)
		FROM '%s'
		WITH (FORMAT csv, HEADER TRUE, DELIMITER ',', QUOTE '"', ESCAPE '\', NULL '');
	`, tableNameStocks, cleanedFilePath)

	_, err = tx.Exec(query)

	if err != nil {
		log.Fatalf("Error copying CSV file %s: %v\n", filePath, err)
	}
}

func addIndex(db *sql.DB) {
	startTime := time.Now()
	// Create indexes for the stocks table
	_, err := db.Exec(`CREATE INDEX idx_stocks_symbol_date ON ` + tableNameStocks + ` (symbol, date DESC);`)
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
