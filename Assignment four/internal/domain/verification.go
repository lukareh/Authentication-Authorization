package domain

// VerificationCheck represents a single verification check result
type VerificationCheck struct {
	Check   string `json:"check"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// VerificationResult represents the complete verification result
type VerificationResult struct {
	Verified bool                 `json:"verified"`
	Checks   []VerificationCheck  `json:"checks"`
	Message  string               `json:"message"`
}

// NewVerificationResult creates a new verification result
func NewVerificationResult() *VerificationResult {
	return &VerificationResult{
		Verified: true,
		Checks:   make([]VerificationCheck, 0),
	}
}

// AddCheck adds a verification check to the result
func (vr *VerificationResult) AddCheck(check, status, message string) {
	vr.Checks = append(vr.Checks, VerificationCheck{
		Check:   check,
		Status:  status,
		Message: message,
	})
	if status == "failed" {
		vr.Verified = false
	}
}

// SetMessage sets the final verification message
func (vr *VerificationResult) SetMessage() {
	if vr.Verified {
		vr.Message = "Token Verified Successfully"
	} else {
		vr.Message = "Invalid Token"
	}
}
