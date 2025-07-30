package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	reapInterval time.Duration
	CacheMap     map[string]cacheEntry
	mu           sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	return &Cache{
		reapInterval: interval,
		CacheMap:     make(map[string]cacheEntry),
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.CacheMap[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.CacheMap[key]
	if exists {
		return entry.val, true
	}
	return nil, false
}

// func (c *Cache) reapLoop() {

// }
