package application_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/application"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/infrastructure"
)

func TestAuthServiceLoginLogout(t *testing.T) {
	t.Parallel()
	name := strings.ReplaceAll(t.Name(), "/", "_")
	client, db, err := infrastructure.Open("sqlite:file:" + name + "?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.Close()
		_ = db.Close()
	})
	ctx := context.Background()
	if err := infrastructure.Migrate(ctx, client); err != nil {
		t.Fatal(err)
	}
	users := infrastructure.NewUserRepository(client)
	sessions := infrastructure.NewSessionRepository(client)
	auth := application.NewAuthService(users, sessions, time.Hour)

	if _, err := users.CreateUserWithPassword(ctx, "auth@example.com", "super-secret-pass"); err != nil {
		t.Fatal(err)
	}
	res, err := auth.Login(ctx, "auth@example.com", "super-secret-pass", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	user, _, err := auth.CurrentUser(ctx, res.SessionToken)
	if err != nil || user.Email != "auth@example.com" {
		t.Fatalf("user=%+v err=%v", user, err)
	}
	if err := auth.Logout(ctx, res.SessionToken); err != nil {
		t.Fatal(err)
	}
	_, _, err = auth.CurrentUser(ctx, res.SessionToken)
	if e, ok := shared.AsError(err); !ok || e.Code != shared.CodeUnauthorized {
		t.Fatalf("err=%v", err)
	}
}
