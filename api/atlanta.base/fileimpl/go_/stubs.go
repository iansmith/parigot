//go:build js
// +build js

// This provides the "stubs" that get filled in when wasmtime is used to run the WASM code. This codepath is
// ONLY for the WASM world.  The true implementations are in this same directory, but with build constraints
// so they are only built on the local machine (not in WASM).  The conect is provided by the func.go
// file in command/runner/
package go_

import _ "unsafe"

//go:noinline
//go:linkname FileSvcOpen filesvc.open
func FileSvcOpen(int32)
