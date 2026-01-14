package services

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP")
	ErrOTPExpired = errors.New("OTP has expired")
)

// OTPService handles OTP generation and validation
type OTPService struct {
	otpStore map[string]*OTPData
	mu       sync.RWMutex
}

// OTPData stores OTP information
type OTPData struct {
	Code      string
	ExpiresAt time.Time
}

// NewOTPService creates a new OTP service
func NewOTPService() *OTPService {
	return &OTPService{
		otpStore: make(map[string]*OTPData),
	}
}

// GenerateOTP generates a 6-digit OTP for a username
func (s *OTPService) GenerateOTP(username string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate random 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(rune('0' + rand.Intn(10)))
	}

	// Store OTP with 5 minute expiration
	s.otpStore[username] = &OTPData{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return otp
}

// ValidateOTP checks if the provided OTP is valid
func (s *OTPService) ValidateOTP(username, otp string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	otpData, exists := s.otpStore[username]
	if !exists {
		return ErrInvalidOTP
	}

	// Check if OTP has expired
	if time.Now().After(otpData.ExpiresAt) {
		delete(s.otpStore, username)
		return ErrOTPExpired
	}

	// Validate OTP
	if otpData.Code != otp {
		return ErrInvalidOTP
	}

	// Delete OTP after successful validation
	delete(s.otpStore, username)
	return nil
}
