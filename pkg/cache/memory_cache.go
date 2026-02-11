package cache

import (
	"sync"
)

// MemoryCache is a thread-safe generic in-memory cache.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

// NewMemoryCache creates a new instance of MemoryCache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]interface{}),
	}
}

// Set adds or updates an item in the cache.
func (c *MemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// Get retrieves an item from the cache.
// Returns (value, found).
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, found := c.items[key]
	return val, found
}

// Delete removes an item from the cache.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]interface{})
}
