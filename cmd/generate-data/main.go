package main

import (
	"fmt"
	"os"
	"stockgame/internal/dataaccess"
	"stockgame/internal/database"
	"stockgame/internal/logic"
	"stockgame/internal/service"
	"stockgame/internal/util"
	"strconv"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide an argument that is the number of files to generate")
		fmt.Println("Usage: go run main.go <number_of_files>")
		return
	}
	filesToGenerate, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error converting argument to integer:", err)
		return
	}
	_, dbHost, dbPort, dbUser, dbPassword, dbName, _ := util.GetDBEnv()

	database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	db := database.GetDB()
	startTime := time.Now()
	stockDataAccess := &dataaccess.StockDataAccessImpl{
		DB: db,
	}
	mockStockLogic := &logic.StockLogicImpl{}
	stockService := &service.StockServiceImpl{
		StockDataAccess: stockDataAccess,
		StockLogic:      mockStockLogic,
	}
	generateService := &service.GeneratorServiceImpl{
		StockService: stockService,
	}

	generateService.GenerateFiles(filesToGenerate)
	fmt.Printf("Total time taken: %v for %v files\n", time.Since(startTime), filesToGenerate)
}
