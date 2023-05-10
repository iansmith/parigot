//go:build js
// +build js

// This provides the "stubs" that get filled in when wasmtime is used to run the WASM code. This codepath is
// ONLY for the WASM world.  The true implementations are in this same directory, but with build constraints
// so they are only built on the local machine (not in WASM).  The connection is provided by the func.go
// file in command/runner/
package go_

import _ "unsafe"

//go:noinline
//go:linkname QueueSvcCreateHandler queuesvc.create_handler
func QueueSvcCreateHandler(int32)

//go:noinline
//go:linkname QueueSvcDeleteHandler queuesvc.delete_handler
func QueueSvcDeleteHandler(int32)

//go:noinline
//go:linkname QueueSvcMarkDoneHandler queuesvc.mark_done_handler
func QueueSvcMarkDoneHandler(int32)

//go:noinline
//go:linkname QueueSvcLengthHandler queuesvc.length_handler
func QueueSvcLengthHandler(int32)

//go:noinline
//go:linkname QueueSvcSendHandler queuesvc.send_handler
func QueueSvcSendHandler(int32)

//go:noinline
//go:linkname QueueSvcReceiveHandler queuesvc.receive_handler
func QueueSvcReceiveHandler(int32)

//go:noinline
//go:linkname QueueSvcLocateHandler queuesvc.locate_handler
func QueueSvcLocateHandler(int32)
