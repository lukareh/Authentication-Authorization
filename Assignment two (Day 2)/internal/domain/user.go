package domain

// User represents a user in the system
type User struct {
	Username string
	Password string
	Email    string
	Name     string
}

// Credentials represents login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate checks if credentials are valid (non-empty)
func (c *Credentials) Validate() bool {
	return c.Username != "" && c.Password != ""
}
