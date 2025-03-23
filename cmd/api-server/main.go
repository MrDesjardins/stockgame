package main

import (
	"net/http"

	"stockgame/internal/database"
	"stockgame/internal/service"

	"github.com/gin-gonic/gin"
)

func getStocks(c *gin.Context) {
	stock := service.GetRandomStockWithRandomDayRange(1)
	c.IndentedJSON(http.StatusOK, stock)
}

func main() {
	database.ConnectDB()
	router := gin.Default()
	router.GET("/stocks", getStocks)

	router.Run("localhost:8080")
}
