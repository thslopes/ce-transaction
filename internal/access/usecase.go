package access

import (
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo     UserRepository
	signer   *TokenSigner
	clock    func() time.Time
	baseURL  string
	tokenTTL time.Duration
	authTTL  time.Duration
}

func NewService(repo UserRepository, signer *TokenSigner, baseURL string, clock func() time.Time) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{
		repo:     repo,
		signer:   signer,
		clock:    clock,
		baseURL:  strings.TrimRight(baseURL, "/"),
		tokenTTL: 10 * time.Minute,
		authTTL:  24 * time.Hour,
	}
}

func (s *Service) Authenticate(email, password string) (string, *User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil || user.PasswordHash != hashPassword(password) {
		return "", nil, ErrInvalidCredentials
	}
	if user.Status != StatusActive {
		return "", nil, ErrInactiveUser
	}
	if len(user.Profiles) == 0 {
		return "", nil, ErrMissingProfiles
	}
	token, err := s.signer.Sign(tokenClaims{
		Subject: user.ID,
		Expiry:  s.clock().UTC().Add(s.authTTL),
		Purpose: "auth",
	})
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *Service) AuthenticateToken(token string) (*User, error) {
	claims, err := s.signer.Verify(token)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if claims.Purpose != "auth" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.FindByID(claims.Subject)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Status != StatusActive || len(user.Profiles) == 0 {
		return nil, ErrUnauthorized
	}
	return user, nil
}

func (s *Service) GenerateRegistrationLink(user *User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, ErrUnauthorized
	}
	if user.Status != StatusActive || !hasProfile(user.Profiles, AdminProfile) {
		return "", time.Time{}, ErrForbidden
	}
	expiresAt := s.clock().UTC().Add(s.tokenTTL)
	token, err := s.signer.Sign(tokenClaims{
		Subject: user.ID,
		Expiry:  expiresAt,
		Purpose: "signup",
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return fmt.Sprintf("%s/signup?token=%s", s.baseURL, token), expiresAt, nil
}

func (s *Service) RegisterUser(token, name, email, phone, password, passwordConfirmation string) (*User, error) {
	claims, err := s.signer.Verify(token)
	if err != nil {
		return nil, err
	}
	if claims.Purpose != "signup" {
		return nil, ErrInvalidToken
	}
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(phone) == "" || password == "" {
		return nil, ErrValidation
	}
	if password != passwordConfirmation {
		return nil, ErrValidation
	}
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}
	user := &User{
		ID:           s.repo.NextID(),
		Name:         name,
		Email:        strings.ToLower(email),
		Phone:        phone,
		Status:       StatusInactive,
		Profiles:     []string{},
		CreatedAt:    s.clock().UTC(),
		PasswordHash: hashPassword(password),
	}
	if err := s.repo.Save(user); err != nil {
		return nil, err
	}
	return cloneUser(user), nil
}
