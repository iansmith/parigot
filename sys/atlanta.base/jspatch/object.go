package jspatch

import (
	"fmt"
	"math"
	"unsafe"
)

type jsKind int
type binKind uint32

const (
	typeFlagNone binKind = iota
	typeFlagObject
	typeFlagString
	typeFlagSymbol
	typeFlagFunction
)
const (
	TypeUndefined jsKind = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypeFunction
	TypeObject
)

var undefined = &jsObj{typeFlag: TypeUndefined}

func floatValue(f float64) jsObject {
	if f == 0 {
		return object.get(predefinedZero)
	}
	if f != f {
		return object.get(predefinedNan)
	}
	return newJSObjNum(f)
}

type jsObject interface {
	id() int32
	Index(int) jsObject
	SetIndex(int, jsObject)
	String() string
	Length() int
	InstanceOf(jsObject) bool
	Call(string, jsObject /*[]jsObject*/) jsObject
	Get(string) jsObject
	Invoke(jsObject /*[]jsObject*/) jsObject
	Set(string, jsObject) jsObject /* returns self for chaining*/
	Delete(string)
	// for our internal use
	binaryRep() (uint32, uint32)
	isNumber() bool
}

var object = newObjMap()

const (
	predefinedNan    = 0
	predefinedZero   = 1
	predefinedNull   = 2
	predefinedTrue   = 3
	predefinedFalse  = 4
	predefinedWindow = 5 // globals
	predefinedGo     = 6
)

func init() {
	// order here is carefully constructed
	object.predefined(newJSObjNum(math.NaN()), predefinedNan)

	obj := newJSObjNum(0.0)
	obj.zero = true
	object.predefined(obj, predefinedZero)

	obj = newJSObjGeneric(TypeNull)
	obj.null = true
	object.predefined(obj, predefinedNull)

	obj = newJSObjBool(true)
	object.predefined(obj, predefinedTrue)

	obj = newJSObjBool(false)
	object.predefined(obj, predefinedFalse)

	glob := newJSObjObject(true)
	object.predefined(glob, predefinedWindow)

	obj = newJSObjObject(false)
	object.predefined(obj, predefinedGo)

	obj = newJSObjFunc(makeObject)
	object.put(obj)
	glob.Set("Object", obj)

	obj = newJSObjFunc(makeArray)
	glob.Set("Array", obj)
	object.put(obj)

	fs := newJSObjGeneric(TypeObject)
	glob.Set("fs", fs)
	object.put(fs)

	obj = newJSObjGeneric(TypeObject)
	glob.Set("process", obj)
	object.put(obj)

	constants := object.put(newJSObjGeneric(TypeObject))
	minus1 := newJSObjNum(-1)
	//{ O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused	constants.Set()
	constants.Set("O_WRONLY", minus1).Set("O_RDWR", minus1).Set("O_TRUNC", minus1).
		Set("O_APPEND", minus1).Set("O_EXCL", minus1).Set("O_CREAT", minus1)
	fs.Set("constants", constants)

	glob.Set("fs", fs)
	glob.Set("Uint8Array", object.put(newJSObjFunc(makeUint8Array)))
}

func makeUint8Array() {
	panic("makeUint8Array called")
}

func makeObject() {
	panic("makeObject from global called")
}
func makeArray() {
	panic("makeArray from global called")
}

func goToJS(x any) jsObject {
	switch x := x.(type) {
	case jsObject:
		return x
	//case Func:
	//	return x.Value
	case nil:
		return object.get(predefinedNull)
	case bool:
		if x {
			return object.get(predefinedTrue)
		} else {
			return object.get(predefinedFalse)
		}
	case int:
		obj := floatValue(float64(x))
		obj.(*jsObj).int_ = x
		return obj
	case int8:
		obj := floatValue(float64(x))
		obj.(*jsObj).i8 = x
		return obj
	case int16:
		obj := floatValue(float64(x))
		obj.(*jsObj).i16 = x
		return obj
	case int32:
		obj := floatValue(float64(x))
		obj.(*jsObj).i32 = x
		return obj
	case int64:
		obj := floatValue(float64(x))
		obj.(*jsObj).i64 = x
		return obj
	case uint:
		obj := floatValue(float64(x))
		obj.(*jsObj).uint_ = x
		return obj
	case uint8:
		obj := floatValue(float64(x))
		obj.(*jsObj).u8 = x
		return obj
	case uint16:
		obj := floatValue(float64(x))
		obj.(*jsObj).u16 = x
		return obj
	case uint32:
		obj := floatValue(float64(x))
		obj.(*jsObj).u32 = x
		return obj
	case uint64:
		obj := floatValue(float64(x))
		obj.(*jsObj).u64 = x
		return obj
	case uintptr:
		obj := floatValue(float64(x))
		obj.(*jsObj).uptr = x
		return obj
	case unsafe.Pointer:
		obj := floatValue(float64(uintptr(x)))
		obj.(*jsObj).uptr = uintptr(x)
		return obj
	case float32:
		return floatValue(float64(x))
	case float64:
		return floatValue(x)
	case string:
		return object.put(newJSObjString(x))
	case []any:
		panic("don't know how to convert an array")
		//a := arrayConstructor.New(len(x))
		//for i, s := range x {
		//	a.SetIndex(i, s)
		//}
		//return a
	case map[string]any:
		val := newJSObjGeneric(TypeObject)
		for k, v := range x {
			val.Set(k, goToJS(v))
		}
		return object.put(val)
	default:
		panic("ValueOf: invalid value")
	}
}

