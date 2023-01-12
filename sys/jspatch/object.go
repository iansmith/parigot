package jspatch

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"unsafe"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"

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

//
// InternalLogger and the logger object are a workaround for import cycles because of
// the fact that we need some "extra" places to have the ability to send to the same
// log endpoint that the "normal" logger does.
//

type InternalLogger interface {
	ProccessLogRequest(*logmsg.LogRequest, bool, bool, bool, []byte)
}

var logger InternalLogger

func SetInternalLogger(il InternalLogger) {
	logger = il
}

//
// End of logger crap.  xxxfixme this is probably not a good solution.
//

var undefined = &jsObj{typeFlag: TypeUndefined}

func floatValue(f float64) jsObject {
	if f == 0 {
		return jsObjectMap.get(predefinedZero)
	}
	if f != f {
		return jsObjectMap.get(predefinedNan)
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
	Call(string /*name*/, []jsObject) jsObject
	Get(string) jsObject
	Invoke([]jsObject) jsObject
	Set(string, jsObject) jsObject /* returns self for chaining*/
	Delete(string)
	IsFunc() bool
	IsClassObject() bool
	IsInstance() bool
	AsNumber() float64
	// from a class, create an instance
	newInstance([]jsObject) jsObject
	// for our internal use
	binaryRep() uint64
	isNumber() bool
	this() interface{}
}

var jsObjectMap = newObjMap()

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
	NewInstance([]jsObject) jsObject
	Name() string
}

var badWrappedFunc jsObject

func init() {
	// this is caused because this code runs on linux but works on behalf of WASM
	// and so it ends up being loaded (and init called) twice, once for the runner
	// in linux and once for the module loaded into the WASM side
	print("xxx init called on object.go ... ", runtime.GOOS, "\n")
	// order here is carefully constructed
	jsObjectMap.predefined(newJSObjNum(math.NaN()), predefinedNan)

	obj := newJSObjGeneric(typeFlagNone, TypeObject, "ZERO")
	obj.zero = true
	obj.typeFlag = jsKind(typeFlagNone)
	jsObjectMap.predefined(obj, predefinedZero)

	obj = newJSObjGeneric(typeFlagNone, TypeNull, "NULL")
	obj.null = true
	jsObjectMap.predefined(obj, predefinedNull)

	obj = newJSObjBool(true)
	jsObjectMap.predefined(obj, predefinedTrue)

	obj = newJSObjBool(false)
	jsObjectMap.predefined(obj, predefinedFalse)

	glob := newJSObjObject(true, "WINDOW")
	jsObjectMap.predefined(glob, predefinedWindow)

	goObj := newJSObjObject(false, "GO")
	jsObjectMap.predefined(goObj, predefinedGo)

	obj = newJSObjFunc(makeFuncWrapper, "_makeFuncWrapper")
	jsObjectMap.put(obj)
	goObj.Set("_makeFuncWrapper", obj)

	// obj = newJSObjFunc(makeArray, "makeArray")
	// glob.Set("Array", obj)
	// jsObjectMap.put(obj)

	dateObj = newJSObjClass(&dateClass{offsetFromUTCInMins: 0}, "class(Date)")
	glob.Set("Date", dateObj)
	jsObjectMap.put(dateObj)

	fs := newJSObjGeneric(typeFlagObject, TypeObject, "fs")
	glob.Set("fs", fs)
	jsObjectMap.put(fs)
	writeFn := newJSObjFunc(fsWrite, "fs.Write")
	jsObjectMap.put(writeFn)
	fs.Set("write", writeFn)

	obj = newJSObjGeneric(typeFlagObject, TypeObject, "process")
	glob.Set("process", obj)
	jsObjectMap.put(obj)

	constants := jsObjectMap.put(newJSObjGeneric(typeFlagObject, TypeObject, "constants"))
	minus1 := newJSObjNum(-1)

	//{ O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused	constants.Set()
	constants.Set("O_WRONLY", minus1).Set("O_RDWR", minus1).Set("O_TRUNC", minus1).
		Set("O_APPEND", minus1).Set("O_EXCL", minus1).Set("O_CREAT", minus1)
	fs.Set("constants", constants)

	glob.Set("fs", fs)

	uint8ArrayObj = newJSObjClass(&uint8ArrayClass{}, "class(Uint8Array)")
	jsObjectMap.put(uint8ArrayObj)
	glob.Set("Uint8Array", uint8ArrayObj)

	objectObj = newJSObjClass(&objectClass{}, "class(Object)")
	jsObjectMap.put(objectObj)
	glob.Set("Object", objectObj)

	badWrappedFunc = newJSObjFunc(func(_ []jsObject) jsObject {
		print("calling wrapped go functions from javascript makes no sense when there is no javascript")
		return jsObjectMap.get(predefinedNull)
	}, "bad function call")
	jsObjectMap.put(badWrappedFunc)
}

