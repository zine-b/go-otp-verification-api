package memory

import(
	"sync"
	"prepareGo/internal/domain"
	"context"

)

type OTPRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.OTPVerification
}

func NewOTPRepository() *OTPRepository {
	return &OTPRepository{
		store: make(map[string]*domain.OTPVerification),
	}
}

func (r *OTPRepository) Save(ctx context.Context, verification *domain.OTPVerification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[verification.ID] = verification
	return nil
}

func (r *OTPRepository) FindByID(ctx context.Context, id string) (*domain.OTPVerification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	verification, ok := r.store[id]
	if !ok {
		return nil, domain.ErrVerificationNotFound
	}

	return verification, nil
}