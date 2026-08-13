package infrastructure_test

import (
	"context"
	"testing"

	"github.com/andersonlmarchi/client-manager/services/configuration/internal/domain"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/infrastructure"
)

func openTestRepo(t *testing.T) (*infrastructure.SettingsRepository, func()) {
	t.Helper()
	client, db, err := infrastructure.Open("sqlite:file:configuration?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := infrastructure.Migrate(ctx, client); err != nil {
		t.Fatal(err)
	}
	repo := infrastructure.NewSettingsRepository(client)
	cleanup := func() {
		_ = client.Close()
		_ = db.Close()
	}
	return repo, cleanup
}

func TestSettingsRepositoryRoundTrip(t *testing.T) {
	repo, cleanup := openTestRepo(t)
	defer cleanup()
	ctx := context.Background()

	if err := repo.EnsureDefaults(ctx); err != nil {
		t.Fatal(err)
	}

	pw, err := repo.GetPasswordPolicy(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	if pw.MinLength != 12 || !pw.RequireSymbol {
		t.Fatalf("password policy = %+v", pw)
	}

	pw.MinLength = 16
	if err := repo.SavePasswordPolicy(ctx, pw); err != nil {
		t.Fatal(err)
	}
	pw2, err := repo.GetPasswordPolicy(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	if pw2.MinLength != 16 {
		t.Fatalf("min_length = %d", pw2.MinLength)
	}

	oidc, err := repo.GetOIDCSettings(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	oidc.Issuer = "https://auth.example.com"
	if err := repo.SaveOIDCSettings(ctx, oidc); err != nil {
		t.Fatal(err)
	}

	rl, err := repo.GetRateLimitSettings(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	rl.LoginRPS = 3
	if err := repo.SaveRateLimitSettings(ctx, rl); err != nil {
		t.Fatal(err)
	}

	branding, err := repo.GetBrandingSettings(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	color := "#112233"
	branding.PrimaryColor = &color
	if err := repo.SaveBrandingSettings(ctx, branding); err != nil {
		t.Fatal(err)
	}
	branding2, err := repo.GetBrandingSettings(ctx, domain.DefaultSettingsID)
	if err != nil {
		t.Fatal(err)
	}
	if branding2.PrimaryColor == nil || *branding2.PrimaryColor != color {
		t.Fatalf("branding = %+v", branding2)
	}
}

func TestSettingsRepositoryNotFound(t *testing.T) {
	repo, cleanup := openTestRepo(t)
	defer cleanup()
	_, err := repo.GetPasswordPolicy(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected not found")
	}
}
