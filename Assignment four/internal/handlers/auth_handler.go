package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sso-mock/internal/domain"
	"sso-mock/internal/services"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService  *services.AuthService
	tokenService *services.TokenService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
	}
}

// HandleRedirect handles the initial redirect to IdP
func (h *AuthHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Step 1: Redirect to Identity Provider")

	response := map[string]string{
		"message":      "Redirecting to Identity Provider...",
		"redirect_url": "/login.html",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleLogin handles user login and generates auth code
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds domain.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate credentials
	if !h.authService.ValidateCredentials(&creds) {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Generate authorization code
	authCode, err := h.authService.GenerateAuthCode(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate auth code", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":    "Credentials validated successfully",
		"auth_code":  authCode,
		"expires_in": 300,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleTokenExchange handles auth code exchange for tokens
func (h *AuthHandler) HandleTokenExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tokenReq struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate auth code
	username, err := h.authService.ValidateAuthCode(tokenReq.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate tokens
	tokens := h.tokenService.GenerateTokens(username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

// HandleVerifyToken handles token verification
func (h *AuthHandler) HandleVerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var verifyReq struct {
		IDToken     string `json:"id_token"`
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&verifyReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify tokens
	result := h.tokenService.VerifyToken(verifyReq.IDToken, verifyReq.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
