package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"stockgame/internal/database"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) {

	// Create the stocks table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS stocks (
    symbol VARCHAR NULL,
    date VARCHAR NOT NULL,
    open FLOAT NOT NULL,
    high FLOAT NOT NULL,
    low FLOAT NOT NULL,
    "close" FLOAT NOT NULL,
    adj_close FLOAT NOT NULL,
    volume INTEGER NOT NULL
);`)
	if err != nil {
		println("Cannot create table")
		panic(err)
	}

}

func insertStocks(db *sql.DB) {
	dirPath := "./data/raw/stocks/"

	startTime := time.Now()

	// Delete existing records before inserting
	_, err := db.Exec("DELETE FROM stocks;")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted existing records")

	// Read all CSV files
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			symbol := strings.TrimSuffix(file.Name(), ".csv")
			filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())

			// Import with filename as symbol
			query := fmt.Sprintf(`
							COPY stocks (date, open, high, low, close, adj_close, volume)
							FROM '%s'
							WITH (HEADER TRUE, DELIMITER ',', QUOTE '"', ESCAPE '\', NULL '');
					`, filePath)

			_, err = db.Exec(query)
			if err != nil {
				// Preprocess the CSV file to remove rows with missing values
				cleanedFilePath, err := preprocessCSV(filePath)
				if err != nil {
					fmt.Printf("Error preprocessing CSV file %s: %v\n", file.Name(), err)
					continue
				}
				query := fmt.Sprintf(`
				COPY stocks (date, open, high, low, close, adj_close, volume)
				FROM '%s'
				WITH (HEADER TRUE, DELIMITER ',', QUOTE '"', ESCAPE '\', NULL '');
				`, cleanedFilePath)

				_, err = db.Exec(query)
				if err != nil {
					fmt.Printf("Error copying CSV file %s: %v\n", file.Name(), err)
					continue
				}
				err = os.Remove(cleanedFilePath)
				if err != nil {
					fmt.Printf("Error removing temp file %s: %v\n", cleanedFilePath, err)
				}
			}

			// Update symbol column
			_, err = db.Exec("UPDATE stocks SET symbol = ? WHERE symbol IS NULL;", symbol)
			if err != nil {
				fmt.Printf("Error updating symbol for %s: %v\n", file.Name(), err)
			}

			fmt.Printf("Inserted data from %s\n", file.Name())
		}
	}

	fmt.Println("Data insertion completed.")
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
}

// Preprocess the CSV file to remove rows with missing values
func preprocessCSV(filePath string) (string, error) {
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

	// Read and process each row
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Skip rows with missing values
		if len(row) == 7 && allColumnsHaveValues(row) {
			writer.Write(row)
		}
	}

	return tempFilePath, nil
}

// Helper function to check if all columns have values
func allColumnsHaveValues(row []string) bool {
	for _, col := range row {
		if strings.TrimSpace(col) == "" {
			return false
		}
	}
	return true
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
