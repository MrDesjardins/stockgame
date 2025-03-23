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
	stock := service.GetRandomStockWithRandomDayRange(1)
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
	router.GET("/stocks", getStocks)

	router.Run(fmt.Sprintf("localhost:%s", port))
}
