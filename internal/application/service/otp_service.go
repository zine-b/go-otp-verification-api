package service

import (
	"context"
	"fmt"
	"time"
	"math/rand"
	"prepareGo/internal/application/port/in"
	"prepareGo/internal/application/port/out"
	"prepareGo/internal/domain"
)

type OTPService struct{
	provider      out.OTPProvider
	repository    out.OTPRepository
	rateLimiter   out.RateLimiter
	idempotencyStore out.IdempotencyStore
	generateCode func() string
}

func NewOTPService(provider out.OTPProvider, repository out.OTPRepository, rateLimiter out.RateLimiter,idempotencyStore out.IdempotencyStore, generateCode func() string) *OTPService{
	return &OTPService{
		provider: provider,
		repository: repository,
		rateLimiter: rateLimiter,
		idempotencyStore: idempotencyStore,
		generateCode: generateCode,

	}
}

func (s *OTPService)SendOTP(ctx context.Context, cmd in.SendOTPCommand) (*in.SendOTPResult, error){
	if cmd.Phone == "" {
		return nil, domain.ErrPhoneRequired
	}

	idempotencyKey := buildIdempotencyKey(cmd.Phone, cmd.IdempotencyKey)
	fmt.Printf("raw idempotency key=%q, built key=%q\n", cmd.IdempotencyKey, idempotencyKey)
	
	if idempotencyKey != "" {
		verificationID, found, err := s.idempotencyStore.Find(ctx, idempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("find idempotency key: %w", err)
		}

		if found {
			verification, err := s.repository.FindByID(ctx, verificationID)
			if err != nil {
				return nil, fmt.Errorf("find idempotent verification: %w", err)
			}

			return &in.SendOTPResult{
				VerificationID: verification.ID,
				Status:         string(verification.Status),
				Message:        "otp already sent",
			}, nil
		}
	}

	allowed, err := s.rateLimiter.Allow(ctx, cmd.Phone)
	if err != nil {
		return nil, fmt.Errorf("check rate limit: %w", err)
	}

	if !allowed {
		return nil, domain.ErrRateLimitExceeded
	}


	//generer le code
	//code := s.generateCode()
	code:="123456"

	// je creer mon objet domain à partir de l'entree
	verification := domain.NewOTPVerification(
		generateID(),
		cmd.Phone,
		code,
		5*time.Minute,
	)


	//appel la methode send du provider out
	if err:=s.provider.Send(ctx, verification.Phone, verification.Code); err != nil {
		return nil, fmt.Errorf("send otp failed: %w", err)
	}
	
	if err := s.repository.Save(ctx, verification); err != nil {
		return nil, fmt.Errorf("save otp verification: %w", err)
	}

	if idempotencyKey != "" {
		if err := s.idempotencyStore.Save(ctx, idempotencyKey, verification.ID); err != nil {
			return nil, fmt.Errorf("save idempotency key: %w", err)
		}
	}
	
	return &in.SendOTPResult{
		VerificationID: verification.ID,
		Status: string(verification.Status),
		Message: "otp sent",
	},nil
}

func (s *OTPService) VerifyOTP(ctx context.Context, cmd in.VerifyOTPCommand) error {
	verification, err := s.repository.FindByID(ctx, cmd.VerificationID)
	if err != nil {
		return err
	}

	if ok := verification.Verify(cmd.Code); !ok {
		return domain.ErrInvalidCode
	}

	if err := s.repository.Save(ctx, verification); err != nil {
		return fmt.Errorf("save verified otp: %w", err)
	}

	return nil
}

func GenerateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func generateID() string {
	return fmt.Sprintf("verif_%d", time.Now().UnixNano())
}

func buildIdempotencyKey(phone string, idempotencyKey string) string {
	if idempotencyKey == "" {
		return ""
	}

	return phone + ":" + idempotencyKey
}