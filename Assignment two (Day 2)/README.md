Task: Build a mock SSO flow with authorization-code exchange
Write a program that:
Simulates redirect step
Print: "Redirecting to Identity Provider..."

Simulates IdP Login
User enters username + password
Program pretends login is valid
Program prints: "Auth Code: XYZ123"

Simulates Token Exchange
User enters the auth code
Program returns a mock ID token + access token, e.g.

{ id_token: "abc.id.sig", access_token: "xyz.access.sig" }

Simulates Token Verification
Program verifies:
Token is not empty
Token has valid structure
Token contains expected claims
Print "Token Verified" or "Invalid Token"

Implemented Features ==>

- OAuth 2.0 Authorization Code grant flow implementation
- RESTful API backend server with Go
- JWT-format token generation (ID token and Access token)
- Authorization code generation with expiration and single-use validation
- Token verification with claim validation (sub, email, name, iat, exp)
- Web-based frontend with step-by-step SSO flow visualization
- In-memory storage for authorization codes with cleanup
- CORS-enabled API for frontend-backend communication
- JWT-format token generation
- Authorization code validation
- Token structure verification
- Claims validation
- Comprehensive verification workflow

 Notes

- This is a mock implementation for educational purposes
- Credentials are not actually validated against a real identity provider
- Tokens are not cryptographically signed
- All authorization codes and tokens are time-stamped but expiration is not enforced
