package handler

import (
	"context"
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/scheduler"
	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
)

var _ scheduler.Task = (*CurrencyTaskHandler)(nil)

type CurrencyTaskHandler struct {
	currencySync service.CurrencySync
}

func NewCurrencyTaskHandler(currencySync service.CurrencySync) CurrencyTaskHandler {
	return CurrencyTaskHandler{
		currencySync: currencySync,
	}
}

func (h CurrencyTaskHandler) Name() string {
	return "currency"
}

func (h CurrencyTaskHandler) Schedule() string {
	return "0 0,6,12,18 * * *"
}

func (h CurrencyTaskHandler) Run(ctx context.Context) error {
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return h.currencySync.Run(runCtx)
}
