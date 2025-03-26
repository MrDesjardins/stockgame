package logic

import (
	"math"
	"stockgame/internal/model"
)

func GetScore(userPrices []model.DayPrice, actualStockInfo []model.Stock, bollinger20Days map[string]model.BollingerBand) model.UserScoreResponse {
	println("userPrices: ", len(userPrices))
	println("actualStockInfo: ", len(actualStockInfo))
	scoreObj := model.UserScoreResponse{
		Total:       0,
		InLowHigh:   0,
		InOpenClose: 0,
		InBollinger: 0,
		InDirection: 0,
	} // Return empty struct if no data

	if len(userPrices) == 0 || len(actualStockInfo) == 0 {
		return scoreObj
	}
	for i := range userPrices {
		if i >= len(actualStockInfo) { // In case
			break
		}
		actualStock := actualStockInfo[i]
		// Check if user price is within the actual stock low/high of the day
		if userPrices[i].Price >= actualStock.Low && userPrices[i].Price <= actualStock.High {
			scoreObj.InLowHigh += 10 + 2*i // Bonus if the prediction is accurate the farther in the future
		}
		// Additional point if between open/close (harder)
		// First check is if open is lower than close
		if userPrices[i].Price >= actualStock.Open && userPrices[i].Price <= actualStock.Close {
			scoreObj.InOpenClose += 10 + 2*i // Bonus if the prediction is accurate the farther in the future
		}
		// Second check is if open is higher than close
		if userPrices[i].Price >= actualStock.Close && userPrices[i].Price <= actualStock.Open {
			scoreObj.InOpenClose += 10 + 2*i // Bonus if the prediction is accurate the farther in the future
		}
		// Between Bollinger bands
		if bollingerBand, found := bollinger20Days[actualStock.Date]; found {
			if userPrices[i].Price >= bollingerBand.LowerBand && userPrices[i].Price <= bollingerBand.UpperBand {
				scoreObj.InBollinger += 5
			}
		}

	}

	// Small bonus if the user was in the right direction
	isUserThinkBullish := userPrices[0].Price < userPrices[len(userPrices)-1].Price
	isStockBullish := actualStockInfo[0].Open < actualStockInfo[len(actualStockInfo)-1].Close
	if isUserThinkBullish == isStockBullish {
		scoreObj.InDirection += 10
	}
	scoreObj.Total = scoreObj.InLowHigh + scoreObj.InOpenClose + scoreObj.InBollinger + scoreObj.InDirection
	return scoreObj
}

func CalculateBollingerBands(stockInfo []model.Stock, day int) map[string]model.BollingerBand {
	if len(stockInfo) < day {
		return map[string]model.BollingerBand{} // Return empty map if not enough data
	}

	mapDayPrices := make(map[string]model.BollingerBand)

	for i := day; i < len(stockInfo); i++ { // Start at 'day' for correct rolling window
		// Compute moving average
		sum := 0.0
		for j := i - day; j < i; j++ { // Window spans 'i-day' to 'i-1'
			sum += stockInfo[j].Close
		}
		average := sum / float64(day)

		// Compute standard deviation
		sumSquares := 0.0
		for j := i - day; j < i; j++ {
			diff := stockInfo[j].Close - average
			sumSquares += diff * diff
		}
		standardDeviation := math.Sqrt(sumSquares / float64(day)) // Use day (population) or day-1 (sample)

		// Compute bands
		upperBand := average + 2*standardDeviation
		lowerBand := average - 2*standardDeviation

		// Store result
		mapDayPrices[stockInfo[i].Date] = model.BollingerBand{
			Date:      stockInfo[i].Date,
			UpperBand: upperBand,
			Average:   average,
			LowerBand: lowerBand,
		}
	}

	return mapDayPrices
}
