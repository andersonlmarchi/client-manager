package domain_test

import (
	"testing"

	"github.com/andersonlmarchi/client-manager/services/configuration/internal/domain"
)

func TestPasswordPolicyValidate(t *testing.T) {
	t.Parallel()
	p := domain.DefaultPasswordPolicy()
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p.MinLength = 4
	if err := p.Validate(); err == nil {
		t.Fatal("expected invalid min_length")
	}
}

func TestOIDCSettingsValidate(t *testing.T) {
	t.Parallel()
	s := domain.DefaultOIDCSettings()
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Issuer = ""
	if err := s.Validate(); err == nil {
		t.Fatal("expected invalid issuer")
	}
}

func TestRateLimitSettingsValidate(t *testing.T) {
	t.Parallel()
	s := domain.DefaultRateLimitSettings()
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.DefaultBurst = 1
	s.DefaultRPS = 10
	if err := s.Validate(); err == nil {
		t.Fatal("expected invalid burst")
	}
}

func TestBrandingSettingsValidate(t *testing.T) {
	t.Parallel()
	s := domain.DefaultBrandingSettings()
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.AppName = ""
	if err := s.Validate(); err == nil {
		t.Fatal("expected invalid app name")
	}
}
