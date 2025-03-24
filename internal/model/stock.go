package model

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
type UserSolution struct {
	Symbol    string     `json:"symbol"`
	AfterDate string     `json:"afterDate"`
	DayPrice  []DayPrice `json:"dayPrice"`
}
