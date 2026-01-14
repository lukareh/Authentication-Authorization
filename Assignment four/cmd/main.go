package main

import (
	"fmt"
	"log"
	"net/http"
	"sso-mock/internal/handlers"
	"sso-mock/internal/repository"
	"sso-mock/internal/services"
)

func main() {
	// Initialize repository
	authRepo := repository.NewAuthCodeRepository()

	// Initialize services
	authService := services.NewAuthService(authRepo)
	tokenService := services.NewTokenService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	// Setup routes
	setupRoutes(authHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}

func setupRoutes(handler *handlers.AuthHandler) {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/api/auth/redirect", handler.HandleRedirect)
	http.HandleFunc("/api/auth/login", handler.HandleLogin)
	http.HandleFunc("/api/auth/token", handler.HandleTokenExchange)
	http.HandleFunc("/api/auth/verify", handler.HandleVerifyToken)
}
