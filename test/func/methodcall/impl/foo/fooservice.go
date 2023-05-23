package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/g/methodcall/v1"

	"github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.ServerWasmContext(context.Background(), "foo:main")
	defer pcontext.Dump(ctx)

	pcontext.Logf(ctx, pcontext.Info, "started main on")

	//log.Printf("xxxmain of foo")
	// syscall.RegisterExport(parigot_main, apiwasm.NewReturnDataWithBuffer,
	// 	apiwasm.NewString)
	ch := make(chan struct{})
	pcontext.Logf(ctx, pcontext.Info, "about to call wasm export")
	print("XXX HERE\n")
	closure := syscall.WasmExport("example", example_)
	pcontext.Logf(ctx, pcontext.Info, "closure1")
	go closure(ctx, ch)
	pcontext.Logf(ctx, pcontext.Info, "closure2")
	for range ch {
		log.Printf("got the signal from the closure")
	}

}

// example is a function that is an example of a function to be exported. From the standpoint of example,
// it receives this call and performs its return in the normal way.  It is running on the goroutine
// that was used to call WasmExport.
func example(s string, i int) int {
	print("example called with parameters '", s, "' and ", i, "\n")
	return len(s) + i //don't make i too big!
}

// example_ is an example wrapper function that understand how to translate the uint64s and dynamic
// buffers provided to it into the proper parameters for the function it wraps.  In this example,
// it pulls the string data and and length from parameters 0 and 1, and another integer from parameter 2.
// The string value is housed in buf[6], just because we can.  Only the caller and this wrapper
// have to know how the params and variable data are encoded.
// param[0] = string length
// param[1] = pointer to string (variable buffer 6)
//
//	WasmExport("example", example_, [][]byte{s})
func example_(raw []uint64, buffer [][]byte) uint64 {

	// sanity check
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buffer[6]))

	if uint32(raw[0]) != uint32(sh.Len) {
		panic("bad encoding in example (arg 1 len != buffer 6 len)")
	}
	if uint64(uintptr(unsafe.Pointer(sh))) != raw[1] {
		panic("bad encoding in example (arg 1 data != buffer 6 data)")
	}
	careful := readStringFromVariableBuffer(buffer, 6, int(raw[0]))

	//try it more direct way, note this way makes no copy of the data!
	//so the string is "gone" after this call returns... don't make a copy of it!
	aggressive := string(buffer[6])
	if aggressive != careful {
		panic("bad encoding for example (aggressive not the same as careful)")
	}
	i := int(raw[3])
	rawReturn := example(careful, i)
	ret := uint64(rawReturn)
	return ret
}

func readStringFromVariableBuffer(buffer [][]byte, index int, length int) string {
	// have to make a copy because str is immutable
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buffer[index]))
	source := (*byte)(unsafe.Pointer(sh.Data))
	copyBuf := make([]byte, length)
	for i := 0; i < length; i++ {
		bp := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(source)) + uintptr(i)*unsafe.Sizeof(byte(0))))
		copyBuf[i] = *bp
	}
	return string(copyBuf)
}

//go:export parigot_main
//go:linkname parigot_main
func parigot_main(argv []string, envp map[string]string) {
	//lib.FlagParseCreateEnv()
	panic("here")
	bg := context.Background()
	methodcall.ExportFooServiceOrPanic()
	methodcall.RequireFooServiceOrPanic(bg)
	s := &fooServer{}
	methodcall.RunFooService(s)
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *methodcallmsg.AddMultiplyRequest) (*methodcallmsg.AddMultiplyResponse, error) {
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")
	resp := &methodcallmsg.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, nil
}

func (f *fooServer) LucasSequence(ctx context.Context) (*methodcallmsg.LucasSequenceResponse, error) {
	pcontext.Debugf(ctx, "LucasSequence", "received call for fooServer.LucasSequence")
	resp := &methodcallmsg.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, nil
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *methodcallmsg.WritePiRequest) error {
	pcontext.Debugf(ctx, "WritePi", "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 {
		return fmt.Errorf("number of terms in WritePi must be a positive integer")
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
	}
	pcontext.Debugf(ctx, "WritePi", "%f", runningTotal)
	return nil
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (f *fooServer) Ready(ctx context.Context) bool {
	methodcall.WaitFooServiceOrPanic()

	return true
}
