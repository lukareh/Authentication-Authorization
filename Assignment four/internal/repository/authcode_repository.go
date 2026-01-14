package repository

import (
	"sso-mock/internal/domain"
	"sync"
	"time"
)

// AuthCodeRepository handles storage and retrieval of authorization codes
type AuthCodeRepository struct {
	codes map[string]domain.AuthCode
	mu    sync.RWMutex
}

// NewAuthCodeRepository creates a new authorization code repository
func NewAuthCodeRepository() *AuthCodeRepository {
	return &AuthCodeRepository{
		codes: make(map[string]domain.AuthCode),
	}
}

// Store saves an authorization code
func (r *AuthCodeRepository) Store(code, username string, expiresIn time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.codes[code] = domain.AuthCode{
		Code:      code,
		Username:  username,
		ExpiresAt: time.Now().Add(expiresIn),
	}
	return nil
}

// Get retrieves an authorization code
func (r *AuthCodeRepository) Get(code string) (*domain.AuthCode, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	authCode, exists := r.codes[code]
	if !exists {
		return nil, false
	}
	return &authCode, true
}

// Delete removes an authorization code (used after exchange)
func (r *AuthCodeRepository) Delete(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.codes, code)
	return nil
}

// CleanExpired removes expired authorization codes
func (r *AuthCodeRepository) CleanExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for code, authCode := range r.codes {
		if now.After(authCode.ExpiresAt) {
			delete(r.codes, code)
		}
	}
}