func goToJS(x any) jsObject {
	switch x := x.(type) {
	case jsObject:
		return x
	case func(interface{}):
		print("got to func case in goToJS\n")
		return nil
	//	return x.Value
	case nil:
		return jsObjectMap.get(predefinedNull)
	case bool:
		if x {
			return jsObjectMap.get(predefinedTrue)
		} else {
			return jsObjectMap.get(predefinedFalse)
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
		return jsObjectMap.put(newJSObjString(x))
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
			val.Set(k, goToJS(v))
		}
		return jsObjectMap.put(val)
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
		prop:      make(map[string]jsObject),
		typeFlag:  TypeNumber,
		f64:       f,
		debugName: fmt.Sprint(f),
	}
}

// class
func newJSObjClass(class class, debugName string) *jsObj {
	return &jsObj{
		n:         -7007,
		prop:      make(map[string]jsObject),
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
		prop:      make(map[string]jsObject),
		typeFlag:  TypeObject,
		global:    isGlobal,
		binFlag:   typeFlagNone,
		debugName: debugName,
	}
}

// func
func newJSObjFunc(fn func([]jsObject) jsObject, debugName string) *jsObj {
	return &jsObj{
		n:         -4004,
		prop:      make(map[string]jsObject),
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
		prop:      make(map[string]jsObject),
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
		prop:      make(map[string]jsObject),
		typeFlag:  t,
		binFlag:   b,
		debugName: debugName,
	}
}
func newJSObjString(s string) *jsObj {
	return &jsObj{
		n:        -6006,
		prop:     make(map[string]jsObject),
		typeFlag: TypeString,
		s:        s,
		binFlag:  typeFlagString,
	}
}
func newJSObjArr(size int) *jsObj {
	o := &jsObj{
		n:        -8008,
		prop:     make(map[string]jsObject),
		typeFlag: TypeObject,
	}
	o.arr = make([]jsObject, size)
	return o
}

type jsObj struct {
	n         int32
	prop      map[string]jsObject
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
	arr       []jsObject
	fn        func([]jsObject) jsObject
	_this     interface{} // only INSTANCES have this
	class     class       //for classes
}

func (j *jsObj) id() int32 {
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
	return fmt.Sprintf("[%s:%s:%d:0x%x,%d props]", j.debugName, j.typeInfo(), j.n, j.binaryRep(), len(j.prop))
}

func (j *jsObj) IsFunc() bool {
	return j.fn != nil
}

