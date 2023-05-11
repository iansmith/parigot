package syscall

import (
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
)

// var envVerbose = os.Getenv("PARIGOT_VERBOSE")

// // Flip this switch for debug output.
// var libparigotVerbose = false || envVerbose != ""

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
//go:wasm-module parigot
//go:export locate
func Locate(*syscallmsg.LocateRequest) *syscallmsg.LocateResponse

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
//
//go:wasm-module parigot
//go:export dispatch
func Dispatch(*syscallmsg.DispatchRequest) *syscallmsg.DispatchResponse

// BlockUntilCall is used to block a process until a request is received from another process.  Even when
// all the "processes" are in a single process for debugging, the BlockUntilCall is for the same purpose.
// go:wasm-module parigot
//
//go:wasm-module parigot
//go:export blockUntilCall
func BlockUntilCall(*syscallmsg.BlockUntilCallRequest) *syscallmsg.BlockUntilCallResponse

// BindMethod is the way that a particular service gets associated with
// a given method id. This is normally not needed by user code because the
// generated code for any service will call this automatically.
//
//go:wasm-module parigot
//go:export bindMethod
func BindMethod(*syscallmsg.BindMethodRequest) *syscallmsg.BindMethodResponse

// Run is a request to start running. Note that this may not return
// immediately and may fail entirely.  For most user code this is not
// used because user code usually uses file.WaitFileServiceOrPanic() to
// block service File until it is cleared to run.
//
//go:wasm-module parigot
//go:export run
func Run(*syscallmsg.RunRequest) *syscallmsg.RunResponse

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
//go:wasm-module parigot
//go:export export
func Export(*syscallmsg.ExportRequest) *syscallmsg.ExportResponse

// ReturnValue is not a call that user code should be using. It is the
// mechanism for a return value to be communicated back to the caller
// from the caller.  User code will typically use the wrappers around
// this that make the method calls looking synchronous.
//
//go:wasm-module parigot
//go:export return_value
func ReturnValue(*syscallmsg.ReturnValueRequest) *syscallmsg.ReturnValueResponse

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
//go:wasm-module parigot
//go:export require
func Require(*syscallmsg.RequireRequest) *syscallmsg.RequireResponse

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
//
//go:wasm-module parigot
//go:export exit
func Exit(*syscallmsg.ExitRequest) *syscallmsg.ExitResponse
