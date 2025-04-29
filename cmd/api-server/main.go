package main

import (
	"fmt"
	"log"
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

	real, err := h.StockService.GetStockInfo(symbolUUID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cannot find the stock information"})
		return
	}
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
	env := os.Getenv("GO_ENV")
	println("env", env)
	isProduction := env == "production"
	if !isProduction {
		err := godotenv.Load(".env")
		if err != nil {

			panic("Error loading .env file")
		}
	}
	port := os.Getenv("VITE_API_PORT")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if isProduction {
		dbUrl := os.Getenv("DATABASE_URL")
		database.ConnectDBFullPath(dbUrl)
	} else {
		database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	}
	// Create dependencies in the correct order
	stockDataAccess := &dataaccess.StockDataAccessImpl{
		DB: database.GetDB(),
	}
	stockLogic := &logic.StockLogicImpl{}
	stockService := &service.StockServiceImpl{
		StockDataAccess: stockDataAccess,
		StockLogic:      stockLogic,
	}

	// Use stockService in your handler initialization
	handler := &SolutionHandler{
		StockService: stockService,
		ScoringLogic: &logic.ScoringLogicImpl{},
	}

	router := SetupRouter(handler, isProduction)
	log.Printf("Listening on 0.0.0.0:%s", port)
	router.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

func SetupRouter(handler *SolutionHandler, isProduction bool) *gin.Engine {
	router := gin.Default()
	// Read the HTML file
	path := ""
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
		path = "/usr/local/bin/public/"
	} else {
		path = "./cmd/api-server/public/"
	}

	// Cors
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/stocks", handler.getStocks)
	router.POST("/solution", handler.postSolution)
	router.Static("/assets", path+"assets")
	router.GET("/", func(c *gin.Context) {

		htmlContent, err := os.ReadFile(path + "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading %sindex.html", path))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlContent))
	})
	return router
}
