package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Document represents a document in the system
type Document struct {
	ID      int    `json:"id"`
	Owner   string `json:"owner"`
	Content string `json:"content"`
}

// User represents a user with a role
type User struct {
	Username string `json:"username"`
	Role     string `json:"role"` // "user" or "admin"
}

// AccessRequest represents the request body for document access
type AccessRequest struct {
	Username   string `json:"username"`
	DocumentID int    `json:"documentId"`
}

// AccessResponse represents the response for document access
type AccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Content string `json:"content,omitempty"`
}

var (
	documents []Document
	users     map[string]User
)

func init() {
	// Initialize sample documents
	documents = []Document{
		{ID: 1, Owner: "userA", Content: "A's secret"},
		{ID: 2, Owner: "userB", Content: "B's secret"},
		{ID: 3, Owner: "userA", Content: "Another secret from A"},
	}

	// Sample users with roles
	users = map[string]User{
		"userA": {Username: "userA", Role: "user"},
		"userB": {Username: "userB", Role: "user"},
		"admin": {Username: "admin", Role: "admin"},
	}
}

func main() {
	// Enable CORS
	http.HandleFunc("/api/access", corsMiddleware(handleAccess))

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleAccess handles document access requests
func handleAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, AccessResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Verify user exists
	user, exists := users[req.Username]
	if !exists {
		respondJSON(w, AccessResponse{
			Success: false,
			Message: "User not found!",
		})
		return
	}

	// Find the document
	var document *Document
	for i := range documents {
		if documents[i].ID == req.DocumentID {
			document = &documents[i]
			break
		}
	}

	// Check if document exists
	if document == nil {
		respondJSON(w, AccessResponse{
			Success: false,
			Message: "Document not found!",
		})
		return
	}

	// Check authorization
	if canAccessDocument(user, document) {
		respondJSON(w, AccessResponse{
			Success: true,
			Message: "Access Granted!",
			Content: document.Content,
		})
	} else {
		respondJSON(w, AccessResponse{
			Success: false,
			Message: "Access Denied",
		})
	}
}

// canAccessDocument checks if a user can access a document
func canAccessDocument(user User, document *Document) bool {
	// Admins can access all documents
	if user.Role == "admin" {
		return true
	}

	// Regular users can only access their own documents
	return user.Username == document.Owner
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data AccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
