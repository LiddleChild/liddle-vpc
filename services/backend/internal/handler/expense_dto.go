package handler

import (
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
)

type choiceResponse struct {
	Value     string `json:"value"`
	TextColor string `json:"textColor,omitempty"`
	FillColor string `json:"fillColor,omitempty"`
}

type accountResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type catalogResponse struct {
	Categories []choiceResponse  `json:"categories"`
	Tags       []choiceResponse  `json:"tags"`
	Accounts   []accountResponse `json:"accounts"`
}

type accountExpenseResponse struct {
	AccountID        int     `json:"accountId"`
	AccountName      string  `json:"accountName"`
	Expenses         float64 `json:"expenses"`
	TransactionCount int     `json:"transactionCount"`
}

type monthlySummaryResponse struct {
	Month            string                   `json:"month"`
	Expenses         float64                  `json:"expenses"`
	TransactionCount int                      `json:"transactionCount"`
	AccountCount     int                      `json:"accountCount"`
	AccountExpenses  []accountExpenseResponse `json:"accountExpenses"`
	SyncedAt         time.Time                `json:"syncedAt"`
}

func newCatalogResponse(catalog model.Catalog) catalogResponse {
	categories := make([]choiceResponse, 0, len(catalog.Categories))

	for _, choice := range catalog.Categories {
		categories = append(categories, choiceResponse{Value: choice.Value, TextColor: choice.TextColor, FillColor: choice.FillColor})
	}

	tags := make([]choiceResponse, 0, len(catalog.Tags))
	for _, choice := range catalog.Tags {
		tags = append(tags, choiceResponse{Value: choice.Value, TextColor: choice.TextColor, FillColor: choice.FillColor})
	}

	accounts := make([]accountResponse, 0, len(catalog.Accounts))
	for _, account := range catalog.Accounts {
		accounts = append(accounts, accountResponse{ID: account.ID, Name: account.Name})
	}

	return catalogResponse{
		Categories: categories,
		Tags:       tags,
		Accounts:   accounts,
	}
}

func newMonthlySummaryResponse(summary model.MonthlySummary) monthlySummaryResponse {
	accountExpenses := make([]accountExpenseResponse, 0, len(summary.AccountExpenses))

	for _, expense := range summary.AccountExpenses {
		accountExpenses = append(accountExpenses, accountExpenseResponse{
			AccountID:        expense.AccountID,
			AccountName:      expense.AccountName,
			Expenses:         expense.Expenses,
			TransactionCount: expense.TransactionCount,
		})
	}

	return monthlySummaryResponse{
		Month:            summary.Month,
		Expenses:         summary.Expenses,
		TransactionCount: summary.TransactionCount,
		AccountCount:     summary.AccountCount,
		AccountExpenses:  accountExpenses,
		SyncedAt:         summary.SyncedAt,
	}
}
