package cache

import (
	"GeeCache/byteview"
	"GeeCache/lru"
)

func (c *Cache) Add(key string, value byteview.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.CacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value byteview.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(byteview.ByteView), ok
	}
	return
}
