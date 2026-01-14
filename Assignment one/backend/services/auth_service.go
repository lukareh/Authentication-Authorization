package services

import (
	"errors"

	"authentication/domain"
	"authentication/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWeakPassword      = errors.New("password must be at least 6 characters")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register creates a new user with hashed password
func (s *AuthService) Register(username, password string) error {
	// Validate password strength
	if len(password) < 6 {
		return ErrWeakPassword
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := domain.NewUser(username, string(hashedPassword))

	// Save to repository
	return s.userRepo.Create(user)
}

// Login verifies user credentials
func (s *AuthService) Login(username, password string) error {
	// Fetch user from repository
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return ErrInvalidCredentials
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}
