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
	columns, err := repo.gristClient.ListTableColumns[transactionColumnFields](ctx, repo.documentID, "Transactions")
	if err != nil {
		return model.Catalog{}, err
	}

	accounts, err := repo.gristClient.ListTableRecords[accountFields](ctx, repo.documentID, "Accounts")
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
		if strings.TrimSpace(record.Fields.Name) == "" || record.ID < 1 {
			return model.Catalog{}, fmt.Errorf("account record has invalid name")
		}
		catalog.Accounts = append(catalog.Accounts, model.Account{ID: record.ID, Name: record.Fields.Name})
	}

	return catalog, nil
}

func (repo Repository) GetMonthlyExpenseSummary(ctx context.Context, filters model.ExpenseFilters) (model.MonthlySummary, error) {
	monthStart, err := time.Parse("2006-01", filters.Month)
	if err != nil || monthStart.Format("2006-01") != filters.Month {
		return model.MonthlySummary{}, fmt.Errorf("invalid month")
	}

	query := `
	select
		t."Account" as account_id,
		coalesce(sum(t."Amount_In_Account_Currency" * t."Inflow_Direction"), 0) as expenses,
		count(*) as transaction_count
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
	accountRecords, err := repo.gristClient.ListTableRecords[accountFields](ctx, repo.documentID, "Accounts")
	if err != nil {
		return model.MonthlySummary{}, err
	}

	accountNames := make(map[int]string, len(accountRecords))
	for _, account := range accountRecords {
		if strings.TrimSpace(account.Fields.Name) != "" && account.ID > 0 {
			accountNames[account.ID] = account.Fields.Name
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

func getColumnChoices(column gristapi.TableColumn[transactionColumnFields]) ([]model.Choice, error) {
	widgetOptions := column.Fields.WidgetOptions
	if widgetOptions == "" {
		return nil, fmt.Errorf("column %s has no choice options", column.ID)
	}

	var options choiceOptions

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

func (repo Repository) ListCurrencies(ctx context.Context, tableID string) ([]model.Currency, error) {
	records, err := repo.gristClient.ListTableRecords[currencyFields](ctx, repo.documentID, tableID)
	if err != nil {
		return nil, err
	}

	rows := make([]model.Currency, 0, len(records))
	for _, record := range records {
		updatedAt := time.Time{}
		if record.Fields.UpdatedAt > 0 {
			updatedAt = time.Unix(record.Fields.UpdatedAt, 0)
		}

		rows = append(rows, model.Currency{
			ID:        record.ID,
			Name:      record.Fields.Name,
			RateToTHB: record.Fields.RateToTHB,
			UpdatedAt: updatedAt,
		})
	}

	return rows, nil
}

func (repo Repository) UpdateCurrencies(ctx context.Context, tableID string, currencies []model.Currency) error {
	updates := make([]gristapi.RecordUpdate[currencyFields], 0, len(currencies))
	for _, currency := range currencies {
		updates = append(updates, gristapi.RecordUpdate[currencyFields]{
			ID: currency.ID,
			Fields: currencyFields{
				Name:      currency.Name,
				RateToTHB: currency.RateToTHB,
				UpdatedAt: currency.UpdatedAt.UTC().Unix(),
			},
		})
	}

	return repo.gristClient.UpdateTableRecords(ctx, repo.documentID, tableID, updates)
}
