package http

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/application"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/domain"
)

const maxBodyBytes = 64 << 10

type Server struct {
	Settings  *application.SettingsService
	AdminKey  string
}

func NewServer(settings *application.SettingsService, adminKey string) *Server {
	return &Server{Settings: settings, AdminKey: adminKey}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", HandleHealth)
	mux.HandleFunc("GET /v1/password-policy", s.handleGetPasswordPolicy)
	mux.HandleFunc("PUT /v1/password-policy", s.requireAdmin(s.handlePutPasswordPolicy))
	mux.HandleFunc("GET /v1/oidc-settings", s.handleGetOIDCSettings)
	mux.HandleFunc("PUT /v1/oidc-settings", s.requireAdmin(s.handlePutOIDCSettings))
	mux.HandleFunc("GET /v1/rate-limits", s.handleGetRateLimits)
	mux.HandleFunc("PUT /v1/rate-limits", s.requireAdmin(s.handlePutRateLimits))
	mux.HandleFunc("GET /v1/branding", s.handleGetBranding)
	mux.HandleFunc("PUT /v1/branding", s.requireAdmin(s.handlePutBranding))
}

func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.authorized(r) {
			_ = shared.WriteErrorProblem(w, r, shared.NewError(shared.CodeUnauthorized, "admin authentication required"))
			return
		}
		next(w, r)
	}
}

func (s *Server) authorized(r *http.Request) bool {
	if s.AdminKey == "" {
		return false
	}
	key := r.Header.Get("X-Admin-Key")
	if key == "" {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			key = strings.TrimSpace(auth[7:])
		}
	}
	if key == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(key), []byte(s.AdminKey)) == 1
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

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

type passwordPolicyDTO struct {
	MinLength     int  `json:"min_length"`
	RequireUpper  bool `json:"require_upper"`
	RequireLower  bool `json:"require_lower"`
	RequireNumber bool `json:"require_number"`
	RequireSymbol bool `json:"require_symbol"`
	MaxAgeDays    int  `json:"max_age_days"`
	HistoryCount  int  `json:"history_count"`
}

func passwordPolicyFromDomain(p domain.PasswordPolicy) passwordPolicyDTO {
	return passwordPolicyDTO{
		MinLength:     p.MinLength,
		RequireUpper:  p.RequireUpper,
		RequireLower:  p.RequireLower,
		RequireNumber: p.RequireNumber,
		RequireSymbol: p.RequireSymbol,
		MaxAgeDays:    p.MaxAgeDays,
		HistoryCount:  p.HistoryCount,
	}
}

func (d passwordPolicyDTO) toDomain() domain.PasswordPolicy {
	return domain.PasswordPolicy{
		MinLength:     d.MinLength,
		RequireUpper:  d.RequireUpper,
		RequireLower:  d.RequireLower,
		RequireNumber: d.RequireNumber,
		RequireSymbol: d.RequireSymbol,
		MaxAgeDays:    d.MaxAgeDays,
		HistoryCount:  d.HistoryCount,
	}
}

func (s *Server) handleGetPasswordPolicy(w http.ResponseWriter, r *http.Request) {
	p, err := s.Settings.GetPasswordPolicy(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, passwordPolicyFromDomain(p))
}

func (s *Server) handlePutPasswordPolicy(w http.ResponseWriter, r *http.Request) {
	var body passwordPolicyDTO
	if err := decodeJSON(w, r, &body); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	if err := s.Settings.UpdatePasswordPolicy(r.Context(), body.toDomain()); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	p, err := s.Settings.GetPasswordPolicy(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, passwordPolicyFromDomain(p))
}

type oidcSettingsDTO struct {
	Issuer                 string `json:"issuer"`
	AccessTokenTTLSeconds  int    `json:"access_token_ttl_seconds"`
	RefreshTokenTTLSeconds int    `json:"refresh_token_ttl_seconds"`
	IDTokenTTLSeconds      int    `json:"id_token_ttl_seconds"`
}

func oidcFromDomain(s domain.OIDCSettings) oidcSettingsDTO {
	return oidcSettingsDTO{
		Issuer:                 s.Issuer,
		AccessTokenTTLSeconds:  s.AccessTokenTTLSeconds,
		RefreshTokenTTLSeconds: s.RefreshTokenTTLSeconds,
		IDTokenTTLSeconds:      s.IDTokenTTLSeconds,
	}
}

