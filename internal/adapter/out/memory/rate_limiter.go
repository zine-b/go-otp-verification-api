package memory

import (
	"context"
	"sync"
	"time"
)

// Implementation du port out RateLimiter
type RateLimiter struct{
	mu       sync.Mutex
	limit    int
	window   time.Duration
	requests map[string][]time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter{
	return &RateLimiter{
		limit: limit,
		window: window,
		requests: make(map[string][]time.Time),
	}
}

// key c le numero de telephone 
func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// avant 10min de maintenant
	cutoff := now.Add(-r.window)

	timestamps := r.requests[key]

	valid := make([]time.Time, 0, len(timestamps))
	
	// je garde que les valeurs pendant les derniers 10min
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= r.limit {
		r.requests[key] = valid
		return false, nil
	}

	valid = append(valid, now)
	r.requests[key] = valid

	return true, nil


}