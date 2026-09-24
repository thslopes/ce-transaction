package access

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MemoryUserRepository struct {
	mu          sync.RWMutex
	usersByID   map[string]*User
	usersByMail map[string]*User
	nextID      int
	clock       func() time.Time
}

func NewMemoryUserRepository(clock func() time.Time, bootstrap BootstrapAdmin) *MemoryUserRepository {
	if clock == nil {
		clock = time.Now
	}
	repo := &MemoryUserRepository{
		usersByID:   make(map[string]*User),
		usersByMail: make(map[string]*User),
		nextID:      1,
		clock:       clock,
	}
	if bootstrap.Email != "" {
		repo.mustCreateBootstrapAdmin(bootstrap)
	}
	return repo
}

func (r *MemoryUserRepository) mustCreateBootstrapAdmin(bootstrap BootstrapAdmin) {
	now := r.clock().UTC()
	user := &User{
		ID:           r.NextID(),
		Name:         bootstrap.Name,
		Email:        strings.ToLower(bootstrap.Email),
		Phone:        bootstrap.Phone,
		Status:       StatusActive,
		Profiles:     []string{AdminProfile},
		CreatedAt:    now,
		ValidatedAt:  &now,
		ValidatedBy:  "bootstrap",
		PasswordHash: hashPassword(bootstrap.Password),
	}
	_ = r.Save(user)
}

func (r *MemoryUserRepository) FindByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneUser(r.usersByMail[strings.ToLower(email)]), nil
}

func (r *MemoryUserRepository) FindByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneUser(r.usersByID[id]), nil
}

func (r *MemoryUserRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cloned := cloneUser(user)
	r.usersByID[cloned.ID] = cloned
	r.usersByMail[strings.ToLower(cloned.Email)] = cloned
	return nil
}

func (r *MemoryUserRepository) NextID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("usr_%03d", r.nextID)
	r.nextID++
	return id
}

func (r *MemoryUserRepository) CreateUserForTest(name, email, phone, password string, status UserStatus, profiles []string) *User {
	user := &User{
		ID:           r.NextID(),
		Name:         name,
		Email:        strings.ToLower(email),
		Phone:        phone,
		Status:       status,
		Profiles:     append([]string(nil), profiles...),
		CreatedAt:    r.clock().UTC(),
		PasswordHash: hashPassword(password),
	}
	_ = r.Save(user)
	return cloneUser(user)
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
