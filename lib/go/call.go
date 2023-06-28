package lib

import (
	"context"
	"time"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var cidToCompleter = make(map[string]future.Completer)

// MatchCompleter is a utility for adding a new
// cid and completer to the tables used to look up
// the location where response values should be sent.
func MatchCompleter(cid id.CallId, comp future.Completer) {
	if cidToCompleter[cid.String()] != nil {
		panic("unexpected duplicate call id for matching to client side")
	}
	cidToCompleter[cid.String()] = comp
}

// CompleteCall is called from the CallOne handler to cause a
// prior dispatch call to be completed. The matching is done
// based on the cid.
func CompleteCall(ctx context.Context, cid id.CallId, result *anypb.Any, resultErr int32) syscall.KernelErr {
	comp, ok := cidToCompleter[cid.String()]
	if !ok {
		return syscall.KernelErr_NotFound
	}
	delete(cidToCompleter, cid.String())
	comp.CompleteMethod(ctx, result, resultErr)
	return syscall.KernelErr_NoError
}

var internalFuture = make(map[string]time.Time)
var keyToRealFuture = make(map[string]*future.Method[proto.Message, int32])

// ExpireMethod() checks the internal list of guest side futures
// that have no call id associated with them.  These futures come about
// when a implementation of a server function returns a future
// that is not completed.  This future likely exists because the
// implementation of the server function called another service
// and the result of the server function thus cannot be calculated
// immediately.  When the call is completed, the Success or Failure
// functions will be called on the original future.  This function
// exists to maintain a list so that we can expire and cancel futures
// that have waiting longer than the timeout time.
func ExpireMethod(ctx context.Context) {
	dead := make([]string, 0)
	for key, elem := range internalFuture {
		future := keyToRealFuture[key]
		if !future.Completed() {
			curr := pcontext.CurrentTime(ctx)
			diff := curr.Sub(elem)
			if diff.Milliseconds() > apishared.FunctionTimeoutInMillis {
				dead = append(dead, key)
			}
		} else {
			// just clean up finished ones
			dead = append(dead, key)
		}
	}
	// don't want delete as we iterate on internalFuture
	for _, key := range dead {
		f := keyToRealFuture[key]
		delete(internalFuture, key)
		f.Cancel()
	}
}

// AddServerReturn is called to register a server side function
// result as a future.
func AddServerReturn(fut future.Method[proto.Message, int32]) {

}
