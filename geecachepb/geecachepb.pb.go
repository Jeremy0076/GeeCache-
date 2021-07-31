package geecachepb


type Request struct {
	Group string   `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	Key   string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

type Response struct {
	Value []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}