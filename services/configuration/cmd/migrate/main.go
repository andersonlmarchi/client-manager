package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersonlmarchi/client-manager/services/configuration/internal/infrastructure"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "up" {
		fmt.Fprintln(os.Stderr, "usage: migrate up")
		os.Exit(2)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, db, err := infrastructure.Open(databaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer client.Close()
	defer db.Close()

	driver := "pgx"
	if strings.HasPrefix(databaseURL, "sqlite:") || strings.HasPrefix(databaseURL, "file:") {
		driver = "sqlite"
	}
	if err := infrastructure.EnsureSchema(ctx, db, driver); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := infrastructure.Migrate(ctx, client); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	repo := infrastructure.NewSettingsRepository(client)
	if err := repo.EnsureDefaults(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
