package jspatch

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

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

func (w *WasmMem) TrueAddr(addr int32) uintptr {
	return w.memPtr + uintptr(addr)
}

func (w *WasmMem) LoadSliceOfValues(addr int32) jsObject {
	array := w.GetInt64(addr)
	l := w.GetInt64(addr + 8)
	arr := make([]jsObj, l)
	a := goToJS(arr)
	for i := int64(0); i < l; i++ {
		a.SetIndex(int(i), w.LoadValue(int32(array+i*8))) //xxx why?why give me a 64 bit ptr?
	}
	return a
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
	result := make([]byte, l)
	for i := int64(0); i < l; i++ {
		ptr := w.memPtr + uintptr(array) + uintptr(i)
		result[i] = *((*byte)(unsafe.Pointer(ptr)))
	}
	return result
}
func (w *WasmMem) GetFloat32(addr int32) float32 {
	ptr := (*uint32)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float32frombits(*ptr)
}

func (w *WasmMem) GetFloat64(addr int32) float64 {
	ptr := (*uint64)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float64frombits(*ptr)
}

// stupid 64 bit trick load side
func (w *WasmMem) LoadValue(addr int32) jsObject {
	f := w.GetFloat64(addr)
	if !math.IsNaN(f) { // is all zeros a valid float?
		if math.Float64bits(f) == 0 {
			return undefined
		}
		return floatValue(f)
	}
	// maybe it's not a valid float...
	if math.Float64bits(f) == 0 {
		return undefined
	}
	// normal procedure
	t := (math.Float64bits(f) >> 32) & 7
	enter("LoadValue", fmt.Sprint("binary type ", t))
	id := w.GetInt32(addr)
	return object.get(id)
}

// stupid 64 bit trick, save side... we are assuming v is a small int
func (w *WasmMem) StoreValue(addr int32, obj jsObject) {
	if !obj.isNumber() && obj.id() < 0 {
		panic("attempt store a value that isn't in the global table: " + fmt.Sprint(obj.id()))
	}

	highOrder, lowOrder := obj.binaryRep()
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr+4))
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	binary.LittleEndian.PutUint32(buf, highOrder)
	ptr -= 4
	header.Data = ptr
	binary.LittleEndian.PutUint32(buf, lowOrder)
}
