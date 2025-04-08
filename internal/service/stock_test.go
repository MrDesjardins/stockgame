package service

import (
	"slices"
	"stockgame/internal/dataaccess"
	"stockgame/internal/model"
	"testing"
)

func TestGetRandomStock(t *testing.T) {
	choices := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "META"}
	mockDataAccess := &dataaccess.StockDataAccessImpl{}
	mockService := &StockServiceImpl{
		StockDataAccess: mockDataAccess,
	}
	stock := mockService.GetRandomStock(choices)
	if slices.Contains(choices, stock) == false {
		t.Errorf("Expected to find %s", stock)
	}
}

type StockDataAccessMockImpl struct {
	GetPricesForStockFunc            func(symbol string) []model.StockPublic
	GetUniqueStockSymbolsFunc        func() []string
	GetUniqueStockSymbolsFuncCall    int
	GetPricesForStockInTimeRangeFunc func(symbol, startDate, endDate string) []model.Stock
	GetStocksAfterDateFunc           func(symbol, afterDate string) []model.Stock
	GetStocksBeforeEqualDateFunc     func(symbol, beforeDate string) []model.Stock
	GetStockInfoFunc                 func(symbolUUID string) model.StockInfo
}

func (s *StockDataAccessMockImpl) GetPricesForStock(symbol string) []model.StockPublic {
	if s.GetPricesForStockFunc != nil {
		return s.GetPricesForStockFunc(symbol)
	}
	return nil
}

func (s *StockDataAccessMockImpl) GetUniqueStockSymbols() []string {
	if s.GetUniqueStockSymbolsFunc != nil {
		s.GetUniqueStockSymbolsFuncCall++
		return s.GetUniqueStockSymbolsFunc()
	}
	return nil
}

func (s *StockDataAccessMockImpl) GetPricesForStockInTimeRange(symbol, startDate, endDate string) []model.Stock {
	if s.GetPricesForStockInTimeRangeFunc != nil {
		return s.GetPricesForStockInTimeRangeFunc(symbol, startDate, endDate)
	}
	return nil
}

func (s *StockDataAccessMockImpl) GetStocksAfterDate(symbol, afterDate string) []model.Stock {
	if s.GetStocksAfterDateFunc != nil {
		return s.GetStocksAfterDateFunc(symbol, afterDate)
	}
	return nil
}

func (s *StockDataAccessMockImpl) GetStocksBeforeEqualDate(symbol, beforeDate string) []model.Stock {
	if s.GetStocksBeforeEqualDateFunc != nil {
		return s.GetStocksBeforeEqualDateFunc(symbol, beforeDate)
	}
	return nil
}

func (s *StockDataAccessMockImpl) GetStockInfo(symbolUUID string) model.StockInfo {
	if s.GetStockInfoFunc != nil {
		return s.GetStockInfoFunc(symbolUUID)
	}
	return model.StockInfo{}
}

func TestGetRandomStockFromPersistence(t *testing.T) {
	// Create a mockDataAccess object with a function GetUniqueStockSymbols that return fake symboles
	mockDataAccess := &StockDataAccessMockImpl{
		GetUniqueStockSymbolsFunc: func() []string {
			return []string{"AAPL", "GOOGL"}
		},
		GetPricesForStockFunc: func(symbol string) []model.StockPublic {
			return []model.StockPublic{
				{SymbolUUID: "AAPL", Volume: 1000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				{SymbolUUID: "AAPL", Volume: 2000, Date: "2023-01-02", Open: 101, High: 111, Low: 89, Close: 107, AdjClose: 1550},
			}
		},
	}
	mockService := &StockServiceImpl{
		StockDataAccess: mockDataAccess,
		GetRandomStockSelectorFunc: func(symbols []string) string {
			return "AAPL"
		},
	}
	stocks := mockService.GetRandomStockFromPersistence()
	if len(stocks) != 2 {
		t.Errorf("Expected the same amount of stock found in the database but found %d", len(stocks))
	}
	if mockDataAccess.GetUniqueStockSymbolsFuncCall != 1 {
		t.Errorf("Expected to call GetUniqueStockSymbols once but called %d times", mockDataAccess.GetUniqueStockSymbolsFuncCall)
	}
}

func TestGetRandomStockWithRandomDayRange(t *testing.T) {
	t.Run("No stock from persistence", func(t *testing.T) {
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				return []model.StockPublic{}
			},
		}
		stocks := mockService.GetRandomStockWithRandomDayRange(2)
		if len(stocks) > 0 {
			t.Errorf("Expected to find some stocks no stocks but found %d", len(stocks))
		}
	})

	t.Run("Found Stocks with high volume, open price above zero, enough to cover the time period", func(t *testing.T) {
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				return []model.StockPublic{
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-03", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				}
			},
		}
		stocks := mockService.GetRandomStockWithRandomDayRange(2)
		if len(stocks) == 0 {
			t.Errorf("Expected to find some stocks but found 0")
		}
	})

	t.Run("Found Stocks but not enough to cover the time period", func(t *testing.T) {
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				return []model.StockPublic{
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-03", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				}
			},
		}
		stocks := mockService.GetRandomStockWithRandomDayRange(20)
		if len(stocks) > 0 {
			t.Errorf("Expected to find some stocks but found 0")
		}
	})

	t.Run("Found Stocks but open price at zero", func(t *testing.T) {
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				return []model.StockPublic{
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 0, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-03", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				}
			},
		}
		stocks := mockService.GetRandomStockWithRandomDayRange(2)
		if len(stocks) > 0 {
			t.Errorf("Expected to find some stocks but found 0")
		}
	})

	t.Run("Found Stocks but average volume too low", func(t *testing.T) {
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				return []model.StockPublic{
					{SymbolUUID: "AAPL", Volume: 24000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 24000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					{SymbolUUID: "AAPL", Volume: 24000, Date: "2023-01-03", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				}
			},
		}
		stocks := mockService.GetRandomStockWithRandomDayRange(2)
		if len(stocks) > 0 {
			t.Errorf("Expected to find some stocks but found 0")
		}
	})
	t.Run("Will Retry 15 times maximum", func(t *testing.T) {
		count := 0
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				count++
				return []model.StockPublic{
					{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
				}
			},
		}
		mockService.GetRandomStockWithRandomDayRange(2)
		if count != 15 {
			t.Errorf("Expected to try 15 times but tried %d", count)
		}
	})

	t.Run("Will Retry 2 times if the second time is good data", func(t *testing.T) {
		count := 0
		mockDataAccess := &StockDataAccessMockImpl{}
		mockService := &StockServiceImpl{
			StockDataAccess: mockDataAccess,
			GetRandomStockFromPersistenceSelectorFunc: func() []model.StockPublic {
				count++
				if count == 1 {
					// Bad data
					return []model.StockPublic{
						{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					}
				} else {
					// Good data
					return []model.StockPublic{
						{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-01", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
						{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
						{SymbolUUID: "AAPL", Volume: 50000, Date: "2023-01-02", Open: 100, High: 110, Low: 90, Close: 105, AdjClose: 1233},
					}
				}

			},
		}
		mockService.GetRandomStockWithRandomDayRange(2)
		if count != 2 {
			t.Errorf("Expected to try 2 times but tried %d", count)
		}
	})
}
