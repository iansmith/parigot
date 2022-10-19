package jspatch

import (
	"encoding/binary"
	"log"
	"math"
	"os"
	"unsafe"
)

var memPtr *uintptr

func SetMemPtr(p *uintptr) {
	memPtr = p
}

func ValueSetIndex(int64, int32, int64, int32) {
	print("not implemented: js.valueSetIndex")
	os.Exit(1)
}

func ValuePrepareString(int32, int64, int32) {
	print("not implemented: js.valuePrepareString")
	os.Exit(1)
}

func ValueLoadString(int64, int32, int32, int32, int32) {
	print("not implemented: js.ValueLoadString")
	os.Exit(1)
}

func FinalizeRef(int64, int32) {
	print("not implemented: js.finalizeRef")
	os.Exit(1)
}
func ValueNew(int32, int64, int32, int32, int32, int32) {
	print("not implemented: js.valueNew")
	os.Exit(1)
}
func ValueCall(int32, int64, int32, int32, int32, int32, int32, int32) {
	print("not implemented: js.valueCall")
	os.Exit(1)
}

//func StringVal() {
//	print("not implemented: js.stringVal")
//	os.Exit(1)
//}

func strConvert(memPtr uintptr, ptr int32, length int32) string {
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
var null = &window{jsObj: &jsObj{n: 2}}
var trueBool = &window{jsObj: &jsObj{n: 3}}
var falseBool = &window{jsObj: &jsObj{n: 4}}
var win = &window{jsObj: &jsObj{n: 5}}

func init() {
	object[5] = win
	object[4] = falseBool
	object[3] = trueBool
	object[2] = null
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
	for i := int32(0); i < 8; i++ {
		buff[i] = *((*byte)(unsafe.Pointer(*memPtr + uintptr(addr+i))))
		log.Printf("loading byte %d, %x", i, buff[i])
	}
	// get this as a little indian 64 bit unsigned int
	bits := binary.LittleEndian.Uint64(buff)
	// convert from float
	f := math.Float64frombits(bits)
	if !math.IsNaN(f) {
		panic("unexpected float")
	}
	id := binary.LittleEndian.Uint32(buff)
	log.Printf("id requested: %d", id)
	return int64(id) //should be the object from the table
}

// stupid 64 bit trick, save side
func storeValue(addr int32, v int64) {
	u := uint64(v)
	f := math.Float64frombits(u)
	if !math.IsNaN(f) {
		panic("unexpected float")
	}
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[4:], nanHead)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(v))
	for i := int32(0); i < 8; i++ {
		ptr := (*byte)(unsafe.Pointer(*memPtr + (uintptr(addr + i))))
		*ptr = buf[i]
	}
}

func ValueGet(retval int32, vAddr int32, propLen int32, propPtr int32, what_is_this int32) {
	//const v = loadValue(v_addr);
	//const p = loadString(p_ptr, p_len);
	//const x = loadValue(x_addr);
	//Reflect.set(v, p, x);
	prop := strConvert(*memPtr, propPtr, propLen)
	log.Printf("Value Get: %s", prop)
	value := loadValue(vAddr)
	result := value // this should be a Reflect.get(value,prop)
	storeValue(retval, result)
}
