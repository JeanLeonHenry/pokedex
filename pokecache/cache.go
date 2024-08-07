package pokecache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{interval: interval, data: make(map[string]cacheEntry), mu: &sync.Mutex{}}
	go cache.reapLoop()
	return cache
}

func (c Cache) String() (result string) {
	result = "{"
	for key, entry := range c.data {
		offset := strings.Split(strings.Split(key, "?")[1], "&")[0]
		result += fmt.Sprintln(offset, entry.createdAt.Format(time.DateTime))
	}
	return result + "}"
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{createdAt: time.Now(), val: val}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.data[key]
	return entry.val, ok
}
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		for key, entry := range c.data {
			if time.Now().Sub(entry.createdAt) > c.interval {
				c.mu.Lock()
				delete(c.data, key)
				c.mu.Unlock()
			}
		}
	}
}
