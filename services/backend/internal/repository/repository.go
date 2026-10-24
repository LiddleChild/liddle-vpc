package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

func (repo Repository) ListCurrencyRows(ctx context.Context, tableID string) ([]model.Currency, error) {
	records, err := repo.gristClient.ListTableRecords(ctx, repo.documentID, tableID)
	if err != nil {
		return nil, err
	}

	rows := make([]model.Currency, 0, len(records))
	for _, record := range records {
		if record.ID < 1 {
			return nil, fmt.Errorf("currency table contains an invalid row ID")
		}
		name, ok := record.Fields["Currency"].(string)
		if !ok {
			name = ""
		}
		rateToTHB, err := gristNumber(record.Fields["Rate To THB"])
		if err != nil {
			return nil, fmt.Errorf("currency row %d has an invalid rate", record.ID)
		}
		updatedAtSeconds, err := gristNumber(record.Fields["Updated At"])
		if err != nil {
			return nil, fmt.Errorf("currency row %d has an invalid updated date", record.ID)
		}
		updatedAt := time.Time{}
		if updatedAtSeconds != 0 {
			seconds := int64(updatedAtSeconds)
			nanoseconds := int64((updatedAtSeconds - float64(seconds)) * float64(time.Second))
			updatedAt = time.Unix(seconds, nanoseconds).UTC()
		}

		rows = append(rows, model.Currency{
			ID:        record.ID,
			Name:      name,
			RateToTHB: rateToTHB,
			UpdatedAt: updatedAt,
		})
	}

	return rows, nil
}

func gristNumber(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}
	var number float64
	switch value := value.(type) {
	case float64:
		number = value
	case int:
		number = float64(value)
	case int64:
		number = float64(value)
	case json.Number:
		parsed, err := value.Float64()
		if err != nil {
			return 0, err
		}
		number = parsed
	default:
		return 0, fmt.Errorf("expected number, got %T", value)
	}
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, fmt.Errorf("number must be finite")
	}
	return number, nil
}

func (repo Repository) UpdateCurrencyRates(ctx context.Context, tableID string, updates []gristapi.RecordUpdate) error {
	return repo.gristClient.UpdateTableRecords(ctx, repo.documentID, tableID, updates)
}
