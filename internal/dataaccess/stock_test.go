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
	EXPECTED := 1000
	stock := GetUniqueStockSymbols()
	if len(stock) < EXPECTED {
		t.Errorf("Expected at last %d", EXPECTED)
	}
}
