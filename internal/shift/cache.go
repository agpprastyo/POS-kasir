package shift

import (
	"POS-kasir/pkg/cache"
	"fmt"

	"github.com/google/uuid"
)

// Cache handles caching for shift status.
type Cache struct {
	cache cache.Cache
}

// NewCache results a new instance of Shift Cache.
func NewCache(c cache.Cache) *Cache {
	return &Cache{
		cache: c,
	}
}

func (c *Cache) key(userID uuid.UUID) string {
	return fmt.Sprintf("shift:open:%s", userID.String())
}

// SetOpen marks the user as having an open shift or not.
func (c *Cache) SetOpen(userID uuid.UUID, open bool) {
	val := "0"
	if open {
		val = "1"
	}
	c.cache.Set(c.key(userID), []byte(val), 0)
}

// GetOpen checks if the user has an open shift.
// Returns (isOpen, found).
func (c *Cache) GetOpen(userID uuid.UUID) (bool, bool) {
	valBytes, err := c.cache.Get(c.key(userID))
	if err != nil || valBytes == nil {
		return false, false
	}
	isOpen := string(valBytes) == "1"
	return isOpen, true
}

// Clear removes the user's shift status from cache.
func (c *Cache) Clear(userID uuid.UUID) {
	c.cache.Delete(c.key(userID))
}
