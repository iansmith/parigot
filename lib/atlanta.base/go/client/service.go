package client

import (
	"github.com/iansmith/parigot/lib/id"

	"google.golang.org/protobuf/proto"
)

type ClientSideService struct {
	svc id.Service
}

func NewClientSideService(svc id.Service) *ClientSideService {
	return &ClientSideService{
		svc: svc,
	}
}
func (b *ClientSideService) Dispatch(method string, in proto.Message, out proto.Message) error {
	return nil
}
