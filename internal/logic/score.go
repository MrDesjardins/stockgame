package logic

import "stockgame/internal/model"

func GetScore(userPrices []model.DayPrice, actualStockInfo []model.Stock) int {

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

	}

	// Small bonus if the user was in the right direction
	isUserThinkBullish := userPrices[0].Price < userPrices[len(userPrices)-1].Price
	isStockBullish := actualStockInfo[0].Open < actualStockInfo[len(actualStockInfo)-1].Close
	if isUserThinkBullish == isStockBullish {
		score += 10
	}
	return score
}
