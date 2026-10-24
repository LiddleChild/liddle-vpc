package frankfurterapi

import "time"

type Rate struct {
	ToTHB float64
	Date  time.Time
}

type frankfurterRate struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}
