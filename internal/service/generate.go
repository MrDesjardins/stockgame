package service

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"stockgame/internal/database"
	"stockgame/internal/model"
	"strconv"
	"sync"
	"time"
)

const FOLDER_GENERATED_FILE = "generated_files"

var maxWorkers = runtime.NumCPU() * 2 // Double the CPU cores
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
	s.DeletePreviousFiles()
	stockSet := make(map[string]bool)
	var mu sync.Mutex

	// Further limit concurrent workers to avoid DB connection issues
	maxWorkers := maxWorkers // int(float64(maxConnection) * 0.75)

	fmt.Printf("Using %d workers based on available database connections\n", maxWorkers)

	// Add periodic stats logging
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			database.LogDBStats()
		}
	}()

	jobs := make(chan int, numberOfFiles)
	var wg sync.WaitGroup

	// Start worker pool with limited goroutines
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for i := range jobs {
				fmt.Printf("Worker %d processing job: %d\n", workerId, i)

				// Process job with proper error handling
				func() {
					// Make sure to cancel at the end of this inner function

					stockPublics := s.StockService.GetRandomStockWithRandomDayRange(model.Number_initial_stock_shown)
					if len(stockPublics) == 0 {
						fmt.Printf("Worker %d: Got empty stock list for job %d\n", workerId, i)
						return
					}

					symbolUUID := stockPublics[0].SymbolUUID

					mu.Lock()
					if stockSet[symbolUUID] {
						mu.Unlock()
						return
					}
					stockSet[symbolUUID] = true
					mu.Unlock()

					stockInfo, err := s.StockService.GetStockInfo(symbolUUID)
					if err != nil {
						fmt.Printf("Error getting stock info for %s: %v\n", symbolUUID, err)
						return
					}

					s.WriteToFile(strconv.Itoa(i), stockPublics, stockInfo)
				}()
			}
		}(w)
	}

	// Submit jobs
	for i := 0; i < numberOfFiles; i++ {
		jobs <- i
	}
	close(jobs)

	// Add timeout for the whole process
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All jobs completed successfully")
	case <-time.After(10 * time.Minute): // Overall timeout
		fmt.Println("Generation process timed out after 10 minutes!")
	}
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
func (s *GeneratorServiceImpl) WriteToFile(fileName string, stocks []model.StockPublic, stockInfo model.StockInfo) {
	pathFolders := fmt.Sprintf("./%s/%s.json", FOLDER_GENERATED_FILE, fileName)

	// Combine the data
	data := FileData{
		Info:   stockInfo,
		Prices: stocks,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling JSON for %s: %v\n", fileName, err)
		return
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll("./"+FOLDER_GENERATED_FILE, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Write the JSON data to the file
	if err := os.WriteFile(pathFolders, jsonData, 0644); err != nil {
		fmt.Printf("Error writing to file %s: %v\n", fileName, err)
		return
	}
}
