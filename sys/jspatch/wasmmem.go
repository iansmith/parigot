package jspatch

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"unsafe"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/jspatch/jsemul"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch to get debug output from WasmMem.  It does the lowest level work
// of grabbing or pushing memory in the other processes address space.  It works on
// behalf of the kernel as part of every process.
var wasmmemVerbose = true

type WasmMem struct {
	memPtr uintptr
}

func NewWasmMem(memPtr uintptr) *WasmMem {
	return &WasmMem{
		memPtr: memPtr,
	}
}

func (w *WasmMem) SetUint8(addr int32, v byte) {
	ptr := (*byte)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	*ptr = v
}
func (w *WasmMem) String() string {
	return fmt.Sprintf("mem-%x", w.memPtr)
}

func (w *WasmMem) TrueAddr(addr int32) uintptr {
	return w.memPtr + uintptr(addr)
}

func (w *WasmMem) TestSliceIsZeroLen(addr int32) bool {
	return w.GetInt64(addr+8) == 0
}

func (w *WasmMem) LoadSliceOfValues(addr int32) []jsemul.JsObject {
	array := w.GetInt64(addr)
	l := w.GetInt64(addr + 8)
	if l == 0 {
		return nil
	}
	slice := make([]jsemul.JsObject, l)
	for i := int64(0); i < l; i++ {
		slice[int(i)] = w.LoadValue(int32(array + i*8))
	}
	return slice
}

func (w *WasmMem) SetFloat64(addr int32, value float64) {
	floatBits := math.Float64bits(value)
	print(fmt.Sprintf("xxx SetFloat64 called %x\n", floatBits))
	w.SetInt64(addr, int64(floatBits))
}

func (w *WasmMem) SetInt64(addr int32, value int64) {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := w.memPtr + uintptr(addr)
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	binary.LittleEndian.PutUint64(buf, uint64(value))
}
func (w *WasmMem) SetInt32(addr int32, value int32) {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := w.memPtr + uintptr(addr)
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	binary.LittleEndian.PutUint32(buf, uint32(value))
}

func (w *WasmMem) GetInt64(addr int32) int64 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := w.memPtr + uintptr(addr)
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	value := binary.LittleEndian.Uint64(buf)
	return int64(value)
}
func (w *WasmMem) GetInt32(addr int32) int32 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr))
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	value := binary.LittleEndian.Uint32(buf)
	return int32(value)
}

func (w *WasmMem) LoadStringWithLen(dataAddr int32, lenAddr int32) string {
	ptr := w.GetInt64(dataAddr)
	len_ := w.GetInt64(lenAddr)
	buf := make([]byte, len_)
	for i := int64(0); i < len_; i++ {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(int32(ptr+i))))
		buf[i] = *str
	}
	return string(buf)
}

func (w *WasmMem) LoadCString(ptr int32) string {
	var buf bytes.Buffer
	i := int32(0)
	for {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(int32(ptr+i))))
		buf.WriteByte(*str)
		if *str == 0 {
			break
		}
		i++
	}
	asBytes := buf.Bytes()
	return string(asBytes[:len(asBytes)-1])
}

func (w *WasmMem) CopyToMemAddr(memAddr int32, content []byte) {
	len_ := int32(len(content))
	for i := int32(0); i < len_; i++ {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(int32(memAddr)+i)))
		*str = content[i]
	}
}
func (w *WasmMem) CopyToPtr(dataAddr int32, content []byte) {
	//print("CopyToPtr: data addr ", uintptr(dataAddr), " len of content ", len(content), "\n")
	ptr := w.GetInt32(dataAddr)
	len_ := int32(len(content))
	for i := int32(0); i < len_; i++ {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(int32(ptr)+i)))
		*str = content[i]
	}
}

func (w *WasmMem) LoadStringTwoPtrs(addr int32, l int32) string {
	data := w.GetInt32(addr)
	size := w.GetInt32(l)
	buf := make([]byte, size)
	for i := int32(0); i < size; i++ {
		ptr := (*byte)(unsafe.Pointer(w.memPtr + uintptr(int32(data)+i)))
		buf[i] = *ptr
	}
	return string(buf)
}

