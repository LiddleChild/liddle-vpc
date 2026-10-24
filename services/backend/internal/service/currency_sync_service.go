package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LiddleChild/liddle-vpc/backend/internal/frankfurterapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
	"github.com/LiddleChild/liddle-vpc/backend/internal/timeutil"
)

type CurrencyRepository interface {
	ListCurrencies(ctx context.Context, tableID string) ([]model.Currency, error)
	UpdateCurrencies(ctx context.Context, tableID string, currencies []model.Currency) error
}

type CurrencySyncService struct {
	tableID     string
	repo        CurrencyRepository
	frankfurter frankfurterapi.Client
}

func NewCurrencySyncService(
	tableID string,
	repo CurrencyRepository,
	frankfurter frankfurterapi.Client,
) CurrencySyncService {
	return CurrencySyncService{
		tableID:     tableID,
		repo:        repo,
		frankfurter: frankfurter,
	}
}

func (svc CurrencySyncService) Run(ctx context.Context) error {
	currencies, err := svc.repo.ListCurrencies(ctx, svc.tableID)
	if err != nil {
		return fmt.Errorf("read currency table: %w", err)
	}

	codes := make([]string, 0, len(currencies))
	for _, currency := range currencies {
		code := strings.ToUpper(strings.TrimSpace(currency.Name))
		if code == "" {
			continue
		}

		codes = append(codes, code)
	}

	rates, err := svc.frankfurter.RatesToTHB(ctx, codes)
	if err != nil {
		return fmt.Errorf("fetch currency rates: %w", err)
	}

	updatedCurrencies := make([]model.Currency, 0, len(currencies))
	for _, currency := range currencies {
		code := strings.ToUpper(strings.TrimSpace(currency.Name))
		rate, ok := rates[code]
		if !ok {
			continue
		}

		updatedCurrencies = append(updatedCurrencies, currency.UpdateRate(
			1/rate.Rate,
			timeutil.DateOnly(rate.UpdatedAt),
		))
	}

	if len(updatedCurrencies) > 0 {
		if err := svc.repo.UpdateCurrencies(ctx, svc.tableID, updatedCurrencies); err != nil {
			return fmt.Errorf("write currency rates: %w", err)
		}
	}

	updatedCodes := make([]string, 0, len(updatedCurrencies))
	for _, currency := range updatedCurrencies {
		updatedCodes = append(updatedCodes, currency.Name)
	}

	slog.Info(
		"currency rate sync completed",
		"currencies", updatedCodes,
	)

	return nil
}
