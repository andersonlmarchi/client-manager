package shared_test

import (
	"testing"

	"github.com/t-code/platform-core/packages/shared"
)

func TestModuleName(t *testing.T) {
	t.Parallel()
	if shared.ModuleName != "packages/shared" {
		t.Fatalf("ModuleName = %q, want packages/shared", shared.ModuleName)
	}
}

func TestVersionNonEmpty(t *testing.T) {
	t.Parallel()
	if shared.Version == "" {
		t.Fatal("Version must not be empty")
	}
}
