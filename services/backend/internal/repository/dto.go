package repository

type accountFields struct {
	Name string `json:"Name"`
}

type transactionColumnFields struct {
	WidgetOptions string `json:"widgetOptions"`
}

type choiceOptions struct {
	Choices       []string                    `json:"choices"`
	ChoiceOptions map[string]choiceAppearance `json:"choiceOptions"`
}

type choiceAppearance struct {
	TextColor string `json:"textColor"`
	FillColor string `json:"fillColor"`
}

type currencyFields struct {
	Name      string  `json:"currency"`
	RateToTHB float64 `json:"rate_to_thb"`
	UpdatedAt int64   `json:"updated_at"`
}

type monthlySummaryRecord struct {
	AccountID        int     `json:"account_id"`
	Expenses         float64 `json:"expenses"`
	TransactionCount int     `json:"transaction_count"`
}
