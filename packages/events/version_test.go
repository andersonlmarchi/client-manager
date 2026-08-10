package events_test

import (
	"testing"

	"github.com/andersonlmarchi/client-manager/packages/events"
)

func TestModuleName(t *testing.T) {
	t.Parallel()
	if events.ModuleName != "packages/events" {
		t.Fatalf("ModuleName = %q, want packages/events", events.ModuleName)
	}
}

func TestVersionNonEmpty(t *testing.T) {
	t.Parallel()
	if events.Version == "" {
		t.Fatal("Version must not be empty")
	}
}
