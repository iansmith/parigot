package jspatch

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/sys/jspatch/jsemul"
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

const jsEmulVerbose = true

func jsLog(spec string, x ...interface{}) {
	if jsEmulVerbose {
		s := fmt.Sprintf(spec, x...)
		log.Print(s)
	}
}

func (j *JSPatch) ValueIndex(sp int32) {
	jsLog("ValueIndex")
	index := j.mem.GetInt64(sp + 16)
	value := j.mem.LoadValue(sp + 8)
	result := value.Index(int(index))
	wrapped := jsemul.GoToJS(result)
	j.mem.StoreValue(sp+32, wrapped)
}

func (j *JSPatch) ValueSetIndex(sp int32) {
	jsLog("ValueSetIndex")
	index := j.mem.GetInt64(sp + 16)
	value := j.mem.LoadValue(sp + 8)
	newValue := j.mem.LoadValue(sp + 24)
	value.SetIndex(int(index), newValue)
}

func (j *JSPatch) ValueLoadString(sp int32) {
	jsLog("ValueLoadIndex")
	str := j.mem.LoadValue(sp + 8)
	slice := j.mem.LoadSlice(sp + 16)
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	content := str.String()
	// XXX this has to be changed between tinygo and gc
	// XXX uintptr(len(content)) -> len(content)
	l := uintptr(len(content))
	ptr.Len = l
	ptr.Cap = l
	ptr.Data = uintptr(unsafe.Pointer((&content)))
}

func (j *JSPatch) ValueLength(sp int32) {
	jsLog("ValueLength")
	v := j.mem.LoadValue(sp + 8)
	l := v.Length()
	j.mem.SetInt64(sp+16, int64(l))
}

func (j *JSPatch) ValueInstanceOf(sp int32) {
	jsLog("ValueInstanceOf")
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
	//jsLog("FinalizeRef")
	panic("finalize not implemented (multiple encodings)")
}

//go:noinline
func (j *JSPatch) ValueCall(sp int32) {
	//jsLog("ValueCall")
	// can go switch the stack on us here?
	inst := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	args := j.mem.LoadSliceOfValues(sp + 32)
	v := inst.Call(prop, args)
	j.mem.StoreValue(sp+56, v) //xxx should we be doing this without any conversion?
	j.mem.SetUint8(sp+64, 1)
}

func (j *JSPatch) StringVal(sp int32) {
	//jsLog("StringVal")
	s := j.mem.LoadString(sp + 8)
	obj := jsemul.GoToJS(s)
	j.mem.StoreValue(sp+24, obj)
}

func (j *JSPatch) ValueGet(sp int32) {
	//jsLog("ValueGet")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	v := value.Get(prop)
	j.mem.StoreValue(sp+32, v)
}

func (j *JSPatch) ValueInvoke(sp int32) {
	//jsLog("ValueInvoke")
	// can go switch the stack on us here?
	callee := j.mem.LoadValue(sp + 8)
	args := j.mem.LoadSliceOfValues(sp + 16)
	result := callee.Invoke(args)
	j.mem.StoreValue(sp+40, jsemul.GoToJS(result))
	j.mem.SetUint8(sp+48, 1)
}

func (j *JSPatch) ValueNew(sp int32) {
	// can go switch the stack on us here?
	//jsLog("ValueNew")
	v := j.mem.LoadValue(sp + 8)
	if !v.IsClassObject() {
		panic(fmt.Sprintf("attempt to call valueNew() on something that is not a class, %s", v))
	}
	var args []jsemul.JsObject
	if !j.mem.TestSliceIsZeroLen(sp + 16) {
		args = j.mem.LoadSliceOfValues(sp + 16)
	}
	result := v.NewInstance(args)
	j.mem.StoreValue(sp+40, result)
	j.mem.SetUint8(sp+48, 1) // should we be setting just 1 byte? on a stack?
}
func (j *JSPatch) ValueSet(sp int32) {
	//jsLog("ValueSet")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	newPropVal := j.mem.LoadValue(sp + 32)
	value.Set(prop, newPropVal)
}

func (j *JSPatch) ValueDelete(sp int32) {
	//jsLog("ValueDelete")
	value := j.mem.LoadValue(sp + 8)
	prop := j.mem.LoadString(sp + 16)
	value.Delete(prop)
}

func (j *JSPatch) CopyBytesToGo(sp int32) {
	panic("CopyBytesToGo not implemented")
}
func (j *JSPatch) CopyBytesToJS(sp int32) {
	dst := j.mem.LoadValue(sp + 8)
	var src []byte
	l := j.mem.GetInt64(sp + 16)
	if l != 0 {
		src = j.mem.LoadSliceWithKnownLength(sp+16, l)
	}
	if !dst.InstanceOf(jsemul.Uint8ArrayObj) {
		panic("CopyBytesToJS called, but dst is not a uint8ArrayObj")
	}
	trueDest := dst.This().(*jsemul.Uint8ArrayInstance)
	shorter := trueDest.CopyBytesToJS(src)
	j.mem.SetInt64(sp+40, int64(shorter))
	j.mem.SetUint8(sp+48, 1)
}
