package application

import (
	"context"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
)

type UserAuthStore interface {
	Authenticate(ctx context.Context, email, password string) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	CreateUserWithPassword(ctx context.Context, email, password string) (domain.User, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time, ip, userAgent *string) (domain.Session, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error)
	Revoke(ctx context.Context, sessionID string, at time.Time) error
}

type AuthService struct {
	users      UserAuthStore
	sessions   SessionStore
	sessionTTL time.Duration
	now        func() time.Time
}

func NewAuthService(users UserAuthStore, sessions SessionStore, sessionTTL time.Duration) *AuthService {
	if sessionTTL <= 0 {
		sessionTTL = 168 * time.Hour
	}
	return &AuthService{
		users:      users,
		sessions:   sessions,
		sessionTTL: sessionTTL,
		now:        time.Now,
	}
}

type LoginResult struct {
	User         domain.User
	SessionID    string
	SessionToken string
	ExpiresAt    time.Time
}

func (s *AuthService) Login(ctx context.Context, email, password string, ip, userAgent *string) (LoginResult, error) {
	user, err := s.users.Authenticate(ctx, email, password)
	if err != nil {
		if e, ok := shared.AsError(err); ok && (e.Code == shared.CodeNotFound || e.Code == shared.CodeUnauthorized) {
			return LoginResult{}, shared.NewError(shared.CodeUnauthorized, "invalid credentials")
		}
		return LoginResult{}, err
	}
	plain, hash, err := domain.NewSessionToken()
	if err != nil {
		return LoginResult{}, shared.Wrap(shared.CodeInternal, "session token", err)
	}
	now := s.now().UTC()
	expires := now.Add(s.sessionTTL)
	sess, err := s.sessions.Create(ctx, user.ID, hash, expires, ip, userAgent)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		User:         user,
		SessionID:    sess.ID,
		SessionToken: plain,
		ExpiresAt:    expires,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, plainToken string) error {
	if plainToken == "" {
		return shared.NewError(shared.CodeUnauthorized, "session required")
	}
	sess, err := s.sessions.GetByTokenHash(ctx, domain.HashSessionToken(plainToken))
	if err != nil {
		return err
	}
	if !sess.Active(s.now().UTC()) {
		return shared.NewError(shared.CodeUnauthorized, "session expired")
	}
	return s.sessions.Revoke(ctx, sess.ID, s.now().UTC())
}

func (s *AuthService) CurrentUser(ctx context.Context, plainToken string) (domain.User, domain.Session, error) {
	if plainToken == "" {
		return domain.User{}, domain.Session{}, shared.NewError(shared.CodeUnauthorized, "session required")
	}
	sess, err := s.sessions.GetByTokenHash(ctx, domain.HashSessionToken(plainToken))
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}
	if !sess.Active(s.now().UTC()) {
		return domain.User{}, domain.Session{}, shared.NewError(shared.CodeUnauthorized, "session expired")
	}
	user, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return domain.User{}, domain.Session{}, err
	}
	if user.Status != domain.UserStatusActive {
		return domain.User{}, domain.Session{}, shared.NewError(shared.CodeForbidden, "user is disabled")
	}
	return user, sess, nil
}
