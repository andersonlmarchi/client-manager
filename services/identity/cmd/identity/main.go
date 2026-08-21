package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/application"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/infrastructure"
	identityhttp "github.com/andersonlmarchi/client-manager/services/identity/internal/transport/http"
)

func main() {
	logger := shared.NewLogger(shared.ParseLogLevel(os.Getenv("LOG_LEVEL")))
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	client, db, err := infrastructure.Open(databaseURL)
	if err != nil {
		logger.Error("database open failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer client.Close()
	defer db.Close()

	ttlHours := 168
	if v := os.Getenv("SESSION_TTL_HOURS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			logger.Error("invalid SESSION_TTL_HOURS")
			os.Exit(1)
		}
		ttlHours = n
	}
	cookieSecure := os.Getenv("COOKIE_SECURE") == "true" || os.Getenv("COOKIE_SECURE") == "1"

	users := infrastructure.NewUserRepository(client)
	sessions := infrastructure.NewSessionRepository(client)
	auth := application.NewAuthService(users, sessions, time.Duration(ttlHours)*time.Hour)
	api := identityhttp.NewServer(auth, cookieSecure)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("identity listening", slog.String("addr", addr))
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
