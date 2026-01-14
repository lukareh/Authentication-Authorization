package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"authentication/services"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService *services.AuthService
	otpService  *services.OTPService
}

func NewAuthHandler(authService *services.AuthService, otpService *services.OTPService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		otpService:  otpService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyOTPRequest struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Log incoming registration request
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	fmt.Println("\nREGISTER REQUEST: ")
	fmt.Printf("Username: %s\n", req.Username)
	fmt.Printf("Password (received): %s\n", req.Password)
	fmt.Printf("Password (bcrypt hash): %s\n", string(hashedPwd))
	fmt.Println("-------------------")

	err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "User registered successfully",
	})
}

// Login handles user login and generates OTP
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Log incoming login request
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	fmt.Println("\nLOGIN REQUEST: ")
	fmt.Printf("Username: %s\n", req.Username)
	fmt.Printf("Password (received): %s\n", req.Password)
	fmt.Printf("Password (bcrypt hash): %s\n", string(hashedPwd))
	fmt.Println("\n")

	// Verify password
	err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		sendResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Generate OTP
	otp := h.otpService.GenerateOTP(req.Username)
	
	// Print OTP to terminal (server console)
	fmt.Printf("OTP for user '%s': %s\n", req.Username, otp)

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Password verified. OTP sent to terminal.",
	})
}

// VerifyOTP handles OTP verification
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Log incoming OTP verification request
	fmt.Println("\nVERIFY OTP REQUEST: ")
	fmt.Printf("Username: %s\n", req.Username)
	fmt.Printf("OTP: %s\n", req.OTP)

	// Verify OTP
	err := h.otpService.ValidateOTP(req.Username, req.OTP)
	if err != nil {
		sendResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Login successful!",
	})
}

func sendResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
