package service

import (
	"os"
	"slices"
	"stockgame/internal/database"
	"testing"
)

func TestMain(m *testing.M) {
	database.ConnectDB()
	defer database.CloseDB()
	code := m.Run()
	os.Exit(code)
}

func TestGetRandomStock(t *testing.T) {
	choices := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "META"}

	stock := GetRandomStock(choices)
	if slices.Contains(choices, stock) == false {
		t.Errorf("Expected to find %s", stock)
	}
}

func TestGetRandomStockFromPersistence(t *testing.T) {
	stocks := GetRandomStockFromPersistence()
	if len(stocks) > 0 {
		t.Errorf("Expected a stock symbol")
	}
}
