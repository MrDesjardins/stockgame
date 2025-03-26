package service

import (
	"math/rand/v2"
	"stockgame/internal/dataaccess"
	"stockgame/internal/model"
)

func GetRandomStockWithRandomDayRange(numberOfDays int) []model.Stock {
	for numberOfTry := 0; numberOfTry < 100; numberOfTry++ {
		stocks := GetRandomStockFromPersistence()

		// Check if there is activity (volume) for the days of the stock
		volume := 0
		for _, stock := range stocks {

			volume += stock.Volume
		}
		volumeAverage := volume / len(stocks)
		println("Volume average: ", volumeAverage)
		if volumeAverage < 25000 {
			continue // Try again
		}

		if len(stocks) < numberOfDays {
			continue // Try again
		}
		index := rand.IntN(len(stocks) - numberOfDays)
		return stocks[index : index+numberOfDays] // Found a good candidate
	}
	return []model.Stock{}
}

func GetStockPriceForTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	stocks := dataaccess.GetPricesForStockInTimeRange(symbol, startDate, endDate)
	return stocks
}

func GetStockBeforeEqualDate(symbol string, beforeDate string) []model.Stock {
	stocks := dataaccess.GetStocksBeforeEqualDate(symbol, beforeDate)
	return stocks
}

func GetStocksAfterDate(symbol string, afterDate string) []model.Stock {
	stocks := dataaccess.GetStocksAfterDate(symbol, afterDate)
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
