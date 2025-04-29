package main

import (
	"fmt"
	"os"
	"stockgame/internal/dataaccess"
	"stockgame/internal/database"
	"stockgame/internal/logic"
	"stockgame/internal/service"
	"time"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
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
	generateService.GenerateFiles(10)
	fmt.Printf("Total time taken: %v\n", time.Since(startTime))
	defer db.Close()
}
