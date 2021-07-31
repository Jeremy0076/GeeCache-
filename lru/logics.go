package lru

import "container/list"

// NewCache :Init a Cache
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		MaxBytes:  maxBytes,
		OnEvicted: onEvicted,
		Ll:        list.New(),
		Cache:     make(map[string]*list.Element),
	}
}

// Get :get key-value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.Cache[key]; ok {
		c.Ll.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest :Remove the oldest Entry
func (c *Cache) RemoveOldest() {
	ele := c.Ll.Back()
	if ele != nil {
		c.Ll.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(c.Cache, kv.key)
		c.NBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add :Add or Update an entry
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.Cache[key]; ok {
		c.Ll.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		c.NBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		ele := c.Ll.PushFront(&Entry{key, value})
		c.Cache[key] = ele
		c.NBytes += int64(len(key)) + int64(value.Len())
	}
	for c.MaxBytes != 0 && c.MaxBytes < c.NBytes {
		c.RemoveOldest()
	}
}

// Len :return how many Entry in Cache
func (c *Cache) Len() int {
	return c.Ll.Len()
}