func (w *WasmMem) LoadString(addr int32) string {
	ptr := w.GetInt64(addr + 0)
	l := w.GetInt64(addr + 8)
	if l > 4096 {
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
			Stamp:   timestamppb.Now(), //xxx fixme(iansmith) should be using kernel time
			Message: fmt.Sprintf("LOADSTRING: xxx!!!! wasmmem refusing to load string because length is too large: 0x%x ", l),
		}, true, false, false, nil)
		return ""
	}
	buf := make([]byte, l)
	for i := int64(0); i < l; i++ {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(ptr+i)))
		buf[i] = *str
	}
	return string(buf)
}

func (w *WasmMem) LoadSlice(addr int32) []byte {
	array := w.GetInt64(addr)
	l := w.GetInt64(addr + 8)
	if l < 4096 {
		print(fmt.Sprintf("LoadSlice called but len is %d\n", l))
		debug.PrintStack()
		print("END OF STACK\n")
	}
	return w.LoadSliceWithLenAddr(int32(array), int32(l))
}

func (w *WasmMem) LoadSliceWithKnownLength(addr int32, l int64) []byte {
	array := w.GetInt64(addr)
	result := make([]byte, int(l))
	for i := int64(0); i < l; i++ {
		ptr := w.memPtr + uintptr(array) + uintptr(i)
		result[i] = *((*byte)(unsafe.Pointer(ptr)))
	}
	return result
}

func (w *WasmMem) LoadSliceWithLenAddr(addr, lenAddr int32) []byte {
	if lenAddr < 4096 {
		wasmmemPrint("LoadSliceWithAddr", "Returning empty slice 0x%x,0x%x because len addr is too small (0x%d) \n",
			addr, lenAddr, lenAddr)
		return []byte{}

	}
	l := w.GetInt64(lenAddr)
	if l == 0 {
		wasmmemPrint("LoadSliceWithAddr ", "Returning empty slice %x,%x because len is zero\n",
			addr, lenAddr)
		return []byte{}
	}
	return w.LoadSliceWithKnownLength(addr, l)
}

func (w *WasmMem) GetFloat32(addr int32) float32 {
	ptr := (*uint32)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float32frombits(*ptr)
}

func (w *WasmMem) GetFloat64(addr int32) float64 {
	ptr := (*uint64)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float64frombits(*ptr)
}

// stupid 64 bit trick load side.  It's not as clean as it should be because it
// has to use a lot from jsemul.
func (w *WasmMem) LoadValue(addr int32) jsemul.JsObject {
	f := w.GetFloat64(addr)
	if !math.IsNaN(f) { // is all zeros a valid float?
		if math.Float64bits(f) == 0 {
			return jsemul.Undefined
		}
		return jsemul.FloatValue(f)
	}
	// maybe it's not a valid float...
	if math.Float64bits(f) == 0 {
		return jsemul.Undefined
	}
	// normal procedure
	//t := (math.Float64bits(f) >> 32) & 7
	id := w.GetInt32(addr)
	if id < 0 || id > 100 {
		print(fmt.Sprintf("bad id for value %x, 64 bit version %x\n", id, w.GetInt64(addr)))
	}
	return jsemul.JsObjectMap.Get(id)
}

// stupid 64 bit trick, save side... we are assuming the id is a small int on obj
// if this is not something that represents itself.  This should not be so
// dependent on the jsemul package.
func (w *WasmMem) StoreValue(addr int32, obj jsemul.JsObject) {
	if !obj.IsNumber() && obj.Id() < 0 {
		panic(fmt.Sprintf("attempt to store a value that isn't in the global table: %s\n", obj))
	}
	bits := obj.BinaryRep()
	// print(fmt.Sprintf("setting binary rep for number? %v, %x\n", obj.isNumber(), bits))
	// this conversion to int from uint depends on the nanHead not having the first bit set!
	w.SetInt64(addr, int64(bits))

	// buf := []byte{}
	// header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	// ptr := (w.memPtr + uintptr(addr+4))
	// header.Data = ptr
	// header.Len = 4
	// header.Cap = 4
	// binary.LittleEndian.PutUint32(buf, highOrder)
	// ptr -= 4
	// header.Data = ptr
	// binary.LittleEndian.PutUint32(buf, lowOrder)
}

func wasmmemPrint(method string, spec string, arg ...interface{}) {
	if wasmmemVerbose {
		print(method, fmt.Sprintf(spec, arg...))
	}
}
