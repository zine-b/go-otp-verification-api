package domain

import (
	"errors"
)

var (
	ErrPhoneRequired 	 	= errors.New("phone number is required")
	ErrInvalidCode       	= errors.New("invalid otp code")
	ErrVerificationNotFound = errors.New("verification not found")

	// Maximum 3 OTP par numéro toutes les 10 minutes.
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
)