package wheeler

import (
	"reflect"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

// CallMatcher is an internal data structure object that
// connects calls to Dispatch (the call) with the response
// which are created by ReturnValue requests.
type CallMatcher interface {
	// Response is called when a return value is
	// being processed. Any value that
	// is returned is NOT from the execution but from
	// the Response call itself.  Be aware that the
	// Response call is likely to be from a different
	// host than the original Dispatch call.
	Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr
	// Dispatch creates the necessary entries to handle
	// a future call to Response.  The value returned is
	// related to the Dispatch itself, it is not related
	// to the execution of the call being registered.
	Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr
	// Ready returns a resolved call or nil if no Promises are
	// resolved for the given host.
	Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr)
}

// Listener is a type that is used when we are running
type Listener interface {
	Case() []reflect.SelectCase
	// Handle sets up the return value whose pointer is
	// given in the params.  The first parameter is the value
	// sent through the channel.  The second parameter is the choice
	// number of the case selected--this probably only matters
	// a methodRequestListener.
	Handle(reflect.Value, int, *syscall.ReadOneResponse)
}

type ExitListener struct {
	ch chan int32
}

// MakeSidMidCombo is a utility for construction of a key (string) that is
// derived from the sid and mid given.
func MakeSidMidCombo(sid id.ServiceId, mid id.MethodId) string {
	return sid.String() + "," + mid.String()
}

// CallInfo is sent to the channels that represent service/method calls.
type CallInfo struct {
	cid   id.CallId
	param *anypb.Any
}
