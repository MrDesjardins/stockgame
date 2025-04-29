package service

import (
	"encoding/json"
	"fmt"
	"os"
	"stockgame/internal/model"
	"sync"
)

const FOLDER_GENERATED_FILE = "generated_files"

// Create a struct to hold the combined data
type FileData struct {
	Info   model.StockInfo     `json:"info"`
	Prices []model.StockPublic `json:"prices"`
}

type GeneratorService interface {
	GenerateFiles(numberOfFiles int)
	WriteToFile(symbolUUID string, stocks []model.StockPublic, stockInfo model.StockInfo)
}
type GeneratorServiceImpl struct {
	StockService StockService
}

func (s *GeneratorServiceImpl) GenerateFiles(numberOfFiles int) {
	// Set to know if the stock already taken
	s.DeletePreviousFiles()
	stockSet := make(map[string]bool)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Launch goroutines based on fixed count
	for i := 0; i < numberOfFiles; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Get a random stock symbol
			stockPublics := s.StockService.GetRandomStockWithRandomDayRange(model.Number_initial_stock_shown)
			symbolUUID := stockPublics[0].SymbolUUID

			mu.Lock()
			// Check if we've already processed this stock
			if stockSet[symbolUUID] {
				mu.Unlock()
				return
			}
			// Mark this stock as processed
			stockSet[symbolUUID] = true
			mu.Unlock()

			// Get the stock info
			stockInfo, err := s.StockService.GetStockInfo(symbolUUID)
			if err != nil {
				fmt.Printf("Error getting stock info for %s: %v\n", symbolUUID, err)
				return
			}

			s.WriteToFile(symbolUUID, stockPublics, stockInfo)
		}()
	}

	wg.Wait()
}

func (s *GeneratorServiceImpl) DeletePreviousFiles() {
	// Check if the directory exists
	if _, err := os.Stat("./" + FOLDER_GENERATED_FILE); os.IsNotExist(err) {
		// Directory does not exist, no need to delete
		return
	}
	// Remove the directory and its contents
	err := os.RemoveAll("./" + FOLDER_GENERATED_FILE)
	if err != nil {
		fmt.Printf("Error deleting directory: %v\n", err)
		return
	}
	// Create the directory again
	err = os.MkdirAll("./"+FOLDER_GENERATED_FILE, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
}
func (s *GeneratorServiceImpl) WriteToFile(symbolUUID string, stocks []model.StockPublic, stockInfo model.StockInfo) {
	pathFolders := fmt.Sprintf("./%s/%s.json", FOLDER_GENERATED_FILE, symbolUUID)

	// Combine the data
	data := FileData{
		Info:   stockInfo,
		Prices: stocks,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling JSON for %s: %v\n", symbolUUID, err)
		return
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll("./"+FOLDER_GENERATED_FILE, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Write the JSON data to the file
	if err := os.WriteFile(pathFolders, jsonData, 0644); err != nil {
		fmt.Printf("Error writing to file %s: %v\n", symbolUUID, err)
		return
	}
}
