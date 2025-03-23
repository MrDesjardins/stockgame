package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"stockgame/internal/database"
	"stockgame/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getStocks(c *gin.Context) {
	stock := service.GetRandomStockWithRandomDayRange(20)
	c.IndentedJSON(http.StatusOK, stock)
}

func getStockInTimeRange(c *gin.Context) {

	stockSymbol := c.Query("symbol")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	if stockSymbol == "" || startDate == "" || endDate == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "symbol, startDate, and endDate are required query parameters"})
		return
	}
	// Make sure we are not querying for more than 30 days
	startDateGo, err := time.Parse("2006-01-20", startDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "startDate is not in the correct format"})
		return
	}
	endDateGo, err := time.Parse("2006-01-20", endDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "endDate is not in the correct format"})
		return
	}
	if endDateGo.Sub(startDateGo).Hours()/24 > 30 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cannot query for more than 30 days"})
		return
	}
	stock := service.GetStockPriceForTimeRange(stockSymbol, startDate, endDate)
	c.IndentedJSON(http.StatusOK, stock)
}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	port := os.Getenv("API_PORT")
	database.ConnectDB()
	router := gin.Default()
	// Cors
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})
	router.GET("/stocks", getStocks)
	router.GET("/stocksInTime", getStockInTimeRange)

	router.Run(fmt.Sprintf("localhost:%s", port))
}
