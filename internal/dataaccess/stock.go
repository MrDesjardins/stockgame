package dataaccess

import (
	"fmt"
	"stockgame/internal/database"
	"stockgame/internal/model"
)

func GetPricesForStock(symbol string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT rowid, symbol, date, open, high, low, close, adj_close, volume
		FROM stocks
		WHERE symbol = ?
		ORDER BY date ASC
	`
	rows, err := db.Query(query, symbol)
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

func GetPricesForStockInTimeRange(symbol string, startDate string, endDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT rowid, symbol, date, open, high, low, close, adj_close, volume
		FROM stocks
		WHERE symbol = ?
		AND date >= ?
		AND date <= ?
		ORDER BY date ASC
	`
	rows, err := db.Query(query, symbol, startDate, endDate)
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

func GetStocksAfterDate(symbol string, afterDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT rowid, symbol, date, open, high, low, close, adj_close, volume
		FROM stocks
		WHERE symbol = ?
		AND date > ?
		ORDER BY date ASC
		LIMIT ?
	`
	rows, err := db.Query(query, symbol, afterDate, model.User_stock_to_guess)
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

func GetStocksBeforeEqualDate(symbol string, beforeDate string) []model.Stock {
	db := database.GetDB()
	query := `
		SELECT rowid, symbol, date, open, high, low, close, adj_close, volume
		FROM stocks
		WHERE symbol = ?
		AND date <= ?
		ORDER BY date DESC
		LIMIT ?
	`
	rows, err := db.Query(query, symbol, beforeDate, model.Number_initial_stock_shown)
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
