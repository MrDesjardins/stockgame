package service

import (
	"context"
	"math/rand/v2"
	"stockgame/internal/dataaccess"
	"stockgame/internal/database"
	"stockgame/internal/logic"
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
	StockLogic                                logic.StockLogic
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
		isValid, upperBound := s.StockLogic.IsStocksValid(stocks, numberOfDays)
		if !isValid {
			continue OuterLoop // Try again
		}
		lowerBound := rand.IntN(upperBound)
		return stocks[lowerBound : lowerBound+numberOfDays] // Found a good candidate
	}
	return []model.StockPublic{}
}

func (s *StockServiceImpl) GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock {

	ctx, cancel := context.WithTimeout(context.Background(), database.CONTEXT_TIMEOUT)
	defer cancel()
	stocks := s.StockDataAccess.GetPricesForStockInTimeRange(ctx, symbol, startDate, endDate)
	return stocks
}
func (s *StockServiceImpl) GetStocksBeforeEqualDate(symbol string, beforeDate string) []model.Stock {
	ctx, cancel := context.WithTimeout(context.Background(), database.CONTEXT_TIMEOUT)
	defer cancel()
	stocks := s.StockDataAccess.GetStocksBeforeEqualDate(ctx, symbol, beforeDate)
	return stocks
}
func (s *StockServiceImpl) GetStockInfo(symbolUUID string) (model.StockInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.CONTEXT_TIMEOUT)
	defer cancel()
	stock, err := s.StockDataAccess.GetStockInfo(ctx, symbolUUID)
	return stock, err
}
func (s *StockServiceImpl) GetStocksAfterDate(symbolUUID string, afterDate string) []model.Stock {
	ctx, cancel := context.WithTimeout(context.Background(), database.CONTEXT_TIMEOUT)
	defer cancel()
	stocks := s.StockDataAccess.GetStocksAfterDate(ctx, symbolUUID, afterDate)
	return stocks
}

func (s *StockServiceImpl) GetRandomStockFromPersistence() []model.StockPublic {
	ctx, cancel := context.WithTimeout(context.Background(), database.CONTEXT_TIMEOUT)
	defer cancel()
	syms := s.StockDataAccess.GetUniqueStockSymbols(ctx)

	var symbol string
	if s.GetRandomStockSelectorFunc != nil {
		symbol = s.GetRandomStockSelectorFunc(syms)
	} else {
		symbol = s.GetRandomStock(syms) // fallback to actual implementation
	}

	stocks := s.StockDataAccess.GetPricesForStock(ctx, symbol)
	return stocks
}

func (s *StockServiceImpl) GetRandomStock(symbol []string) string {
	index := rand.IntN(len(symbol))
	return symbol[index]
}
