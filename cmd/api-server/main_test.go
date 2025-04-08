package main

import (
	"net/http"
	"net/http/httptest"
	"stockgame/internal/model"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type StockServiceMockImpl struct {
	mock.Mock
}

type ScoringLogicMockImpl struct {
	mock.Mock
}

func (m *StockServiceMockImpl) GetRandomStockWithRandomDayRange(n int) []model.StockPublic {
	args := m.Called(n)
	return args.Get(0).([]model.StockPublic)
}
func (m *StockServiceMockImpl) GetRandomStock(symbol []string) string {
	args := m.Called()
	return args.Get(0).(string)
}
func (m *StockServiceMockImpl) GetStockInfo(symbolUUID string) model.StockInfo {
	args := m.Called(symbolUUID)
	return args.Get(0).(model.StockInfo)
}
func (m *StockServiceMockImpl) GetStocksAfterDate(symbol, date string) []model.Stock {
	args := m.Called(symbol, date)
	return args.Get(0).([]model.Stock)
}
func (m *StockServiceMockImpl) GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	args := m.Called(symbol, startDate, endDate)
	return args.Get(0).([]model.Stock)
}
func (m *StockServiceMockImpl) GetRandomStockFromPersistence() []model.StockPublic {
	args := m.Called()
	return args.Get(0).([]model.StockPublic)
}
func (m *StockServiceMockImpl) GetStocksBeforeEqualDate(symbol, date string) []model.Stock {
	args := m.Called(symbol, date)
	return args.Get(0).([]model.Stock)
}

func (m *ScoringLogicMockImpl) CalculateBollingerBands(stockInfo []model.Stock, day int) map[string]model.BollingerBand {
	args := m.Called(stockInfo, day)
	return args.Get(0).(map[string]model.BollingerBand)
}
func (m *ScoringLogicMockImpl) GetScore(userPrices []model.DayPrice, actualStockInfo []model.Stock, bollinger20Days map[string]model.BollingerBand) model.UserScoreResponse {
	args := m.Called(userPrices, actualStockInfo, bollinger20Days)
	return args.Get(0).(model.UserScoreResponse)
}

func TestApiServerRequestGetStocks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("GetStocks", func(t *testing.T) {
		// Arrange
		mockService := new(StockServiceMockImpl)
		expectedStocks := []model.StockPublic{
			{
				SymbolUUID: "AAPL",
				Date:       "2023-10-01",
				Open:       float64(150.0),
				High:       float64(155.0),
				Low:        float64(148.0),
				Close:      float64(152.0),
				AdjClose:   float64(152.0),
				Volume:     1000,
			},
			{
				SymbolUUID: "GOOG",
				Date:       "2024-05-03",
				Open:       float64(50.0),
				High:       float64(55.0),
				Low:        float64(48.0),
				Close:      float64(52.0),
				AdjClose:   float64(52.0),
				Volume:     200,
			},
		}

		mockService.On("GetRandomStockWithRandomDayRange", 40).Return(expectedStocks)

		handler := &SolutionHandler{
			StockService: mockService,
			ScoringLogic: nil, // Assuming you don't need this for the test
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		router := SetupRouter(handler)

		req, _ := http.NewRequest(http.MethodGet, "/stocks", nil)
		c.Request = req

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		assert.Contains(t, body, "AAPL")
		assert.Contains(t, body, "GOOG")

		mockService.AssertExpectations(t)
	})
}

func TestApiServerRequestPostSolution(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run(("InvalidBody"), func(t *testing.T) {
		// Arrange
		mockService := new(StockServiceMockImpl)
		mockScoringLogic := new(ScoringLogicMockImpl)
		handler := &SolutionHandler{
			StockService: mockService,
			ScoringLogic: mockScoringLogic,
		}
		body := "{invalid json}"
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		bodyIoReader := strings.NewReader(body)

		router := SetupRouter(handler)

		req, _ := http.NewRequest(http.MethodPost, "/solution", bodyIoReader)
		c.Request = req
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "Invalid JSON")
	})
	t.Run(("Missing symbolUUID"), func(t *testing.T) {
		// Arrange
		mockService := new(StockServiceMockImpl)
		mockScoringLogic := new(ScoringLogicMockImpl)
		handler := &SolutionHandler{
			StockService: mockService,
			ScoringLogic: mockScoringLogic,
		}
		body := `{"afterDate": "2023-10-01", "estimatedDayPrices": []}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		bodyIoReader := strings.NewReader(body)

		router := SetupRouter(handler)

		req, _ := http.NewRequest(http.MethodPost, "/solution", bodyIoReader)
		c.Request = req
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "symbolUUID and afterDate are required query parameters")
	})
	t.Run(("Missing afterDate"), func(t *testing.T) {
		// Arrange
		mockService := new(StockServiceMockImpl)
		mockScoringLogic := new(ScoringLogicMockImpl)
		handler := &SolutionHandler{
			StockService: mockService,
			ScoringLogic: mockScoringLogic,
		}
		body := `{"symbolUUID": "AAPL", "estimatedDayPrices": []}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		bodyIoReader := strings.NewReader(body)

		router := SetupRouter(handler)

		req, _ := http.NewRequest(http.MethodPost, "/solution", bodyIoReader)
		c.Request = req
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "symbolUUID and afterDate are required query parameters")
	})
	t.Run("NoUserEstimatedPrice", func(t *testing.T) {
		mockService := new(StockServiceMockImpl)
		mockScoringLogic := new(ScoringLogicMockImpl)
		mockService.On("GetStockInfo", "AAPL").Return(model.StockInfo{
			Symbol: "AAPL",
			Name:   "Apple Inc.",
		})
		mockService.On("GetStocksBeforeEqualDate", "AAPL", "2023-10-01").Return([]model.Stock{
			{
				Symbol:   "AAPL",
				Date:     "2023-09-30",
				Open:     150.0,
				High:     155.0,
				Low:      148.0,
				Close:    152.0,
				AdjClose: 152.0,
				Volume:   1000,
			},
		})
		mockService.On("GetStocksAfterDate", "AAPL", "2023-10-01").Return([]model.Stock{
			{
				Symbol:   "AAPL",
				Date:     "2023-10-01",
				Open:     150.0,
				High:     155.0,
				Low:      148.0,
				Close:    152.0,
				AdjClose: 152.0,
				Volume:   1000,
			},
		})
		mockScoringLogic.On("CalculateBollingerBands", mock.Anything, 20).Return(map[string]model.BollingerBand{
			"2023-10-01": {
				LowerBand: 145.0,
				UpperBand: 155.0,
			},
		})
		mockScoringLogic.On("GetScore", mock.Anything, mock.Anything, mock.Anything).Return(model.UserScoreResponse{
			Total:       0,
			InLowHigh:   0,
			InOpenClose: 0,
			InBollinger: 0,
			InDirection: 0,
		})

		handler := &SolutionHandler{
			StockService: mockService,
			ScoringLogic: mockScoringLogic,
		}
		body := `{"symbolUUID": "AAPL", "afterDate": "2023-10-01", "estimatedDayPrices": []}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		bodyIoReader := strings.NewReader(body)
		req, _ := http.NewRequest(http.MethodPost, "/solution", bodyIoReader)
		c.Request = req
		c.Request.Header.Set("Content-Type", "application/json")
		handler.postSolution(c)
		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "score")
		assert.Contains(t, responseBody, "total")
		assert.Contains(t, responseBody, "inLowHigh")
		assert.Contains(t, responseBody, "inOpenClose")
		assert.Contains(t, responseBody, "inDirection")
	})
}
