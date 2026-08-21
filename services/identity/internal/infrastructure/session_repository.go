package infrastructure

import (
	"context"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/ent"
	"github.com/andersonlmarchi/client-manager/services/identity/ent/session"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
)

type SessionRepository struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time, ip, userAgent *string) (domain.Session, error) {
	id, err := shared.NewUUID()
	if err != nil {
		return domain.Session{}, shared.Wrap(shared.CodeInternal, "session id", err)
	}
	create := r.client.Session.Create().
		SetID(id).
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		SetUserID(userID)
	if ip != nil {
		create.SetIP(*ip)
	}
	if userAgent != nil {
		create.SetUserAgent(*userAgent)
	}
	row, err := create.Save(ctx)
	if err != nil {
		return domain.Session{}, shared.Wrap(shared.CodeInternal, "create session", err)
	}
	return domain.Session{
		ID:        row.ID,
		UserID:    userID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		IP:        row.IP,
		UserAgent: row.UserAgent,
	}, nil
}

func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error) {
	row, err := r.client.Session.Query().Where(session.TokenHashEQ(tokenHash)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.Session{}, shared.NewError(shared.CodeUnauthorized, "invalid session")
		}
		return domain.Session{}, shared.Wrap(shared.CodeInternal, "get session", err)
	}
	userID := ""
	if u, err := row.QueryUser().Only(ctx); err == nil {
		userID = u.ID
	} else {
		return domain.Session{}, shared.Wrap(shared.CodeInternal, "session user", err)
	}
	return domain.Session{
		ID:        row.ID,
		UserID:    userID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: row.RevokedAt,
		IP:        row.IP,
		UserAgent: row.UserAgent,
	}, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, sessionID string, at time.Time) error {
	n, err := r.client.Session.Update().
		Where(session.IDEQ(sessionID), session.RevokedAtIsNil()).
		SetRevokedAt(at).
		Save(ctx)
	if err != nil {
		return shared.Wrap(shared.CodeInternal, "revoke session", err)
	}
	if n == 0 {
		return shared.NewError(shared.CodeUnauthorized, "invalid session")
	}
	return nil
}
