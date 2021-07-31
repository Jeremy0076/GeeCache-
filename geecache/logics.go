package geecache

import (
	"GeeCache/byteview"
	"GeeCache/cache"
	pb "GeeCache/geecachepb"
	"GeeCache/peers"
	"GeeCache/singleflight"
	"fmt"
	"log"
	"sync"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create new group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.Cache{CacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup get a group by name
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get get value from mainCache
func (g *Group) Get(key string) (byteview.ByteView, error) {
	if key == "" {
		return byteview.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

// load from locally or peer
func (g *Group) load(key string) (value byteview.ByteView, err error) {
	// each key is only fetched once (either locally or remotely)
	// regardless of the number of concurrent callers.
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFormPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})

	if err != nil {
		return viewi.(byteview.ByteView), nil
	}
	return
}

// getLocally use callback func
func (g *Group) getLocally(key string) (byteview.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return byteview.ByteView{}, err
	}
	value := byteview.ByteView{B: byteview.CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache add to mainCache
func (g *Group) populateCache(key string, value byteview.ByteView) {
	g.mainCache.Add(key, value)
}

// RegsiterPeers registers a PeerPicker(HTTPPool) for choosing remote peer
func (g *Group) RegsiterPeers(peers peers.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFormPeer(peer peers.PeerGetter, key string) (byteview.ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return byteview.ByteView{}, err
	}
	return byteview.ByteView{B: res.Value}, nil
}