func (d oidcSettingsDTO) toDomain() domain.OIDCSettings {
	return domain.OIDCSettings{
		Issuer:                 d.Issuer,
		AccessTokenTTLSeconds:  d.AccessTokenTTLSeconds,
		RefreshTokenTTLSeconds: d.RefreshTokenTTLSeconds,
		IDTokenTTLSeconds:      d.IDTokenTTLSeconds,
	}
}

func (s *Server) handleGetOIDCSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.Settings.GetOIDCSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, oidcFromDomain(settings))
}

func (s *Server) handlePutOIDCSettings(w http.ResponseWriter, r *http.Request) {
	var body oidcSettingsDTO
	if err := decodeJSON(w, r, &body); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	if err := s.Settings.UpdateOIDCSettings(r.Context(), body.toDomain()); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	settings, err := s.Settings.GetOIDCSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, oidcFromDomain(settings))
}

type rateLimitDTO struct {
	DefaultRPS   int `json:"default_rps"`
	DefaultBurst int `json:"default_burst"`
	LoginRPS     int `json:"login_rps"`
	LoginBurst   int `json:"login_burst"`
}

func rateLimitFromDomain(s domain.RateLimitSettings) rateLimitDTO {
	return rateLimitDTO{
		DefaultRPS:   s.DefaultRPS,
		DefaultBurst: s.DefaultBurst,
		LoginRPS:     s.LoginRPS,
		LoginBurst:   s.LoginBurst,
	}
}

func (d rateLimitDTO) toDomain() domain.RateLimitSettings {
	return domain.RateLimitSettings{
		DefaultRPS:   d.DefaultRPS,
		DefaultBurst: d.DefaultBurst,
		LoginRPS:     d.LoginRPS,
		LoginBurst:   d.LoginBurst,
	}
}

func (s *Server) handleGetRateLimits(w http.ResponseWriter, r *http.Request) {
	settings, err := s.Settings.GetRateLimitSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, rateLimitFromDomain(settings))
}

func (s *Server) handlePutRateLimits(w http.ResponseWriter, r *http.Request) {
	var body rateLimitDTO
	if err := decodeJSON(w, r, &body); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	if err := s.Settings.UpdateRateLimitSettings(r.Context(), body.toDomain()); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	settings, err := s.Settings.GetRateLimitSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, rateLimitFromDomain(settings))
}

type brandingDTO struct {
	AppName      string  `json:"app_name"`
	LogoURL      *string `json:"logo_url"`
	PrimaryColor *string `json:"primary_color"`
	SupportEmail *string `json:"support_email"`
	SMTPHostRef  *string `json:"smtp_host_ref"`
	SMTPPortRef  *string `json:"smtp_port_ref"`
	SMTPUserRef  *string `json:"smtp_user_ref"`
	FromEmail    *string `json:"from_email"`
}

func brandingFromDomain(s domain.BrandingSettings) brandingDTO {
	return brandingDTO{
		AppName:      s.AppName,
		LogoURL:      s.LogoURL,
		PrimaryColor: s.PrimaryColor,
		SupportEmail: s.SupportEmail,
		SMTPHostRef:  s.SMTPHostRef,
		SMTPPortRef:  s.SMTPPortRef,
		SMTPUserRef:  s.SMTPUserRef,
		FromEmail:    s.FromEmail,
	}
}

func (d brandingDTO) toDomain() domain.BrandingSettings {
	return domain.BrandingSettings{
		AppName:      d.AppName,
		LogoURL:      d.LogoURL,
		PrimaryColor: d.PrimaryColor,
		SupportEmail: d.SupportEmail,
		SMTPHostRef:  d.SMTPHostRef,
		SMTPPortRef:  d.SMTPPortRef,
		SMTPUserRef:  d.SMTPUserRef,
		FromEmail:    d.FromEmail,
	}
}

func (s *Server) handleGetBranding(w http.ResponseWriter, r *http.Request) {
	settings, err := s.Settings.GetBrandingSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, brandingFromDomain(settings))
}

func (s *Server) handlePutBranding(w http.ResponseWriter, r *http.Request) {
	var body brandingDTO
	if err := decodeJSON(w, r, &body); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	if err := s.Settings.UpdateBrandingSettings(r.Context(), body.toDomain()); err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	settings, err := s.Settings.GetBrandingSettings(r.Context())
	if err != nil {
		_ = shared.WriteErrorProblem(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, brandingFromDomain(settings))
}
