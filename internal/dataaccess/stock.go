package dataaccess

import (
	"fmt"
	"stockgame/internal/database"
	"stockgame/internal/model"
)

func GetPricesForStock(symbol string) []model.StockPublic {
	db := database.GetDB()
	query := `
		SELECT stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume, stocks_info.symbol_uuid
		FROM stocks
		INNER JOIN stocks_info 
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks.symbol = ?
		ORDER BY date ASC
	`
	rows, err := db.Query(query, symbol)
	if err != nil {
		fmt.Println("Error querying stock: ", err)
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
		continue

	}
	return stocks
}

func GetUniqueStockSymbols() []string {
	db := database.GetDB()
	query := `
		SELECT DISTINCT(symbol)
		FROM stocks
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying stock symbols: ", err)
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

func GetPricesForStockInTimeRange(symbolUUID string, startDate string, endDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT stocks.rowid, stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		INNER JOIN stocks_info
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks_info.symbol_uuid = ?
		AND stocks.date >= ?
		AND stocks.date <= ?
		ORDER BY stocks.date ASC
	`
	rows, err := db.Query(query, symbolUUID, startDate, endDate)
	if err != nil {
		fmt.Println("Error querying stock: ", err)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Id, &stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)
		continue

	}
	return stocks
}

func GetStocksAfterDate(symbolUUID string, afterDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT stocks.rowid, stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		INNER JOIN stocks_info
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks_info.symbol_uuid = ?
		AND stocks.date > ?
		ORDER BY stocks.date ASC
		LIMIT ?
	`
	rows, err := db.Query(query, symbolUUID, afterDate, model.User_stock_to_guess)
	if err != nil {
		fmt.Println("Error querying stock: ", err)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Id, &stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)
		continue

	}
	return stocks
}

func GetStocksBeforeEqualDate(symbolUUID string, beforeDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT stocks.rowid, stocks.symbol, stocks.date, stocks.open, stocks.high, stocks.low, stocks.close, stocks.adj_close, stocks.volume
		FROM stocks
		INNER JOIN stocks_info
			ON stocks.symbol = stocks_info.symbol
		WHERE stocks_info.symbol_uuid = ?
		AND stocks.date <= ?
		ORDER BY stocks.date DESC
		LIMIT ?
	`
	rows, err := db.Query(query, symbolUUID, beforeDate, model.Number_initial_stock_shown)
	if err != nil {
		fmt.Println("Error querying stock: ", err)
		return []model.Stock{}
	}
	defer rows.Close()
	var stocks = []model.Stock{}
	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.Id, &stock.Symbol, &stock.Date, &stock.Open, &stock.High, &stock.Low, &stock.Close, &stock.AdjClose, &stock.Volume)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		stocks = append(stocks, stock)
		continue

	}
	return stocks
}

func GetStockInfo(symbolUUID string) model.StockInfo {
	db := database.GetDB()
	query := `
		SELECT symbol, name, symbol_uuid
		FROM stocks_info
		WHERE symbol_uuid = ?
		LIMIT 1
	`
	rows, err := db.Query(query, symbolUUID)
	if err != nil {
		fmt.Println("Error querying stock: ", err)
		return model.StockInfo{
			SymbolUUID: symbolUUID,
			Symbol:     "",
			Name:       "",
		}
	}
	defer rows.Close()
	var stock = model.StockInfo{}
	for rows.Next() {
		err := rows.Scan(&stock.Symbol, &stock.Name, &stock.SymbolUUID)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			continue
		}
		break

	}
	return stock
}
