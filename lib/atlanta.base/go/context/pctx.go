package context

import (
	"github.com/iansmith/parigot/g/parigot/log"
	logC "github.com/iansmith/parigot/lib/log"
	"google.golang.org/protobuf/proto"
)

type Pctx interface {
	Log() log.Log
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
	ToBytes() ([]byte, error)
}

type pctx struct {
	logger log.Log
	line   []*log.LogRequest
	entry  map[string]string
}

func NewPctx() Pctx {
	line := []*log.LogRequest{}
	entry := make(map[string]string)
	logger := logC.NewProtoLogger(line)
	return &pctx{
		logger: logger,
		line:   line,
		entry:  entry,
	}
}

func NewPctxWithLog(l log.Log) Pctx {
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

func (p *pctx) Log() log.Log {
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
