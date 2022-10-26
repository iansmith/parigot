package jspatch

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"unsafe"
)

var memPtr *uintptr

func SetMemPtr(p *uintptr) {
	memPtr = p
}
func MemPtr() uintptr {
	return *memPtr
}

func ValueSetIndex(v int32, i int32, x int32, wtf int32) {
	log.Printf("not implemented: js.valueSetIndex (%x,%d,%x,%x)",
		v, i, x, wtf)
	os.Exit(1)
}

func ValueSet(vaddr int32, ptr int32, len int32, xaddr int32, wtf int32) {
	print("not implemented: js.valueSet %x,%x,%d,%x,%x",
		vaddr, ptr, len, xaddr, wtf)
	os.Exit(1)
}

func ValuePrepareString(v int32, x int32, i int32) {
	log.Printf("Value prepare string %x,%x %d", v, x, i)
	os.Exit(1)
}

func ValueLoadString(result int32, v int32, slice int32, len int32, cap int32) {
	log.Printf("Value Load String %x,(%x %d,%d)", v, slice, len, cap)
	os.Exit(1)
}

func FinalizeRef(x int32, y int32) {
	log.Printf("FinalizeRef %x,%x", x, y)
	os.Exit(1)
}
func ValueNew(result int32, v int32, args int32, args_len int32, args_cap int32, qqq int32) {
	log.Printf("ValueNew(%x,%x,%x,%d,%d) ???%x", result, v, args, args_len, args_cap, qqq)

	os.Exit(1)
}
func ValueCall(result int32, v int32, m int32, m_len int32, args int32, args_len int32, args_cap int32, qqq int32) {
	log.Printf("ValueCall(%x,%x,%x,%d,%x,%d,%d) ???%x", result, v, m, m_len, args, args_len, args_cap, qqq)
	os.Exit(1)
}

func StringVal(retVal int32, ptr int32, len int32, wtf int32) {
	log.Printf("js.stringVal called %x,%x,%d,%x",
		retVal, ptr, len, wtf)
	os.Exit(1)
}

func StrConvert(memPtr uintptr, ptr int32, length int32) string {
	// we could probably go bytesConvert and claim our cap was equal to our len but...
	buf := make([]byte, length)
	for i := int32(0); i < length; i++ {
		b := (*byte)(unsafe.Pointer(memPtr + uintptr(ptr+i)))
		buf[i] = *b
	}
	s := string(buf)
	return s
}

type jsObject interface {
	id() int
}

var object = make(map[int]jsObject)
var null = &jsObj{n: 2}
var trueBool = &jsObj{n: 3}
var falseBool = &jsObj{n: 4}
var win = &window{jsObj: &jsObj{n: 5}}
var zero = &jsObj{n: 1}
var notANumber = &jsObj{n: 0}

func init() {
	object[5] = win
	object[4] = falseBool
	object[3] = trueBool
	object[2] = null
	object[1] = zero
	object[0] = notANumber
}

type jsObj struct {
	n int
}

func (j *jsObj) id() int {
	return j.n
}

type window struct {
	*jsObj
}

const nanHead = 0x7FF80000

// stupid 64 bit trick load side
func loadValue(addr int32) int64 {
	buff := make([]byte, 8)
	log.Printf("loadValue(%x) memPtr=%x", addr, *memPtr)
	for i := int32(0); i < 8; i++ {
		buff[i] = *((*byte)(unsafe.Pointer(*memPtr + uintptr(addr+i))))
	}
	// get this as a little indian 64 bit unsigned int
	bits := binary.LittleEndian.Uint64(buff)
	// convert from float
	f := math.Float64frombits(bits)
	if !math.IsNaN(f) {
		panic("unexpected float" + fmt.Sprint(f))
	}
	id := binary.LittleEndian.Uint32(buff)
	return int64(id) //should be the object from the table
}

// stupid 64 bit trick, save side... we are assuming v is a small int
func storeValue(addr int32, v int64) {
	log.Printf("storeValue(%x,%x)-- assuming int", addr, v)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[4:], nanHead)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(v))
	for i := int32(0); i < 8; i++ {
		ptr := (*byte)(unsafe.Pointer(*memPtr + (uintptr(addr + i))))
		*ptr = buf[i]
	}
}

func ValueGet(retval int32, vAddr int32, propLen int32, propPtr int32, what_is_this int32) {
	prop := StrConvert(*memPtr, propLen, propPtr)
	log.Printf("ValueGet(%s,%x,%x)", prop, vAddr, retval)
	value := loadValue(vAddr)
	result := value // this should be a Reflect.get(value,prop)
	storeValue(retval, result)
}

func SetInt64(memptr uintptr, addr int32, value int64) {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (*memPtr + uintptr(addr))
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	binary.LittleEndian.PutUint64(buf, uint64(value))
}

func GetInt64(memptr uintptr, addr int32) int64 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (memptr + uintptr(addr))
	header.Data = ptr
	header.Len = 8
	header.Cap = 8
	value := binary.LittleEndian.Uint64(buf)
	return int64(value)
}
func GetInt32(memptr uintptr, addr int32) int32 {
	buf := []byte{}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	ptr := (memptr + uintptr(addr))
	header.Data = ptr
	header.Len = 4
	header.Cap = 4
	value := binary.LittleEndian.Uint32(buf)
	return int32(value)
}

func LoadSlice(memptr uintptr, addr int32) []byte {
	array := GetInt64(memptr, addr)
	l := GetInt64(memptr, addr+8)
	result := make([]byte, l)
	for i := int64(0); i < l; i++ {
		ptr := memptr + uintptr(array) + uintptr(i)
		result[i] = *((*byte)(unsafe.Pointer(ptr)))
	}
	return result
}
