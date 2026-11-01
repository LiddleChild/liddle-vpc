package service

type IngestionRequest struct {
	Transactions []IngestionItem
}

type IngestionItem struct {
	Version      string
	Notification string
	Transaction  IngestedTransaction
}

type IngestedTransaction struct {
	App      string
	Type     string
	Date     string
	Amount   string
	Currency string
}
