package model

type StockPublic struct {
	Date       string  `json:"date"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	AdjClose   float64 `json:"adj_close"`
	Volume     int     `json:"volume"`
	SymbolUUID string  `json:"symbol_uuid"`
}
type Stock struct {
	Id       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Date     string  `json:"date"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	AdjClose float64 `json:"adj_close"`
	Volume   int     `json:"volume"`
}

type DayPrice struct {
	Day   int     `json:"day"`
	Price float64 `json:"price"`
}
type UserSolutionRequest struct {
	SymbolUUID string     `json:"symbolUUID"`
	AfterDate  string     `json:"afterDate"`
	DayPrice   []DayPrice `json:"estimatedDayPrices"`
}

type UserScoreResponse struct {
	Total       int `json:"total"`
	InLowHigh   int `json:"inLowHigh"`
	InOpenClose int `json:"inOpenClose"`
	InBollinger int `json:"inBollinger"`
	InDirection int `json:"inDirection"`
}
type UserSolutionResponse struct {
	Symbol string                   `json:"symbol"`
	Name   string                   `json:"name"`
	Score  UserScoreResponse        `json:"score"`
	Stocks []Stock                  `json:"stocks"`
	BB20   map[string]BollingerBand `json:"bb20"`
}

type BollingerBand struct {
	Date      string  `json:"date"`
	UpperBand float64 `json:"upperBand"`
	Average   float64 `json:"average"`
	LowerBand float64 `json:"lowerBand"`
}

type StockInfo struct {
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	SymbolUUID string `json:"symbol_uuid"`
}
