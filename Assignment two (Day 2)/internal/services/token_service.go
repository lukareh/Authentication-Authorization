package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"sso-mock/internal/domain"
	"strings"
	"time"
)

// TokenService handles token generation and verification
type TokenService struct{}

// NewTokenService creates a new token service
func NewTokenService() *TokenService {
	return &TokenService{}
}

// GenerateTokens creates ID and access tokens for a user
func (s *TokenService) GenerateTokens(username string) *domain.Token {
	claims := domain.TokenClaims{
		Sub:   "user_" + generateRandomString(6),
		Email: username + "@example.com",
		Name:  username,
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(1 * time.Hour).Unix(),
	}

	idToken := s.createJWT(claims, "id_token")
	accessToken := s.createJWT(claims, "access_token")

	return &domain.Token{
		IDToken:     idToken,
		AccessToken: accessToken,
		ExpiresIn:   3600,
		TokenType:   "Bearer",
	}
}

// VerifyToken verifies token structure and claims
func (s *TokenService) VerifyToken(idToken, accessToken string) *domain.VerificationResult {
	result := domain.NewVerificationResult()

	// Check 1: Tokens not empty
	if idToken == "" || accessToken == "" {
		result.AddCheck("Token not empty", "failed", "Token is empty")
		result.SetMessage()
		log.Println("Step 4: Token Verification Failed - Empty tokens")
		return result
	}
	result.AddCheck("Token not empty", "passed", "")

	// Check 2: ID Token JWT structure
	if !s.isValidJWTStructure(idToken) {
		result.AddCheck("ID Token JWT structure", "failed", "Invalid JWT structure")
	} else {
		result.AddCheck("ID Token JWT structure", "passed", "")
	}

	// Check 3: Access Token JWT structure
	if !s.isValidJWTStructure(accessToken) {
		result.AddCheck("Access Token JWT structure", "failed", "Invalid JWT structure")
	} else {
		result.AddCheck("Access Token JWT structure", "passed", "")
	}

	// Check 4: Token claims
	if !s.verifyTokenClaims(idToken) {
		result.AddCheck("Token claims validation", "failed", "Missing or invalid claims")
	} else {
		result.AddCheck("Token claims validation", "passed", "")
	}

	result.SetMessage()
	log.Printf("Step 4: Token Verification - Overall Status: %v\n", result.Verified)

	return result
}

// createJWT creates a mock JWT token
func (s *TokenService) createJWT(claims domain.TokenClaims, tokenType string) string {
	// Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Payload
	payload := map[string]interface{}{
		"sub":   claims.Sub,
		"email": claims.Email,
		"name":  claims.Name,
		"iat":   claims.Iat,
		"exp":   claims.Exp,
		"type":  tokenType,
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Signature (mock)
	signature := generateRandomString(43)

	return headerB64 + "." + payloadB64 + "." + signature
}

// isValidJWTStructure checks if token has valid JWT format
func (s *TokenService) isValidJWTStructure(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	// Verify each part can be decoded from base64
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		// Try to decode with padding
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			// Try with padding
			padded := part
			switch len(part) % 4 {
			case 2:
				padded += "=="
			case 3:
				padded += "="
			}
			if _, err := base64.StdEncoding.DecodeString(padded); err != nil {
				return false
			}
		}
	}

	return true
}

// verifyTokenClaims verifies token claims are present and valid
func (s *TokenService) verifyTokenClaims(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	// Decode payload
	payloadB64 := parts[1]
	padding := 4 - len(payloadB64)%4
	if padding != 4 {
		payloadB64 += strings.Repeat("=", padding)
	}

	payloadJSON, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		return false
	}

	var claims domain.TokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return false
	}

	// Verify required claims exist
	if claims.Sub == "" || claims.Email == "" || claims.Name == "" {
		return false
	}

	if claims.Exp == 0 {
		return false
	}

	return true
}

// generateRandomString creates a random string for codes/signatures
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}
