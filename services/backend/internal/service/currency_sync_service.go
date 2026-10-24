package service

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/LiddleChild/liddle-vpc/backend/internal/frankfurterapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/gristapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/model"
	"github.com/LiddleChild/liddle-vpc/backend/internal/timeutil"
)

type CurrencyRepository interface {
	ListCurrencyRows(ctx context.Context, tableID string) ([]model.Currency, error)
	UpdateCurrencyRates(ctx context.Context, tableID string, updates []gristapi.RecordUpdate) error
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
	currencies, err := svc.repo.ListCurrencyRows(ctx, svc.tableID)
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

	updates := make([]gristapi.RecordUpdate, 0, len(currencies))
	updatedCurrencies := make(map[string]struct{})
	for _, currency := range currencies {
		code := strings.ToUpper(strings.TrimSpace(currency.Name))
		rate, ok := rates[code]
		if !ok {
			continue
		}

		updatedCurrencies[code] = struct{}{}

		fields := map[string]any{
			"Rate_To_THB": 1 / rate.Rate,
			"Updated_At":  timeutil.DateOnly(rate.UpdatedAt.In(timeutil.LocationBangkok)).Unix(),
		}

		updates = append(updates, gristapi.RecordUpdate{
			ID:     currency.ID,
			Fields: fields,
		})
	}

	if len(updates) > 0 {
		if err := svc.repo.UpdateCurrencyRates(ctx, svc.tableID, updates); err != nil {
			return fmt.Errorf("write currency rates: %w", err)
		}
	}

	slog.Info(
		"Currency rate sync completed",
		"table", svc.tableID,
		"updated_currencies", slices.Collect(maps.Keys(updatedCurrencies)),
	)

	return nil
}
