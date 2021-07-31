package httppool

import (
	"GeeCache/consistenthash"
	"sync"
)

const (
	defaultBasePath = "_/geeCache/"
	defaultReplicas = 3
)

type server int

// HTTPPool 承载 httppool 通信的核心数据结构 self: self[Url]
// example : "http://example.com/_geecache/"
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex // guards peers and httpGetters
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter //keyed by e.g "http:10.0.0.2:8080"
}

// httpGetter http 客户端
type httpGetter struct {
	baseURL string
}
