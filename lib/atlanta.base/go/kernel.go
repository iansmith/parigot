// This package is a thin wrapper around kernel functionality so you can cacd .ll
// that functionality easily and without having to create a KernelClient each time.
// It allows you to use the notation kernel.Exit(0), for example, if you import this
// under the default name.
package lib

import (
	"fmt"
	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/lib/k"
)

func Exit(in *kernel.ExitRequest) error {
	if err := k.Exit(in); err != 0 {
		return fmt.Errorf("exit failed locally, error code: %d", err) // probably wont happen
	}
	return nil // wont happen
}

func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) error {
	if err := k.Register(in, out); err != 0 {
		return fmt.Errorf("register failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}

func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) error {
	if err := k.Locate(in, out); err != 0 {
		return fmt.Errorf("locate failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}

func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) error {
	if err := k.Dispatch(in, out); err != 0 {
		return fmt.Errorf("dispatch failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}
