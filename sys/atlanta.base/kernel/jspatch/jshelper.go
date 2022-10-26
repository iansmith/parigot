package jspatch

import (
	"log"
	"reflect"
	"unsafe"
)

type JSPatch struct {
	mem *wasmMem
}

func NewJSPatchWithMemPtr(memPtr uintptr) *JSPatch {
	return &JSPatch{
		mem: newWasmMem(memPtr),
	}
}

func NewJSPatch() *JSPatch {
	return &JSPatch{}
}

func (j *JSPatch) SetMemPtr(m uintptr) {
	j.mem = newWasmMem(m)
}

const verbose = true
const nanHead = 0x7FF80000

func enter(funcName string, rest ...string) {
	if !verbose {
		return
	}
	result := ""
	for i, s := range rest {
		result += s
		if i != len(rest)-1 {
			result += ","
		}
	}
	if len(rest) > 0 {
		log.Printf("---- entering  %s ---- [%s]", funcName, result)
	} else {
		log.Printf("---- entering %s ----", funcName)
	}
}

func (j *JSPatch) ValueIndex(sp int32) {
	enter("ValueIndex")
	index := j.mem.getInt64(sp + 16)
	value := j.mem.loadValue(sp + 8)
	result := value.getIndex(index)
	j.mem.storeValue(sp+32, result)
}

func (j *JSPatch) ValueSetIndex(sp int32) {
	enter("ValueSetIndex")
	index := j.mem.getInt64(sp + 16)
	value := j.mem.loadValue(sp + 8)
	newValue := j.mem.loadValue(sp + 24)
	value.setIndex(index, newValue)
}

func (j *JSPatch) ValueLoadString(sp int32) {
	enter("ValueLoadIndex")
	str := j.mem.loadValue(sp + 8)
	slice := j.mem.loadSlice(sp + 16)
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	content := str.string()
	ptr.Data = uintptr(unsafe.Pointer((&content)))
	ptr.Len = int(str.length())
	ptr.Cap = int(str.length())
}

func (j *JSPatch) ValueLength(sp int32) {
	enter("ValueLength")
	v := j.mem.loadValue(sp + 8)
	j.mem.setInt64(sp+16, v.length())
}

func (j *JSPatch) ValueInstanceOf(sp int32) {
	enter("ValueInstanceOf")
	example := j.mem.loadValue(sp + 16)
	candidate := j.mem.loadValue(sp + 8)
	if (candidate.isArray() && example.isArray()) ||
		(candidate.isCallable() && example.isCallable()) ||
		(candidate.isString() && example.isString()) {
		j.mem.setUint8(sp+24, 1)
	} else {
		j.mem.setUint8(sp+24, 0)
	}
}

func (j *JSPatch) ValuePrepareString(sp int32) {
	panic("value prepare string not implemented (multiple encodings)")
}

func (j *JSPatch) FinalizeRef(sp int32) {
	enter("FinalizeRef")
	v := j.mem.getInt32(sp + 8)
	log.Printf("finalize ref on obj", v)
	count, ok := refCount[v]
	if !ok {
		log.Printf("unable to decrement refcount on %d", v)
	} else {
		refCount[v] = count - 1
	}
}

func (j *JSPatch) ValueCall(sp int32) {
	enter("ValueCall")
	// can go switch the stack on us here?
	recvr := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	method := recvr.getProp(prop)
	args := j.mem.loadSliceOfValues(sp + 32)
	result := recvr.apply(method, args)
	j.mem.storeValue(sp+56, result)
	j.mem.setUint8(sp+64, 1)

}

func (j *JSPatch) StringVal(sp int32) {
	enter("StringVal")
	s := j.mem.loadString(sp + 8)
	obj := newJSObjString(nextId(), s)
	j.mem.storeValue(sp+24, obj)
}

func (j *JSPatch) StrConvert(memPtr uintptr, ptr int32, length int32) string {
	// we could probably go bytesConvert and claim our cap was equal to our len but...
	enter("StrConvert")
	buf := make([]byte, length)
	for i := int32(0); i < length; i++ {
		b := (*byte)(unsafe.Pointer(memPtr + uintptr(ptr+i)))
		buf[i] = *b
	}
	s := string(buf)
	return s
}

func (j *JSPatch) ValueGet(sp int32) {
	enter("ValueGet")
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	v := value.getProp(prop)
	log.Printf("value get of %d.%s=>%v", value.id(), prop, v.id())
	j.mem.storeValue(sp+32, v)
}

func (j *JSPatch) ValueInvoke(sp int32) {
	enter("ValueInvoke")
	// can go switch the stack on us here?
	callee := j.mem.loadValue(sp + 8)
	args := j.mem.loadSliceOfValues(sp + 16)
	result := callee.call(args)
	j.mem.storeValue(sp+40, result)
	j.mem.setUint8(sp+48, 1)
}
func (j *JSPatch) ValueNew(sp int32) {
	// can go switch the stack on us here?
	enter("ValueNew")
	v := j.mem.loadValue(sp + 8)
	args := j.mem.loadSliceOfValues(sp + 16)
	result := v.construct(args)
	j.mem.storeValue(sp+40, result)
	j.mem.setUint8(sp+48, 1)

}
func (j *JSPatch) ValueSet(sp int32) {
	enter("ValueSet")
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	newPropVal := j.mem.loadValue(sp + 32)
	value.setProp(prop, newPropVal)
}

func (j *JSPatch) ValueDelete(sp int32) {
	enter("ValueDelete")
	value := j.mem.loadValue(sp + 8)
	prop := j.mem.loadString(sp + 16)
	value.deleteProp(prop)
}

func (j *JSPatch) CopyBytesToGo(sp int32) {
	panic("CopyBytesToGo not implemented")
}
func (j *JSPatch) CopyBytesToJS(sp int32) {
	panic("CopyBytesToJS not implemented")
}
