package lib

import (
	glog "github.com/iansmith/parigot/g/pb/log"
	"google.golang.org/protobuf/proto"
)

type Pctx interface {
	Log() Log
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
	ToBytes() ([]byte, error)
}

type pctx struct {
	logger Log
	line   []*glog.LogRequest
	entry  map[string]string
}

func NewPctx() Pctx {
	line := []*glog.LogRequest{}
	entry := make(map[string]string)
	logger := NewProtoLogger(line)
	return &pctx{
		logger: logger,
		line:   line,
		entry:  entry,
	}
}

func NewPctxWithLog(l Log) Pctx {
	return &pctx{
		logger: l,
		entry:  make(map[string]string),
	}
}

func (p *pctx) ToBytes() ([]byte, error) {
	c := &glog.LogCollection{
		Req:   p.line,
		Entry: p.entry,
	}
	return proto.Marshal(c)
}

func (p *pctx) Log() Log {
	return p.logger
}

func (p *pctx) Entry(group, name string) (string, bool) {
	key := group + "." + name
	result, found := p.entry[key]
	return result, found
}
func (p *pctx) SetEntry(group, name, value string) bool {
	key := group + "." + name
	_, present := p.entry[key]
	p.entry[key] = value
	return present
}
