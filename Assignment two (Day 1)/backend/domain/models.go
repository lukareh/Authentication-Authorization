package domain

import "github.com/golang-jwt/jwt/v5"

// User represents a registered user in the system
type User struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	Designation string `json:"designation"`
	Age         int    `json:"age"`
}

// Claims represents the JWT claims structure
type Claims struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	Designation string `json:"designation"`
	Age         int    `json:"age"`
	jwt.RegisteredClaims
}

// RegisterRequest represents the request to register a user
type RegisterRequest struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	Designation string `json:"designation"`
	Age         int    `json:"age"`
}

// RegisterResponse represents the response after registering a user
type RegisterResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// GenerateRequest represents the request to generate a JWT
type GenerateRequest struct {
	Username string `json:"username"`
}

// GenerateResponse represents the response after generating a JWT
type GenerateResponse struct {
	Token string `json:"token"`
}

// ValidateRequest represents the request to validate a JWT
type ValidateRequest struct {
	Token string `json:"token"`
}

// ValidateResponse represents the response after validating a JWT
type ValidateResponse struct {
	Valid   bool        `json:"valid"`
	Message string      `json:"message"`
	Claims  interface{} `json:"claims,omitempty"`
}

// ClaimsData represents the claims data returned in validation response
type ClaimsData struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	Designation string `json:"designation"`
	Age         int    `json:"age"`
}
