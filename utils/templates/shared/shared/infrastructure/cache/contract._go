package cache

import (
	"context"
	"time"
)

// Cache We only use 4 base function in redis
type Cache interface {

	// Set put the initial value
	Set(ctx context.Context, key string, value []byte, exp time.Duration) error

	// Reset update the existing value but keep continue the exp time
	Reset(ctx context.Context, key string, value []byte) error

	// Get the value
	Get(ctx context.Context, key string) (string, error)

	// Del Delete the value
	Del(ctx context.Context, key string) error
}