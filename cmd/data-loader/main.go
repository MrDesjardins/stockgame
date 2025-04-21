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

	// UNLOGGED for performance
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

	fmt.Printf("Drop Table + Creating Table Completed: %v\n", time.Since(startTime))
}

func insertCompanyInfo(db *sql.DB) {
	relativePath := "./data/raw/symbols_valid_meta.csv"
	absolutePath := filepath.Join(util.GetProjectRoot(), relativePath)
	startTime := time.Now()

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
	fmt.Printf("Time to insert data into %s table: %v\n", tableNameStocksInfo, time.Since(startTime))
}
func insertStocksParallel(db *sql.DB) {
	dirPath := "./data/raw/stocks/"
	startTime := time.Now()

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}
	fmt.Println("Found", len(files), "files in directory")

	fileChan := make(chan string, len(files))
	var wg sync.WaitGroup

	// Shared slice to store cleaned file paths
	var cleanedFiles []string
	var mu sync.Mutex

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				cleanedFilePath, err := preprocessCSV(filepath.Join(util.GetProjectRoot(), filePath), getSymbol(filePath))
				if err != nil {
					fmt.Printf("Error preprocessing file %s: %v\n", filePath, err)
					continue
				}

				mu.Lock()
				cleanedFiles = append(cleanedFiles, cleanedFilePath)
				mu.Unlock()
			}
		}()
	}

	// Feed the file paths
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			fileChan <- filepath.Join(dirPath, file.Name())
		}
	}
	close(fileChan)
	wg.Wait()
	fmt.Printf("Insert - PreProcessing took: %v\n", time.Since(startTime))
	buckInsertTime := time.Now()

	// Performance
	_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s SET (autovacuum_enabled = false)", tableNameStocks))
	if err != nil {
		log.Fatal("Failed to disable autovacuum:", err)
	}
	_, err = db.Exec("SET session_replication_role = 'replica'") // Disables triggers
	if err != nil {
		log.Fatal("Failed to disable triggers:", err)
	}

	// Create a channel for file insertion tasks and error reporting
	type insertResult struct {
		filePath string
		err      error
	}
	resultChan := make(chan insertResult, len(cleanedFiles))

	// Use a semaphore pattern to limit concurrent transactions, too many can cause contention
	concurrentTxLimit := 4
	if concurrentTxLimit > maxWorkers {
		concurrentTxLimit = maxWorkers
	}

	sem := make(chan struct{}, concurrentTxLimit)

	// Launch workers to insert files in parallel
	for _, cleanedFile := range cleanedFiles {
		sem <- struct{}{} // Acquire semaphore
		go func(file string) {
			defer func() { <-sem }() // Release semaphore when done

			// Create transaction for this file
			tx, err := db.Begin()
			if err != nil {
				resultChan <- insertResult{file, fmt.Errorf("begin transaction failed: %v", err)}
				return
			}

			err = bulkInsertStockFile(file, tx)
			if err != nil {
				tx.Rollback()
				resultChan <- insertResult{file, err}
				return
			}

			if err := tx.Commit(); err != nil {
				resultChan <- insertResult{file, fmt.Errorf("commit failed: %v", err)}
				return
			}

			_ = os.Remove(file)
			resultChan <- insertResult{file, nil}
		}(cleanedFile)
	}

	// Collect results
	var insertErrors []string
	for i := 0; i < len(cleanedFiles); i++ {
		result := <-resultChan
		if result.err != nil {
			insertErrors = append(insertErrors, fmt.Sprintf("Error with file %s: %v", result.filePath, result.err))
		}
	}

	// Check for any errors
	if len(insertErrors) > 0 {
		for _, err := range insertErrors {
			fmt.Println(err)
		}
		log.Fatal("Failed to insert one or more files")
	}

	// Performance - restore settings
	_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s SET (autovacuum_enabled = true)", tableNameStocks))
	if err != nil {
		log.Fatal("Failed to enable autovacuum:", err)
	}
	_, err = db.Exec("SET session_replication_role = 'origin'") // Re-enables triggers
	if err != nil {
		log.Fatal("Failed to enable triggers:", err)
	}

	fmt.Printf("Insert - Parallel data insertion (bulk) completed in %v\n", time.Since(buckInsertTime))
	fmt.Printf("Insert time taken: %v\n", time.Since(startTime))
}

func bulkInsertStockFile(filePath string, tx *sql.Tx) error {
	query := fmt.Sprintf(`
		COPY %s (date, open, high, low, close, adj_close, volume, symbol)
		FROM '%s'
		WITH (FORMAT csv, HEADER TRUE, DELIMITER ',', QUOTE '"', ESCAPE '\', NULL '');
	`, tableNameStocks, filePath)

	_, err := tx.Exec(query)
	return err
}

func getSymbol(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), ".csv")
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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	db := database.GetDB()
	startTime := time.Now()
	createTables(db)
	insertStocksParallel(db)
	insertCompanyInfo(db)
	addIndex(db)
	fmt.Printf("Total time taken: %v\n", time.Since(startTime))
	defer db.Close()
}
