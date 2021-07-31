package peers

import pb "GeeCache/geecachepb"

// PeerPicker 根据key查找对应的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter HTTP客户端, 根据group, key 返回缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}