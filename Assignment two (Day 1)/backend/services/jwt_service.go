package services

import (
	"fmt"
	"jwt-auth-system/backend/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT operations
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a new JWT token for a given user
func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	
	claims := &domain.Claims{
		Username:    user.Username,
		Role:        user.Role,
		Designation: user.Designation,
		Age:         user.Age,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	// Print to console
	fmt.Printf("Generated Token: %s\n", tokenString)
	fmt.Printf("For User: %s\n", user.Username)
	fmt.Printf("Claims: {\"username\": \"%s\", \"role\": \"%s\", \"designation\": \"%s\", \"age\": %d}\n", 
		user.Username, user.Role, user.Designation, user.Age)
	fmt.Println("---")

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		fmt.Println("Token Invalid:", err.Error())
		return nil, err
	}

	if !token.Valid {
		fmt.Println("Token Invalid: token is not valid")
		return nil, fmt.Errorf("invalid token")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		fmt.Println("Token Invalid: token expired")
		return nil, fmt.Errorf("token expired")
	}

	// Print to console
	fmt.Println("Token Valid")
	fmt.Printf("Claims: {\"username\": \"%s\", \"role\": \"%s\", \"designation\": \"%s\", \"age\": %d}\n", 
		claims.Username, claims.Role, claims.Designation, claims.Age)
	fmt.Println("---")

	return claims, nil
}
