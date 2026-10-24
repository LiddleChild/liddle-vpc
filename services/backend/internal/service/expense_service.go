package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
)

var ErrInvalidExpenseRequest = errors.New("invalid expense request")

type Repository interface {
	GetExpenseCatalog(ctx context.Context) (model.Catalog, error)
	GetMonthlyExpenseSummary(ctx context.Context, filters model.ExpenseFilters) (model.MonthlySummary, error)
}

type ExpenseService struct {
	repo Repository
}

func NewExpenseService(repo Repository) ExpenseService {
	return ExpenseService{
		repo: repo,
	}
}

func (svc ExpenseService) GetExpenseCatalog(ctx context.Context) (model.Catalog, error) {
	return svc.repo.GetExpenseCatalog(ctx)
}

func (svc ExpenseService) GetMonthlyExpenseSummary(ctx context.Context, filters model.ExpenseFilters) (model.MonthlySummary, error) {
	monthStart, err := time.Parse("2006-01", filters.Month)
	if err != nil || monthStart.Format("2006-01") != filters.Month {
		return model.MonthlySummary{}, fmt.Errorf("%w: month must use YYYY-MM format", ErrInvalidExpenseRequest)
	}

	if filters.Accounts == nil || filters.Categories == nil || filters.Tags == nil {
		return model.MonthlySummary{}, fmt.Errorf("%w: accounts, categories, and tags are required", ErrInvalidExpenseRequest)
	}

	for _, accountID := range filters.Accounts {
		if accountID < 1 {
			return model.MonthlySummary{}, fmt.Errorf("%w: account IDs must be positive", ErrInvalidExpenseRequest)
		}
	}

	return svc.repo.GetMonthlyExpenseSummary(ctx, filters)
}
