package wheeler

import (
	"log"

	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type callMatcher struct {
	// waiting and ready are maps from hostId to call id to matchingInfo
	// waiting is only a partially filled structure, it is completed
	// and then moved to ready by Response()
	waiting map[string]*syscall.ResolvedCall
	ready   map[string]map[string]*syscall.ResolvedCall
}

// newCallMatcher returns a  new instance of CallMatcher. It should not be
// called by outside (user level) code.  There is only one CallMatcher in
// the system.
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
		return syscall.KernelErr_BadId
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
		return syscall.KernelErr_BadId
	}
	cidMap[cid.String()] = rc
	return syscall.KernelErr_NoError
}

func (c *callMatcher) Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr {

	if hid.IsZeroOrEmptyValue() || cid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}

	_, ok := c.waiting[cid.String()]
	if ok {
		log.Printf("xxx --- got a call id that we have seen before: %s", cid.Short())
		return syscall.KernelErr_BadCallId
	}
	rc := &syscall.ResolvedCall{
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
	if !ok || len(cid) == 0 {
		return nil, syscall.KernelErr_NoError
	}
	for key, rc := range cid {
		delete(cid, key)
		return rc, syscall.KernelErr_NoError
	}
	return nil, syscall.KernelErr_BadCallId
}
