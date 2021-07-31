package cache

import (
	"GeeCache/lru"
	"sync"
)

type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	CacheBytes int64
}
