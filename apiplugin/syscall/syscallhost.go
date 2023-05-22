package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"
	"unsafe"

	"github.com/iansmith/parigot/eng"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/sharedconst"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

type syscallPlugin struct{}

var ParigiotInitialize = syscallPlugin{}

// xxx global vairable kinda sucks
var currentEng eng.Engine

func (*syscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "parigot", "locate_", locate)
	e.AddSupportedFunc(ctx, "parigot", "dispatch_", dispatch)
	e.AddSupportedFunc(ctx, "parigot", "block_until_call_", blockUntilCall)
	e.AddSupportedFunc(ctx, "parigot", "bind_method_", bindMethod)
	e.AddSupportedFunc(ctx, "parigot", "run_", run)
	e.AddSupportedFunc(ctx, "parigot", "export_", export)
	e.AddSupportedFunc(ctx, "parigot", "return_value_", returnValue)
	e.AddSupportedFunc(ctx, "parigot", "require_", require)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)
	e.AddSupportedFunc_7i32_v(ctx, "parigot", "register_export_", registerExport)

	currentEng = e
	return true
}

func locate(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("locate %s 0x%x", m.Name(), stack)
}

func dispatch(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("dispatch 0x%x", stack)
}

func blockUntilCall(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("blockUntilCall 0x%x", stack)
}
func bindMethod(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("bindMethod 0x%x", stack)
}
func run(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("run 0x%x", stack)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("export 0x%ax", stack)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	instance, err := currentEng.InstanceByName(ctx, m.Name())
	log.Printf("XXX>>> %s:require got %p, %v", m.Name(), instance, err)
	if err != nil {
		panic(fmt.Sprintf("attempt to find module that called require failed, module is: %s", m.Name()))
	}
	length := eng.Util.DecodeU32(stack[0])
	ptr := eng.Util.DecodeU32(stack[1])
	memList, err := instance.Memory(ctx)
	if err != nil {
		panic(fmt.Sprintf("retreiving memory object from instance failed: %v ", err))
	}
	mem := memList[0]
	rawPb, err := mem.ReadBytes(ptr, length)
	if err != nil {
		panic(fmt.Sprintf("retreiving guest bytes of memory failed: %v ", err))
	}
	in := &syscallmsg.LocateRequest{}
	if err := proto.Unmarshal(rawPb, in); err != nil {
		panic(fmt.Sprintf("unable to unmarshal guest data: %v", err))
	}
	log.Printf("got the sent data in require! %s,%s", in.PackageName, in.ServiceName)
}
func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}

func readStringFromGuest(mem api.Memory, nameOffset int32) string {
	l, ok := mem.ReadUint32Le(uint32(nameOffset))
	if !ok {
		panic("unable to read the length of a string from the guest")
	}
	data := uint32(nameOffset + sizeofGuestInt)
	ptr, ok := mem.ReadUint32Le(data)
	if !ok {
		panic("unable to read the data pointer of a string from the guest")
	}
	result := make([]byte, int(l))
	for i := uint32(0); i < l; i++ {
		b, ok := mem.ReadByte(ptr + i)
		if !ok {
			panic("unable to read the data of a string from the guest")
		}
		result[int(i)] = b
	}
	return string(result)
}

// registerExport is call by the guest to tell the host that a guest function can be run (really "started")
// by the host when desired.
func registerExport(ctx context.Context, mod api.Module, param []uint64) {
	//nameHeaderRaw := int32(param[0])
	paramHeaderRaw := int32(param[1])
	//is32Bit := int32(param[2])
	bufferRaw := int32(param[3])
	exclusiveBufferSizePtrRaw := uint32(param[4])
	flagPtrRaw := uintptr(int32(param[5]))
	turnPtrRaw := uintptr(int32(param[6]))

	mem := mod.Memory()
	closure := func() {
		exampleHost(uintptr(paramHeaderRaw), uintptr(bufferRaw), uintptr(exclusiveBufferSizePtrRaw), mem)
	}
	go func(fn func(), flagPtrRaw, turnPtrRaw, exclusiveBufferSizePtrRaw uintptr) {
		time.Sleep(time.Duration(2) * time.Second)
		callSingleGuestFunction(mem, fn, flagPtrRaw, turnPtrRaw, exclusiveBufferSizePtrRaw)
	}(closure, flagPtrRaw, turnPtrRaw, uintptr(exclusiveBufferSizePtrRaw))
	return
}

// callSingleGuestFunction does exactly what it says.  It expects a closure in fn to actually do the call and gather result.  The other parameters
// are so this function can use a mutual exclusion algorithm with th guest.
func callSingleGuestFunction(mem api.Memory, fn func(), flagPtrRaw, turnPtrRaw, exclusiveBufferSizePtrRaw uintptr) {

	print("HOST peterson0\n")
	//peterson lock
	flag0 := flagPtrRaw
	flag1 := flagPtrRaw + sizeofGuestInt

	// ptr1 = 1 b/c we want to enter crit sect
	if ok := mem.WriteUint32Le(uint32(flag1), uint32(1)); !ok {
		panic("out of bounds write of flag ptr 1 = 1")
	}
	// set the turn to be the guest side
	if ok := mem.WriteUint32Le(uint32(turnPtrRaw), uint32(0)); !ok {
		panic("out of bounds write of turn ptr [1]")
	}
	print("HOST peterson1 ... \n")
	for { // while flag[0]==true and turn==0
		fl0, ok := mem.ReadUint32Le(uint32(flag0))
		if !ok {
			panic("unable to read the flag[0] value from the guest system")
		}

		turn, ok := mem.ReadUint32Le(uint32(turnPtrRaw))
		if !ok {
			panic("unable to read the turn value from the guest system")
		}
		if fl0 == 1 && turn == 0 {
			continue
		}
		break
	}

	print("HOST peterson2\n")
	// peterson critical section
	value, ok := mem.ReadUint32Le(uint32(exclusiveBufferSizePtrRaw))
	if !ok {
		panic("out of bounds read of exclusive buffer")
	}

	if value != sharedconst.ParamsNotReady {
		print("no work to do, client has not picked it up yet")
		time.Sleep(time.Duration(1) * time.Millisecond)
	} else {
		fn()
		print("HOST peterson3B\n")

	}
	print("HOST peterson6\n")
	//peterson unlock
	if ok := mem.WriteUint32Le(uint32(flag1), 0); !ok {
		panic("unable to write peterson flag 1 back to guest")
	}
	print("HOST peterson7\n")

}

