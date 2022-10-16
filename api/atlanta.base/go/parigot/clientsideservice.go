package parigot

import (
	"fmt"
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
	respBlob, err := abi.Dispatch(int64(b.svc), method, blob)
	if err != nil {
		return NewFromError("dispatch infrastructure (not dispatch itself) failed", err)
	}
	err = proto.Unmarshal(respBlob, out)
	if err != nil {
		id := NewDispatchErrorFromBytes(respBlob)
		return NewError(fmt.Sprintf("dispatch failed with code: %0d", int64(id)))
	}
	return nil
}
