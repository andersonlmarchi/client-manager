package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/application"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
)

const maxBodyBytes = 64 << 10

type healthResponse struct {
	Status string `json:"status"`
}

type Server struct {
	Auth         *application.AuthService
	CookieSecure bool
}

func NewServer(auth *application.AuthService, cookieSecure bool) *Server {
	return &Server{Auth: auth, CookieSecure: cookieSecure}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", HandleHealth)
	mux.HandleFunc("POST /v1/login", s.handleLogin)
	mux.HandleFunc("POST /v1/logout", s.handleLogout)
	mux.HandleFunc("GET /v1/me", s.handleMe)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		if err == io.EOF {
			return shared.NewError(shared.CodeInvalid, "request body is required")
		}
		return shared.Wrap(shared.CodeInvalid, "invalid json body", err)
	}
	return nil
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userDTO struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type loginResponse struct {
	User         userDTO   `json:"user"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func userFromDomain(u domain.User) userDTO {
	return userDTO{ID: u.ID, Email: u.Email, Status: string(u.Status)}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := decodeJSON(w, r, &body); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	ip := clientIP(r)
	ua := r.UserAgent()
	var ipPtr, uaPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	if ua != "" {
		uaPtr = &ua
	}
	result, err := s.Auth.Login(r.Context(), body.Email, body.Password, ipPtr, uaPtr)
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    result.SessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  result.ExpiresAt,
	})
	writeJSON(w, http.StatusOK, loginResponse{
		User:         userFromDomain(result.User),
		SessionToken: result.SessionToken,
		ExpiresAt:    result.ExpiresAt,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := sessionTokenFromRequest(r)
	if err := s.Auth.Logout(r.Context(), token); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	token := sessionTokenFromRequest(r)
	user, _, err := s.Auth.CurrentUser(r.Context(), token)
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, userFromDomain(user))
}

func sessionTokenFromRequest(r *http.Request) string {
	if c, err := r.Cookie(domain.SessionCookieName); err == nil && c.Value != "" {
		return c.Value
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return ""
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i >= 0 {
		return host[:i]
	}
	return host
}
