package lib

import (
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"google.golang.org/protobuf/proto"
)

func NewFromUnmarshal(b []byte) (*protosupport.Pctx, error) {
	p := protosupport.Pctx{}
	err := proto.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
