//go:build parigot_abi
// +build parigot_abi

// This file has been machine generated from proto/abi.proto.  DO NOT EDIT OR YOU WILL LOSE.

// This file is here to tell tinygo what functions to elide when
// it is generating code for this package.  Most IDEs get confused when
// you have "non defined" functions, although it's perfectly valid
// go.  The functions that start with _ are potentially restricted
// (dangerous).  The functions that are the real implementations
// are implemented by the container, are not in the source or address
// space of a client, and have a trailing _.  Functions with more
// complex signatures (e.g. if they use a string) have to binary
// patched into the wasm module because we are using tinygo and it
// doesn't understand that we want to control the wasm-level type
// signatures.

package abi

// These functions have the true function name at the WASM level.
func OutputString_(_ string)
func _Exit_(_ int64) int64
func Now_()
func _SetNow_(_ int64, _ bool)
func TinygoNotImplemented_(_ string)
func JSNotImplemented_(_ string)

// These functions are wrappers around the true ABI functions, as is done with libc.
// We have to patch functions that take arguments outside of the basic four, and this
// is done with binary editing after compilation is finished.
func OutputString(s string) {
	return OutputString_(string)

}
func _Exit(code int64) int64 {
	return _Exit_(int64)

}
func Now() {
	Now_()
}
func _SetNow(now int64, freeze_clock bool) {
	return _SetNow_(int64, bool)

}
func TinygoNotImplemented(message string) {
	return TinygoNotImplemented_(string)

}
func JSNotImplemented(message string) {
	return JSNotImplemented_(string)

}
