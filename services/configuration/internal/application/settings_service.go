package application

import (
	"context"

	"github.com/andersonlmarchi/client-manager/services/configuration/internal/domain"
)

type SettingsStore interface {
	GetPasswordPolicy(ctx context.Context, id string) (domain.PasswordPolicy, error)
	SavePasswordPolicy(ctx context.Context, policy domain.PasswordPolicy) error
	GetOIDCSettings(ctx context.Context, id string) (domain.OIDCSettings, error)
	SaveOIDCSettings(ctx context.Context, settings domain.OIDCSettings) error
	GetRateLimitSettings(ctx context.Context, id string) (domain.RateLimitSettings, error)
	SaveRateLimitSettings(ctx context.Context, settings domain.RateLimitSettings) error
	GetBrandingSettings(ctx context.Context, id string) (domain.BrandingSettings, error)
	SaveBrandingSettings(ctx context.Context, settings domain.BrandingSettings) error
}

type SettingsService struct {
	store SettingsStore
}

func NewSettingsService(store SettingsStore) *SettingsService {
	return &SettingsService{store: store}
}

func (s *SettingsService) GetPasswordPolicy(ctx context.Context) (domain.PasswordPolicy, error) {
	return s.store.GetPasswordPolicy(ctx, domain.DefaultSettingsID)
}

func (s *SettingsService) UpdatePasswordPolicy(ctx context.Context, policy domain.PasswordPolicy) error {
	policy.ID = domain.DefaultSettingsID
	return s.store.SavePasswordPolicy(ctx, policy)
}

func (s *SettingsService) GetOIDCSettings(ctx context.Context) (domain.OIDCSettings, error) {
	return s.store.GetOIDCSettings(ctx, domain.DefaultSettingsID)
}

func (s *SettingsService) UpdateOIDCSettings(ctx context.Context, settings domain.OIDCSettings) error {
	settings.ID = domain.DefaultSettingsID
	return s.store.SaveOIDCSettings(ctx, settings)
}

func (s *SettingsService) GetRateLimitSettings(ctx context.Context) (domain.RateLimitSettings, error) {
	return s.store.GetRateLimitSettings(ctx, domain.DefaultSettingsID)
}

func (s *SettingsService) UpdateRateLimitSettings(ctx context.Context, settings domain.RateLimitSettings) error {
	settings.ID = domain.DefaultSettingsID
	return s.store.SaveRateLimitSettings(ctx, settings)
}

func (s *SettingsService) GetBrandingSettings(ctx context.Context) (domain.BrandingSettings, error) {
	return s.store.GetBrandingSettings(ctx, domain.DefaultSettingsID)
}

func (s *SettingsService) UpdateBrandingSettings(ctx context.Context, settings domain.BrandingSettings) error {
	settings.ID = domain.DefaultSettingsID
	return s.store.SaveBrandingSettings(ctx, settings)
}
