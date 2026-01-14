package domain

// Token represents the OAuth 2.0 token response
type Token struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	Sub   string `json:"sub"`   // Subject (user ID)
	Email string `json:"email"` // Email address
	Name  string `json:"name"`  // Full name
	Iat   int64  `json:"iat"`   // Issued at timestamp
	Exp   int64  `json:"exp"`   // Expiration timestamp
}
