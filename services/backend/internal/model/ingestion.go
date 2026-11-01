package model

import "time"

type IngestedTransaction struct {
	App          string
	Type         string
	Date         string
	Amount       float64
	Currency     string
	Notification string
	CreatedAt    time.Time
}
