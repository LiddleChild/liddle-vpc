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

type CurrencySync struct {
	tableID     string
	repo        CurrencyRepository
	frankfurter frankfurterapi.Client
}

func NewCurrencySync(
	tableID string,
	repo CurrencyRepository,
	frankfurter frankfurterapi.Client,
) CurrencySync {
	return CurrencySync{
		tableID:     tableID,
		repo:        repo,
		frankfurter: frankfurter,
	}
}

func (service CurrencySync) Run(ctx context.Context) error {
	currencies, err := service.repo.ListCurrencyRows(ctx, service.tableID)
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

	rates, err := service.frankfurter.RatesToTHB(ctx, codes)
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
			"Rate To THB": rate.ToTHB,
			"Updated At":  timeutil.DateOnly(rate.Date.In(timeutil.LocationBangkok)).Unix(),
		}

		updates = append(updates, gristapi.RecordUpdate{
			ID:     currency.ID,
			Fields: fields,
		})
	}

	if len(updates) > 0 {
		if err := service.repo.UpdateCurrencyRates(ctx, service.tableID, updates); err != nil {
			return fmt.Errorf("write currency rates: %w", err)
		}
	}

	slog.Info(
		"Currency rate sync completed",
		"table", service.tableID,
		"updated_currencies", slices.Collect(maps.Keys(updatedCurrencies)),
	)

	return nil
}
