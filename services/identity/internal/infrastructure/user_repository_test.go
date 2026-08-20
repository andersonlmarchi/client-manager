package infrastructure_test

import (
	"context"
	"strings"
	"testing"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/infrastructure"
)

func openRepo(t *testing.T) *infrastructure.UserRepository {
	t.Helper()
	name := strings.ReplaceAll(t.Name(), "/", "_")
	client, db, err := infrastructure.Open("sqlite:file:" + name + "?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.Close()
		_ = db.Close()
	})
	if err := infrastructure.Migrate(context.Background(), client); err != nil {
		t.Fatal(err)
	}
	return infrastructure.NewUserRepository(client)
}

func TestCreateUserWithPasswordAndAuthenticate(t *testing.T) {
	t.Parallel()
	repo := openRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUserWithPassword(ctx, "User@Example.com", "super-secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "user@example.com" || user.Status != domain.UserStatusActive {
		t.Fatalf("user=%+v", user)
	}

	cred, err := repo.GetCredentialByUserID(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cred.Algorithm != domain.PasswordAlgorithm || cred.PasswordHash == "super-secret-pass" {
		t.Fatalf("credential=%+v", cred)
	}

	authed, err := repo.Authenticate(ctx, "user@example.com", "super-secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if authed.ID != user.ID {
		t.Fatalf("authed=%+v", authed)
	}

	_, err = repo.Authenticate(ctx, "user@example.com", "bad-password-xx")
	if err == nil {
		t.Fatal("expected unauthorized")
	}
	if e, ok := shared.AsError(err); !ok || e.Code != shared.CodeUnauthorized {
		t.Fatalf("err=%v", err)
	}
}

func TestCreateUserConflict(t *testing.T) {
	t.Parallel()
	repo := openRepo(t)
	ctx := context.Background()
	if _, err := repo.CreateUserWithPassword(ctx, "dup@example.com", "super-secret-pass"); err != nil {
		t.Fatal(err)
	}
	_, err := repo.CreateUserWithPassword(ctx, "dup@example.com", "super-secret-pass")
	if err == nil {
		t.Fatal("expected conflict")
	}
	if e, ok := shared.AsError(err); !ok || e.Code != shared.CodeConflict {
		t.Fatalf("err=%v", err)
	}
}
