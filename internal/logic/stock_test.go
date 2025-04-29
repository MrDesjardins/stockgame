package logic

import (
	"stockgame/internal/model"
	"testing"
)

func TestIsStocksValid(t *testing.T) {
	t.Run("WithValidStocks", func(t *testing.T) {
		mockService := &StockLogicImpl{}
		stocks := []model.StockPublic{
			{Date: "2022-01-01", Open: 100, High: 105, Low: 95, Close: 100, Volume: 30000},
			{Date: "2022-01-02", Open: 102, High: 107, Low: 97, Close: 102, Volume: 50000},
			{Date: "2022-01-03", Open: 105, High: 110, Low: 100, Close: 105, Volume: 60000},
		}
		numberOfDays := 2
		isValid, upperBound := mockService.IsStocksValid(stocks, numberOfDays)
		if !isValid {
			t.Errorf("Expected stocks to be valid")
		}
		if upperBound != len(stocks)-numberOfDays {
			t.Errorf("Expected upper bound to be %d and not %d", len(stocks)-numberOfDays, upperBound)
		}
	})
	t.Run("with low volumne", func(t *testing.T) {
		mockService := &StockLogicImpl{}
		stocks := []model.StockPublic{
			{Date: "2022-01-01", Open: 100, High: 105, Low: 95, Close: 100, Volume: 1},
			{Date: "2022-01-02", Open: 102, High: 107, Low: 97, Close: 102, Volume: 1},
			{Date: "2022-01-03", Open: 105, High: 110, Low: 100, Close: 105, Volume: 1},
		}
		numberOfDays := 2
		isValid, upperBound := mockService.IsStocksValid(stocks, numberOfDays)
		if isValid {
			t.Errorf("Expected stocks to be valid because of low volume")
		}
		if upperBound != 0 {
			t.Errorf("Expected upper bound to be 0, not %d because it is invalid enough volume", upperBound)
		}
	})
	t.Run("with not enought day", func(t *testing.T) {
		mockService := &StockLogicImpl{}
		stocks := []model.StockPublic{
			{Date: "2022-01-01", Open: 100, High: 105, Low: 95, Close: 100, Volume: 100000},
			{Date: "2022-01-02", Open: 102, High: 107, Low: 97, Close: 102, Volume: 100000},
			{Date: "2022-01-03", Open: 105, High: 110, Low: 100, Close: 105, Volume: 100000},
		}
		numberOfDays := 4
		isValid, upperBound := mockService.IsStocksValid(stocks, numberOfDays)
		if isValid {
			t.Errorf("Expected stocks to be invalid because not enough days")
		}
		if upperBound != 0 {
			t.Errorf("Expected upper bound to be 0, not %d because it is invalid enough day", upperBound)
		}
	})
	t.Run("with zero amount at open", func(t *testing.T) {
		mockService := &StockLogicImpl{}
		stocks := []model.StockPublic{
			{Date: "2022-01-01", Open: 0, High: 105, Low: 95, Close: 100, Volume: 100000},
			{Date: "2022-01-02", Open: 102, High: 107, Low: 97, Close: 102, Volume: 100000},
			{Date: "2022-01-03", Open: 105, High: 110, Low: 100, Close: 105, Volume: 100000},
		}
		numberOfDays := 2
		isValid, upperBound := mockService.IsStocksValid(stocks, numberOfDays)
		if isValid {
			t.Errorf("Expected stocks to have all stock open above zero")
		}
		if upperBound != 0 {
			t.Errorf("Expected upper bound to be 0, not %d.", upperBound)
		}
	})
}
