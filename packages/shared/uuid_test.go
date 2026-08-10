package shared_test

import (
	"strings"
	"testing"

	"github.com/t-code/client-manager/packages/shared"
)

func TestNewUUIDAndParseUUID(t *testing.T) {
	t.Parallel()
	id, err := shared.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 36 {
		t.Fatalf("len = %d, want 36", len(id))
	}
	parsed, err := shared.ParseUUID(id)
	if err != nil {
		t.Fatal(err)
	}
	if parsed != strings.ToLower(id) {
		t.Fatalf("parsed = %q, want %q", parsed, id)
	}
}

func TestParseUUIDRejectsInvalid(t *testing.T) {
	t.Parallel()
	cases := []string{"", "not-a-uuid", "12345678-1234-1234-1234-1234567890zz"}
	for _, in := range cases {
		if _, err := shared.ParseUUID(in); err == nil {
			t.Fatalf("ParseUUID(%q) expected error", in)
		}
	}
}

func TestNewUUIDUniqueness(t *testing.T) {
	t.Parallel()
	a, err := shared.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := shared.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("expected distinct UUIDs")
	}
}
