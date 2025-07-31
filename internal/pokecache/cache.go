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
	c := &Cache{
		reapInterval: interval,
		CacheMap:     make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
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

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.CacheMap {
			t0 := entry.createdAt
			t1 := time.Now()
			elapsedTime := t1.Sub(t0)
			if elapsedTime >= interval {
				delete(c.CacheMap, key)
			}
		}
		c.mu.Unlock()
	}
}
