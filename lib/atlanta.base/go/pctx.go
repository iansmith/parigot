package lib

import (
	glog "github.com/iansmith/parigot/g/parigot/log"
	"github.com/iansmith/parigot/lib/interface_"
	"google.golang.org/protobuf/proto"
)

type pctx struct {
	logger interface_.Log
	line   []*glog.LogRequest
	entry  map[string]string
}

func NewPctx() interface_.Pctx {
	line := []*glog.LogRequest{}
	entry := make(map[string]string)
	logger := NewProtoLogger(line)
	return &pctx{
		logger: logger,
		line:   line,
		entry:  entry,
	}
}

func NewPctxWithLog(l interface_.Log) interface_.Pctx {
	return &pctx{
		logger: l,
		entry:  make(map[string]string),
	}
}

func (p *pctx) ToBytes() ([]byte, error) {
	c := &log.LogCollection{
		Req:   p.line,
		Entry: p.entry,
	}
	return proto.Marshal(c)
}

func (p *pctx) Log() interface_.Log {
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