func (j *jsObj) Index(int) jsObject {
	panic("Index")
}
func (j *jsObj) SetIndex(index int, obj jsObject) {
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

func (j *jsObj) InstanceOf(class jsObject) bool {
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

func (j *jsObj) this() interface{} {
	return j._this
}

// Call is used to make METHOD calls, not function calls. We convert this to a method
// call by stuffing the j here in as the first param.
func (j *jsObj) Call(name string, param []jsObject) jsObject {
	fn, ok := j.prop[name]
	if !ok {
		panic(fmt.Sprintf("unable to find prop %s on object %s", name, j))
	}
	if !fn.IsFunc() {
		panic(fmt.Sprintf("%s is not a function on object %d", name, j.id()))
	}
	result := fn.Invoke(append([]jsObject{j}, param...))
	return result
}

func (j *jsObj) Get(propName string) jsObject {
	v, ok := j.prop[propName]
	if ok {
		return v
	}
	return undefined
}

// Invoke is used to make FUNCTION calls.  If the call is a method, then Call()
// will have transformed the args appropriately.
func (j *jsObj) Invoke(arg []jsObject) jsObject {
	if !j.IsFunc() {
		panic(fmt.Sprintf("unable to make function invocation on object %d (not a func or methodh)", j.id()))
	}
	return j.fn(arg)
}

func (j *jsObj) Set(propName string, obj jsObject) jsObject {
	j.prop[propName] = obj
	return j
}
func (j *jsObj) Delete(propName string) {
	delete(j.prop, propName)
}

func (j *jsObj) isNumber() bool {
	switch j.id() {
	case predefinedNan, predefinedZero, -1001:
		return true
	default:
		return false
	}
}

func (j *jsObj) binaryRep() uint64 {
	// normal case is the silly encoding business with NaN...
	//note that this depends on NanHead not having the low bits set!
	if j.isNumber() {
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
			if j.id() != predefinedNan && j.id() != predefinedZero {
				panic("can't understand number type of id " + fmt.Sprint(j.id()))
			}
		}
		if f == 0 || j.id() == predefinedZero {
			return (uint64(nanHead) << 32) | predefinedZero
		}
		if math.IsNaN(f) || j.id() == predefinedNan {
			return (uint64(nanHead) << 32) | predefinedNan
		}
		return math.Float64bits(f)
	}
	high := uint64(nanHead|j.binFlag) << 32
	low := uint64(j.id())
	result := high | low

	return result
}
func (js *jsObj) AsNumber() float64 {
	return math.Float64frombits(js.binaryRep())
}
func (j *jsObj) newInstance(arg []jsObject) jsObject {
	if j.class == nil {
		panic(fmt.Sprintf("attempt to call newInstance() on %s but not a class", j))
	}
	return j.class.NewInstance(arg)
}

func (j *jsObj) IsInstance() bool {
	return j.this() != nil
}

func makeFuncWrapper(_ []jsObject) jsObject {
	return badWrappedFunc
}

func fsWrite(arg []jsObject) jsObject {
	if len(arg) < 6 {
		print("JSCONSOLE: unable to understand arguments to fsWrite (console) because number of parameters is too small\n")
		return jsObjectMap.get(predefinedZero)
	}
	output := 0
	needsHeader := true
	for i := 2; arg[i].id() != predefinedNull; i += 3 {
		if needsHeader {
			print("JSCONSOLE:")
			needsHeader = false
		}
		var s string
		if arg[i].IsInstance() && arg[i].InstanceOf(uint8ArrayObj) {
			s = string(arg[i].this().(*uint8ArrayInstance).data)
			logJSConsoleMessage(s)
			if strings.HasSuffix(s, "\n") {
				needsHeader = true
			}
		} else {
			print(fmt.Sprintf("unable to understand parameter %d: %s", i, arg[i]))
		}
		if arg[i+1].id() != predefinedZero {
			print(fmt.Sprintf("unable to understand parameter %d: expected 0 but got %s", i+1, arg[i+1]))
		}
		if !arg[i+2].isNumber() {
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
		return jsObjectMap.get(predefinedZero)
	}
	return newJSObjNum(float64(output))
}

func logJSConsoleMessage(s string) {
	req := &logmsg.LogRequest{
		Stamp:   timestamppb.Now(), //xxx should use kernel now
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: s,
	}
	logger.ProccessLogRequest(req, false, false, true, nil)
}

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
	if expected >= o.next {
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
