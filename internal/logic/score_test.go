package logic

import (
	"fmt"
	"math"
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
	bb20 := make(map[string]model.BollingerBand)
	stock := GetScore(userPrice, actualStockInfo, bb20)
	if stock.Total != 76 {
		t.Errorf("Expected score to be 76 and not %d", stock.Total)
	}
}

func TestCalculateBollingBands(t *testing.T) {
	stockInfo := []model.Stock{}
	for i := 0; i < 40; i++ {
		stockInfo = append(stockInfo, model.Stock{
			Id:       10002,
			Date:     fmt.Sprintf("2022-01-%02d", i+1),
			Open:     100,
			High:     100,
			Low:      100,
			Close:    float64(i + 1),
			AdjClose: 100,
			Volume:   100,
		})
	}

	bb20 := CalculateBollingerBands(stockInfo, 20)
	for k, v := range bb20 {
		fmt.Printf("Date: %s, LowerBand: %f, UpperBand: %f\n", k, v.LowerBand, v.UpperBand)
	}
	if len(bb20) != 20 {
		t.Errorf("Expected 1 Bollinger band and not %d", len(bb20))
	}
	if !almostEqual(-1.032, bb20["2022-01-21"].LowerBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-21"].LowerBand)
	}
	if !almostEqual(22.032, bb20["2022-01-21"].UpperBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-21"].UpperBand)
	}

	if !almostEqual(17.967, bb20["2022-01-40"].LowerBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-40"].LowerBand)
	}
	if !almostEqual(41.032, bb20["2022-01-40"].UpperBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-40"].UpperBand)
	}
}

func TestCalculateBollingBands2(t *testing.T) {
	stockInfo := []model.Stock{}
	for i := 0; i < 10; i++ {
		stockInfo = append(stockInfo, model.Stock{
			Id:       10002,
			Date:     fmt.Sprintf("2022-01-%02d", i+1),
			Open:     100,
			High:     100,
			Low:      100,
			Close:    float64(9 + (i%2)*2), // 9 or 11
			AdjClose: 100,
			Volume:   100,
		})
	}

	bb20 := CalculateBollingerBands(stockInfo, 4)
	for k, v := range bb20 {
		fmt.Printf("Date: %s, LowerBand: %f, Average: %f, UpperBand: %f\n", k, v.LowerBand, v.Average, v.UpperBand)
	}
	if len(bb20) != 6 {
		t.Errorf("Expected 1 Bollinger band and not %d", len(bb20))
	}
	if !almostEqual(8.000, bb20["2022-01-05"].LowerBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-05"].LowerBand)
	}
	if !almostEqual(12.000, bb20["2022-01-05"].UpperBand, 0.001) {
		t.Errorf("Expected LowerBand must not be %f", bb20["2022-01-05"].UpperBand)
	}
}

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
