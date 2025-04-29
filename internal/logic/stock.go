package logic

import "stockgame/internal/model"

type StockLogic interface {
	IsStocksValid(stocks []model.StockPublic, numberOfDays int) (isValid bool, upperBound int)
}

type StockLogicImpl struct {
	StockLogic
}

func (s *StockLogicImpl) IsStocksValid(stocks []model.StockPublic, numberOfDays int) (isValid bool, upperBound int) {
	isValid = true
	upperBound = 0
	volume := 0
	for _, stock := range stocks {
		volume += stock.Volume
	}
	if len(stocks) == 0 {
		isValid = false
		return
	}
	volumeAverage := volume / len(stocks)
	if volumeAverage < 25000 {
		isValid = false
		return
	}
	upperBound = len(stocks) - numberOfDays
	if upperBound <= 0 {
		isValid = false
		upperBound = 0
		return
	}
	// Check if some stock in the slice has an open price to zero
	for _, stock := range stocks {
		if stock.Open == 0 {
			upperBound = 0
			isValid = false
			return
		}
	}
	return
}
