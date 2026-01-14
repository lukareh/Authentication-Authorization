package main

import (
	"fmt"
	"jwt-auth-system/backend/handlers"
	"jwt-auth-system/backend/repo"
	"jwt-auth-system/backend/services"
	"log"
	"net/http"
)

// enableCORS is a middleware to enable CORS for all routes
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

func main() {
	// Initialize repository
	userRepo := repo.NewUserRepository()

	// Initialize services
	secretKey := "my-secret-key-change-in-production"
	jwtService := services.NewJWTService(secretKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtService, userRepo)

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/register", enableCORS(authHandler.RegisterUser))
	http.HandleFunc("/generate", enableCORS(authHandler.GenerateToken))
	http.HandleFunc("/validate", enableCORS(authHandler.ValidateToken))

	// Start server
	fmt.Println("JWT Authentication Server")
	fmt.Println("Server starting on http://localhost:8080")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
