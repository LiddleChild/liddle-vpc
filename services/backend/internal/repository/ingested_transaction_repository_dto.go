package repository

type ingestedTransactionFields struct {
	App          string  `json:"app"`
	Type         string  `json:"type"`
	Date         string  `json:"date"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Notification string  `json:"notification"`
}
