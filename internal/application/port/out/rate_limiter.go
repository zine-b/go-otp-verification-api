package out

import "context"

// pk ici dans ce package ?
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}