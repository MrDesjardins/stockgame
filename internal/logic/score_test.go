package logic

import (
	"stockgame/internal/model"
	"testing"
)

func TestGetScoreWithPartialDayBetweenLowAndHigh(t *testing.T) {
	userPrice := []model.DayPrice{
		{Day: 41,
			Price: 100},
		{Day: 42,
			Price: 102},
		{Day: 43,
			Price: 103},
	}
	actualStockInfo := []model.Stock{
		{Id: 10002,
			Symbol:   "AAPL",
			Date:     "2022-01-01",
			Open:     100,
			High:     105,
			Low:      95,
			Close:    100,
			AdjClose: 100,
			Volume:   100},
		{Id: 10003,
			Symbol:   "AAPL",
			Date:     "2022-01-02",
			Open:     102,
			High:     107,
			Low:      97,
			Close:    102,
			AdjClose: 102,
			Volume:   102},
		{Id: 10004,
			Symbol:   "AAPL",
			Date:     "2022-01-02",
			Open:     105,
			High:     110,
			Low:      105,
			Close:    110,
			AdjClose: 102,
			Volume:   102},
	}
	stock := GetScore(userPrice, actualStockInfo)
	if stock != 70 {
		t.Errorf("Expected score to be 70 and not %d", stock)
	}
}
