package out

import (
	"context"
	"prepareGo/internal/domain"

)

type OTPRepository interface {
	Save(ctx context.Context, verification *domain.OTPVerification) error
	FindByID(ctx context.Context, id string) (*domain.OTPVerification, error)
}