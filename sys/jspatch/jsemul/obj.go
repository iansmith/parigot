package jsemul

var objectObj JsObject

type objectClass struct {
}

type objectInstance struct {
	class *objectClass
}

func (oClass *objectClass) Name() string {
	return "object"
}

func (oClass *objectClass) NewInstance(arg []JsObject) JsObject {
	if len(arg) != 0 {
		panic("wrong number of args to uint8Array constructor")
	}
	obj := newJSObjGeneric(typeFlagObject, TypeObject, "inst(object)")
	obj.n = -10101012
	obj.prop = make(map[string]JsObject)
	obj.class = objectObj.(*jsObj).class
	inst := &objectInstance{class: oClass}

	obj._this = inst // only instances have _this
	JsObjectMap.Put(obj)

	return obj
}
