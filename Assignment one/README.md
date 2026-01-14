Task: Build a minimal login system with secure password hashing
Write a program in any language that does the following:

Register User

Input: username + password

Hash the password using bcrypt or argon2
Store {username, hashedPassword} in a simple in-memory array or file
Login User
Input: username + password

Fetch the stored hash
Verify using the hashing library
Print "Login Successful" or "Invalid Credentials"

Add MFA Layer (basic simulation)
After password is correct, generate a 6-digit OTP
Print it on console (pretend it's an SMS)
User enters OTP
Validate OTP and allow login

==>
How to Use

1. Start the backend server by running `go run cmd/main.go` from the `backend/` directory
2. Open `http://localhost:8080` in your browser to access the web interface
3. Register a new user with username and password (password is hashed with bcrypt)
4. Login with credentials - OTP will be displayed in the server terminal
5. Enter the OTP from terminal to complete authentication and login successfully

API Endpoints

- POST `/api/register` - Register new user
- POST `/api/login` - Verify credentials and generate OTP (displayed in terminal)
- POST `/api/verify-otp` - Validate OTP and complete login

Import `Postman_Collection.json` into Postman to test the API endpoints.