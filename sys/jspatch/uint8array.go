package jspatch

var uint8ArrayObj jsObject

type uint8ArrayClass struct {
}

type uint8ArrayInstance struct {
	class *uint8ArrayClass
	data  []byte
}

func (uClass *uint8ArrayClass) Name() string {
	return "uint8Array"
}

func (uClass *uint8ArrayClass) NewInstance(arg []jsObject) jsObject {
	if len(arg) != 1 {
		panic("wrong number of args to uint8Array constructor")
	}
	if !arg[0].isNumber() {
		panic("argument to constructor of uint8arg is not a number")
	}
	size := int(arg[0].AsNumber())
	obj := newJSObjGeneric(typeFlagObject, TypeObject, "inst(uint8Array)")
	obj.n = -10101011
	obj.prop = make(map[string]jsObject)
	obj.class = uint8ArrayObj.(*jsObj).class
	inst := &uint8ArrayInstance{class: uClass, data: make([]byte, size)}

	obj._this = inst // only instances have _this
	jsObjectMap.put(obj)

	return obj
}
