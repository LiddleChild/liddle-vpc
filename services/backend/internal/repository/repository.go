package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/gristapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
)

type Repository struct {
	documentID  string
	gristClient gristapi.Client
}

func New(documentID string, gristClient gristapi.Client) Repository {
	return Repository{
		documentID:  documentID,
		gristClient: gristClient,
	}
}

func (repo Repository) GetExpenseCatalog(ctx context.Context) (model.Catalog, error) {
	columns, err := repo.gristClient.ListTableColumns(ctx, repo.documentID, "Transactions")
	if err != nil {
		return model.Catalog{}, err
	}

	accounts, err := repo.gristClient.ListTableRecords(ctx, repo.documentID, "Accounts")
	if err != nil {
		return model.Catalog{}, err
	}

	catalog := model.Catalog{
		Categories: []model.Choice{},
		Tags:       []model.Choice{},
		Accounts:   make([]model.Account, 0, len(accounts)),
	}

	for _, column := range columns {
		switch column.ID {
		case "Category":
			catalog.Categories, err = getColumnChoices(column)
		case "Tags":
			catalog.Tags, err = getColumnChoices(column)
		}
		if err != nil {
			return model.Catalog{}, err
		}
	}

	if catalog.Categories == nil || catalog.Tags == nil {
		return model.Catalog{}, fmt.Errorf("transaction catalog columns are missing")
	}

	for _, record := range accounts {
		name, ok := record.Fields["Name"].(string)
		if !ok || strings.TrimSpace(name) == "" || record.ID < 1 {
			return model.Catalog{}, fmt.Errorf("account record has invalid name")
		}
		catalog.Accounts = append(catalog.Accounts, model.Account{ID: record.ID, Name: name})
	}

	return catalog, nil
}

type monthlySummaryRecord struct {
	AccountID        int     `json:"accountId"`
	Expenses         float64 `json:"expenses"`
	TransactionCount int     `json:"transactionCount"`
}

func (repo Repository) GetMonthlyExpenseSummary(ctx context.Context, filters model.ExpenseFilters) (model.MonthlySummary, error) {
	monthStart, err := time.Parse("2006-01", filters.Month)
	if err != nil || monthStart.Format("2006-01") != filters.Month {
		return model.MonthlySummary{}, fmt.Errorf("invalid month")
	}

	query := `
	select
		t."Account" as accountId,
		coalesce(sum(t."Amount_In_Account_Currency" * t."Inflow_Direction"), 0) as expenses,
		count(*) as transactionCount
	from "Transactions" t
	where t."Bill_Cycle" >= ?
	  and t."Bill_Cycle" < ?
	  and t."Account" in (select value from json_each(?))
	  and (
	    t."Category" = ''
	    or t."Category" in (select value from json_each(?))
	  )
	  and (
	    t."Tags" is null
	    or exists (
	      select 1
	      from json_each(t."Tags") transaction_tag
	      join json_each(?) selected_tag
	        on selected_tag.value = transaction_tag.value
	    )
	  )
	group by t."Account"
	`
	accountsJSON, _ := json.Marshal(filters.Accounts)
	categoriesJSON, _ := json.Marshal(filters.Categories)
	tagsJSON, _ := json.Marshal(filters.Tags)
	response, err := repo.gristClient.QueryWithParams[monthlySummaryRecord](ctx, gristapi.QueryWithParamsRequest{
		DocID: repo.documentID,
		SQL:   query,
		Args: []any{
			monthStart.Unix(),
			monthStart.AddDate(0, 1, 0).Unix(),
			string(accountsJSON),
			string(categoriesJSON),
			string(tagsJSON),
		},
	})
	if err != nil {
		return model.MonthlySummary{}, err
	}
	accountRecords, err := repo.gristClient.ListTableRecords(ctx, repo.documentID, "Accounts")
	if err != nil {
		return model.MonthlySummary{}, err
	}

	accountNames := make(map[int]string, len(accountRecords))
	for _, account := range accountRecords {
		if name, ok := account.Fields["Name"].(string); ok && strings.TrimSpace(name) != "" && account.ID > 0 {
			accountNames[account.ID] = name
		}
	}

	groupedExpenses := make(map[int]model.AccountExpense, len(response.Records))
	summary := model.MonthlySummary{
		Month:           filters.Month,
		AccountExpenses: make([]model.AccountExpense, 0, len(filters.Accounts)),
		SyncedAt:        time.Now().UTC(),
	}
	for _, record := range response.Records {
		row := record.Fields
		accountName, ok := accountNames[row.AccountID]
		if !ok {
			return model.MonthlySummary{}, fmt.Errorf("expense returned for unknown account")
		}
		expense := model.AccountExpense{
			AccountID:        row.AccountID,
			AccountName:      accountName,
			Expenses:         row.Expenses,
			TransactionCount: row.TransactionCount,
		}
		groupedExpenses[row.AccountID] = expense
		summary.Expenses += expense.Expenses
		summary.TransactionCount += expense.TransactionCount
	}
	summary.AccountCount = len(groupedExpenses)
	for _, accountID := range filters.Accounts {
		accountName, ok := accountNames[accountID]
		if !ok {
			return model.MonthlySummary{}, fmt.Errorf("selected account does not exist")
		}
		if expense, ok := groupedExpenses[accountID]; ok {
			summary.AccountExpenses = append(summary.AccountExpenses, expense)
		} else {
			summary.AccountExpenses = append(summary.AccountExpenses, model.AccountExpense{
				AccountID:   accountID,
				AccountName: accountName,
			})
		}
	}

	return summary, nil
}

func getColumnChoices(column gristapi.TableColumn) ([]model.Choice, error) {
	widgetOptions, ok := column.Fields["widgetOptions"].(string)
	if !ok {
		return nil, fmt.Errorf("column %s has no choice options", column.ID)
	}
	var options struct {
		Choices       []string `json:"choices"`
		ChoiceOptions map[string]struct {
			TextColor string `json:"textColor"`
			FillColor string `json:"fillColor"`
		} `json:"choiceOptions"`
	}
	if err := json.Unmarshal([]byte(widgetOptions), &options); err != nil {
		return nil, fmt.Errorf("column %s has invalid choice options", column.ID)
	}
	if options.Choices == nil {
		return nil, fmt.Errorf("column %s has no choices", column.ID)
	}
	choices := make([]model.Choice, 0, len(options.Choices))
	for _, value := range options.Choices {
		option := options.ChoiceOptions[value]
		choices = append(choices, model.Choice{Value: value, TextColor: option.TextColor, FillColor: option.FillColor})
	}
	return choices, nil
}
