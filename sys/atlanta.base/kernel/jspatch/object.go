package jspatch

import "log"

type jsObject interface {
	id() int32
	getProp(prop string) jsObject
	setProp(prop string, obj jsObject)
	deleteProp(prop string)
	getIndex(int64) jsObject
	setIndex(int64, jsObject)
}

var refCount = make(map[int32]int)
var object = make(map[int32]jsObject)

var null = newJSObj(2)
var trueBool = newJSObj(3)
var falseBool = newJSObj(4)
var win = newJSObj(5)
var zero = newJSObj(1)
var notANumber = newJSObj(0)
var nextId = int32(6)

func init() {
	object[5] = win
	object[4] = falseBool
	object[3] = trueBool
	object[2] = null
	object[1] = zero
	object[0] = notANumber

	refCount[0] = 1
	refCount[1] = 1
	refCount[2] = 1
	refCount[3] = 1
	refCount[4] = 1
	refCount[5] = 1

	objFuncId := nextId
	object[objFuncId] = newJSObjCallable(objFuncId, func() {
		log.Printf("Object function called()")
	})
	refCount[objFuncId] = 1
	nextId++
	arrayFuncId := nextId
	object[arrayFuncId] = newJSObjCallable(arrayFuncId, func() {
		log.Printf("Object function called()")
	})
	refCount[arrayFuncId] = 1
	nextId++
	win.setProp("Object", object[objFuncId])
	win.setProp("Array", object[arrayFuncId])
}

type jsObj struct {
	n         int32
	prop      map[string]jsObject
	callable  bool
	fn        func()
	array     bool
	arrayData []jsObject
}

func newJSObj(n int32) *jsObj {
	return &jsObj{
		n:        n,
		prop:     make(map[string]jsObject),
		callable: false,
	}
}

func newJSObjCallable(n int32, fn func()) *jsObj {
	return &jsObj{
		n:        n,
		prop:     make(map[string]jsObject),
		callable: true,
		fn:       fn,
	}
}

func (j *jsObj) id() int32 {
	return j.n
}

func (j *jsObj) getProp(p string) jsObject {
	v, ok := j.prop[p]
	if !ok {
		log.Printf("unable to find prop '%s', returning undefined", p)
		return null
	}
	return v
}

func (j *jsObj) setProp(p string, v jsObject) {
	j.prop[p] = v.(*jsObj)
}

func (j *jsObj) setIndex(i int64, obj jsObject) {
	if !j.array {
		log.Printf("attempt to set index non array (%d), ignoring", j.id())
	}
	if int(i) >= len(j.arrayData) {
		log.Printf("attempt to set index past end of array (%d) at index (%d) , ignoring",
			j.id(), i)
	}
	j.arrayData[int(i)] = obj
}

func (j *jsObj) getIndex(i int64) jsObject {
	if !j.array {
		log.Printf("attempt to index non array (%d), returning null", j.id())
		return null
	}
	if int(i) >= len(j.arrayData) {
		log.Printf("attempt to index past end of array (%d) at index (%d), returning null",
			j.id(), i)
		return null
	}
	return j.arrayData[int(i)]
}

func (j *jsObj) deleteProp(p string) {
	// old := j.prop[p]
	delete(j.prop, p)
	// xxx?
	// refCount[old.id()]--
}
