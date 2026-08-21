package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andersonlmarchi/client-manager/services/identity/internal/application"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/infrastructure"
	identityhttp "github.com/andersonlmarchi/client-manager/services/identity/internal/transport/http"
)

func newTestAPI(t *testing.T) (http.Handler, *infrastructure.UserRepository) {
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
	users := infrastructure.NewUserRepository(client)
	sessions := infrastructure.NewSessionRepository(client)
	auth := application.NewAuthService(users, sessions, time.Hour)
	api := identityhttp.NewServer(auth, false)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	return mux, users
}

func TestHandleHealth(t *testing.T) {
	t.Parallel()
	mux, _ := newTestAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestHandleHealthRejectsPost(t *testing.T) {
	t.Parallel()
	mux, _ := newTestAPI(t)
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestLoginSuccessAndMeAndLogout(t *testing.T) {
	t.Parallel()
	mux, users := newTestAPI(t)
	ctx := context.Background()
	if _, err := users.CreateUserWithPassword(ctx, "login@example.com", "super-secret-pass"); err != nil {
		t.Fatal(err)
	}

	loginBody := []byte(`{"email":"login@example.com","password":"super-secret-pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", rec.Code, rec.Body.String())
	}
	var loginResp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &loginResp); err != nil {
		t.Fatal(err)
	}
	token, _ := loginResp["session_token"].(string)
	if token == "" {
		t.Fatal("missing session_token")
	}
	cookie := rec.Result().Cookies()
	foundCookie := false
	for _, c := range cookie {
		if c.Name == domain.SessionCookieName && c.Value == token && c.HttpOnly {
			foundCookie = true
		}
	}
	if !foundCookie {
		t.Fatal("expected httpOnly session cookie")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+token)
	meRec := httptest.NewRecorder()
	mux.ServeHTTP(meRec, meReq)
	if meRec.Code != http.StatusOK {
		t.Fatalf("me status=%d body=%s", meRec.Code, meRec.Body.String())
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/v1/logout", nil)
	logoutReq.AddCookie(&http.Cookie{Name: domain.SessionCookieName, Value: token})
	logoutRec := httptest.NewRecorder()
	mux.ServeHTTP(logoutRec, logoutReq)
	if logoutRec.Code != http.StatusNoContent {
		t.Fatalf("logout status=%d", logoutRec.Code)
	}

	meReq2 := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	meReq2.Header.Set("Authorization", "Bearer "+token)
	meRec2 := httptest.NewRecorder()
	mux.ServeHTTP(meRec2, meReq2)
	if meRec2.Code != http.StatusUnauthorized {
		t.Fatalf("me after logout status=%d", meRec2.Code)
	}
}

func TestLoginFailure(t *testing.T) {
	t.Parallel()
	mux, users := newTestAPI(t)
	if _, err := users.CreateUserWithPassword(context.Background(), "fail@example.com", "super-secret-pass"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader([]byte(`{"email":"fail@example.com","password":"wrong-password-x"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
