package access

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrMissingProfiles    = errors.New("user has no profiles")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrValidation         = errors.New("validation error")
)

type contextKey string

const userContextKey contextKey = "access-user"

type UserStatus string

const (
	StatusInactive UserStatus = "inactive"
	StatusActive   UserStatus = "active"
)

const AdminProfile = "admin"

type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Status       UserStatus `json:"status"`
	Profiles     []string   `json:"profiles"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	ValidatedAt  *time.Time `json:"validated_at,omitempty"`
	ValidatedBy  string     `json:"validated_by,omitempty"`
	PasswordHash string     `json:"-"`
}

type BootstrapAdmin struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func cloneUser(user *User) *User {
	if user == nil {
		return nil
	}
	copy := *user
	copy.Profiles = append([]string(nil), user.Profiles...)
	return &copy
}

func hasProfile(profiles []string, wanted string) bool {
	for _, profile := range profiles {
		if profile == wanted {
			return true
		}
	}
	return false
}