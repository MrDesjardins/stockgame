package service

import (
	"math/rand/v2"
	"stockgame/internal/dataaccess"
	"stockgame/internal/model"
)

func GetRandomStockWithRandomDayRange(numberOfDays int) []model.Stock {
	stocks := GetRandomStockFromPersistence()
	if len(stocks) < numberOfDays {
		return stocks
	}
	index := rand.IntN(len(stocks) - numberOfDays)
	return stocks[index : index+numberOfDays]
}

func GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	stocks := dataaccess.GetPricesForStockInTimeRange(symbol, startDate, endDate)
	return stocks
}

func GetRandomStockFromPersistence() []model.Stock {
	syms := dataaccess.GetUniqueStockSymbols()
	symbol := GetRandomStock(syms)
	stocks := dataaccess.GetPricesForStock(symbol)
	return stocks
}

func GetRandomStock(symbol []string) string {
	index := rand.IntN(len(symbol))
	return symbol[index]
}
