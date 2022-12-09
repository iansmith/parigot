package lib

import (
	"fmt"
	"time"

	"github.com/iansmith/parigot/g/pb/call"
	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc         Id
	caller      string
	pctxEnabled bool
	currentPctx Pctx
}

func NewClientSideService(id Id, caller string) *ClientSideService {
	return &ClientSideService{
		svc:    id,
		caller: caller,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

func (c *ClientSideService) SetPctx(pctx *protosupport.PCtx) {
	print("client side service: setting Pctx\n")
	c.currentPctx = newPctxWithTime(time.Now(), pctx)
}
func (c *ClientSideService) EnablePctx() {
	c.pctxEnabled = true
}
func (c *ClientSideService) DisablePctx() {
	c.pctxEnabled = false
}
func (c *ClientSideService) Log(level pblog.LogLevel, message string) {
	if !c.pctxEnabled {
		// we can't print a warning, we are not running in an environment with any type of output
		panic("attempt to use client side logging without pctx enabled")
	}
	if c.pctxEnabled {
		if c.currentPctx == nil {
			c.currentPctx = NewPctx()
		}
		c.currentPctx.Log(level, message)
	}
}
func (c *ClientSideService) DumpLog() {
	if !c.pctxEnabled {
		panic("attempt to dump client side log without pctx enabled")
	}
	if c.currentPctx == nil {
		fmt.Printf("[no log output generated")
	} else {
		for i, e := range c.currentPctx.(*pctx).GetEvent() {
			if e == nil {
				print(fmt.Sprintf("NO EVENT FOUND #%05d\n", i))
				continue
			}
			print(fmt.Sprintf("%05d:%s (%v,%v)\n", i, e.GetMessage(), e.GetLine() == nil, len(e.GetLine())))
			for j, m := range e.GetLine() {
				if m == nil {
					print(fmt.Sprintf("\tNO MESSAGE FOUND #%05d\n", j))
					continue
				}
				print(fmt.Sprintf("\t%05d:%d:%s:%s\n", j, m.GetLevel(),
					m.GetStamp().AsTime().Format(time.Kitchen), m.GetMessage()))
			}
		}
	}

}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*call.DispatchResponse, error) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			return nil, NewPerrorFromError("unable to convert param for dispatch into Any", err)
		}
	}

	var ughPuke *protosupport.PCtx
	if c.currentPctx != nil {
		ughPuke = c.currentPctx.(*pctx).PCtx
	}

	in := &call.DispatchRequest{
		ServiceId: MarshalServiceId(c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    ughPuke,
		Param:     a,
	}
	// xxx this should be going through dispatch anyway
	return CallConnection().Dispatch(in)
}
