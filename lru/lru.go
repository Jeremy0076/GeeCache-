package lru

import "container/list"

type Cache struct {
	MaxBytes int64
	NBytes   int64
	Ll       *list.List
	Cache    map[string]*list.Element
	// Entry is perged func
	OnEvicted func(key string, value Value)
}

type Entry struct {
	key   string
	value Value
}

// Value must have Len() to calculate bytes
type Value interface {
	Len() int
}
