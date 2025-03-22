package main

import (
	"net/http"

	"stockgame/internal/service"

	"github.com/gin-gonic/gin"
)

func getStock(c *gin.Context) {
	stock := service.GetStockFromPersistence()
	c.IndentedJSON(http.StatusOK, stock)
}

func main() {
	router := gin.Default()
	router.GET("/stock", getStock)

	router.Run("localhost:8080")
}
