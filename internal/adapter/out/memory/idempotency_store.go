package memory

import (
	"context"
	"sync"
)
type IdempotencyStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		store: make(map[string]string),
	}
}

func (s *IdempotencyStore) Save(ctx context.Context, key string, verificationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = verificationID
	return nil
}

func (s *IdempotencyStore) Find(ctx context.Context, key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	verificationID, ok := s.store[key]
	return verificationID, ok, nil
}