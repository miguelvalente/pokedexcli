package pokecache

import (
	"fmt"
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

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val_ []byte) {
	// fmt.Println("Added: ", key)
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

	if found, exists := c.cache[key]; exists {
		return found.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.mutex.Lock()

		for key, value := range c.cache {
			elapsed := time.Since(value.createdAt)
			if elapsed > interval {
				delete(c.cache, key)
				fmt.Printf("Removed key: %s, elapsed time: %s\n", key, elapsed)
			}
		}
		c.mutex.Unlock()
	}
}
