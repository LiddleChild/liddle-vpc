package model

import "time"

type Currency struct {
	ID        int
	Name      string
	RateToTHB float64
	UpdatedAt time.Time
}
