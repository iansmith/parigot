package lib

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/protosupport"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Pctx is an interface that wraps protosupport.PCtx at the protobuf level. It has methods to allow it to
// collect logs that are connected to events and to manage a simple key/value store of strings.
type Pctx interface {
	Log(pblog.LogLevel, string)
	EventStart(string)
	EventFinish()
	Now() time.Time
	Marshal() ([]byte, error)
	Dump() string

	// for passing values between different subsystems
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
	DeleteEntry(group string, name string) bool
}

type Item interface {
	StringWithIndent(indent int, buf *bytes.Buffer) string
	IsComposite() bool
	// Marshal([]byte) (int, error)
}

// type CompositeItem interface {
// 	Item
// 	Item() []Item
// 	AddItem(Item)
// }

type msgWrapper struct {
	msg *protosupport.PCtxMessage
}

// func (m *msgWrapper) Marshal(buffer []byte) (int, error) {
// 	encoded, err := proto.MarshalOptions{}.MarshalAppend(buffer, m.msg)
// 	if err != nil {
// 		return -812, err
// 	}
// 	return len(encoded), nil
// }

func (m *msgWrapper) StringWithIndent(indent int, buf *bytes.Buffer) string {
	for i := 0; i < indent; i++ {
		buf.WriteString("\t")
	}

	buf.WriteString(fmt.Sprintf("%s:", m.msg.GetStamp().AsTime().Format(time.RFC3339)))
	switch m.msg.GetLevel() {
	case pblog.LogLevel_LOGLEVEL_DEBUG:
		buf.WriteString("DEBUG:")
	case pblog.LogLevel_LOGLEVEL_INFO:
		buf.WriteString(" INFO:")
	case pblog.LogLevel_LOGLEVEL_WARNING:
		buf.WriteString(" WARN:")
	case pblog.LogLevel_LOGLEVEL_ERROR:
		buf.WriteString("ERROR:")
	case pblog.LogLevel_LOGLEVEL_FATAL:
		buf.WriteString("FATAL:")
	default:
		buf.WriteString("UNKNOWN:")
	}
	s := m.msg.GetMessage()
	buf.WriteString(s)
	if !strings.HasSuffix(s, "\n") {
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m *msgWrapper) IsComposite() bool {
	return false
}

func NewLineItem(msg *protosupport.PCtxMessage) *msgWrapper {
	return &msgWrapper{msg: msg}
}

type eventWrapper struct {
	// for speed of marshalling we keep the empty fields of PCtxMessage here, even though we only
	// use the event field
	event *protosupport.PCtxMessage
}

func (m *eventWrapper) IsComposite() bool {
	return true
}

// func (m *eventWrapper) Marshal(buffer []byte) (int, error) {
// 	encoded, err := proto.MarshalOptions{}.MarshalAppend(buffer, m.event)
// 	if err != nil {
// 		return -812, err
// 	}
// 	return len(encoded), nil
// }

func messageToItem(msg *protosupport.PCtxMessage) Item {
	if msg.GetEvent() != nil {
		return NewEventItem(msg)
	}
	return NewLineItem(msg)
}

func (m *eventWrapper) StringWithIndent(indent int, buf *bytes.Buffer) string {
	for i := 0; i < indent; i++ {
		buf.WriteString("\t")
	}
	buf.WriteString(fmt.Sprintf("--> %s <---", m.event.GetMessage()))
	if m.event.GetEvent() == nil {
		buf.WriteString("[no log messages]\n")
	} else {
		for _, m := range m.event.GetEvent().Line {
			var item Item
			if m.Event != nil {
				item = NewEventItem(m)
			} else {
				item = NewLineItem(m)
			}
			if item.IsComposite() {
				item.StringWithIndent(indent+1, buf)
				continue
			}
			item.StringWithIndent(indent, buf)
		}
	}
	return buf.String()
}

func NewEventItem(event *protosupport.PCtxMessage) *eventWrapper {
	return &eventWrapper{event: event}
}

func NewPctxFromBytes(pctxSlice []byte) (Pctx, error) {
	pctxWire := protosupport.PCtx{}
	err := proto.Unmarshal(pctxSlice, &pctxWire)
	if err != nil {
		return nil, err
	}
	root := &pctxWire
	if root == nil {
		panic("unable to find any log info inside a pctxWire")
	}
	return &pctx{now: time.Now() /*xxxfixme, should ask kernel*/, root: root}, nil
}

func NewPctxFromProtosupport(pctxWire *protosupport.PCtx) Pctx {
	root := pctxWire
	if root == nil {
		panic("unable to find any log info inside a pctxWire")
	}
	return &pctx{now: time.Now() /*xxxfixme, should ask kernel*/, root: root}
}

type pctx struct {
	now  time.Time
	root *protosupport.PCtx
}

func (p *pctx) Dump() string {
	var buf bytes.Buffer
	e := NewEventItem(p.root.GetEvent())
	e.StringWithIndent(0, &buf)
	return buf.String()
}

// NewPctx creates a new Pctx for use on the client side.
func NewPctx() Pctx {
	// xxx fixme ... this should be doing a system call to get the time from kernel
	now := time.Now()

	pbEvent := &protosupport.PCtxMessage{
		Event: &protosupport.PCtxEvent{
			Message: fmt.Sprintf("Call %s started", NewCallId().Short()),
		},
	}

	return &pctx{
		now: now,
		root: &protosupport.PCtx{
			Event: pbEvent,
			Open:  pbEvent,
		},
	}
}

func (p *pctx) Now() time.Time {
	return p.now
}

func (p *pctx) Log(level pblog.LogLevel, msg string) {
	t := p.Now()
	pbMsg := &protosupport.PCtxMessage{
		Stamp:   timestamppb.New(t),
		Level:   level,
		Message: msg,
	}
	if p.root.Open == nil {
		panic("open event is nil in, cannot log")
	}
	if p.root.Open.Event == nil {
		panic("open message is not an event")
	}
	p.root.Open.Event.Line = append(p.root.Open.Event.GetLine(), pbMsg)
}

func (p *pctx) EventStart(message string) {
	if p.root.Open == nil {
		panic("open event not found")
	}
	evt := &protosupport.PCtxMessage{
		Event: &protosupport.PCtxEvent{Message: message, Parent: p.root.Open},
	}
	p.root.Open.Event.Line = append(p.root.Open.Event.GetLine(), evt)
	p.root.Open = evt
}

func (p *pctx) EventFinish() {
	p.root.Open = p.root.Open.Event.Parent
}

func (p *pctx) Marshal() ([]byte, error) {
	return proto.Marshal(p.root)
}

func NewFromUnmarshal(b []byte) (Pctx, error) {
	p := pctx{}
	err := proto.Unmarshal(b, p.root)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *pctx) Entry(group, name string) (string, bool) {
	key := group + "." + name
	m := p.root.GetEntry()
	v, ok := m[key]
	return v, ok
}

func (p *pctx) DeleteEntry(group, name string) bool {
	key := group + "." + name
	m := p.root.GetEntry()
	_, ok := m[key]
	delete(m, key)
	return ok
}

func (p *pctx) SetEntry(group, name, value string) bool {
	key := group + "." + name
	_, ok := p.root.GetEntry()[key]
	p.root.GetEntry()[key] = value
	return ok
}
