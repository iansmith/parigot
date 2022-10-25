package k

import "github.com/iansmith/parigot/g/pb/kernel"

func Locate(_ *kernel.LocateRequest, _ *kernel.LocateResponse) int32 {
	return 0
}

func Register(_ *kernel.RegisterRequest, _ *kernel.RegisterResponse) int32 {
	return 0
}

func Dispatch(_ *kernel.DispatchRequest, _ *kernel.DispatchResponse) int32 {
	return 0
}

func Exit(_ *kernel.ExitRequest) int32 {
	return 0
}
