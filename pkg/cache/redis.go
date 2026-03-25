package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new instance of RedisCache implementing the Cache interface
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves an item from the cache. Returns nil, nil if the key does not exist.
func (r *RedisCache) Get(key string) ([]byte, error) {
	return r.GetWithContext(context.Background(), key)
}

// GetWithContext retrieves an item from the cache with a context.
func (r *RedisCache) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set adds or updates an item in the cache. Set exp to 0 for no expiration.
func (r *RedisCache) Set(key string, val []byte, exp time.Duration) error {
	return r.SetWithContext(context.Background(), key, val, exp)
}

// SetWithContext adds or updates an item in the cache with a context.
func (r *RedisCache) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	return r.client.Set(ctx, key, val, exp).Err()
}

// Delete removes an item from the cache.
func (r *RedisCache) Delete(key string) error {
	return r.DeleteWithContext(context.Background(), key)
}

// DeleteWithContext removes an item from the cache with a context.
func (r *RedisCache) DeleteWithContext(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Reset clears the current database.
func (r *RedisCache) Reset() error {
	return r.ResetWithContext(context.Background())
}

// ResetWithContext clears the current database with a context.
func (r *RedisCache) ResetWithContext(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Close closes the underlying Redis client connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Exists checks if a key exists in the cache.
func (r *RedisCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	val, err := r.client.Exists(ctx, key).Result()
	return val > 0, err
}
