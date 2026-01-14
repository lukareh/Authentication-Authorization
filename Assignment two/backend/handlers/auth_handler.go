package handlers

import (
	"encoding/json"
	"fmt"
	"jwt-auth-system/backend/domain"
	"jwt-auth-system/backend/repo"
	"jwt-auth-system/backend/services"
	"net/http"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	jwtService *services.JWTService
	userRepo   *repo.UserRepository
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(jwtService *services.JWTService, userRepo *repo.UserRepository) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// RegisterUser handles user registration requests
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Role == "" || req.Designation == "" || req.Age <= 0 {
		http.Error(w, "All fields are required (username, role, designation, age)", http.StatusBadRequest)
		return
	}

	// Create user
	user := &domain.User{
		Username:    req.Username,
		Role:        req.Role,
		Designation: req.Designation,
		Age:         req.Age,
	}

	// Store user in repository
	if err := h.userRepo.RegisterUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	fmt.Printf("User Registered: %s (Role: %s, Designation: %s, Age: %d)\n", 
		user.Username, user.Role, user.Designation, user.Age)
	fmt.Println("---")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.RegisterResponse{
		Message:  "User registered successfully",
		Username: user.Username,
	})
}

// GenerateToken handles token generation requests
func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// Get user from repository
	user, err := h.userRepo.GetUser(req.Username)
	if err != nil {
		http.Error(w, "User not found. Please register first.", http.StatusNotFound)
		return
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.GenerateResponse{Token: token})
}

// ValidateToken handles token validation requests
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	claims, err := h.jwtService.ValidateToken(req.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.ValidateResponse{
			Valid:   false,
			Message: "Token Invalid: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.ValidateResponse{
		Valid:   true,
		Message: "Token Valid",
		Claims: domain.ClaimsData{
			Username:    claims.Username,
			Role:        claims.Role,
			Designation: claims.Designation,
			Age:         claims.Age,
		},
	})
}
