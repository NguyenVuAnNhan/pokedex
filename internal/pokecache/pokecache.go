package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	value []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	return entry.value, true
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(value) == 0 {
		return
	}

	c.cache[key] = cacheEntry{
		value: value,
		createdAt: time.Now(),
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		now := time.Now()
		for key, entry := range c.cache {
			if now.Sub(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}

		c.mu.Unlock()
	}
}