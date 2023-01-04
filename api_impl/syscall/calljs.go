//go:build js
// +build js

package syscall

// This provides the "stubs" that get filled in when wasmtime is used to run the WASM code. This is path
// ONLY for the WASM world.  The true implementations are in the package sys because that is *host* code
// that gets bound with the WrapFunc() machinery of wasmtime.

import _ "unsafe"

//go:noinline
//go:linkname locate parigot.locate_
func locate(int32)

//go:noinline
//go:linkname dispatch parigot.dispatch_
func dispatch(int32)

//go:noinline
//go:linkname bindMethod parigot.bind_method_
func bindMethod(int32)

//go:noinline
//go:linkname exit parigot.exit_
func exit(int32)

//go:noinline
//go:linkname blockUntilCall parigot.block_until_call_
func blockUntilCall(int32)

//go:noinline
//go:linkname returnValue parigot.return_value_
func returnValue(int32)

//go:noinline
//go:linkname export parigot.export_
func export(int32)

//go:noinline
//go:linkname require parigot.require_
func require(int32)

//go:noinline
//go:linkname run parigot.run_
func run(int32)

//go:noinline
//go:linkname backdoorLog parigot.backdoor_log_
func backdoorLog(int32)
