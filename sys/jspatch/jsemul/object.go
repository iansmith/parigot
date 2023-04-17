package jsemul

import (
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"unsafe"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/jspatch/jsgo"

	"google.golang.org/protobuf/types/known/timestamppb"
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

var Undefined = &jsObj{typeFlag: TypeUndefined}

func FloatValue(f float64) JsObject {
	if f == 0 {
		return JsObjectMap.Get(predefinedZero)
	}
	if f != f {
		return JsObjectMap.Get(predefinedNan)
	}
	return newJSObjNum(f)
}

type JsObject interface {
	Id() int32
	Index(int) JsObject
	SetIndex(int, JsObject)
	String() string
	Length() int
	InstanceOf(JsObject) bool
	Call(string /*name*/, []JsObject) JsObject
	Get(string) JsObject
	Invoke([]JsObject) JsObject
	Set(string, JsObject) JsObject /* returns self for chaining*/
	Delete(string)
	IsFunc() bool
	IsClassObject() bool
	IsInstance() bool
	IsNumber() bool
	AsNumber() float64
	BinaryRep() uint64
	// from a class, create an instance
	NewInstance([]JsObject) JsObject
	// returns the underlying object from an instance (only!) of a class
	This() interface{}
}

var JsObjectMap = newObjMap()

const (
	predefinedNan    = 0
	predefinedZero   = 1
	predefinedNull   = 2
	predefinedTrue   = 3
	predefinedFalse  = 4
	predefinedWindow = 5 // globals
	predefinedGo     = 6
)

type class interface {
	NewInstance([]JsObject) JsObject
	Name() string
}

var badWrappedFunc JsObject

func init() {
	// this is caused because this code runs on linux but works on behalf of WASM
	// and so it ends up being loaded (and init called) twice, once for the runner
	// in linux and once for the module loaded into the WASM side
	//print("xxx init called on object.go ... ", runtime.GOOS, "\n")
	// order here is carefully constructed
	JsObjectMap.predefined(newJSObjNum(math.NaN()), predefinedNan)

	obj := newJSObjGeneric(typeFlagNone, TypeObject, "ZERO")
	obj.zero = true
	obj.typeFlag = jsKind(typeFlagNone)
	JsObjectMap.predefined(obj, predefinedZero)

	obj = newJSObjGeneric(typeFlagNone, TypeNull, "NULL")
	obj.null = true
	JsObjectMap.predefined(obj, predefinedNull)

	obj = newJSObjBool(true)
	JsObjectMap.predefined(obj, predefinedTrue)

	obj = newJSObjBool(false)
	JsObjectMap.predefined(obj, predefinedFalse)

	glob := newJSObjObject(true, "WINDOW")
	JsObjectMap.predefined(glob, predefinedWindow)

	goObj := newJSObjObject(false, "GO")
	JsObjectMap.predefined(goObj, predefinedGo)

	obj = newJSObjFunc(makeFuncWrapper, "_makeFuncWrapper")
	JsObjectMap.Put(obj)
	goObj.Set("_makeFuncWrapper", obj)

	// obj = newJSObjFunc(makeArray, "makeArray")
	// glob.Set("Array", obj)
	// jsObjectMap.put(obj)

	dateObj = newJSObjClass(&dateClass{offsetFromUTCInMins: 0}, "class(Date)")
	glob.Set("Date", dateObj)
	JsObjectMap.Put(dateObj)

	fs := newJSObjGeneric(typeFlagObject, TypeObject, "fs")
	glob.Set("fs", fs)
	JsObjectMap.Put(fs)
	writeFn := newJSObjFunc(fsWrite, "fs.Write")
	JsObjectMap.Put(writeFn)
	fs.Set("write", writeFn)

	obj = newJSObjGeneric(typeFlagObject, TypeObject, "process")
	glob.Set("process", obj)
	JsObjectMap.Put(obj)

	constants := JsObjectMap.Put(newJSObjGeneric(typeFlagObject, TypeObject, "constants"))
	minus1 := newJSObjNum(-1)

	//{ O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused	constants.Set()
	constants.Set("O_WRONLY", minus1).Set("O_RDWR", minus1).Set("O_TRUNC", minus1).
		Set("O_APPEND", minus1).Set("O_EXCL", minus1).Set("O_CREAT", minus1)
	fs.Set("constants", constants)

	glob.Set("fs", fs)

	Uint8ArrayObj = newJSObjClass(&uint8ArrayClass{}, "class(Uint8Array)")
	JsObjectMap.Put(Uint8ArrayObj)
	glob.Set("Uint8Array", Uint8ArrayObj)

	objectObj = newJSObjClass(&objectClass{}, "class(Object)")
	JsObjectMap.Put(objectObj)
	glob.Set("Object", objectObj)

	badWrappedFunc = newJSObjFunc(func(_ []JsObject) JsObject {
		print("calling wrapped go functions from javascript makes no sense when there is no javascript")
		return JsObjectMap.Get(predefinedNull)
	}, "bad function call")
	JsObjectMap.Put(badWrappedFunc)
}

func GoToJS(x any) JsObject {
	switch x := x.(type) {
	case JsObject:
		return x
	case func(interface{}):
		print("got to func case in goToJS\n")
		return nil
	//	return x.Value
	case nil:
		return JsObjectMap.Get(predefinedNull)
	case bool:
		if x {
			return JsObjectMap.Get(predefinedTrue)
		} else {
			return JsObjectMap.Get(predefinedFalse)
		}
	case int:
		obj := FloatValue(float64(x))
		obj.(*jsObj).int_ = x
		return obj
	case int8:
		obj := FloatValue(float64(x))
		obj.(*jsObj).i8 = x
		return obj
	case int16:
		obj := FloatValue(float64(x))
		obj.(*jsObj).i16 = x
		return obj
	case int32:
		obj := FloatValue(float64(x))
		obj.(*jsObj).i32 = x
		return obj
	case int64:
		obj := FloatValue(float64(x))
		obj.(*jsObj).i64 = x
		return obj
	case uint:
		obj := FloatValue(float64(x))
		obj.(*jsObj).uint_ = x
		return obj
	case uint8:
		obj := FloatValue(float64(x))
		obj.(*jsObj).u8 = x
		return obj
	case uint16:
		obj := FloatValue(float64(x))
		obj.(*jsObj).u16 = x
		return obj
	case uint32:
		obj := FloatValue(float64(x))
		obj.(*jsObj).u32 = x
		return obj
	case uint64:
		obj := FloatValue(float64(x))
		obj.(*jsObj).u64 = x
		return obj
	case uintptr:
		obj := FloatValue(float64(x))
		obj.(*jsObj).uptr = x
		return obj
	case unsafe.Pointer:
		obj := FloatValue(float64(uintptr(x)))
		obj.(*jsObj).uptr = uintptr(x)
		return obj
	case float32:
		return FloatValue(float64(x))
	case float64:
		return FloatValue(x)
	case string:
		return JsObjectMap.Put(newJSObjString(x))
	case []any:
		panic("don't know how to convert an array")
		//a := arrayConstructor.New(len(x))
		//for i, s := range x {
		//	a.SetIndex(i, s)
		//}
		//return a
	case map[string]any:
		val := newJSObjGeneric(typeFlagNone, TypeObject, "converted props")
		for k, v := range x {
			val.Set(k, GoToJS(v))
		}
		return JsObjectMap.Put(val)
	default:
		print("--------------DUMP FROM JS LAND--------------\n")
		debug.PrintStack()
		print("--------------END FROM JS LAND--------------\n")
		v := reflect.ValueOf(x)
		switch v.Kind() {
		case reflect.Slice:
			//ev := reflect.ValueOf(v.Interface())
			sl := v.Interface().([]jsObj)
			print("found a slice,xxx,", v.Len(), " and ", v.Index(0).Kind().String(), " -- ", len(sl))
			arr := newJSObjArr(v.Len())
			for i := 0; i < v.Len(); i++ {
				arr.SetIndex(i, &sl[i])
			}
			return arr
		default:
			panic("we don't know what to do with type " + v.Kind().String())
		}
	}
}

// num
func newJSObjNum(f float64) *jsObj {
	return &jsObj{
		n:         -1001,
		prop:      make(map[string]JsObject),
		typeFlag:  TypeNumber,
		f64:       f,
		debugName: fmt.Sprint(f),
	}
}

// class
func newJSObjClass(class class, debugName string) *jsObj {
	return &jsObj{
		n:         -7007,
		prop:      make(map[string]JsObject),
		typeFlag:  TypeObject,
		class:     class,
		binFlag:   typeFlagNone,
		debugName: debugName,
	}
}

// flag object
func newJSObjObject(isGlobal bool, debugName string) *jsObj {
	return &jsObj{
		n:         -5005,
		prop:      make(map[string]JsObject),
		typeFlag:  TypeObject,
		global:    isGlobal,
		binFlag:   typeFlagNone,
		debugName: debugName,
	}
}

// func
func newJSObjFunc(fn func([]JsObject) JsObject, debugName string) *jsObj {
	return &jsObj{
		n:         -4004,
		prop:      make(map[string]JsObject),
		typeFlag:  TypeFunction,
		fn:        fn,
		binFlag:   typeFlagFunction,
		debugName: debugName,
	}
}

// bool
func newJSObjBool(b bool) *jsObj {
	return &jsObj{
		n:         -3003,
		prop:      make(map[string]JsObject),
		typeFlag:  TypeBoolean,
		b:         b,
		binFlag:   typeFlagNone,
		debugName: fmt.Sprint(b),
	}
}

// generic object
func newJSObjGeneric(b binKind, t jsKind, debugName string) *jsObj {
	return &jsObj{
		n:         -2002,
		prop:      make(map[string]JsObject),
		typeFlag:  t,
		binFlag:   b,
		debugName: debugName,
	}
}
func newJSObjString(s string) *jsObj {
	return &jsObj{
		n:        -6006,
		prop:     make(map[string]JsObject),
		typeFlag: TypeString,
		s:        s,
		binFlag:  typeFlagString,
	}
}
func newJSObjArr(size int) *jsObj {
	o := &jsObj{
		n:        -8008,
		prop:     make(map[string]JsObject),
		typeFlag: TypeObject,
	}
	o.arr = make([]JsObject, size)
	return o
}

type jsObj struct {
	n         int32
	prop      map[string]JsObject
	typeFlag  jsKind
	binFlag   binKind
	b         bool
	f32       float32
	f64       float64
	i8        int8
	i16       int16
	i32       int32
	i64       int64
	uint_     uint
	u8        uint8
	u16       uint16
	u32       uint32
	u64       uint64
	uptr      uintptr
	int_      int
	null      bool
	zero      bool
	global    bool
	s         string
	debugName string
	arr       []JsObject
	fn        func([]JsObject) JsObject
	_this     interface{} // only INSTANCES have this
	class     class       //for classes
}

func (j *jsObj) Id() int32 {
	return j.n
}

func (j *jsObj) typeInfo() string {
	part1, part2 := "***UNKNOWN***", "***UNKNOWN***"
	switch j.binFlag {
	case typeFlagNone:
		part1 = "None"
	case typeFlagObject:
		part1 = "Obj"
	case typeFlagString:
		part1 = "Str"
	case typeFlagSymbol:
		part1 = "Sym"
	case typeFlagFunction:
		part1 = "Func"
	}
	switch j.typeFlag {
	case TypeUndefined:
		part2 = "undef"
	case TypeNull:
		part2 = "null"
	case TypeBoolean:
		part2 = "bool"
	case TypeNumber:
		part2 = "num"
	case TypeString:
		part2 = "str"
	case TypeSymbol:
		part2 = "sym"
	case TypeFunction:
		part2 = "func"
	case TypeObject:
		part2 = "obj"
	}
	return fmt.Sprintf("%s/%s", part1, part2)
}
func (j *jsObj) String() string {
	return fmt.Sprintf("[%s:%s:%d:0x%x,%d props]", j.debugName, j.typeInfo(), j.n, j.BinaryRep(), len(j.prop))
}

func (j *jsObj) IsFunc() bool {
	return j.fn != nil
}

func (j *jsObj) Index(int) JsObject {
	panic("Index")
}
func (j *jsObj) SetIndex(index int, obj JsObject) {
	print(fmt.Sprintf("xxx set index on %s: index is %d, param is %s\n", j, index, obj))
	j.arr[index] = obj
	//panic("SetIndex")

}
func (j *jsObj) Length() int {
	panic("Length")
}
func (j *jsObj) IsClassObject() bool {
	return j.class != nil
}

func (j *jsObj) InstanceOf(class JsObject) bool {
	if !j.IsInstance() {
		panic(fmt.Sprintf("attempt to call IsClass() on %s but not an instance", j))
	}
	if !class.IsClassObject() {
		panic(fmt.Sprintf("attempt to call IsClass() but parameter %s is not a class", class))
	}
	// print(fmt.Sprintf("xxx instance of comparison: is %s instance of %s [%v,%v]",
	// 	j, class, j.class == nil, class.(*jsObj).class == nil))
	return class.(*jsObj).class == j.class //pointer equality, because class objects are singletons
}

func (j *jsObj) This() interface{} {
	return j._this
}

// Call is used to make METHOD calls, not function calls. We convert this to a method
// call by stuffing the j here in as the first param.
func (j *jsObj) Call(name string, param []JsObject) JsObject {
	fn, ok := j.prop[name]
	if !ok {
		panic(fmt.Sprintf("unable to find prop %s on object %s", name, j))
	}
	if !fn.IsFunc() {
		panic(fmt.Sprintf("%s is not a function on object %d", name, j.Id()))
	}
	result := fn.Invoke(append([]JsObject{j}, param...))
	return result
}

func (j *jsObj) Get(propName string) JsObject {
	v, ok := j.prop[propName]
	if ok {
		return v
	}
	return Undefined
}

// Invoke is used to make FUNCTION calls.  If the call is a method, then Call()
// will have transformed the args appropriately.
func (j *jsObj) Invoke(arg []JsObject) JsObject {
	if !j.IsFunc() {
		panic(fmt.Sprintf("unable to make function invocation on object %d (not a func or methodh)", j.Id()))
	}
	return j.fn(arg)
}

func (j *jsObj) Set(propName string, obj JsObject) JsObject {
	j.prop[propName] = obj
	return j
}
func (j *jsObj) Delete(propName string) {
	delete(j.prop, propName)
}

func (j *jsObj) IsNumber() bool {
	switch j.Id() {
	case predefinedNan, predefinedZero, -1001:
		return true
	default:
		return false
	}
}

func (j *jsObj) BinaryRep() uint64 {
	// normal case is the silly encoding business with NaN...
	//note that this depends on NanHead not having the low bits set!
	if j.IsNumber() {
		var f float64
		switch {
		case j.f32 != 0:
			f = float64(j.f32)
		case j.f64 != 0:
			f = j.f64
		case j.i8 != 0:
			f = float64(j.i8)
		case j.i16 != 0:
			f = float64(j.i16)
		case j.i32 != 0:
			f = float64(j.i32)
		case j.i64 != 0:
			f = float64(j.i64)
		case j.int_ != 0:
			f = float64(j.int_)
		case j.uint_ != 0:
			f = float64(j.uint_)
		default:
			if j.Id() != predefinedNan && j.Id() != predefinedZero && j != Undefined {
				panic("can't understand number type of id " + fmt.Sprint(j.Id()))
			}
		}
		if j == Undefined { // that is how the js interface defines it
			return 0
		}
		if f == 0 || j.Id() == predefinedZero {
			return (uint64(jsgo.NanHead) << 32) | predefinedZero
		}
		if math.IsNaN(f) || j.Id() == predefinedNan {
			return (uint64(jsgo.NanHead) << 32) | predefinedNan
		}
		return math.Float64bits(f)
	}
	high := uint64(jsgo.NanHead|j.binFlag) << 32
	low := uint64(j.Id())
	result := high | low

	return result
}
func (js *jsObj) AsNumber() float64 {
	return math.Float64frombits(js.BinaryRep())
}
func (j *jsObj) NewInstance(arg []JsObject) JsObject {
	if j.class == nil {
		panic(fmt.Sprintf("attempt to call NewInstance() on %s but not a class", j))
	}
	return j.class.NewInstance(arg)
}

func (j *jsObj) IsInstance() bool {
	return j.This() != nil
}

func makeFuncWrapper(_ []JsObject) JsObject {
	return badWrappedFunc
}

func fsWrite(arg []JsObject) JsObject {
	if len(arg) < 6 {
		print("JSCONSOLE: unable to understand arguments to fsWrite (console) because number of parameters is too small\n")
		return JsObjectMap.Get(predefinedZero)
	}
	output := 0
	for i := 2; arg[i].Id() != predefinedNull; i += 3 {
		var s string
		if arg[i].IsInstance() && arg[i].InstanceOf(Uint8ArrayObj) {
			s = string(arg[i].This().(*Uint8ArrayInstance).data)
			logJSConsoleMessage(s)
		} else {
			print(fmt.Sprintf("unable to understand parameter %d: %s", i, arg[i]))
		}
		if arg[i+1].Id() != predefinedZero {
			print(fmt.Sprintf("unable to understand parameter %d: expected 0 but got %s", i+1, arg[i+1]))
		}
		if !arg[i+2].IsNumber() {
			print(fmt.Sprintf("unable to understand parameter %d: expected a number for length but got %s", i+2, arg[i+2]))
		} else {
			if int(arg[i+2].AsNumber()) != len(s) {
				print(fmt.Sprintf("unable to understand parameter %d: mismatched length, expected %d but got %s", i+2, len(s), arg[i+2]))
			} else {
				output += len(s)
			}
		}
	}
	if output == 0 {
		return JsObjectMap.Get(predefinedZero)
	}
	return newJSObjNum(float64(output))
}

func logJSConsoleMessage(s string) {
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	req := &logmsg.LogRequest{
		Stamp:   timestamppb.Now(), //xxx should use kernel now
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: s,
	}
	backdoor.Log(req, false, false, true, nil)
}

//
// OBJMAP
//

type objMap struct {
	lock sync.Mutex
	next int32
	data map[int32]JsObject
}

func newObjMap() *objMap {
	return &objMap{
		next: 0,
		data: make(map[int32]JsObject),
	}
}

// returns the object it just Put in the map
func (o *objMap) Put(obj JsObject) JsObject {
	if obj.Id() >= 0 {
		panic("attempt to put object in map second time " + fmt.Sprint(obj.Id()))
	}
	t := obj.(*jsObj).typeFlag
	if t == TypeNumber || t == TypeUndefined {
		panic("should not be putting numbers or undefined into the object map")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	obj.(*jsObj).n = o.next
	o.next++
	o.data[obj.Id()] = obj
	return obj
}

func (o *objMap) predefined(obj JsObject, expected int32) JsObject {
	obj.(*jsObj).n = expected
	o.lock.Lock()
	defer o.lock.Unlock()
	o.data[expected] = obj
	if expected >= o.next {
		o.next = expected + 1
	}
	return obj
}

func (o *objMap) Get(id int32) JsObject {
	if id < 0 {
		panic("attempt to get object from map that was never entered " + fmt.Sprint(id))
	}
	o.lock.Lock()
	defer o.lock.Unlock()

	v, ok := o.data[id]
	if !ok {
		panic("unknown object id " + fmt.Sprint(id))
	}
	return v
}
