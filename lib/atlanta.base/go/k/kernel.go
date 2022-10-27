//go:build parigot

package k

import (
	"github.com/iansmith/parigot/g/pb/kernel"
	_ "unsafe"
)

//go:noinline
//go:linkname locate parigot.locate_
func locate(in *kernel.LocateRequest, out *kernel.LocateResponse) int32

func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) int32 {
	return locate(in, out)
}

//go:noinline
//go:linkname register parigot.register_
func register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) int32

func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) int32 {
	return register(in, out)
}

//go:noinline
//go:linkname dispatch parigot.dispatch_
func dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) int32

func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) int32 {
	return dispatch(in, out)
}

//go:noinline
//go:linkname exit parigot.exit_
func exit(in *kernel.ExitRequest) int32

func Exit(in *kernel.ExitRequest) int32 {
	return exit(in)
}
