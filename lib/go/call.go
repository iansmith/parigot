package lib

import (
	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
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
func CompleteCall(cid id.CallId, result *anypb.Any, resultErr int32) syscall.KernelErr {
	comp, ok := cidToCompleter[cid.String()]
	if !ok {
		return syscall.KernelErr_NotFound
	}
	delete(cidToCompleter, cid.String())
	comp.CompleteMethod(result, resultErr)
	return syscall.KernelErr_NoError
}
