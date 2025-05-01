package dataaccess

import (
	"context"
	"fmt"
	"stockgame/internal/database"
	"stockgame/internal/model"
)

type StockDataAccess interface {
	GetPricesForStock(ctx context.Context, symbol string) []model.StockPublic
	GetUniqueStockSymbols(ctx context.Context) []string
	GetPricesForStockInTimeRange(ctx context.Context, symbol string, startDate string, endDate string) []model.Stock
	GetStocksAfterDate(ctx context.Context, symbol string, afterDate string) []model.Stock
	GetStocksBeforeEqualDate(ctx context.Context, symbol string, beforeDate string) []model.Stock
	GetStockInfo(ctx context.Context, symbolUUID string) (model.StockInfo, error)
}

type StockDataAccessImpl struct {
	DB database.DBInterface
	StockDataAccess
}

func (s *StockDataAccessImpl) GetPricesForStock(ctx context.Context, symbol string) []model.StockPublic {
	query := `
		SELECT stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume, stocks_info.symbol_uuid
		FROM stocks
		INNER JOIN stocks_info 
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks.symbol = $1
		ORDER BY date ASC
	`
	rows, err := s.DB.QueryContext(ctx, query, symbol)
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
	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
	return stocks
}

func (s *StockDataAccessImpl) GetUniqueStockSymbols(ctx context.Context) []string {
	query := `
		SELECT DISTINCT(symbol)
		FROM stocks
	`
	rows, err := s.DB.QueryContext(ctx, query)
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

func (s *StockDataAccessImpl) GetPricesForStockInTimeRange(ctx context.Context, symbol string, startDate string, endDate string) []model.Stock {
	query := `
		SELECT stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date >= $2
		AND stocks.date <= $3
		ORDER BY stocks.date ASC
	`
	rows, err := s.DB.QueryContext(ctx, query, symbol, startDate, endDate)
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
	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
	return stocks
}

func (s *StockDataAccessImpl) GetStocksAfterDate(ctx context.Context, symbol string, afterDate string) []model.Stock {
	query := `
		SELECT stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date > $2
		ORDER BY stocks.date ASC
		LIMIT $3
	`
	rows, err := s.DB.QueryContext(ctx, query, symbol, afterDate, model.User_stock_to_guess)
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
	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
	return stocks
}

func (s *StockDataAccessImpl) GetStocksBeforeEqualDate(ctx context.Context, symbol string, beforeDate string) []model.Stock {
	query := `
		SELECT  stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		WHERE stocks.symbol = $1
		AND stocks.date <= $2
		ORDER BY stocks.date DESC
		LIMIT $3
	`
	rows, err := s.DB.QueryContext(ctx, query, symbol, beforeDate, model.Number_initial_stock_shown)
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
	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
	return stocks
}

func (s *StockDataAccessImpl) GetStockInfo(ctx context.Context, symbolUUID string) (result model.StockInfo, err error) {
	result = model.StockInfo{
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
		err = fmt.Errorf("database connection is nil")
		return
	}
	rows, err2 := s.DB.QueryContext(ctx, query, symbolUUID)
	if err2 != nil {
		err = fmt.Errorf("error querying stock: %v", err2)
		return
	}
	defer rows.Close()
	if rows.Next() {
		err2 := rows.Scan(&result.Symbol, &result.Name, &result.SymbolUUID)
		if err2 != nil {
			err = fmt.Errorf("error scanning stock: %v", err2)
		}
		return
	} else {
		err = fmt.Errorf("no data found for symbolUUID: %s", symbolUUID)
		return
	}
}
