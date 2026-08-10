package shared_test

import (
	"log/slog"
	"testing"

	"github.com/t-code/client-manager/packages/shared"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()
	cases := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"WARN":    slog.LevelWarn,
		"error":   slog.LevelError,
		"info":    slog.LevelInfo,
		"unknown": slog.LevelInfo,
	}
	for in, want := range cases {
		if got := shared.ParseLogLevel(in); got != want {
			t.Fatalf("ParseLogLevel(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()
	logger := shared.NewLogger(slog.LevelInfo)
	if logger == nil {
		t.Fatal("expected logger")
	}
}
