package jspatch

import (
	"encoding/binary"
	"log"
	"math"
	"reflect"
	"unsafe"
)

type wasmMem struct {
	memPtr uintptr
}

func newWasmMem(memPtr uintptr) *wasmMem {
	return &wasmMem{
		memPtr: memPtr,
	}
}

func (w *wasmMem) setInt64(addr int32, value int64) {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr))
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	binary.LittleEndian.PutUint64(buf, uint64(value))
}

func (w *wasmMem) getInt64(addr int32) int64 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr))
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	value := binary.LittleEndian.Uint64(buf)
	return int64(value)
}
func (w *wasmMem) getInt32(addr int32) int32 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr))
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	value := binary.LittleEndian.Uint32(buf)
	return int32(value)
}

func (w *wasmMem) loadString(addr int32) string {
	ptr := w.getInt64(addr + 0)
	l := w.getInt64(addr + 8)
	buf := make([]byte, l)
	for i := int64(0); i < l; i++ {
		str := (*byte)(unsafe.Pointer(w.memPtr + uintptr(ptr+i)))
		buf[i] = *str
	}
	return string(buf)
}

func (w *wasmMem) loadSlice(memptr uintptr, addr int32) []byte {
	array := w.getInt64(addr)
	l := w.getInt64(addr + 8)
	result := make([]byte, l)
	for i := int64(0); i < l; i++ {
		ptr := w.memPtr + uintptr(array) + uintptr(i)
		result[i] = *((*byte)(unsafe.Pointer(ptr)))
	}
	return result
}
func (w *wasmMem) getFloat32(addr int32) float32 {
	ptr := (*uint32)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float32frombits(*ptr)
}

func (w *wasmMem) getFloat64(addr int32) float64 {
	ptr := (*uint64)(unsafe.Pointer(w.memPtr + uintptr(addr)))
	return math.Float64frombits(*ptr)
}

// stupid 64 bit trick load side
func (w *wasmMem) loadValue(addr int32) jsObject {
	f := w.getFloat64(addr)
	log.Printf("loadValue float: %f, %v", f, math.IsNaN(f))
	id := w.getInt32(addr)
	log.Printf("loadValue id: %d", id)
	return object[id]
}

// stupid 64 bit trick, save side... we are assuming v is a small int
func (w *wasmMem) storeValue(addr int32, obj jsObject) {
	log.Printf("storeValue(%x,%x)-- assuming int", addr, obj.id())
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (w.memPtr + uintptr(addr+4))
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	binary.LittleEndian.PutUint32(buf, uint32(nanHead))
	ptr -= 4
	header.Data = ptr
	binary.LittleEndian.PutUint32(buf, uint32(obj.id()))
	refCount[obj.id()]++
}
