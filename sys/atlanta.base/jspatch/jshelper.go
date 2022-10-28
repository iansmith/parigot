package jspatch

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

type JSPatch struct {
	mem *WasmMem
}

func NewJSPatchWithMemPtr(memPtr uintptr) *JSPatch {
	return &JSPatch{
		mem: NewWasmMem(memPtr),
	}
}

func NewJSPatch() *JSPatch {
	return &JSPatch{}
}

func (j *JSPatch) SetMemPtr(m uintptr) {
	j.mem = NewWasmMem(m)
}

const verbose = false
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
	index := j.mem.GetInt64(sp + 16)
	value := j.mem.LoadValue(sp + 8)
	result := value.Index(int(index))
	wrapped := goToJS(result)
	j.mem.StoreValue(sp+32, wrapped)
}

func (j *JSPatch) ValueSetIndex(sp int32) {
	enter("ValueSetIndex")
	index := j.mem.GetInt64(sp + 16)
	value := j.mem.LoadValue(sp + 8)
	newValue := j.mem.LoadValue(sp + 24)
	value.SetIndex(int(index), newValue)
}

func (j *JSPatch) ValueLoadString(sp int32) {
	enter("ValueLoadIndex")
	str := j.mem.LoadValue(sp + 8)
	slice := j.mem.LoadSlice(sp + 16)
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	content := str.String()
	ptr.Data = uintptr(unsafe.Pointer((&content)))
	ptr.Len = len(content)
	ptr.Cap = len(content)
}

func (j *JSPatch) ValueLength(sp int32) {
	enter("ValueLength")
	v := j.mem.LoadValue(sp + 8)
	l := v.Length()
	j.mem.SetInt64(sp+16, int64(l))
}

func (j *JSPatch) ValueInstanceOf(sp int32) {
	enter("ValueInstanceOf")
	example := j.mem.LoadValue(sp + 16)
	candidate := j.mem.LoadValue(sp + 8)
	if example.InstanceOf(candidate) {
		j.mem.SetUint8(sp+24, 1)
	} else {
		j.mem.SetUint8(sp+24, 0)
	}
}

func (j *JSPatch) ValuePrepareString(sp int32) {
	panic("value prepare string not implemented (multiple encodings)")
}

func (j *JSPatch) FinalizeRef(sp int32) {
	enter("FinalizeRef")
	panic("finalize not implemented (multiple encodings)")
}

func (j *JSPatch) ValueCall(sp int32) {
	enter("ValueCall")
	// can go switch the stack on us here?
	recvr := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	args := j.mem.LoadSliceOfValues(sp + 32)
	v := recvr.Call(prop, args)
	j.mem.StoreValue(sp+56, goToJS(v))
	j.mem.SetUint8(sp+64, 1)
}

func (j *JSPatch) StringVal(sp int32) {
	enter("StringVal")
	s := j.mem.LoadString(sp + 8)
	obj := goToJS(s)
	j.mem.StoreValue(sp+24, obj)
}

//func (j *JSPatch) StrConvert(memPtr uintptr, ptr int32, length int32) string {
//	// we could probably go bytesConvert and claim our cap was equal to our len but...
//	enter("StrConvert")
//	buf := make([]byte, length)
//	for i := int32(0); i < length; i++ {
//		b := (*byte)(unsafe.Pointer(memPtr + uintptr(ptr+i)))
//		buf[i] = *b
//	}
//	s := string(buf)
//	return s
//}

func (j *JSPatch) ValueGet(sp int32) {
	enter("ValueGet")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	v := value.Get(prop)
	j.mem.StoreValue(sp+32, v)
}

func (j *JSPatch) ValueInvoke(sp int32) {
	enter("ValueInvoke")
	// can go switch the stack on us here?
	callee := j.mem.LoadValue(sp + 8)
	args := j.mem.LoadSliceOfValues(sp + 16)
	result := callee.Invoke(args)
	j.mem.StoreValue(sp+40, goToJS(result))
	j.mem.SetUint8(sp+48, 1)
}
func (j *JSPatch) ValueNew(sp int32) {
	// can go switch the stack on us here?
	enter("ValueNew")
	v := j.mem.LoadValue(sp + 8)
	args := j.mem.LoadSliceOfValues(sp + 16)
	panic("don't know how to call constructor for ValueNew:" +
		fmt.Sprintf("v=%v,args=%v", v, args))

}
func (j *JSPatch) ValueSet(sp int32) {
	enter("ValueSet")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	newPropVal := j.mem.LoadValue(sp + 32)
	value.Set(prop, newPropVal)
}

func (j *JSPatch) ValueDelete(sp int32) {
	enter("ValueDelete")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	value.Delete(prop)
}

func (j *JSPatch) CopyBytesToGo(sp int32) {
	panic("CopyBytesToGo not implemented")
}
func (j *JSPatch) CopyBytesToJS(sp int32) {
	panic("CopyBytesToJS not implemented")
}
