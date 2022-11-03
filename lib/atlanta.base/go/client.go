package lib

import (
	"github.com/iansmith/parigot/g/pb/kernel"
	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/parigot"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc         Id
	caller      string
	pctxEnabled bool
	currentPctx Pctx
}

func NewClientSideService(id Id) *ClientSideService {
	return &ClientSideService{
		svc: id,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
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

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*kernel.DispatchResponse, error) {
	a, err := anypb.New(param)
	if err != nil {
		return nil, NewPerrorFromError("unable to convert param for dispatch into Any", err)
	}
	var p *parigot.PCtx
	if c.currentPctx != nil {
		p = c.currentPctx.(*pctx).PCtx //xxx fixme, shouldn't be looking inside
	}
	in := &kernel.DispatchRequest{
		ServiceId: MarshalServiceId(c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    p,
		Param:     a,
	}
	// xxx this should be going through dispatch anyway
	return Dispatch(in)
}
