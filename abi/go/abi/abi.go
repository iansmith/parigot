//go:build parigot_abi
// +build parigot_abi

package abi

import (
	"time"
)

//go:wasm-module parigot_abi
//export OutputString
func OutputString(string)

//go:wasm-module parigot_abi
//export JSNotImplemented
func JSNotImplemented()

//go:wasm-module parigot_abi
//export JSHandleEvent
func JSHandleEvent()

//go:wasm-module parigot_abi
//export TinyGoNotImplemented
func TinyGoNotImplemented()

//go:wasm-module parigot_abi
//export Exit
func Exit(exitCode int32)

//go:wasm-module parigot_abi
//export Now
func Now() int64

//go:wasm-module parigot_abi
//export SetNow
func SetNow(time.Time)

//go:wasm-module parigot_abi
//export fd_write
func FdWrite(int32, int32, int32, int32) int32
