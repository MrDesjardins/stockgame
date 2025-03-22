package service

import "stockgame/internal/model"

func GetStockFromPersistence() model.Stock {
	return model.Stock{
		Id:       1,
		Symbol:   "AAPL",
		Date:     "2021-01-01",
		Open:     100.0,
		High:     105.0,
		Low:      95.0,
		Close:    102.0,
		AdjClose: 102.0,
		Volume:   1000000,
	}
}
