package dataaccess

import (
	"os"
	"stockgame/internal/database"
	"testing"
)

func TestMain(m *testing.M) {
	database.ConnectDB()
	defer database.CloseDB()
	code := m.Run()
	os.Exit(code)
}
func TestGetStock(t *testing.T) {
	mockService := &StockDataAccessImpl{}
	EXPECTED := 1000
	stock := mockService.GetUniqueStockSymbols()
	if len(stock) < EXPECTED {
		t.Errorf("Expected atleast %d and not %d", EXPECTED, len(stock))
	}
}
