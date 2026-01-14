package domain

import "time"

// AuthCode represents an authorization code with associated metadata
type AuthCode struct {
	Code      string
	Username  string
	ExpiresAt time.Time
}

// IsExpired checks if the authorization code has expired
func (a *AuthCode) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// IsValid checks if the authorization code is valid
func (a *AuthCode) IsValid() bool {
	return a.Code != "" && !a.IsExpired()
}
