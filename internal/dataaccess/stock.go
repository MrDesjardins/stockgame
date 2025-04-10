package dataaccess

import (
	"fmt"
	"stockgame/internal/database"
	"stockgame/internal/model"
)

type StockDataAccess interface {
	GetPricesForStock(symbol string) []model.StockPublic
	GetUniqueStockSymbols() []string
	GetPricesForStockInTimeRange(symbol string, startDate string, endDate string) []model.Stock
	GetStocksAfterDate(symbol string, afterDate string) []model.Stock
	GetStocksBeforeEqualDate(symbol string, beforeDate string) []model.Stock
	GetStockInfo(symbolUUID string) (model.StockInfo, error)
}
type StockDataAccessImpl struct {
	DB database.DBInterface
	StockDataAccess
}

func (s *StockDataAccessImpl) GetPricesForStock(symbol string) []model.StockPublic {
	query := `
		SELECT stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume, stocks_info.symbol_uuid
		FROM stocks
		INNER JOIN stocks_info 
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks.symbol = $1
		ORDER BY date ASC
	`
	rows, err := s.DB.Query(query, symbol)
	if err != nil {
		fmt.Println("GetPricesForStock Error querying stock: ", err, query)
		return []model.StockPublic{}
	}
	defer rows.Close()
	var stocks = []model.StockPublic{}
	for rows.Next() {
		var stock model.StockPublic
		err := rows.Scan(&stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume, &stock.SymbolUUID)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)

	}
	return stocks
}

func (s *StockDataAccessImpl) GetUniqueStockSymbols() []string {
	query := `
		SELECT DISTINCT(symbol)
		FROM stocks
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println("GetUniqueStockSymbols Error querying stock symbols: ", err, query)
		return []string{}
	}
	defer rows.Close()
	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		symbols = append(symbols, symbol)
	}
	return symbols
}

func (s *StockDataAccessImpl) GetPricesForStockInTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	query := `
		SELECT stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date >= $2
		AND stocks.date <= $3
		ORDER BY stocks.date ASC
	`
	rows, err := s.DB.Query(query, symbol, startDate, endDate)
	if err != nil {
		fmt.Println("GetPricesForStockInTimeRange Error querying stock: ", err, query)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)

	}
	return stocks
}

func (s *StockDataAccessImpl) GetStocksAfterDate(symbol string, afterDate string) []model.Stock {
	query := `
		SELECT stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date > $2
		ORDER BY stocks.date ASC
		LIMIT $3
	`
	rows, err := s.DB.Query(query, symbol, afterDate, model.User_stock_to_guess)
	if err != nil {
		fmt.Println("GetStocksAfterDate Error querying stock: ", err, query)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)

	}
	return stocks
}

func (s *StockDataAccessImpl) GetStocksBeforeEqualDate(symbol string, beforeDate string) []model.Stock {
	query := `
		SELECT  stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date <= $2
		ORDER BY stocks.date DESC
		LIMIT $3
	`
	rows, err := s.DB.Query(query, symbol, beforeDate, model.Number_initial_stock_shown)
	if err != nil {
		fmt.Println("GetStocksBeforeEqualDate Error querying stock: ", err, query)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)

	}
	return stocks
}

func (s *StockDataAccessImpl) GetStockInfo(symbolUUID string) (model.StockInfo, error) {
	errorSymbol := model.StockInfo{
		SymbolUUID: symbolUUID,
		Symbol:     "",
		Name:       "",
	}
	query := `
		SELECT symbol, name, symbol_uuid
		FROM stocks_info
		WHERE symbol_uuid = $1
		LIMIT 1
	`
	if s.DB == nil {
		return errorSymbol, fmt.Errorf("database connection is nil")
	}
	rows, err := s.DB.Query(query, symbolUUID)
	if err != nil {
		return errorSymbol, fmt.Errorf("error querying stock: %v", err)
	}
	if rows.Next() {
		var stock = model.StockInfo{}
		err := rows.Scan(&stock.Symbol, &stock.Name, &stock.SymbolUUID)
		if err != nil {
			return errorSymbol, fmt.Errorf("error scanning stock: %v", err)
		}
		return stock, nil
	} else {
		return errorSymbol, fmt.Errorf("no data found for symbolUUID: %s", symbolUUID)
	}
}
