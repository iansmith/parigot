package apishared

import (
	"hash/crc32"
	"unsafe"
)

// constants in this package are ones that _MUST_ be synchronized
// between the two go systems, HOST and GUEST.

// WasmWith is the size of a uintptr for the wasm virtual machine.
var WasmWidth = int32(unsafe.Sizeof(uintptr(0))) // in bytes
// WasmIs32Bit is true on a 32 bit wasm implementation
var WasmIs32Bit bool // init function

// WasmPageSize is the size of a memory page in wasm. I believe this is
// dictated by the spec.
const WasmPageSize = 4096

func init() {
	if WasmWidth == 4 {
		WasmIs32Bit = true
	}
}

// EntryPointSymbol is what should be used to start up a ready instance.  Note that we are turning
// off the instantiation's normal call to start so that we can control the startup and its entry point.
const EntryPointSymbol = "_start"

// In parigot the 8 byte magic value, when needed, is the date of the french
// and us revolutions, in hex.
var MagicStringOfBytes = uint64(0x1789071417760704)

// GuestReceiveBufferSize is the maximum amount data that the guest expects to
// read in a response back from the server.  Usually the send side is
// known apriori.
const GuestReceiveBufferSize = WasmPageSize

// unix domain socket for talking to the logviewer... note that the SocketEnvVar
// should be "" when you are running an app inside the dev container.  You only
// need SocketEnvVar for things running on the local machine, like the log viewer app.
//const SocketEnvVar = "PARIGOT_SOCKET_DIR"
//const SocketName = "logviewer.sock"

// ExpectedStackDumpSize is used to allocate space so that stack trace
// can be placed in it, then read back line by line.
const ExpectedStackDumpSize = 4096 * 2

// const FrontMatterSize = 12
// const TrailerSize = 4
// const WriteTimeout = 100 * time.Millisecond
// const ReadTimeout = 100 * time.Millisecond
// const LongReadTimeout = 500 * time.Millisecond

// KoopmanTable is the `crc32.Koopman` data in a table ready to use for
// CRC32 computations.
var KoopmanTable = crc32.MakeTable(crc32.Koopman)

// ReadBuffer is the maximum amount of data you can expect to receive in
// single read call with files or the network.
var ReadBufferSize = 8192

// The amount of time we will wait for a function call to be
// completed.
var FunctionTimeoutInMillis = int64(3000)