func writeParam(mem api.Memory, paramHeader uintptr, index int32, value uint64) {
	if ok := mem.WriteUint64Le(uint32(paramHeader+(uintptr(index)*sizeofGuestUint64)), value); !ok {
		panic(fmt.Sprintf("unable to write guest parameter int param slice (%d)\n", index))
	}
}

const sizeofGuestInt = 4
const sizeofGuestSliceHeader = 12
const sizeofGuestByte = 1
const sizeofGuestUint64 = 8

// writeVariableSized data copies the data pointed to by data, with length length into the slice at index index inside buffer. Buffer is an
// slice of []byte. If the data is too large to fit in the buffer chosen by index, abortTooLarge is consulted.  If abortTooLarge is
// true, then the program panics.  This is appropriate for sending data that must be received in full. If abortTooLarge is false,
// the data is truncated to fit.  This latter policy is probably best for things like human readable strings. It returns a value
// suitable for use a parameter that points to the slice that has been altered.
func writeVariableSizedData(mem api.Memory, buffer uintptr, index int32, length uint32, data *byte, abortTooLarge bool) uint64 {
	if length > uint32(sharedconst.DynamicSizedData[index]) {
		if abortTooLarge {
			msg := fmt.Sprintf("size of data (%d) is larger than size of variable sized buffer (%d) at index (%d)",
				length, sharedconst.DynamicSizedData[index], index)
			panic(msg)
		}
		// truncate
		length = uint32(sharedconst.DynamicSizedData[index])
	}
	//get element slice, index * sizeof(sh)
	sh := buffer + uintptr(index*sizeofGuestSliceHeader)
	// write len
	if ok := mem.WriteUint32Le(uint32(sh), length); !ok {
		panic("unable to write data int len field of variable buffer, out of bounds")
	}
	// write cap
	if ok := mem.WriteUint32Le(uint32(sh+sizeofGuestInt), length); !ok {
		panic("unable to write data int cap field of variable buffer, out of bounds")
	}
	// data pointer
	destArea := uint32(sh + (2 * sizeofGuestInt))
	for i := uint32(0); i < length; i++ {
		dest := destArea + (i * sizeofGuestByte)
		source := (*byte)(unsafe.Pointer((uintptr(unsafe.Pointer(data)) + (uintptr(i) * sizeofGuestByte))))
		if ok := mem.WriteByte(dest, *source); !ok {
			panic("unable to write variable sized data, during data copy out of bounds")
		}
	}
	return uint64(sh)
}

func exampleHost(paramHeader uintptr, buffer uintptr, exclusiveBufferSizePtrRaw uintptr, mem api.Memory) {
	s := "hello, parigot"

	log.Printf("example host setting up params: %v", paramHeader == 0)
	setNumberOfParameters(mem, paramHeader, uint32(3))
	writeParam(mem, paramHeader, 0, uint64(len(s)))
	writeParam(mem, paramHeader, 1, writeStringToVariableSizeBuffer(mem, buffer, s, 6))
	writeParam(mem, paramHeader, 2, 42)
	signalNumberOfParameters(mem, uint32(3), exclusiveBufferSizePtrRaw)
}

func signalNumberOfParameters(mem api.Memory, num uint32, exclusiveBufferSizePtrRaw uintptr) {
	if ok := mem.WriteUint32Le(uint32(exclusiveBufferSizePtrRaw), num); !ok {
		panic("signal number of params, out of bounds")
	}

}

func setNumberOfParameters(mem api.Memory, paramHeader uintptr, num uint32) {
	if num >= sharedconst.MaxExportParam {
		panic(fmt.Sprintf("too many parameters for guest function call (max is %d)", sharedconst.MaxExportParam))
	}
	// setup params for guest side
	log.Printf("param header needed by settNumberOfParams %v", paramHeader == 0)
	if ok := mem.WriteUint32Le(uint32(paramHeader), num); !ok {
		panic("unable to write parameter to example guest function")
	}
	// setup params for guest side
	if ok := mem.WriteUint32Le(uint32(paramHeader+sizeofGuestUint64), num); !ok {
		panic("unable to write parameter to example guest function")
	}
}

// writeStringToVariableSizeBuffer is a convenience wrapper for writeVariableSizedData that handles putting a string into
// one of the variable sized buffers.
func writeStringToVariableSizeBuffer(mem api.Memory, buffer uintptr, s string, index int32) uint64 {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sData := (*byte)(unsafe.Pointer(strHeader.Data))
	l := uint32(len(s))
	return writeVariableSizedData(mem, buffer, index, l, sData, false)
}
