// This package is a thin wrapper around kernel functionality so you can call
// that functionality easily and without having to create a KernelClient each time.
// It allows you to use the notation kernel.Exit(0), for example, if you import this
// under the default name.
package lib

import (
	"fmt"
	"github.com/iansmith/parigot/g/parigot/kernel"
)

func Exit(in *kernel.ExitRequest) error {
	if err := exit(in); err != 0 {
		return fmt.Errorf("exit failed locally, error code: %d", err) // probably wont happen
	}
	return nil // wont happen
}

func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) error {
	if err := register(in, out); err != 0 {
		return fmt.Errorf("register failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}

func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) error {
	if err := locate(in, out); err != 0 {
		return fmt.Errorf("locate failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}

func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) error {
	if err := dispatch(in, out); err != 0 {
		return fmt.Errorf("dispatch failed locally, error code: %d", err) // probably wont happen
	}
	return nil
}

//go:export locate parigot.locate_
func locate(in *kernel.LocateRequest, out *kernel.LocateResponse) int32

//go:export register parigot.register_
func register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) int32

//go:export dispatch parigot.dispatch_
func dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) int32

//go:export exit parigot.exit_
func exit(in *kernel.ExitRequest) int32
