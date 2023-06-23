package lib

import (
	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

var cidToClient = make(map[string]*ClientSideService)

// AddMatchingClient is a utility for adding a new
// cid, client pair to the tables used to look up
// the location where response values should be sent.
func AddMatchingClient(cid id.CallId, cs *ClientSideService) {
	if cidToClient[cid.String()] != nil {
		panic("unexpected duplicate call id for matching to client side")
	}
	cidToClient[cid.String()] = cs
}

// CompleteCall is called from the CallOne handler to cause a
// prior dispatch call to be completed. The matching is done
// based on the cid and cidToClient.
func Resolve(_ id.HostId, cid id.CallId, result *anypb.Any, resultErr int32) syscall.KernelErr {
	cs, ok := cidToClient[cid.String()]
	if !ok {
		return syscall.KernelErr_NotFound
	}
	delete(cidToClient, cid.String())
	return cs.CompleteCall(cid, result, resultErr)
}
