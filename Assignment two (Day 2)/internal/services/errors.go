package services

import "errors"

var (
	// ErrInvalidAuthCode is returned when an authorization code is invalid
	ErrInvalidAuthCode = errors.New("invalid authorization code")

	// ErrAuthCodeExpired is returned when an authorization code has expired
	ErrAuthCodeExpired = errors.New("authorization code expired")

	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
)
