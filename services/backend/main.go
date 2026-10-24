package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LiddleChild/liddle-vpc/backend/internal/frankfurterapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/gristapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/handler"
	"github.com/LiddleChild/liddle-vpc/backend/internal/repository"
	"github.com/LiddleChild/liddle-vpc/backend/internal/scheduler"
	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		slog.Error(fmt.Sprintf("error: %s", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return err
	}

	var (
		gristClient        = gristapi.NewClient(cfg.GristEndpoint, cfg.GristAPIKey)
		frankfurtherClient = frankfurterapi.NewClient(cfg.FrankfurterEndpoint)

		repo = repository.New(cfg.GristDocumentID, gristClient)

		expenseSvc      = service.NewExpenseService(repo)
		currencySyncSvc = service.NewCurrencySyncService(cfg.CurrencyTableID, repo, frankfurtherClient)

		expenseHandler = handler.NewExpenseHandler(expenseSvc)
		taskHandler    = handler.NewTaskHandler(currencySyncSvc)

		currencySyncTaskHandler = handler.NewCurrencySyncTaskHandler(currencySyncSvc)
	)

	mux := chi.NewMux()
	mux.Use(LoggerMiddleware)
	mux.Use(chimiddleware.Recoverer)

	mux.Get("/expenses/v1/catalog", expenseHandler.GetExpenseCatalog)
	mux.Get("/expenses/v1/monthly-summary", expenseHandler.GetMonthlyExpenseSummary)
	mux.Post("/tasks/v1/currency-sync", taskHandler.TriggerCurrencyTask)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	taskScheduler := scheduler.New()
	taskScheduler.ScheduleTask(currencySyncTaskHandler)

	notifyCtx, notifyCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer notifyCancel()

	runtimeGroup, runtimeCtx := errgroup.WithContext(notifyCtx)

	runtimeGroup.Go(func() error {
		slog.Info(fmt.Sprintf("listening on :%d", cfg.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	<-runtimeCtx.Done()

	slog.Info("shutting down")

	var shutdownErrors []error

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	slog.Info("stopping task scheduler")
	taskScheduler.Stop(shutdownCtx)

	shutdownCtx, shutdownCancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	slog.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		shutdownErrors = append(shutdownErrors, err)
		slog.Error("error shutting down http server",
			"error", err.Error(),
		)
	}

	slog.Info("gracefully shutdown")

	errs := append(
		[]error{runtimeGroup.Wait()},
		shutdownErrors...,
	)

	return errors.Join(errs...)
}
