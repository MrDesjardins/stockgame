package logic

import (
	"fmt"
	"math"
	"stockgame/internal/model"
)

func GetScore(userPrices []model.DayPrice, actualStockInfo []model.Stock, bollinger20Days map[string]model.BollingerBand) int {

	score := 0

	println("userPrices: ", len(userPrices))
	println("actualStockInfo: ", len(actualStockInfo))

	if len(userPrices) == 0 || len(actualStockInfo) == 0 {
		return score
	}
	for i := 0; i < len(userPrices); i++ {
		if i >= len(actualStockInfo) { // In case
			break
		}
		actualStock := actualStockInfo[i]
		if userPrices[i].Price >= actualStock.Low && userPrices[i].Price <= actualStock.High {
			score += 10
		}
		// Additional point if between open/close (harder)
		// First check is if open is lower than close
		if userPrices[i].Price >= actualStock.Open && userPrices[i].Price <= actualStock.Close {
			score += 10
		}
		// Second check is if open is higher than close
		if userPrices[i].Price >= actualStock.Close && userPrices[i].Price <= actualStock.Open {
			score += 10
		}
		// Between Bollinger bands
		if bollingerBand, found := bollinger20Days[actualStock.Date]; found {
			if userPrices[i].Price >= bollingerBand.LowerBand && userPrices[i].Price <= bollingerBand.UpperBand {
				score += 5
			}
		}

	}

	// Small bonus if the user was in the right direction
	isUserThinkBullish := userPrices[0].Price < userPrices[len(userPrices)-1].Price
	isStockBullish := actualStockInfo[0].Open < actualStockInfo[len(actualStockInfo)-1].Close
	if isUserThinkBullish == isStockBullish {
		score += 10
	}
	return score
}

func CalculateBollingerBands(stockInfo []model.Stock, day int) map[string]model.BollingerBand {
	if len(stockInfo) < day {
		return map[string]model.BollingerBand{} // Return empty map if not enough data
	}

	mapDayPrices := make(map[string]model.BollingerBand)
	firstDayGetBBIndex := len(stockInfo) - day

	for i := firstDayGetBBIndex; i < len(stockInfo); i++ {
		if i-day < 0 {
			continue // Skip if there aren't enough past data points
		}

		// Compute moving average
		sum := 0.0
		for j := i - day; j < i; j++ {
			sum += stockInfo[j].Close
		}
		average := sum / float64(day)

		// Compute standard deviation using sample formula (n-1)
		sum = 0.0
		for j := i - day; j < i; j++ {
			diff := stockInfo[j].Close - average
			sum += diff * diff
		}
		standardDeviation := math.Sqrt(sum / float64(day-1)) // Fix: using n-1

		// Compute bands
		upperBand := average + 2*standardDeviation
		lowerBand := average - 2*standardDeviation

		// Store result
		mapDayPrices[stockInfo[i].Date] = model.BollingerBand{
			Date:      stockInfo[i].Date,
			UpperBand: upperBand,
			LowerBand: lowerBand,
		}

		// Debugging: Print values for verification
		fmt.Printf("Date: %s, Close: %.2f, Avg: %.2f, StdDev: %.2f, Upper: %.2f, Lower: %.2f\n",
			stockInfo[i].Date, stockInfo[i].Close, average, standardDeviation, upperBand, lowerBand)
	}

	return mapDayPrices
}
