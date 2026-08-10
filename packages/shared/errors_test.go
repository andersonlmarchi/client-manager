package shared_test

import (
	"errors"
	"testing"

	"github.com/andersonlmarchi/client-manager/packages/shared"
)

func TestNewErrorAndAsError(t *testing.T) {
	t.Parallel()
	err := shared.NewError(shared.CodeNotFound, "missing")
	got, ok := shared.AsError(err)
	if !ok {
		t.Fatal("expected AsError to succeed")
	}
	if got.Code != shared.CodeNotFound || got.Message != "missing" {
		t.Fatalf("unexpected error: %+v", got)
	}
}

func TestWrapUnwrap(t *testing.T) {
	t.Parallel()
	root := errors.New("root")
	err := shared.Wrap(shared.CodeInvalid, "bad input", root)
	if !errors.Is(err, root) {
		t.Fatal("expected errors.Is to find root")
	}
	if got := err.Error(); got == "" {
		t.Fatal("expected non-empty Error()")
	}
}
