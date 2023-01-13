package jsemul

var Uint8ArrayObj JsObject

type uint8ArrayClass struct {
}

type Uint8ArrayInstance struct {
	class *uint8ArrayClass
	data  []byte
}

func (uClass *uint8ArrayClass) Name() string {
	return "uint8Array"
}

func (uClass *uint8ArrayClass) NewInstance(arg []JsObject) JsObject {
	if len(arg) != 1 {
		panic("wrong number of args to uint8Array constructor")
	}
	if !arg[0].IsNumber() {
		panic("argument to constructor of uint8arg is not a number")
	}
	size := int(arg[0].AsNumber())
	obj := newJSObjGeneric(typeFlagObject, TypeObject, "inst(uint8Array)")
	obj.n = -10101011
	obj.prop = make(map[string]JsObject)
	obj.class = Uint8ArrayObj.(*jsObj).class
	inst := &Uint8ArrayInstance{class: uClass, data: make([]byte, size)}
	// maybe copyBytesToJS should be a method here?
	obj._this = inst // only instances have _this
	JsObjectMap.Put(obj)

	return obj
}

// CopyBytesToJS here isn't really great because we are not implementing it as a proper
// function of the instance object. However, the caller already has all the pieces to
// call this function properly and doing an extra wrap/unwrap sequence seems silly.
func (uInst *Uint8ArrayInstance) CopyBytesToJS(in []byte) int {
	shorter := len(uInst.data)
	var excess int
	if len(in) < shorter {
		excess = shorter - len(in)
		shorter = len(in)
	}
	for i := 0; i < shorter; i++ {
		uInst.data[i] = in[i]
	}
	// just for safety
	for i := 0; i < excess; i++ {
		uInst.data[shorter+i] = 0
	}
	return shorter
}
