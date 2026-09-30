// Command server runs the calculator REST API.
//
// Environment variables:
//
//	PORT        port to listen on (default 8080)
//	STATIC_DIR  optional directory with the built frontend; when set, it is
//	            served at "/" so the whole app runs from a single container.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calculator/internal/api"
)

func main() {
	port := envOr("PORT", "8080")
	mux := api.NewRouter()

	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		mux.Handle("GET /", http.FileServer(http.Dir(dir)))
		slog.Info("serving frontend", "dir", dir)
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           api.Logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	go func() {
		slog.Info("server listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on Ctrl+C / docker stop.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
	slog.Info("server stopped")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
