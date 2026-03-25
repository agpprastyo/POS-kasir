package cache

import (
	"context"
	"time"
)

// Cache defines the generic caching interface.
// It is designed to be fully compatible with the fiber.Storage interface.
type Cache interface {
	Get(key string) ([]byte, error)
	GetWithContext(ctx context.Context, key string) ([]byte, error)
	Set(key string, val []byte, exp time.Duration) error
	SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error
	Delete(key string) error
	DeleteWithContext(ctx context.Context, key string) error
	Reset() error
	ResetWithContext(ctx context.Context) error
	Close() error
	Exists(key string) (bool, error)
}
