package parigot

import (
	"github.com/iansmith/parigot/g/parigot/abi"
	"google.golang.org/protobuf/proto"
)

type ClientSideService struct {
	svc ServiceId
}

func NewClientSideService(svc ServiceId) *ClientSideService {
	return &ClientSideService{
		svc: svc,
	}
}
func (b *ClientSideService) Dispatch(method string, in proto.Message, out proto.Message) Error {
	var blob []byte
	if in != nil {
		var err error
		blob, err = proto.Marshal(in)
		if err != nil {
			return NewFromError("unable to marshall input parameter in dispatch", err)
		}
	}
	respBlob := abi.Dispatch(int64(b.svc), method, blob)
	err := proto.Unmarshal(respBlob, out)
	if err != nil {
		id := NewDispatchErrorFromBytes(respBlob)
		return NewErrorFromId("dispatch failed", AnyId(id))
	}
	return nil
}
