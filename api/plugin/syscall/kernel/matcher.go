package kernel

import (
	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type matcher struct {
	// waiting and ready are maps from hostId to call id to matchingInfo
	// waiting is only a partially filled structure, it is completed
	// and then moved to ready by Response()
	waiting map[string]*syscall.ResolvedCall
	ready   map[string][]*syscall.ResolvedCall
}

// newCallMatcher returns a  new instance of CallMatcher. It should not be
// called by outside (user level) code.  There is only one CallMatcher in
// the system.
func newCallMatcher() *matcher {
	return &matcher{
		waiting: make(map[string]*syscall.ResolvedCall),
		ready:   make(map[string][]*syscall.ResolvedCall),
	}
}

var _ callMatcher = &matcher{}

func (c *matcher) Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr {
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
	cidList, ok := c.ready[hid.String()]
	if !ok {
		c.ready[hid.String()] = []*syscall.ResolvedCall{}
		cidList = c.ready[hid.String()]
	}

	for _, candidate := range cidList {
		cand := id.UnmarshalCallId(candidate.CallId)
		if cand.Equal(cid) {
			return syscall.KernelErr_BadId
		}
	}
	c.ready[hid.String()] = append(cidList, rc)
	return syscall.KernelErr_NoError
}

func (c *matcher) Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr {

	if hid.IsZeroOrEmptyValue() || cid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}

	_, ok := c.waiting[cid.String()]
	if ok {
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

func (c *matcher) Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr) {
	return c.readyImpl(hid, false)
}

func (c *matcher) ReadyLen(hid id.HostId) int {
	cidList, ok := c.ready[hid.String()]
	if !ok {
		return 0
	}
	return len(cidList)

}

func (c *matcher) ReadyPeek(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr) {
	return c.readyImpl(hid, true)
}

// readyImpl is the code for getting the next ready item. It is shared by Ready and ReadyPeek.
func (c *matcher) readyImpl(hid id.HostId, isPeek bool) (*syscall.ResolvedCall, syscall.KernelErr) {
	cidList, ok := c.ready[hid.String()]
	if !ok || len(cidList) == 0 {
		return nil, syscall.KernelErr_NoError
	}
	rc := cidList[0]
	if !isPeek {
		cidList = cidList[1:]
		c.ready[hid.String()] = cidList
	}
	return rc, syscall.KernelErr_NoError

}

// callMatcher is an internal data structure object that
// connects calls to Dispatch (the call) with the response
// which are created by ReturnValue requests.  This data structure
// should not be used outside the nameserver.
type callMatcher interface {
	// Response is called when a return value is
	// being processed. Any value that
	// is returned is NOT from the execution but from
	// the Response call itself.
	Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr
	// Dispatch creates the necessary entries to handle
	// a future call to Response.  The value returned is
	// related to the Dispatch itself, it is not related
	// to the execution of the call being registered.
	Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr
	// Ready returns a resolved call or nil if no outstanding
	// resolutions are ready.
	Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr)
	// ReadyLen returns how many items are ready.  This is usueful
	// primarily for debugging.
	ReadyLen(hid id.HostId) int
	// ReadyLen returns the front item from the ready list, but leaves
	// it still in the list. This is primarily for debugging.
	ReadyPeek(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr)
}
