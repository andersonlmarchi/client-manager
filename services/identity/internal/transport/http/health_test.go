package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	identityhttp "github.com/andersonlmarchi/client-manager/services/identity/internal/transport/http"
)

func TestHandleHealth(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	identityhttp.NewServer().RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandleHealthRejectsPost(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	identityhttp.NewServer().RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
