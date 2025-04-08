package main

import (
	"fmt"
	"net/http"
	"os"

	"stockgame/internal/dataaccess"
	"stockgame/internal/database"
	"stockgame/internal/logic"
	"stockgame/internal/model"
	"stockgame/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SolutionHandler struct {
	StockService service.StockService
	ScoringLogic logic.ScoringLogic
}

func (h *SolutionHandler) getStocks(c *gin.Context) {

	stock := h.StockService.GetRandomStockWithRandomDayRange(model.Number_initial_stock_shown)
	c.IndentedJSON(http.StatusOK, stock)
}

func (h *SolutionHandler) postSolution(c *gin.Context) {
	// Read the body of the request
	// Bind JSON directly to a struct
	userSolution := model.UserSolutionRequest{}
	if err := c.BindJSON(&userSolution); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	// Extract values from the struct
	symbolUUID := userSolution.SymbolUUID
	afterDate := userSolution.AfterDate
	dayPrice := userSolution.DayPrice
	if symbolUUID == "" || afterDate == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "symbolUUID and afterDate are required query parameters"})
		return
	}

	real := h.StockService.GetStockInfo(symbolUUID)
	realStocksBeforeDate := h.StockService.GetStocksBeforeEqualDate(real.Symbol, afterDate)
	realStocksAfterDate := h.StockService.GetStocksAfterDate(real.Symbol, afterDate)
	fullList := append(realStocksBeforeDate, realStocksAfterDate...) // To calculuate Bollinger Bands we need the price before and after the date
	// Score
	bollinger20Days := h.ScoringLogic.CalculateBollingerBands(fullList, 20)
	score := h.ScoringLogic.GetScore(dayPrice, realStocksAfterDate, bollinger20Days)
	solutionResponse := model.UserSolutionResponse{
		Symbol: real.Symbol,
		Name:   real.Name,
		Score:  score,
		Stocks: realStocksAfterDate,
		BB20:   bollinger20Days,
	}
	c.IndentedJSON(http.StatusOK, solutionResponse)
}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	port := os.Getenv("VITE_API_PORT")
	database.ConnectDB()

	// Create dependencies in the correct order
	stockDataAccess := &dataaccess.StockDataAccessImpl{}
	stockService := &service.StockServiceImpl{
		StockDataAccess: stockDataAccess,
	}

	// Use stockService in your handler initialization
	handler := &SolutionHandler{
		StockService: stockService,
		ScoringLogic: &logic.ScoringLogicImpl{},
	}

	router := SetupRouter(handler)

	router.Run(fmt.Sprintf("localhost:%s", port))
}

func SetupRouter(handler *SolutionHandler) *gin.Engine {
	router := gin.Default()
	// Cors
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/stocks", handler.getStocks)
	router.POST("/solution", handler.postSolution)
	return router
}
