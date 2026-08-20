package infrastructure

import (
	"context"
	"strings"

	"github.com/andersonlmarchi/client-manager/packages/shared"
	"github.com/andersonlmarchi/client-manager/services/identity/ent"
	"github.com/andersonlmarchi/client-manager/services/identity/ent/credential"
	"github.com/andersonlmarchi/client-manager/services/identity/ent/user"
	"github.com/andersonlmarchi/client-manager/services/identity/internal/domain"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) CreateUserWithPassword(ctx context.Context, email, password string) (domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return domain.User{}, shared.NewError(shared.CodeInvalid, "user email is required")
	}
	hash, err := domain.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	userID, err := shared.NewUUID()
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "user id", err)
	}
	credID, err := shared.NewUUID()
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "credential id", err)
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "begin tx", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	exists, err := tx.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "check email", err)
	}
	if exists {
		err = shared.NewError(shared.CodeConflict, "email already registered")
		return domain.User{}, err
	}

	row, err := tx.User.Create().
		SetID(userID).
		SetEmail(email).
		SetStatus(user.StatusActive).
		Save(ctx)
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "create user", err)
	}
	_, err = tx.Credential.Create().
		SetID(credID).
		SetPasswordHash(hash).
		SetAlgorithm(domain.PasswordAlgorithm).
		SetUserID(row.ID).
		Save(ctx)
	if err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "create credential", err)
	}
	if err = tx.Commit(); err != nil {
		return domain.User{}, shared.Wrap(shared.CodeInternal, "commit tx", err)
	}
	err = nil
	return domain.User{ID: row.ID, Email: row.Email, Status: domain.UserStatus(row.Status)}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	row, err := r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, shared.NewError(shared.CodeNotFound, "user not found")
		}
		return domain.User{}, shared.Wrap(shared.CodeInternal, "get user by email", err)
	}
	return domain.User{ID: row.ID, Email: row.Email, Status: domain.UserStatus(row.Status)}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	row, err := r.client.User.Query().Where(user.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, shared.NewError(shared.CodeNotFound, "user not found")
		}
		return domain.User{}, shared.Wrap(shared.CodeInternal, "get user by id", err)
	}
	return domain.User{ID: row.ID, Email: row.Email, Status: domain.UserStatus(row.Status)}, nil
}

func (r *UserRepository) GetCredentialByUserID(ctx context.Context, userID string) (domain.Credential, error) {
	row, err := r.client.Credential.Query().Where(credential.HasUserWith(user.IDEQ(userID))).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.Credential{}, shared.NewError(shared.CodeNotFound, "credential not found")
		}
		return domain.Credential{}, shared.Wrap(shared.CodeInternal, "get credential", err)
	}
	return domain.Credential{
		ID:           row.ID,
		UserID:       userID,
		PasswordHash: row.PasswordHash,
		Algorithm:    row.Algorithm,
	}, nil
}

func (r *UserRepository) Authenticate(ctx context.Context, email, password string) (domain.User, error) {
	u, err := r.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if u.Status != domain.UserStatusActive {
		return domain.User{}, shared.NewError(shared.CodeForbidden, "user is disabled")
	}
	cred, err := r.GetCredentialByUserID(ctx, u.ID)
	if err != nil {
		return domain.User{}, err
	}
	ok, err := domain.VerifyPassword(cred.PasswordHash, password)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, shared.NewError(shared.CodeUnauthorized, "invalid credentials")
	}
	return u, nil
}
