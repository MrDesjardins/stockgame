package service

type GeneratorService interface {
	GenerateFiles(numberOfFiles int)
}
type GeneratorServiceImpl struct {
	StockService StockService
}

func (s *GeneratorServiceImpl) GenerateFiles(numberOfFiles int) {
	// Set to know if the stock already taken
	// stockSet := make(map[string]bool)
	// for i := 0; i < numberOfFiles; i++ {
	// 	// Get a random stock symbol
	// 	stockPublics := s.StockService.GetRandomStockWithRandomDayRange(model.Number_initial_stock_shown)
	// 	symbolUUID := stockPublics[0].SymbolUUID
	// 	if stockSet[symbolUUID] {
	// 		continue // Skip if the stock is already taken
	// 	}
	// 	stockSet[symbolUUID] = true
	// 	// Get the stock info
	// 	stockInfo, err := s.StockService.GetStockInfo(symbolUUID)
	// 	if err != nil {
	// 		continue // Skip if the stock info is not found
	// 	}

	// 	// Write the prices to a file
	// 	s.StockService.WriteToFile(stocks, symbolUUID)
	// }
}
