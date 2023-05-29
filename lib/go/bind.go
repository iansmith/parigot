package lib

import (
	"github.com/iansmith/parigot/apiwasm/syscall"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/proto"
)

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func BindMethodIn(in *syscallmsg.BindMethodRequest) (*syscallmsg.BindMethodResponse, error) {
	return bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func BindMethodInNoPctx(in *syscallmsg.BindMethodRequest, _ func(proto.Message) error) (*syscallmsg.BindMethodResponse, error) {
	return bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
// func BindMethodOut(in *syscallmsg.BindMethodRequest) (*syscallmsg.BindMethodResponse, error) {
func BindMethodOut(in *syscallmsg.BindMethodRequest) (*syscallmsg.BindMethodResponse, error) {
	return bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_OUT)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func BindMethodBoth(in *syscallmsg.BindMethodRequest) (*syscallmsg.BindMethodResponse, error) {
	return bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_BOTH)
}

// bindMethodByName is the implementation of all three of the Bind* calls.
func bindMethodByName(in *syscallmsg.BindMethodRequest, dir syscallmsg.MethodDirection) (*syscallmsg.BindMethodResponse, error) {
	in.Direction = dir
	out, err := syscall.BindMethod(in)
	if err != nil {
		return nil, err
	}
	// XXX FIX ME
	// kid := nil //XXX id.Unmarshal((*protosupportmsg.KernelErrorId)(out.MethodId))

	// if kid.IsError() {
	// 	return nil, id.NewPerrorFromId("bind", kid)
	// }
	return out, nil
}
