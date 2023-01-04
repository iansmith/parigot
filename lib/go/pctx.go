package lib

import (
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"

	"google.golang.org/protobuf/proto"
)

func NewFromUnmarshal(b []byte) (*protosupportmsg.Pctx, error) {
	p := protosupportmsg.Pctx{}
	err := proto.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
