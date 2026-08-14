package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/application"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/infrastructure"
	confighttp "github.com/andersonlmarchi/client-manager/services/configuration/internal/transport/http"
)

func main() {
	logger := shared.NewLogger(shared.ParseLogLevel(os.Getenv("LOG_LEVEL")))
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	adminKey := os.Getenv("ADMIN_API_KEY")
	if adminKey == "" {
		logger.Error("ADMIN_API_KEY is required")
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	client, db, err := infrastructure.Open(databaseURL)
	if err != nil {
		logger.Error("database open failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer client.Close()
	defer db.Close()

	repo := infrastructure.NewSettingsRepository(client)
	service := application.NewSettingsService(repo)
	serverAPI := confighttp.NewServer(service, adminKey)

	mux := http.NewServeMux()
	serverAPI.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("configuration listening", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
