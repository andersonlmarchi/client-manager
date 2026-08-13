package infrastructure

import (
	"context"
	"fmt"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent/brandingsettings"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent/oidcsettings"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent/passwordpolicy"
	"github.com/andersonlmarchi/client-manager/services/configuration/ent/ratelimitsettings"
	"github.com/andersonlmarchi/client-manager/services/configuration/internal/domain"
)

type SettingsRepository struct {
	client *ent.Client
}

func NewSettingsRepository(client *ent.Client) *SettingsRepository {
	return &SettingsRepository{client: client}
}

func (r *SettingsRepository) GetPasswordPolicy(ctx context.Context, id string) (domain.PasswordPolicy, error) {
	row, err := r.client.PasswordPolicy.Query().Where(passwordpolicy.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.PasswordPolicy{}, shared.NewError(shared.CodeNotFound, "password policy not found")
		}
		return domain.PasswordPolicy{}, shared.Wrap(shared.CodeInternal, "get password policy", err)
	}
	return domain.PasswordPolicy{
		ID:            row.ID,
		MinLength:     row.MinLength,
		RequireUpper:  row.RequireUpper,
		RequireLower:  row.RequireLower,
		RequireNumber: row.RequireNumber,
		RequireSymbol: row.RequireSymbol,
		MaxAgeDays:    row.MaxAgeDays,
		HistoryCount:  row.HistoryCount,
	}, nil
}

func (r *SettingsRepository) SavePasswordPolicy(ctx context.Context, policy domain.PasswordPolicy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	err := r.client.PasswordPolicy.Create().
		SetID(policy.ID).
		SetMinLength(policy.MinLength).
		SetRequireUpper(policy.RequireUpper).
		SetRequireLower(policy.RequireLower).
		SetRequireNumber(policy.RequireNumber).
		SetRequireSymbol(policy.RequireSymbol).
		SetMaxAgeDays(policy.MaxAgeDays).
		SetHistoryCount(policy.HistoryCount).
		OnConflictColumns(passwordpolicy.FieldID).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return shared.Wrap(shared.CodeInternal, "save password policy", err)
	}
	return nil
}

func (r *SettingsRepository) GetOIDCSettings(ctx context.Context, id string) (domain.OIDCSettings, error) {
	row, err := r.client.OIDCSettings.Query().Where(oidcsettings.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.OIDCSettings{}, shared.NewError(shared.CodeNotFound, "oidc settings not found")
		}
		return domain.OIDCSettings{}, shared.Wrap(shared.CodeInternal, "get oidc settings", err)
	}
	return domain.OIDCSettings{
		ID:                     row.ID,
		Issuer:                 row.Issuer,
		AccessTokenTTLSeconds:  row.AccessTokenTTLSeconds,
		RefreshTokenTTLSeconds: row.RefreshTokenTTLSeconds,
		IDTokenTTLSeconds:      row.IDTokenTTLSeconds,
	}, nil
}

func (r *SettingsRepository) SaveOIDCSettings(ctx context.Context, settings domain.OIDCSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	err := r.client.OIDCSettings.Create().
		SetID(settings.ID).
		SetIssuer(settings.Issuer).
		SetAccessTokenTTLSeconds(settings.AccessTokenTTLSeconds).
		SetRefreshTokenTTLSeconds(settings.RefreshTokenTTLSeconds).
		SetIDTokenTTLSeconds(settings.IDTokenTTLSeconds).
		OnConflictColumns(oidcsettings.FieldID).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return shared.Wrap(shared.CodeInternal, "save oidc settings", err)
	}
	return nil
}

