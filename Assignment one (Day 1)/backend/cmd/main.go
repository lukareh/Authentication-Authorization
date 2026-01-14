package main

import (
	"fmt"
	"log"
	"net/http"

	"authentication/handlers"
	"authentication/repository"
	"authentication/services"
)

func main() {
	// Initialize dependencies
	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	otpService := services.NewOTPService()
	authHandler := handlers.NewAuthHandler(authService, otpService)

	// Setup routes
	http.HandleFunc("/api/register", enableCORS(authHandler.Register))
	http.HandleFunc("/api/login", enableCORS(authHandler.Login))
	http.HandleFunc("/api/verify-otp", enableCORS(authHandler.VerifyOTP))

	// Serve static files (frontend)
	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	fmt.Println("  Secure Login System with MFA - API Server")
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("API Endpoints:")
	fmt.Println("  POST /api/register")
	fmt.Println("  POST /api/login")
	fmt.Println("  POST /api/verify-otp")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// enableCORS adds CORS headers to allow frontend requests
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
