package domain_test

import (
	"testing"

	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
)

func TestHashAndVerifyPassword(t *testing.T) {
	t.Parallel()
	hash, err := domain.HashPassword("correct-horse-battery")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" || hash == "correct-horse-battery" {
		t.Fatal("hash must not equal plaintext")
	}
	ok, err := domain.VerifyPassword(hash, "correct-horse-battery")
	if err != nil || !ok {
		t.Fatalf("verify ok=%v err=%v", ok, err)
	}
	ok, err = domain.VerifyPassword(hash, "wrong-password-value")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected mismatch")
	}
}

func TestHashPasswordRejectsShort(t *testing.T) {
	t.Parallel()
	if _, err := domain.HashPassword("short"); err == nil {
		t.Fatal("expected error")
	}
}

func TestUserValidate(t *testing.T) {
	t.Parallel()
	u := domain.User{ID: "1", Email: "a@b.c", Status: domain.UserStatusActive}
	if err := u.Validate(); err != nil {
		t.Fatal(err)
	}
	u.Status = "weird"
	if err := u.Validate(); err == nil {
		t.Fatal("expected invalid status")
	}
}
