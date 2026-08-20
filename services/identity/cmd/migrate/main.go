package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/andersonlmarchi/client-manager/services/identity/internal/infrastructure"
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

	if err := infrastructure.EnsureSchema(ctx, db, infrastructure.DriverOf(databaseURL)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := infrastructure.Migrate(ctx, client); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
