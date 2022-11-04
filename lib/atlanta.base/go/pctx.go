package lib

import (
	"fmt"
	"time"

	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/parigot"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Pctx is an interface that wraps parigot.PCtx at the protobuf level. It has methods to allow it to
// collect logs that are connected to events and to manage a simple key/value store of strings.
type Pctx interface {
	Log(pblog.LogLevel, string)
	EventStart(string)
	EventFinish()
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
	DeleteEntry(group string, name string) bool
	Now() time.Time
	Marshal() ([]byte, error)
}

func NewPctxFromBytes(pctxSlice []byte) (Pctx, error) {
	pctxWire := parigot.PCtx{}
	err := proto.Unmarshal(pctxSlice, &pctxWire)
	if err != nil {
		return nil, err
	}
	return newPctxWithTime(time.Now(), &pctxWire), nil
}

type pctx struct {
	*parigot.PCtx
	now       time.Time
	openEvent *parigot.PCtxEvent
}

// NewPctx creates a new Pctx for use on the client side.
func NewPctx() Pctx {
	// xxx fixme ... this should be doing a system call to get the time from kernel
	now := time.Now()
	p := &parigot.PCtx{
		Event: []*parigot.PCtxEvent{
			&parigot.PCtxEvent{
				Message: fmt.Sprintf("PCtx created on client side"),
			},
		},
	}
	return &pctx{
		PCtx:      p,
		now:       now,
		openEvent: p.Event[0],
	}
}

func newPctxWithTime(t time.Time, inner *parigot.PCtx) Pctx {
	return &pctx{now: t, PCtx: inner}
}

func (p *pctx) Now() time.Time {
	return p.now
}

func (p *pctx) Log(level pblog.LogLevel, msg string) {
	var t time.Time
	if p.openEvent == nil {
		p.EventStart("anonymous event created")
	}
	p.openEvent.Line = append(p.openEvent.GetLine(), &parigot.PCtxMessage{
		Stamp:   timestamppb.New(t),
		Level:   level,
		Message: msg,
	})
}

func (p *pctx) Entry(group, name string) (string, bool) {
	key := group + "." + name
	m := p.PCtx.GetEntry()
	v, ok := m[key]
	return v, ok
}
func (p *pctx) DeleteEntry(group, name string) bool {
	key := group + "." + name
	m := p.PCtx.GetEntry()
	_, ok := m[key]
	delete(m, key)
	return ok
}

func (p *pctx) SetEntry(group, name, value string) bool {
	key := group + "." + name
	_, ok := p.PCtx.GetEntry()[key]
	p.PCtx.GetEntry()[key] = value
	return ok
}

func (p *pctx) EventStart(message string) {
	if p.openEvent != nil {
		p.PCtx.Event = append(p.PCtx.GetEvent(), p.openEvent)
	}
	p.openEvent = &parigot.PCtxEvent{
		Message: message,
	}
}

func (p *pctx) EventFinish() {
	p.PCtx.Event = append(p.PCtx.GetEvent(), p.openEvent)
}

func (p *pctx) Marshal() ([]byte, error) {
	return proto.Marshal(p.PCtx)
}

func NewFromUnmarshal(b []byte) (Pctx, error) {
	p := pctx{}
	err := proto.Unmarshal(b, p.PCtx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
