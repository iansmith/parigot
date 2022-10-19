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
			return NewErrorFromError(fmt.Sprintf("failed to marshal request in dispatch to '%s'", method), err)
		}
	}
	dispatchResp := abi.Dispatch(int64(b.svc), method, blob)
	if dispatchResp.GetErrorCode() != 0 {
		return NewErrorFromId("dispatch failed on method "+method, AnyId(NewDispatchErrorId(dispatchResp.GetErrorCode())))
	}
	err := proto.Unmarshal(dispatchResp.GetBlob(), out)
	if err != nil {
		return NewErrorFromError(fmt.Sprintf("failed to unmarshal response to dispatch of '%s'", method), err)
	}
	return nil
}
