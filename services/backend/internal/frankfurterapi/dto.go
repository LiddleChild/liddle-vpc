package frankfurterapi

import "time"

type ExchangeRateResponse map[string]ExchangeRate

type ExchangeRate struct {
	Rate      float64
	UpdatedAt time.Time
}

type frankfurterRate struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}
