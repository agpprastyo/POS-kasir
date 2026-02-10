package shift

import (
	"POS-kasir/pkg/cache"
	"fmt"

	"github.com/google/uuid"
)

// Cache handles caching for shift status.
type Cache struct {
	memCache *cache.MemoryCache
}

// NewCache results a new instance of Shift Cache.
func NewCache(memCache *cache.MemoryCache) *Cache {
	return &Cache{
		memCache: memCache,
	}
}

func (c *Cache) key(userID uuid.UUID) string {
	return fmt.Sprintf("shift:open:%s", userID.String())
}

// SetOpen marks the user as having an open shift or not.
func (c *Cache) SetOpen(userID uuid.UUID, open bool) {
	// We only store "true" for open shifts.
	// If closed, we delete the key to save memory, or strictly store false.
	// The validation logic usually checks "is present and true".
	// Let's store the boolean value for explicit clarity.
	c.memCache.Set(c.key(userID), open)
}

// GetOpen checks if the user has an open shift.
// Returns (isOpen, found).
func (c *Cache) GetOpen(userID uuid.UUID) (bool, bool) {
	val, found := c.memCache.Get(c.key(userID))
	if !found {
		return false, false
	}
	isOpen, ok := val.(bool)
	if !ok {
		return false, false
	}
	return isOpen, true
}

// Clear removes the user's shift status from cache.
func (c *Cache) Clear(userID uuid.UUID) {
	c.memCache.Delete(c.key(userID))
}
