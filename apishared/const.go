package apishared

import (
	"hash/crc32"
	"unsafe"
)

// constants in this package are ones that _MUST_ be synchronized
// between the two go systems, HOST and GUEST.

// EntryPoint must be the name of the exported function in the guest where
// parigot will start executing a service.
const EntryPoint = "parigot_main"

// ReturnDataName must be the name of the exported function in the guest
// that allocates space for the host to write results into.
//
// It is expected that the memory returned from this function
// will be marked as "dont collect" for the guest GC and will
// be released by the host code that actually uses the result.
//
// Host code should NOT instantiate a ReturnData struct, because
// the size and layout of the struct in the host machine is useless.
// This function should be called (eng/Instance.ReturnData) with
// the necessary arguements that can then be copied into the
// correct layou on the guest.
const ReturnDataName = "apiwasm.new_return_data_with_buffer"

// ReturnDataSize is how big a ReturnData struct is, as measured
// by the guest.
const ReturnDataSize = 24

// ReturnDataSize is where in a ReturnData struct the 128 bit error
// or id should be placed.  Note this is measured by the guest.
const ReturnDataIdErrOffset = 8

// ParamsNotReady is used between the host and guest to indicate the host has a set of parametrs
// ready to be used for guest execution.
const ParamsNotReady uint32 = 0xfffffffe

// ParamsDie is a signal between the host and guest that indicates that the host wants to
// remove the given function from the set of functions available from the host.  This is
// signal just like ParamsNotReady.
const ParamsDie uint32 = 0xfffffffd

// WasmPageSize is the standard wasm page size.
const WasmPageSize = 4096

// PagesPerExport is the amount of data that is consumed by buffers that are used when
// calling from the host to the guest.
const PagesPerExport = 16 * WasmPageSize //64k
// DynamicSizedData can be any length and it is the size of the various buffers that
// are in use for this function's variable sized data.  All the variable sized data
// is actually resident inside a single allocation that is sized by PagesPerExport and
// default wasm page size.
var DynamicSizedData = []int32{PagesPerExport >> 1, PagesPerExport >> 3, PagesPerExport >> 3, PagesPerExport >> 4, //8,2,2,1 wasm pages
	PagesPerExport >> 4, PagesPerExport>>4 - WasmPageSize>>1, PagesPerExport>>4 - WasmPageSize>>1} //1,1,0.5,0.5 wasm pges

// Durations to sleep between attempts to retrieve the next call to an exported function.
// If the number of misses is larger than the length of SleepSeqMicro then the last
// value in this list will be used.  Since that value is very large (100ms) that number
// will quickly become the average amount of time waited.  The smaller values are just
// to handle busy periods.
var SleepSeqMicro = []int{50, 100, 500, 1000, 5000, 10000, 20000, 20000, 50000, 50000, 100000}

// WasmWith is the size of a uintptr for the wasm virtual machine.
var WasmWidth = int32(unsafe.Sizeof(uintptr(0))) // in bytes
// WasmIs32Bit is true on a 32 bit wasm implementation
var WasmIs32Bit bool // init function

// MaxExportParam is the maximum num of elements in the slice that is passed as
// parameters.  How these are mapped to "true" parameters is up to the coordination
// of the host side and the guest side.
const MaxExportParam = 16 // note: strings or slices typically count as two params, not one

func init() {
	if WasmWidth == 4 {
		WasmIs32Bit = true
	}
	//log.Printf("xxx init of const")
}

// EntryPointSymbol is what should be used to start up a ready instance.  Note that we are turning
// off the instantiation's normal call to start so that we can control the startup and its entry point.
const EntryPointSymbol = "_start"

// if your message doesn't start with this, you have lost sync or had some other encoding
// problem.  if you can, it might be a good time to try reconnect to the sender.
var MagicStringOfBytes = uint64(0x1789071417760704)
var FrontMatterSize = 8 + 4
var TrailerSize = 4

var KoopmanTable = crc32.MakeTable(crc32.Koopman)
var ReadBufferSize = 4096

// unix domain socket for talking to the logviewer... note that the SocketEnvVar
// should be "" when you are running an app inside the dev container.  You only
// need SocketEnvVar for things running on the local machine, like the log viewer app.
const SocketEnvVar = "PARIGOT_SOCKET_DIR"
const SocketName = "logviewer.sock"

const ExpectedStackDumpSize = 4096 * 2
