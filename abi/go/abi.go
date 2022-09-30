//go:build parigot_abi
// +build parigot_abi

package abi

import "time"

//go:wasm-module parigot_abi
//export outputString
func OutputString(string)

//export exit
func Exit(int)

//export now
func Now() time.Time

//export setNow
func SetNow(time.Time)
