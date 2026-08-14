package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andersonlmarchi/client-manager/services/configuration/internal/application"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/infrastructure"
	confighttp "github.com/andersonlmarchi/client-manager/services/configuration/internal/transport/http"
)

func newTestServer(t *testing.T) (*confighttp.Server, http.Handler) {
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
	ctx := context.Background()
	if err := infrastructure.Migrate(ctx, client); err != nil {
		t.Fatal(err)
	}
	repo := infrastructure.NewSettingsRepository(client)
	if err := repo.EnsureDefaults(ctx); err != nil {
		t.Fatal(err)
	}
	svc := application.NewSettingsService(repo)
	api := confighttp.NewServer(svc, "test-admin-key")
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	return api, mux
}

func TestHandleHealth(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetPasswordPolicy(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/password-policy", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["min_length"].(float64) != 12 {
		t.Fatalf("body=%v", body)
	}
}

func TestPutPasswordPolicyRequiresAdmin(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	payload := []byte(`{"min_length":16,"require_upper":true,"require_lower":true,"require_number":true,"require_symbol":true,"max_age_days":0,"history_count":0}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/password-policy", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestPutPasswordPolicyWithAdminKey(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	payload := []byte(`{"min_length":16,"require_upper":true,"require_lower":true,"require_number":true,"require_symbol":true,"max_age_days":90,"history_count":3}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/password-policy", bytes.NewReader(payload))
	req.Header.Set("X-Admin-Key", "test-admin-key")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["min_length"].(float64) != 16 || body["max_age_days"].(float64) != 90 {
		t.Fatalf("body=%v", body)
	}
}

func TestPutPasswordPolicyValidation(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	payload := []byte(`{"min_length":2,"require_upper":true,"require_lower":true,"require_number":true,"require_symbol":true,"max_age_days":0,"history_count":0}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/password-policy", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer test-admin-key")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestOIDCRateLimitAndBrandingRoundTrip(t *testing.T) {
	t.Parallel()
	_, mux := newTestServer(t)

	put := func(path string, body string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader([]byte(body)))
		req.Header.Set("X-Admin-Key", "test-admin-key")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status=%d body=%s", path, rec.Code, rec.Body.String())
		}
	}

	put("/v1/oidc-settings", `{"issuer":"https://auth.example.com","access_token_ttl_seconds":600,"refresh_token_ttl_seconds":86400,"id_token_ttl_seconds":600}`)
	put("/v1/rate-limits", `{"default_rps":30,"default_burst":60,"login_rps":2,"login_burst":4}`)
	put("/v1/branding", `{"app_name":"Acme","primary_color":"#abc123","smtp_host_ref":"secrets/smtp/host"}`)

	for _, path := range []string{"/v1/oidc-settings", "/v1/rate-limits", "/v1/branding"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s status=%d", path, rec.Code)
		}
	}
}
