package syscall

import (
	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

var _matcher CallMatcher = newCallMatcher()

func matcher() CallMatcher {
	return _matcher
}

type callMatcher struct {
	// waiting and ready are maps from hostId to call id to matchingInfo
	// waiting is only a partially filled structure, it is completed
	// and then moved to ready by Response()
	waiting map[string]*syscall.ResolvedCall
	ready   map[string]map[string]*syscall.ResolvedCall
}

func newCallMatcher() *callMatcher {
	return &callMatcher{
		waiting: make(map[string]*syscall.ResolvedCall),
		ready:   make(map[string]map[string]*syscall.ResolvedCall),
	}
}

func (c *callMatcher) Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr {
	if cid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}
	rc, ok := c.waiting[cid.String()]
	if !ok {
		return syscall.KernelErr_BadCallId
	}
	delete(c.waiting, cid.String())
	rc.Result = a
	rc.ResultError = err
	hid := id.UnmarshalHostId(rc.HostId)
	// we are now ready to move to the ready list
	cidMap, ok := c.ready[hid.String()]
	if !ok {
		c.ready[hid.String()] = make(map[string]*syscall.ResolvedCall)
		cidMap = c.ready[hid.String()]
	}
	_, ok = cidMap[cid.String()]
	if ok {
		return syscall.KernelErr_BadCallId
	}
	cidMap[cid.String()] = rc
	return syscall.KernelErr_NoError
}

func (c *callMatcher) Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr {
	rc, ok := c.waiting[cid.String()]
	if ok {
		return syscall.KernelErr_BadCallId
	}
	rc = &syscall.ResolvedCall{
		HostId:      hid.Marshal(),
		CallId:      cid.Marshal(),
		Result:      nil,
		ResultError: 0,
	}
	c.waiting[cid.String()] = rc
	return syscall.KernelErr_NoError
}
func (c *callMatcher) Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr) {
	cid, ok := c.ready[hid.String()]
	if !ok {
		return nil, syscall.KernelErr_NoError
	}
	for _, rc := range cid {
		return rc, syscall.KernelErr_NoError
	}
	return nil, syscall.KernelErr_BadCallId
}
