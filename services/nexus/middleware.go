package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(wrappedWriter, r)

		slog.Info(
			fmt.Sprintf("%s %s", r.Method, r.URL.String()),
			slog.Int("status", wrappedWriter.Status()),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