func (r *SettingsRepository) GetRateLimitSettings(ctx context.Context, id string) (domain.RateLimitSettings, error) {
	row, err := r.client.RateLimitSettings.Query().Where(ratelimitsettings.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.RateLimitSettings{}, shared.NewError(shared.CodeNotFound, "rate limit settings not found")
		}
		return domain.RateLimitSettings{}, shared.Wrap(shared.CodeInternal, "get rate limit settings", err)
	}
	return domain.RateLimitSettings{
		ID:           row.ID,
		DefaultRPS:   row.DefaultRps,
		DefaultBurst: row.DefaultBurst,
		LoginRPS:     row.LoginRps,
		LoginBurst:   row.LoginBurst,
	}, nil
}

func (r *SettingsRepository) SaveRateLimitSettings(ctx context.Context, settings domain.RateLimitSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	err := r.client.RateLimitSettings.Create().
		SetID(settings.ID).
		SetDefaultRps(settings.DefaultRPS).
		SetDefaultBurst(settings.DefaultBurst).
		SetLoginRps(settings.LoginRPS).
		SetLoginBurst(settings.LoginBurst).
		OnConflictColumns(ratelimitsettings.FieldID).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return shared.Wrap(shared.CodeInternal, "save rate limit settings", err)
	}
	return nil
}

func (r *SettingsRepository) GetBrandingSettings(ctx context.Context, id string) (domain.BrandingSettings, error) {
	row, err := r.client.BrandingSettings.Query().Where(brandingsettings.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.BrandingSettings{}, shared.NewError(shared.CodeNotFound, "branding settings not found")
		}
		return domain.BrandingSettings{}, shared.Wrap(shared.CodeInternal, "get branding settings", err)
	}
	return domain.BrandingSettings{
		ID:           row.ID,
		AppName:      row.AppName,
		LogoURL:      row.LogoURL,
		PrimaryColor: row.PrimaryColor,
		SupportEmail: row.SupportEmail,
		SMTPHostRef:  row.SMTPHostRef,
		SMTPPortRef:  row.SMTPPortRef,
		SMTPUserRef:  row.SMTPUserRef,
		FromEmail:    row.FromEmail,
	}, nil
}

func (r *SettingsRepository) SaveBrandingSettings(ctx context.Context, settings domain.BrandingSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	create := r.client.BrandingSettings.Create().
		SetID(settings.ID).
		SetAppName(settings.AppName)
	if settings.LogoURL != nil {
		create.SetLogoURL(*settings.LogoURL)
	}
	if settings.PrimaryColor != nil {
		create.SetPrimaryColor(*settings.PrimaryColor)
	}
	if settings.SupportEmail != nil {
		create.SetSupportEmail(*settings.SupportEmail)
	}
	if settings.SMTPHostRef != nil {
		create.SetSMTPHostRef(*settings.SMTPHostRef)
	}
	if settings.SMTPPortRef != nil {
		create.SetSMTPPortRef(*settings.SMTPPortRef)
	}
	if settings.SMTPUserRef != nil {
		create.SetSMTPUserRef(*settings.SMTPUserRef)
	}
	if settings.FromEmail != nil {
		create.SetFromEmail(*settings.FromEmail)
	}
	err := create.OnConflictColumns(brandingsettings.FieldID).UpdateNewValues().Exec(ctx)
	if err != nil {
		return shared.Wrap(shared.CodeInternal, "save branding settings", err)
	}
	return nil
}

func (r *SettingsRepository) EnsureDefaults(ctx context.Context) error {
	if err := r.SavePasswordPolicy(ctx, domain.DefaultPasswordPolicy()); err != nil {
		return fmt.Errorf("default password policy: %w", err)
	}
	if err := r.SaveOIDCSettings(ctx, domain.DefaultOIDCSettings()); err != nil {
		return fmt.Errorf("default oidc settings: %w", err)
	}
	if err := r.SaveRateLimitSettings(ctx, domain.DefaultRateLimitSettings()); err != nil {
		return fmt.Errorf("default rate limits: %w", err)
	}
	if err := r.SaveBrandingSettings(ctx, domain.DefaultBrandingSettings()); err != nil {
		return fmt.Errorf("default branding: %w", err)
	}
	return nil
}
