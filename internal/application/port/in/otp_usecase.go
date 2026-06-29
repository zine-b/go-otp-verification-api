package in

import "context"
type SendOTPCommand struct{
	Phone          string
	Channel        string 
	IdempotencyKey string //requestId
}

type SendOTPResult struct{
	VerificationID string `json:"verification_id"` 
	Status         string `json:"status"`
	Message        string `json:"message"`
}

type VerifyOTPCommand struct{
	VerificationID string
	Code           string
}

type OTPUseCase interface{
	SendOTP(ctx context.Context, cmd SendOTPCommand) (*SendOTPResult, error)
	VerifyOTP(ctx context.Context, cmd VerifyOTPCommand) error 
}