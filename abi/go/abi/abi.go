//go:build parigot_abi
// +build parigot_abi

package abi

import "time"

//go:wasm-module parigot_abi
//export OutputString
func OutputString(string)

//go:wasm-module parigot_abi
//export JSNotImplemented
func JSNotImplemented()

//go:wasm-module parigot_abi
//export JSNotImplemented1
func JSNotImplemented1(int32)

//go:wasm-module parigot_abi
//export JSHandleEvent
func JSHandleEvent()

//go:wasm-module parigot_abi
//export TinyGoNotImplemented
func TinyGoNotImplemented()

//go:wasm-module parigot_abi
//export TinyGoNotImplemented1
func TinyGoNotImplemented1(int32)

//go:wasm-module parigot_abi
//export Exit
func Exit(int)

//go:wasm-module parigot_abi
//export Now
func Now() time.Time

//go:wasm-module parigot_abi
//export SetNow
func SetNow(time.Time)
