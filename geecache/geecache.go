package geecache

import (
	"GeeCache/cache"
	"GeeCache/peers"
	"GeeCache/singleflight"
)

// Getter 获取接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数
type GetterFunc func(key string) ([]byte, error)

// Get 接口型函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group cache namespace and associate loaded data
type Group struct {
	name      string
	getter    Getter
	mainCache cache.Cache
	peers     peers.PeerPicker
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}
