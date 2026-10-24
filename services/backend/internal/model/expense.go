package model

import "time"

type Choice struct {
	Value     string
	TextColor string
	FillColor string
}

type Account struct {
	ID   int
	Name string
}

type Catalog struct {
	Categories []Choice
	Tags       []Choice
	Accounts   []Account
}

type ExpenseFilters struct {
	Month      string
	Accounts   []int
	Categories []string
	Tags       []string
}

type AccountExpense struct {
	AccountID        int
	AccountName      string
	Expenses         float64
	TransactionCount int
}

type MonthlySummary struct {
	Month            string
	Expenses         float64
	TransactionCount int
	AccountCount     int
	AccountExpenses  []AccountExpense
	SyncedAt         time.Time
}
