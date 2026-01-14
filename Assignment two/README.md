JWT Authentication System

A mini JWT-based authentication system built with Go using clean architecture. The system generates JWT tokens with claims `{"sub": "<user-id>", "role": "user"}`, signs them using HS256 algorithm with a secret key, and validates tokens by verifying signatures and checking expiration (5-minute expiry). Includes a simple HTML frontend for testing and a Postman collection for API testing.

API Endpoints
POST /generate - Generate JWT token
Request: { "user_id": "user123" }
Response: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }

POST /validate - Validate JWT token
Request: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
Response: { "valid": true, "message": "Token Valid", "claims": {"sub": "user123", "role": "user"} }

Running
go run backend/cmd/main.go

Server starts on `http://localhost:8080`
