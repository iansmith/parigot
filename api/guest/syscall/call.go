package syscall

import (
	"context"
	"fmt"
	"log"
	"os"
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
func MatchCompleter(hid id.HostId, cid id.CallId, comp future.Completer) {
	if getCompleter(hid, cid) != nil {
		panic("unexpected duplicate call id for matching to client side")
	}
	setCompleter(hid, cid, comp)
}

func completerKey(hid id.HostId, cid id.CallId) string {
	return fmt.Sprintf("%s;%s", hid.Short(), cid.String())
}
func getCompleter(hid id.HostId, cid id.CallId) future.Completer {
	log.Printf("xxx getCompleter: %s, %s", hid.Short(), cid.Short())
	key := completerKey(hid, cid)
	return cidToCompleter[key]
}
func setCompleter(hid id.HostId, cid id.CallId, f future.Completer) {
	log.Printf("xxx setCompleter: %s, %s", hid.Short(), cid.Short())
	key := completerKey(hid, cid)
	cidToCompleter[key] = f
}
func delCompleted(hid id.HostId, cid id.CallId) {
	log.Printf("xxx delCompleter: %s, %s", hid.Short(), cid.Short())
	key := completerKey(hid, cid)
	delete(cidToCompleter, key)

}

var iter = 0

// CompleteCall is called from the ReadOneAndCall handler to cause a
// prior dispatch call to be completed. The matching is done
// based on the cid.
func CompleteCall(ctx context.Context, hid id.HostId, cid id.CallId, result *anypb.Any, resultErr int32) syscall.KernelErr {
	iter++
	if iter == 20 {
		os.Exit(1)
	}
	comp := getCompleter(hid, cid)
	if comp == nil {
		pcontext.Errorf(ctx, " no way to complete complete call: %s (on host %s)", cid.Short(), CurrentHostId().Short())
		log.Printf("xxx -- cidToCompleter after failure: %+v (on %s)", cidToCompleter, CurrentHostId().Short())
		return syscall.KernelErr_NotFound
	}
	log.Printf("xxx -- about to delete the cidToCompleter entry from %+v", cidToCompleter)
	delCompleted(hid, cid)
	log.Printf("xxxx --- complete method call %s", cid.Short())
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
	log.Printf("xxx -- expire method total size %d", len(internalFuture))
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
