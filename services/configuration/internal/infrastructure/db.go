package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

const SchemaName = "configuration"

func Open(databaseURL string) (*ent.Client, *sql.DB, error) {
	driver, dsn, err := parseDatabaseURL(databaseURL)
	if err != nil {
		return nil, nil, err
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}
	var entDriver dialect.Driver
	switch driver {
	case "pgx":
		entDriver = entsql.OpenDB(dialect.Postgres, db)
	case "sqlite":
		entDriver = entsql.OpenDB(dialect.SQLite, db)
	default:
		_ = db.Close()
		return nil, nil, fmt.Errorf("unsupported driver %q", driver)
	}
	return ent.NewClient(ent.Driver(entDriver)), db, nil
}

func parseDatabaseURL(databaseURL string) (driver string, dsn string, err error) {
	if databaseURL == "" {
		return "", "", fmt.Errorf("DATABASE_URL is required")
	}
	if strings.HasPrefix(databaseURL, "sqlite:") || strings.HasPrefix(databaseURL, "file:") {
		dsn = strings.TrimPrefix(databaseURL, "sqlite:")
		return "sqlite", dsn, nil
	}
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", "", fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	switch u.Scheme {
	case "postgres", "postgresql":
		return "pgx", databaseURL, nil
	case "sqlite":
		return "sqlite", strings.TrimPrefix(databaseURL, "sqlite:"), nil
	default:
		return "", "", fmt.Errorf("unsupported DATABASE_URL scheme %q", u.Scheme)
	}
}

func EnsureSchema(ctx context.Context, db *sql.DB, driver string) error {
	if driver != "pgx" {
		return nil
	}
	_, err := db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS `+SchemaName)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	_, err = db.ExecContext(ctx, `SET search_path TO `+SchemaName)
	if err != nil {
		return fmt.Errorf("set search_path: %w", err)
	}
	return nil
}

func Migrate(ctx context.Context, client *ent.Client) error {
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}
	return nil
}
