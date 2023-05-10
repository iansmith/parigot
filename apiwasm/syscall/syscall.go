package syscall

import (
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
)

// var envVerbose = os.Getenv("PARIGOT_VERBOSE")

// // Flip this switch for debug output.
// var libparigotVerbose = false || envVerbose != ""

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
// go:wasm-module parigot
// go:export locate
func Locate(*syscallmsg.LocateRequest) *syscallmsg.LocateResponse

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
// go:wasm-module parigot
// go:export dispatch
func Dispatch(*syscallmsg.DispatchRequest) *syscallmsg.DispatchResponse

// BlockUntilCall is used to block a process until a request is received from another process.  Even when
// all the "processes" are in a single process for debugging, the BlockUntilCall is for the same purpose.
// go:wasm-module parigot
// go:export blockUntilCall
func BlockUntilCall(*syscallmsg.BlockUntilCallRequest) *syscallmsg.BlockUntilCallResponse

// go:wasm-module parigot
// go:export bindMethod
func BindMethod(*syscallmsg.BindMethodRequest) *syscallmsg.BindMethodResponse

// go:wasm-module parigot
// go:export run
func Run(*syscallmsg.RunRequest) *syscallmsg.RunResponse

// go:wasm-module parigot
// go:export export
func Export(*syscallmsg.ExportRequest) *syscallmsg.ExportResponse

// go:wasm-module parigot
// go:export return_value
func ReturnValue(*syscallmsg.ReturnValueRequest) *syscallmsg.ReturnValueResponse

// go:wasm-module parigot
// go:export require
func Require(*syscallmsg.RequireRequest) *syscallmsg.RequireResponse

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
// go:wasm-module parigot
// go:export exit
func Exit(*syscallmsg.ExitRequest) *syscallmsg.ExitResponse
