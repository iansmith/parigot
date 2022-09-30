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

// OutputString sends the specified string to the current terminal when running in local
// mode.
func OutputString(_ string) {}

// Exit causes the current service to exit and the value provided is the error code.
func Exit(_ int) {}

// Now returns the current time as a time.Time _unless_ the function SetNow has been
// called, in which case it returns whatever value was provided to SetNow.
func Now() time.Time { return time.Time{} }

// SetNow sets the value returned by Now.  Once set, Now will only change its
// value after subsequent calls to SetNow.
func SetNow(_ time.Time) {}
