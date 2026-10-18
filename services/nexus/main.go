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

	"github.com/LiddleChild/liddle-vpc/gristguard/grist"
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
		client = grist.NewClient(cfg.GristEndpoint, cfg.GristAPIKey)
		repo   = NewRepository(cfg.GristDocumentID, client)
		proxy  = NewProxy(cfg.GristEndpoint, repo, cfg.GristAPIKey)
	)

	mux := chi.NewMux()
	mux.Use(LoggerMiddleware)
	mux.Use(chimiddleware.Recoverer)

	endpoints, err := repo.ListAllEndpoints(context.Background())
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		mux.Handle(endpoint, http.HandlerFunc(proxy.Handle))
	}

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
	}()

	<-stopCh
	slog.Info("Shutting down")

	if err := server.Close(); err != nil {
		return err
	}

	slog.Info("Gracefully shutdown")
	return nil
}
