package domain

import (
	"time"
)

type OTPStatus string

const (
	OTPStatusPending  OTPStatus = "pending"
	OTPStatusVerified OTPStatus = "verified"
	OTPStatusExpired  OTPStatus = "expired"
)

type OTPVerification struct{
	ID        string
	Phone     string
	Code      string
	Status    OTPStatus
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewOTPVerification(id string, phone string, code string, ttl time.Duration) *OTPVerification{
	now := time.Now()

	return &OTPVerification{
		ID: 	   id,
		Phone: 	   phone,
		Code: 	   code,
		Status:    OTPStatusPending,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}

}

func (o *OTPVerification) Verify(code string) bool {
	if time.Now().After(o.ExpiresAt){
		o.Status=OTPStatusExpired
		return false
	}
	if code !=o.Code{
		return false
	}
	
	o.Status = OTPStatusVerified
	return true

}