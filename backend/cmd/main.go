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

	"qrgen/internal/bootstrap"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.InitApp(ctx)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			app.AppLogger.Error("close Redis client", "error", err)
		}
	}()

	server := &http.Server{
		Addr:              app.Config.Server.Port,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       time.Duration(app.Config.Server.Timeout) * time.Second,
		WriteTimeout:      time.Duration(app.Config.Server.Timeout) * time.Second,
		IdleTimeout:       time.Duration(app.Config.Server.IdleTimeout) * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		app.AppLogger.Info("HTTP server started", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server: %w", err)

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("graceful shutdown: %w", err)
		}

		return nil
	}
}
