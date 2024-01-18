package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheEntry
	mu   sync.Mutex
}

type cacheEntry struct {
	createdAt int64
	value     []byte
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now().Unix(),
		value:     value,
	}
	c.data[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return data.value, true
}

func NewCache(interval int64) *Cache {
	cache := Cache{
		data: make(map[string]cacheEntry),
		// the mutex zero value is the "unlocked" state, sweet
	}
	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) reapLoop(interval int64) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		<-ticker.C
		c.mu.Lock()
		defer c.mu.Unlock()

		for key, entry := range c.data {
			if time.Now().Unix()-entry.createdAt >= interval {
				delete(c.data, key)
			}
		}
	}
}
