package out

import "context"

type IdempotencyStore interface {
	// abc-123 → verif_123 : idempotency_key --> verifkey
	Save(ctx context.Context, key string, verificationID string) error
	Find(ctx context.Context, key string) (string, bool, error)
}