// num
func newJSObjNum(f float64) *jsObj {
	return &jsObj{
		n:        -1001,
		prop:     make(map[string]jsObject),
		typeFlag: TypeNumber,
		f64:      f,
	}
}

// flag object
func newJSObjObject(isGlobal bool) *jsObj {
	return &jsObj{
		n:        -5005,
		prop:     make(map[string]jsObject),
		typeFlag: TypeObject,
		global:   isGlobal,
	}
}

// func
func newJSObjFunc(fn func()) *jsObj {
	return &jsObj{
		n:        -4004,
		prop:     make(map[string]jsObject),
		typeFlag: TypeFunction,
		fn:       fn,
	}
}

// bool
func newJSObjBool(b bool) *jsObj {
	return &jsObj{
		n:        -3003,
		prop:     make(map[string]jsObject),
		typeFlag: TypeBoolean,
		b:        b,
	}
}

// generic object
func newJSObjGeneric(t jsKind) *jsObj {
	return &jsObj{
		n:        -2002,
		prop:     make(map[string]jsObject),
		typeFlag: t,
	}
}
func newJSObjString(s string) *jsObj {
	return &jsObj{
		n:        -6006,
		prop:     make(map[string]jsObject),
		typeFlag: TypeString,
		s:        s,
	}
}

type jsObj struct {
	n        int32
	prop     map[string]jsObject
	typeFlag jsKind
	b        bool
	f32      float32
	f64      float64
	i8       int8
	i16      int16
	i32      int32
	i64      int64
	uint_    uint
	u8       uint8
	u16      uint16
	u32      uint32
	u64      uint64
	uptr     uintptr
	int_     int
	null     bool
	zero     bool
	global   bool
	s        string
	fn       func()
}

func (j *jsObj) id() int32 {
	return j.n
}

func (j *jsObj) String() string {
	return fmt.Sprintf("[%d,%+v]", j.n, j.prop)
}

func (j *jsObj) Index(int) jsObject {
	panic("Index")
}
func (j *jsObj) SetIndex(int, jsObject) {
	panic("SetIndex")

}
func (j *jsObj) Length() int {
	panic("Length")

}
func (j *jsObj) InstanceOf(o jsObject) bool {
	other := o.(*jsObj)
	if j.id() == other.id() {
		return true
	}
	return j.typeFlag == other.typeFlag
}
func (j *jsObj) Call(string, jsObject /* []jsObject*/) jsObject {
	panic("Call")
}
func (j *jsObj) Get(propName string) jsObject {
	v, ok := j.prop[propName]
	if ok {
		return v
	}
	return undefined
}
func (j *jsObj) Invoke(jsObject /*[]jsObject*/) jsObject {
	panic("Invoke")
}
func (j *jsObj) Set(propName string, obj jsObject) jsObject {
	j.prop[propName] = obj
	return j
}
func (j *jsObj) Delete(propName string) {
	delete(j.prop, propName)
}

func (j *jsObj) isNumber() bool {
	return j.typeFlag == TypeNumber
}

func (j *jsObj) BinaryTypeFlag() uint32 {
	switch j.typeFlag {
	case TypeUndefined, TypeNull, TypeBoolean, TypeNumber:
		return uint32(typeFlagNone)
	case TypeObject:
		return uint32(typeFlagObject)
	case TypeString:
		return uint32(typeFlagString)
	case TypeSymbol:
		return uint32(typeFlagSymbol)
	case TypeFunction:
		return uint32(typeFlagFunction)
	}
	panic("unknown type found in jsObj")
}

func (j *jsObj) binaryRep() (uint32, uint32) {
	switch j.typeFlag {
	case TypeUndefined:
		return 0, 0
	case TypeNumber:
		bits := math.Float64bits(j.f64)
		low := uint32(bits & 0xffffffff)
		high := uint32(bits >> 32)
		return high, low
	}
	// normal case is the silly encoding business with NaN
	high := nanHead | j.BinaryTypeFlag()
	low := uint32(j.n)
	return high, low
}

//
//func (j *jsObj) FuncAddr() func() {
//	if j.typeFlag != TypeFunction {
//		panic("attempt to get function address of non function")
//	}
//	return j.fn
//}
//
//func (j *jsObj) IsFunc() bool {
//	return j.typeFlag == TypeFunction
//}

//
// OBJMAP
//

type objMap struct {
	next int32
	data map[int32]jsObject
}

func newObjMap() *objMap {
	return &objMap{
		next: 0,
		data: make(map[int32]jsObject),
	}
}

// returns the object it just put in the map
func (o *objMap) put(obj jsObject) jsObject {
	if obj.id() >= 0 {
		panic("attempt to put object in map second time " + fmt.Sprint(obj.id()))
	}
	t := obj.(*jsObj).typeFlag
	if t == TypeNumber || t == TypeUndefined {
		panic("should not be putting numbers or undefined into the object map")
	}
	obj.(*jsObj).n = o.next
	o.next++
	o.data[obj.id()] = obj
	return obj
}

func (o *objMap) predefined(obj jsObject, expected int32) jsObject {
	obj.(*jsObj).n = expected
	o.data[expected] = obj
	if expected > o.next {
		o.next = expected + 1
	}
	return obj
}

func (o *objMap) get(id int32) jsObject {
	if id < 0 {
		panic("attempt to get object from map that was never entered " + fmt.Sprint(id))
	}
	v, ok := o.data[id]
	if !ok {
		panic("unknown object id " + fmt.Sprint(id))
	}
	return v
}
