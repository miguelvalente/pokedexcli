package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mutex sync.Mutex
}

func NewCache(inteval time.Duration) Cache {
	return Cache{}
}

func (c *Cache) Add(key string, val_ []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val_,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	found, exists := c.cache[key]
	if found, exists := c.cache[key]; exists {
		return found.val, true
	}
	return nil, false
}
