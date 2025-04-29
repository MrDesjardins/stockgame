package main

import (
	"fmt"
	"stockgame/internal/dataaccess"
	"stockgame/internal/database"
	"stockgame/internal/logic"
	"stockgame/internal/service"
	"stockgame/internal/util"
	"time"
)

func main() {
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
	generateService.GenerateFiles(10)
	fmt.Printf("Total time taken: %v\n", time.Since(startTime))
	defer db.Close()
}
