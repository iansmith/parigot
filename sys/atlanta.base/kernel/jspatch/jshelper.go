package jspatch

import (
	"log"
	"os"
	"unsafe"
)

type JSPatch struct {
	mem *wasmMem
}

func NewJSPatch(memPtr *uintptr) *JSPatch {
	return &JSPatch{
		mem: newWasmMem(*memPtr),
	}
}

const nanHead = 0x7FF80000

func (j *JSPatch) ValueIndex(sp int32) {
	index := j.mem.getInt64(sp + 16)
	value := j.mem.loadValue(sp + 8)
	result := value.getIndex(index)
	j.mem.storeValue(sp+32, result)
}

func (j *JSPatch) ValueSetIndex(sp int32) {
	index := j.mem.getInt64(sp + 16)
	value := j.mem.loadValue(sp + 8)
	newValue := j.mem.loadValue(sp + 24)
	value.setIndex(index, newValue)
}

func (j *JSPatch) ValuePrepareString(v int32, x int32, i int32) {
	log.Printf("Value prepare string %x,%x %d", v, x, i)
	os.Exit(1)
}

func (j *JSPatch) ValueLoadString(result int32, v int32, slice int32, len int32, cap int32) {
	log.Printf("Value Load String %x,(%x %d,%d)", v, slice, len, cap)
	os.Exit(1)
}

func (j *JSPatch) FinalizeRef(sp int32) {
	v := j.mem.getInt32(sp + 8)
	log.Printf("finalize ref on obj", v)
	count, ok := refCount[v]
	if !ok {
		log.Printf("unable to decrement refcount on %d", v)
	} else {
		refCount[v] = count - 1
	}
}
func (j *JSPatch) ValueNew(result int32, v int32, args int32, args_len int32, args_cap int32, qqq int32) {
	log.Printf("ValueNew(%x,%x,%x,%d,%d) ???%x", result, v, args, args_len, args_cap, qqq)

	os.Exit(1)
}
func (j *JSPatch) ValueCall(sp int32) {
	// can go switch the stack on us here?
	recvr := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	method := recvr.getProp(prop)
	args := j.mem.loadSliceOfValues(sp + 32)
	result := recvr.apply(method, args)
	j.storeValue(sp+56, result)
	j.mem.setUint8(sp+64, 1)

}

func (j *JSPatch) StringVal(retVal int32, ptr int32, len int32, wtf int32) {
	log.Printf("js.stringVal called %x,%x,%d,%x",
		retVal, ptr, len, wtf)
	os.Exit(1)
}

func (j *JSPatch) StrConvert(memPtr uintptr, ptr int32, length int32) string {
	// we could probably go bytesConvert and claim our cap was equal to our len but...
	buf := make([]byte, length)
	for i := int32(0); i < length; i++ {
		b := (*byte)(unsafe.Pointer(memPtr + uintptr(ptr+i)))
		buf[i] = *b
	}
	s := string(buf)
	return s
}

func (j *JSPatch) ValueGet(sp int32) {
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	v := value.getProp(prop)
	log.Printf("value get of %d.%s=>%v", value.id(), prop, v.id())
	j.mem.storeValue(sp+32, v)
}

func (j *JSPatch) ValueSet(sp int32) {
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	newPropVal := j.mem.loadValue(sp + 32)
	value.setProp(prop, newPropVal)
}

func (j *JSPatch) ValueDelete(sp int32) {
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	value.deleteProp(prop)
}
