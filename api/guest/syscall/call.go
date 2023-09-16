package syscall

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"

	"google.golang.org/protobuf/types/known/anypb"
)

var cidToCompleter = make(map[string]*completionInfo)

type completionInfo struct {
	completer      future.Completer
	originalSource id.HostId
}

func newCompletionInfo(completer future.Completer, hid id.HostId) *completionInfo {
	return &completionInfo{completer, hid}

}

// MatchCompleter is a utility for adding a new
// cid and completer to the tables used to look up
// the location where response values should be sent.  The time value has
// be passed in from the outside.
func MatchCompleter(ctx context.Context, t time.Time, hid id.HostId, origHost id.HostId, cid id.CallId, comp future.Completer) {
	cand, _ := getCompleter(hid, cid)
	if cand != nil {
		panic("unexpected duplicate call id for matching to client side")
	}
	setCompleter(ctx, t, hid, origHost, cid, newCompletionInfo(comp, origHost))
}

func completerKey(hid id.HostId, cid id.CallId) string {
	return fmt.Sprintf("%s;%s", hid.Short(), cid.Short())
}
func getCompleter(hid id.HostId, cid id.CallId) (future.Completer, id.HostId) {
	key := completerKey(hid, cid)
	result := cidToCompleter[key]
	if result == nil {
		return nil, id.HostIdZeroValue()
	}
	return result.completer, result.originalSource
}
func setCompleter(ctx context.Context, t time.Time, hid id.HostId, origHost id.HostId, cid id.CallId, info *completionInfo) {
	key := completerKey(hid, cid)
	internalFuture[key] = t
	cidToCompleter[key] = info
}

func delCompleted(hid id.HostId, cid id.CallId) {
	key := completerKey(hid, cid)
	delete(cidToCompleter, key)
	delete(internalFuture, key)
}

var iter = 0

// CompleteCall is called from the ReadOneAndCall handler to cause a
// prior dispatch call to be completed. The matching is done
// based on the cid.
func CompleteCall(ctx context.Context, hid id.HostId, cid id.CallId, result *anypb.Any, resultErr int32) syscall.KernelErr {
	comp, orig := getCompleter(hid, cid)
	if comp == nil {
		slog.Warn("CompleteCall received, but no invoke found, ignoring")
		return syscall.KernelErr_NotFound
	}
	delCompleted(hid, cid)
	//log.Printf("xxx --- cmplete call %s", result.TypeUrl)
	//debug.PrintStack()
	comp.CompleteMethod(ctx, result, resultErr, orig)
	return syscall.KernelErr_NoError
}

var internalFuture = make(map[string]time.Time)

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
func ExpireMethod(ctx context.Context, curr time.Time) {
	dead := make([]string, 0)
	for key, elem := range internalFuture {
		comp, ok := cidToCompleter[key]
		if !ok {
			log.Printf("unable to find matching completer for key %s", key)
			continue
		}
		if comp.completer != nil {
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
		f := cidToCompleter[key]
		delete(internalFuture, key)
		f.completer.Cancel()
		delete(cidToCompleter, key)
	}
}
