package domain

import "time"

// User represents a user in the system
type User struct {
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}

// NewUser creates a new user instance
func NewUser(username, hashedPassword string) *User {
	return &User{
		Username:       username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}
}
