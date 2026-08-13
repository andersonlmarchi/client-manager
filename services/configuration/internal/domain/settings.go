package domain

const DefaultSettingsID = "default"

type PasswordPolicy struct {
	ID            string
	MinLength     int
	RequireUpper  bool
	RequireLower  bool
	RequireNumber bool
	RequireSymbol bool
	MaxAgeDays    int
	HistoryCount  int
}

func DefaultPasswordPolicy() PasswordPolicy {
	return PasswordPolicy{
		ID:            DefaultSettingsID,
		MinLength:     12,
		RequireUpper:  true,
		RequireLower:  true,
		RequireNumber: true,
		RequireSymbol: true,
		MaxAgeDays:    0,
		HistoryCount:  0,
	}
}

func (p PasswordPolicy) Validate() error {
	if p.ID == "" {
		return errInvalid("password policy id is required")
	}
	if p.MinLength < 8 || p.MinLength > 128 {
		return errInvalid("password min_length must be between 8 and 128")
	}
	if p.MaxAgeDays < 0 || p.MaxAgeDays > 3650 {
		return errInvalid("password max_age_days must be between 0 and 3650")
	}
	if p.HistoryCount < 0 || p.HistoryCount > 24 {
		return errInvalid("password history_count must be between 0 and 24")
	}
	return nil
}

type OIDCSettings struct {
	ID                      string
	Issuer                  string
	AccessTokenTTLSeconds   int
	RefreshTokenTTLSeconds  int
	IDTokenTTLSeconds       int
}

func DefaultOIDCSettings() OIDCSettings {
	return OIDCSettings{
		ID:                     DefaultSettingsID,
		Issuer:                 "http://localhost:8081",
		AccessTokenTTLSeconds:  900,
		RefreshTokenTTLSeconds: 2592000,
		IDTokenTTLSeconds:      900,
	}
}

func (s OIDCSettings) Validate() error {
	if s.ID == "" {
		return errInvalid("oidc settings id is required")
	}
	if s.Issuer == "" {
		return errInvalid("oidc issuer is required")
	}
	if s.AccessTokenTTLSeconds < 60 || s.AccessTokenTTLSeconds > 86400 {
		return errInvalid("access_token_ttl_seconds out of range")
	}
	if s.RefreshTokenTTLSeconds < 60 || s.RefreshTokenTTLSeconds > 31536000 {
		return errInvalid("refresh_token_ttl_seconds out of range")
	}
	if s.IDTokenTTLSeconds < 60 || s.IDTokenTTLSeconds > 86400 {
		return errInvalid("id_token_ttl_seconds out of range")
	}
	return nil
}

type RateLimitSettings struct {
	ID           string
	DefaultRPS   int
	DefaultBurst int
	LoginRPS     int
	LoginBurst   int
}

func DefaultRateLimitSettings() RateLimitSettings {
	return RateLimitSettings{
		ID:           DefaultSettingsID,
		DefaultRPS:   60,
		DefaultBurst: 120,
		LoginRPS:     5,
		LoginBurst:   10,
	}
}

func (s RateLimitSettings) Validate() error {
	if s.ID == "" {
		return errInvalid("rate limit settings id is required")
	}
	if s.DefaultRPS < 1 || s.DefaultBurst < s.DefaultRPS {
		return errInvalid("invalid default rate limit values")
	}
	if s.LoginRPS < 1 || s.LoginBurst < s.LoginRPS {
		return errInvalid("invalid login rate limit values")
	}
	return nil
}

type BrandingSettings struct {
	ID          string
	AppName     string
	LogoURL     *string
	PrimaryColor *string
	SupportEmail *string
	SMTPHostRef *string
	SMTPPortRef *string
	SMTPUserRef *string
	FromEmail   *string
}

func DefaultBrandingSettings() BrandingSettings {
	return BrandingSettings{
		ID:      DefaultSettingsID,
		AppName: "Client Manager",
	}
}

func (s BrandingSettings) Validate() error {
	if s.ID == "" {
		return errInvalid("branding settings id is required")
	}
	if s.AppName == "" {
		return errInvalid("branding app_name is required")
	}
	return nil
}
