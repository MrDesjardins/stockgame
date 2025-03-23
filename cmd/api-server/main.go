package main

import (
	"fmt"
	"net/http"
	"os"

	"stockgame/internal/database"
	"stockgame/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getStocks(c *gin.Context) {
	stock := service.GetRandomStockWithRandomDayRange(20)
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

	router.Run(fmt.Sprintf("localhost:%s", port))
}
