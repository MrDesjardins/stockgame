package dataaccess

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

//	func TestMain(m *testing.M) {
//		database.ConnectDB()
//		defer database.CloseDB()
//		code := m.Run()
//		os.Exit(code)
//	}
// func TestGetStock(t *testing.T) {
// 	mockService := &StockDataAccessImpl{}
// 	EXPECTED := 1000
// 	stock := mockService.GetUniqueStockSymbols()
// 	if len(stock) < EXPECTED {
// 		t.Errorf("Expected atleast %d and not %d", EXPECTED, len(stock))
// 	}
// }

func TestGetStockInfo_QuerySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"symbol", "name", "symbol_uuid"}).
		AddRow("AAPL", "Apple Inc.", "uuid-123")

	mock.ExpectQuery("SELECT symbol, name, symbol_uuid").
		WithArgs("uuid-123").
		WillReturnRows(rows)

	dao := &StockDataAccessImpl{DB: db}
	result, err := dao.GetStockInfo("uuid-123")

	assert.NoError(t, err)
	assert.Equal(t, "AAPL", result.Symbol)
	assert.Equal(t, "Apple Inc.", result.Name)
	assert.Equal(t, "uuid-123", result.SymbolUUID)
}

func TestGetStockInfo_QueryNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"symbol", "name", "symbol_uuid"}) // no rows

	mock.ExpectQuery("SELECT symbol, name, symbol_uuid").
		WithArgs("uuid-456").
		WillReturnRows(rows)

	dao := &StockDataAccessImpl{DB: db}
	result, err := dao.GetStockInfo("uuid-456")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no data found for symbolUUID: uuid-456")
	assert.Equal(t, "uuid-456", result.SymbolUUID)
}

func TestGetStockInfo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT symbol, name, symbol_uuid").
		WithArgs("uuid-789").
		WillReturnError(fmt.Errorf("database error"))

	dao := &StockDataAccessImpl{DB: db}
	result, err := dao.GetStockInfo("uuid-789")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error querying stock")
	assert.Equal(t, "uuid-789", result.SymbolUUID)
}
