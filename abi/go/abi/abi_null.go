//go:build !parigot_abi
// +build !parigot_abi

// Package abi package defines the interface between client side code and the runtime.  This
// abi is actually called parigot_abi at the wasm level.  This implementation of the
// package is primarily for IDE users, so the IDE will be able to "find" the implementation
// of the package and its functions, but these are not in used in the real implementation.
// Documentation for the golang version of the ABI is provided here, again for the
// convenience of IDE users.

package abi

import "time"

//go:wasm-module parigot_abi
//export OutputString
func OutputString(string) {}

//go:wasm-module parigot_abi
//export OutputStringConvert
func OutputStringConvert(int32, int32) {}

func JSNotImplemented() {}

func JSHandleEvent() {}

func TinyGoNotImplemented() {}

func FdWrite(int32, int32, int32, int32) int32 {
	return 0
}

func Exit(int32) {}

func Now() int64 {
	return 0
}

func SetNow(time.Time) {}
