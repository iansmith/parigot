package sys

import (
	"github.com/iansmith/parigot/apishared/id"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/sys/dep"

	"google.golang.org/protobuf/types/known/anypb"
)

// Flip this switch to get extra debug information from the nameserver when it is doing
// various lookups.
var nameserverVerbose = false

const MaxService = 127

const parigotNameserverRemoteAddress = "parigot_ns:13330"

type callContext struct {
	mid          id.MethodId                         // the method id this call is going to be made TO
	method       string                              // if the call is remote our LOCAL mid wont mean squat, the remote needs the name
	target       dep.DepKey                          // the process/addr this call is going to be made TO
	cid          id.CallId                           // call id that should be be used by the caller to match results
	sender       dep.DepKey                          // the process/addr this call is going to be made FROM
	sid          id.ServiceId                        // service that is being called
	respCh       chan *syscallmsg.ReturnValueRequest // this is where to send the return results
	param        *anypb.Any                          // where to put the param data
	pctx         *protosupportmsg.Pctx               // where to put the previous pctx
	timedOut     bool                                // set to true when we waited on a call for a while and didn't get anything
	exitAfterUse bool                                // this is set to true ONLY when the nscore has requested it AND the inflight queue is empty
}
