//go:build parigot
package k

import "github.com/iansmith/parigot/g/pb/kernel"

//export Locate parigot.locate_
func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) int32

//export Register parigot.register_
func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) int32

//export Dispatch parigot.dispatch_
func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) int32

//export Exit parigot.exit_
func Exit(in *kernel.ExitRequest) int32
