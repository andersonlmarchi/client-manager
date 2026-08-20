package domain

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID     string
	Email  string
	Status UserStatus
}

func (u User) Validate() error {
	if u.ID == "" {
		return errInvalid("user id is required")
	}
	if u.Email == "" {
		return errInvalid("user email is required")
	}
	switch u.Status {
	case UserStatusActive, UserStatusDisabled:
	default:
		return errInvalid("user status is invalid")
	}
	return nil
}

type Credential struct {
	ID           string
	UserID       string
	PasswordHash string
	Algorithm    string
}

func (c Credential) Validate() error {
	if c.ID == "" {
		return errInvalid("credential id is required")
	}
	if c.UserID == "" {
		return errInvalid("credential user_id is required")
	}
	if c.PasswordHash == "" {
		return errInvalid("credential password_hash is required")
	}
	if c.Algorithm == "" {
		return errInvalid("credential algorithm is required")
	}
	return nil
}
