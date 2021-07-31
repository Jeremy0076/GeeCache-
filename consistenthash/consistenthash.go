package consistenthash

type Hash func(data []byte) uint32

// Map
//  hash:哈希函数
//  replicas:复制倍数
//  keys:哈希环 机器序列
//  hashMap:虚拟节点与真实节点的映射表
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}


