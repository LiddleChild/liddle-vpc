package model

import "time"

type Choice struct {
	Value     string `json:"value"`
	TextColor string `json:"textColor,omitempty"`
	FillColor string `json:"fillColor,omitempty"`
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Catalog struct {
	Categories []Choice  `json:"categories"`
	Tags       []Choice  `json:"tags"`
	Accounts   []Account `json:"accounts"`
}

type ExpenseFilters struct {
	Month      string   `json:"month"`
	Accounts   []int    `json:"accounts"`
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
}

type AccountExpense struct {
	AccountID        int     `json:"accountId"`
	AccountName      string  `json:"accountName"`
	Expenses         float64 `json:"expenses"`
	TransactionCount int     `json:"transactionCount"`
}

type MonthlySummary struct {
	Month            string           `json:"month"`
	Expenses         float64          `json:"expenses"`
	TransactionCount int              `json:"transactionCount"`
	AccountCount     int              `json:"accountCount"`
	AccountExpenses  []AccountExpense `json:"accountExpenses"`
	SyncedAt         time.Time        `json:"syncedAt"`
}
