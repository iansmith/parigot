package syscall

import (
	"context"
	"reflect"
	"runtime"
	"time"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/sharedconst"
)

// WasmExport makes the function provided "visible" to the host side with the given name.
// When the host side wants the function to run (or "be called"), it will put the parameters to fn
// into a set a uint64 parameters, and possible one or more data areas for dynamically sized
// data.  Then it will "trigger" (really "call") the guest side function, likely a wrapper
// around a real function.

// from the host side.  When the channel "unblocks" it will have the parameters available
// to it to operate.  Variable sized parameters will use the buffers provided to WasmExport and the
// host will not overfill these buffers.   The particular protocol used for parametrs that do not
// fit into the provided buffers is the decision of the host side (abort, truncate, etc).
// Note that if a buffer is used its address will be provided as one of the "normal" (uint64)
// parameters to the function.
//
// WasmExport should only be called from its own goroutine.  A simple example is:
// go WasmExport(...)
//
// The gorutine created will be used to run the exported function every time the exported
// function is "called" from the host side. The caller on the host side busywaits for
// this function to return, so exported functions should not take a great deal of time
// (like more than 500 milliseconds).
//
// The caller of this function must be on the guest side, and should provided as many f
// fixed size buffers of bytes that are needed for this function.  Each buffer of
// bytes is used by the host side to provide variable sized data.  The use any
// particular is up to the coordination of host and guest.
//
// All of the buffer values provided to this function must only be accessed by this function
// (and thus this goroutine).  They cannot be shared.
func WasmExport(name string, fn func(param []uint64, buffer [][]byte) uint64) func(context.Context, chan struct{}) {

	// this is the crucial piece of information that is changed by the host side
	petersonExclParamSize := uint32(sharedconst.ParamsNotReady)
	petersonExclParamSizePtr := uintptr(unsafe.Pointer(&petersonExclParamSize))

	// for peterson algorithm
	var flag [2]int32
	flagPtr := uintptr(unsafe.Pointer(&flag))
	var turn int32
	turnPtr := uintptr(unsafe.Pointer(&turn))

	// setup the dynamic buffers based on the bufOffset
	pages := make([]byte, sharedconst.PagesPerExport)
	// all the buffers share the same underlying storage (pages)
	buffer := make([][]byte, sharedconst.MaxExportParam)
	sum := 0
	bigHeader := (*reflect.SliceHeader)(unsafe.Pointer(&pages))
	for i := 0; i < len(sharedconst.DynamicSizedData); i++ {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&buffer[i]))
		sh.Len = int(sharedconst.DynamicSizedData[i])
		sh.Cap = int(sharedconst.DynamicSizedData[i])
		sh.Data = bigHeader.Data + uintptr(sum)
		shData := (*byte)(unsafe.Pointer(sh.Data))
		b := unsafe.Slice(shData, sharedconst.DynamicSizedData[i])
		buffer[i] = b
		sum += int(sharedconst.DynamicSizedData[i])
	}
	bufferHeader := uintptr(unsafe.Pointer(&buffer))

	// setup the params
	param := make([]uint64, sharedconst.MaxExportParam)
	paramHeader := uint32(uintptr(unsafe.Pointer(&param)))

	is32Bit := uint32(1) //true
	if sharedconst.WasmWidth != 4 {
		if sharedconst.WasmWidth != 8 {
			panic("unable to understand natural size of wasm implementation")
		}
		is32Bit = 0 //64 bit, false
	}
	nameHeader := uint32(uintptr((unsafe.Pointer(&name))))

	//
	// tell the host side that we are ready
	//

	RegisterExport(nameHeader,
		paramHeader,
		is32Bit,
		uint32(bufferHeader),
		uint32(petersonExclParamSizePtr),
		uint32(flagPtr),
		uint32(turnPtr))

	runner := func(ctx context.Context, ch chan struct{}) {
		wasmExportRun(ctx, ch, param, buffer, &flag, &turn, &petersonExclParamSize, fn)
	}
	return runner
}

func wasmExportRun(ctx context.Context, exitChan chan struct{}, param []uint64, buffer [][]byte, flag *[2]int32, turn *int32, petersonExclParamSize *uint32,
	fn func(param []uint64, buffer [][]byte) uint64) {
	// these are some loop maintenance vars for sleeping between invocations
	miss := 0
	needSleep := false
	pcontext.Debugf(ctx, "wasmExportRun", "arrived in run")
	//
	// Mainloop
	//
	pcontext.Debugf(ctx, "wasmExportRun", "entered the main loop of wasmExportRun")
	for {
		//peterson lock
		(*flag)[0] = 1
		*turn = 1
		for (*flag)[1] == 1 && *turn == 1 {
			runtime.Gosched()
		}
		pcontext.Debugf(ctx, "wasmExportRun", "we have the peterson lock (wasm)")

		// peterson critical section
		if uint32(*petersonExclParamSize) == sharedconst.ParamsNotReady {
			pcontext.Debugf(ctx, "wasmExportRun", "peterson 3A")
			miss++
			needSleep = true
		} else {
			pcontext.Debugf(ctx, "wasmExportRun", "peterson 3A")
			miss = 0
			if uint32(*petersonExclParamSize) == sharedconst.ParamsDie {
				exitChan <- struct{}{}
				return
			}
			// the value we grabbed in exclusiveBuffer is the size of the passed params
			result := fn(param, buffer)

			// reset paramHeaders... not that the host side MUST ignore these when collecting
			// results
			paramHeader := (*reflect.SliceHeader)(unsafe.Pointer(&param))
			bufferHeader := (*reflect.SliceHeader)(unsafe.Pointer(&buffer))

			paramHeader.Len = sharedconst.MaxExportParam
			paramHeader.Cap = sharedconst.MaxExportParam
			// reset buffers
			bufferHeader.Len = len(sharedconst.DynamicSizedData)
			bufferHeader.Cap = len(sharedconst.DynamicSizedData)
			pcontext.Debugf(ctx, "wasmExportRun", "peterson 4")
			for i := 0; i < len(sharedconst.DynamicSizedData); i++ {
				buffSliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&buffer[i]))
				buffSliceHeader.Len = len(sharedconst.DynamicSizedData)
				buffSliceHeader.Cap = len(sharedconst.DynamicSizedData)
			}

			// put the result at slot 0 of the pool
			param[0] = result
			*petersonExclParamSize = sharedconst.ParamsNotReady
		}
		//peterson unlock
		flag[0] = 0
		if needSleep {
			pcontext.Debugf(ctx, "wasmExportRun", "peterson unlock")
			sleepSome(miss)
		}
	}
}

func sleepSome(missCount int) {
	if missCount >= len(sharedconst.SleepSeqMicro) {
		missCount = len(sharedconst.SleepSeqMicro) - 1
	}
	dur := time.Duration(sharedconst.SleepSeqMicro[missCount]) * time.Microsecond
	//before := time.Now()
	time.Sleep(dur)
	//after := time.Now()
	//actual := after.Sub(before)
}
