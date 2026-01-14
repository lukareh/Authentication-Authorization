package services

import (
	"log"
	"sso-mock/internal/domain"
	"sso-mock/internal/repository"
	"time"
)

// AuthService handles authentication business logic
type AuthService struct {
	repo *repository.AuthCodeRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(repo *repository.AuthCodeRepository) *AuthService {
	return &AuthService{repo: repo}
}

// ValidateCredentials validates user credentials
// In this mock, any non-empty credentials are valid
func (s *AuthService) ValidateCredentials(creds *domain.Credentials) bool {
	return creds.Validate()
}

// GenerateAuthCode generates and stores an authorization code
func (s *AuthService) GenerateAuthCode(username string) (string, error) {
	code := generateRandomString(8)
	expiresIn := 5 * time.Minute

	err := s.repo.Store(code, username, expiresIn)
	if err != nil {
		return "", err
	}

	log.Printf("Step 2: IdP Login - User: %s, Auth Code: %s\n", username, code)
	return code, nil
}

// ValidateAuthCode validates an authorization code and returns the username
func (s *AuthService) ValidateAuthCode(code string) (string, error) {
	authCode, exists := s.repo.Get(code)
	if !exists {
		log.Printf("Step 3: Token Exchange Failed - Invalid auth code: %s\n", code)
		return "", ErrInvalidAuthCode
	}

	if authCode.IsExpired() {
		log.Printf("Step 3: Token Exchange Failed - Expired auth code: %s\n", code)
		s.repo.Delete(code)
		return "", ErrAuthCodeExpired
	}

	// Delete code after successful validation (single use)
	s.repo.Delete(code)

	log.Printf("Step 3: Token Exchange Success - User: %s\n", authCode.Username)
	return authCode.Username, nil
}
