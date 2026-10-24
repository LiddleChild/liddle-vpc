package model

import "time"

type Currency struct {
	ID        int
	Name      string
	RateToTHB float64
	UpdatedAt time.Time
}

func (c Currency) UpdateRate(rate float64, updatedAt time.Time) Currency {
	c.RateToTHB = rate
	c.UpdatedAt = updatedAt
	return c
}
