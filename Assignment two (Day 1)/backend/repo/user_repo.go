package repo

import (
	"errors"
	"jwt-auth-system/backend/domain"
	"sync"
)

// UserRepository handles user storage operations
type UserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

// NewUserRepository creates a new user repository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// RegisterUser stores a user in memory
func (r *UserRepository) RegisterUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return errors.New("user already exists")
	}

	r.users[user.Username] = user
	return nil
}

// GetUser retrieves a user by username
func (r *UserRepository) GetUser(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UserExists checks if a user exists
func (r *UserRepository) UserExists(username string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[username]
	return exists
}

// GetAllUsers returns all registered users (for debugging)
func (r *UserRepository) GetAllUsers() map[string]*domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	usersCopy := make(map[string]*domain.User)
	for k, v := range r.users {
		usersCopy[k] = v
	}
	return usersCopy
}
