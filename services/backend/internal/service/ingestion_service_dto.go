package service

import "time"

type IngestionRequest struct {
	Transactions []IngestionItem
}

type IngestionItem struct {
	Version      string
	Notification string
	Transaction  IngestedTransaction
}

type IngestedTransaction struct {
	App       string
	Type      string
	Date      string
	CreatedAt time.Time
	Amount    string
	Currency  string
}
