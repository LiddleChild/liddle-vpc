package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LiddleChild/liddle-vpc/backend/internal/gristapi"
	"github.com/LiddleChild/liddle-vpc/backend/internal/handler"
	"github.com/LiddleChild/liddle-vpc/backend/internal/repository"
	"github.com/LiddleChild/liddle-vpc/backend/internal/service"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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
		client         = gristapi.NewClient(cfg.GristEndpoint, cfg.GristAPIKey)
		repo           = repository.New(cfg.GristDocumentID, client)
		expenseSvc     = service.NewExpenses(repo)
		expenseHandler = handler.NewHandler(expenseSvc)
	)

	mux := chi.NewMux()
	mux.Use(LoggerMiddleware)
	mux.Use(chimiddleware.Recoverer)

	mux.Get("/expenses/v1/catalog", expenseHandler.GetExpenseCatalog)
	mux.Get("/expenses/v1/monthly-summary", expenseHandler.GetMonthlyExpenseSummary)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		slog.Info(fmt.Sprintf("Listening on :%d", cfg.Port))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		close(stopCh)
	}()

	<-stopCh
	slog.Info("Shutting down")

	if err := server.Close(); err != nil {
		return err
	}

	slog.Info("Gracefully shutdown")
	return nil
}
