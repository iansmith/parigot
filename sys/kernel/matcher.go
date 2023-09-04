package kernel

import (
	"io"
	"log"
	"log/slog"
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

// callMatcher is an internal data structure object that
// connects calls to Dispatch (the call) with the response
// which are created by ReturnValue requests.  This data structure
// is perhaps the most critical in Parigot and should not
// be changed without serious consideration.
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
	// The function that is the last param is only use if
	// some host-side call is calling dispatch and wants to
	// receive the result.
	Dispatch(hid, orig id.HostId, cid id.CallId, mid id.MethodId, fn func(*syscall.ResolvedCall), w io.Writer) syscall.KernelErr
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

type waitingCall struct {
	inner        *syscall.ResolvedCall
	writer       io.Writer
	hostCallback func(*syscall.ResolvedCall)
}

type matcher struct {
	lock sync.Mutex
	// waiting and ready are maps from hostId to call id to matchingInfo
	// waiting is only a partially filled structure, it is completed
	// and then moved to ready by Response()
	waiting map[string]*waitingCall
	ready   map[string][]*syscall.ResolvedCall
}

// newCallMatcher returns a  new instance of CallMatcher. It should not be
// called by outside (user level) code.  There is only one CallMatcher in
// the system.
func newCallMatcher() *matcher {
	return &matcher{
		waiting: make(map[string]*waitingCall),
		ready:   make(map[string][]*syscall.ResolvedCall),
	}
}

var _ callMatcher = &matcher{}

func (c *matcher) Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr {
	c.lock.Lock()
	defer c.lock.Unlock()

	if cid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}
	waiter, ok := c.waiting[cid.String()]
	if !ok {
		return syscall.KernelErr_BadId
	}
	slog.Info("waiter retrieved for call", "call", cid.Short(),
		"callback", waiter.hostCallback != nil)

	delete(c.waiting, cid.String())
	waiter.inner.Result = a
	waiter.inner.ResultError = err
	hid := id.UnmarshalHostId(waiter.inner.HostId)
	if hid.IsZeroOrEmptyValue() {
		slog.Info("we just pulled a bad hid from waiting list", "cid", cid.Short())
	}

	// check to see if this is actually a host only receiver
	if waiter.hostCallback != nil {
		slog.Info("callback found, about to do callback")
		waiter.hostCallback(waiter.inner)
		return syscall.KernelErr_NoError
	}

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
	if hid.IsZeroOrEmptyValue() {
		panic("attempt to empty in response")
	}
	c.ready[hid.String()] = append(cidList, waiter.inner)
	return syscall.KernelErr_NoError
}

func (c *matcher) Dispatch(target, source id.HostId, cid id.CallId, mid id.MethodId, fn func(*syscall.ResolvedCall), w io.Writer) syscall.KernelErr {
	c.lock.Lock()
	defer c.lock.Unlock()

	if target.IsZeroOrEmptyValue() || cid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}

	_, ok := c.waiting[cid.String()]
	if ok {
		return syscall.KernelErr_BadCallId
	}
	if source.IsZeroOrEmptyValue() {
		source = target // for launch and exit
	}

	rc := &syscall.ResolvedCall{
		HostId:      source.Marshal(),
		CallId:      cid.Marshal(),
		MethodId:    mid.Marshal(),
		Result:      nil,
		ResultError: 0,
	}
	waiter := &waitingCall{
		hostCallback: fn,
		inner:        rc,
		writer:       w,
	}
	c.waiting[cid.String()] = waiter
	return syscall.KernelErr_NoError
}

func (c *matcher) Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr) {

	return c.readyImpl(hid, false)
}

func (c *matcher) ReadyLen(hid id.HostId) int {
	c.lock.Lock()
	defer c.lock.Unlock()

	cidList, ok := c.ready[hid.String()]
	if !ok {
		return 0
	}
	return len(cidList)

}

func (c *matcher) ReadyPeek(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr) {
	return c.readyImpl(hid, true)
}

func dumpReadyImpl(origHid id.HostId, ready map[string][]*syscall.ResolvedCall) {
	first := true
	for k, rc := range ready {
		if len(rc) > 0 {
			if first {
				log.Printf("--- ready dump-- %s", origHid.Short())
				first = false
			}
			log.Printf("---- > > > %s,(len is %d)", k, len(rc))
			for i := 0; i < len(rc); i++ {
				hid := id.UnmarshalHostId(rc[0].HostId)
				cid := id.UnmarshalCallId(rc[0].CallId)
				mid := id.UnmarshalMethodId(rc[0].MethodId)
				log.Printf("---- > > > > \t %d,(%s,%s,%s) ?result type? %s", i, hid.Short(), cid.Short(), mid.Short(), rc[0].GetResult().TypeUrl)
			}
		}
	}
}

// readyImpl is the code for getting the next ready item. It is shared by Ready and ReadyPeek.
func (c *matcher) readyImpl(hid id.HostId, isPeek bool) (*syscall.ResolvedCall, syscall.KernelErr) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if hid.IsZeroOrEmptyValue() {
		panic("bad value put in matcher")
	}

	dumpReadyImpl(hid, c.ready)
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
