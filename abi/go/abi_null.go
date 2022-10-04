//go:build !parigot_abi
// +build !parigot_abi

// Package abi package defines the interface between client side code and the runtime.  This
// abi is actually called parigot_abi at the wasm level.  This implementation of the
// package is primarily for IDE users, so the IDE will be able to "find" the implementation
// of the package and its functions, but these are not in used in the real implementation.
// Documentation for the golang version of the ABI is provided here, again for the
// convenience of IDE users.

package _go

import "time"

func OutputString(string) {}

func JSNotImplemented() {}

func TinyGoNotImplemented() {}

func Exit(int) {}

func Now() time.Time {
	return time.Time{}
}

func SetNow(time.Time) {}
