package service

import (
	"math/rand/v2"
	"stockgame/internal/dataaccess"
	"stockgame/internal/model"
)

type StockService interface {
	GetStockInfo(symbolUUID string) (model.StockInfo, error)
	GetStocksBeforeEqualDate(symbol, date string) []model.Stock
	GetStocksAfterDate(symbol, date string) []model.Stock
	GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock
	GetRandomStockWithRandomDayRange(numberOfDays int) []model.StockPublic
	GetRandomStockFromPersistence() []model.StockPublic
	GetRandomStock(symbol []string) string
}
type StockServiceImpl struct {
	StockDataAccess                           dataaccess.StockDataAccess
	GetRandomStockSelectorFunc                func(choices []string) string
	GetRandomStockFromPersistenceSelectorFunc func() []model.StockPublic
}

func (s *StockServiceImpl) GetRandomStockWithRandomDayRange(numberOfDays int) []model.StockPublic {
OuterLoop:
	for numberOfTry := 0; numberOfTry < 15; numberOfTry++ {
		var stocks []model.StockPublic
		if s.GetRandomStockFromPersistenceSelectorFunc != nil {
			stocks = s.GetRandomStockFromPersistenceSelectorFunc()
		} else {
			stocks = s.GetRandomStockFromPersistence()
		}
		// Check if there is activity (volume) for the days of the stock
		volume := 0
		for _, stock := range stocks {
			volume += stock.Volume
		}
		if len(stocks) == 0 {
			continue // Try again
		}
		volumeAverage := volume / len(stocks)
		if volumeAverage < 25000 {
			continue // Try again
		}
		upperBound := len(stocks) - numberOfDays
		if upperBound <= 0 {
			continue // Try again
		}
		// Check if some stock in the slice has an open price to zero
		for _, stock := range stocks {
			if stock.Open == 0 {
				continue OuterLoop // Try again
			}
		}
		lowerBound := rand.IntN(upperBound)
		return stocks[lowerBound : lowerBound+numberOfDays] // Found a good candidate
	}
	return []model.StockPublic{}
}

func (s *StockServiceImpl) GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	stocks := s.StockDataAccess.GetPricesForStockInTimeRange(symbol, startDate, endDate)
	return stocks
}
func (s *StockServiceImpl) GetStocksBeforeEqualDate(symbol string, beforeDate string) []model.Stock {
	stocks := s.StockDataAccess.GetStocksBeforeEqualDate(symbol, beforeDate)
	return stocks
}
func (s *StockServiceImpl) GetStockInfo(symbolUUID string) (model.StockInfo, error) {
	stock, err := s.StockDataAccess.GetStockInfo(symbolUUID)
	return stock, err
}
func (s *StockServiceImpl) GetStocksAfterDate(symbolUUID string, afterDate string) []model.Stock {
	stocks := s.StockDataAccess.GetStocksAfterDate(symbolUUID, afterDate)
	return stocks
}

func (s *StockServiceImpl) GetRandomStockFromPersistence() []model.StockPublic {
	syms := s.StockDataAccess.GetUniqueStockSymbols()

	var symbol string
	if s.GetRandomStockSelectorFunc != nil {
		symbol = s.GetRandomStockSelectorFunc(syms)
	} else {
		symbol = s.GetRandomStock(syms) // fallback to actual implementation
	}

	stocks := s.StockDataAccess.GetPricesForStock(symbol)
	return stocks
}

func (s *StockServiceImpl) GetRandomStock(symbol []string) string {
	index := rand.IntN(len(symbol))
	return symbol[index]
}
