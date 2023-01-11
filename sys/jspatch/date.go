package jspatch

import (
	"fmt"
	"time"
)

var dateObj jsObject

type dateClass struct {
	offsetFromUTCInMins int // default for all dates
}

type dateInstance struct {
	class  *dateClass
	offset int
	t      time.Time
}

func (dClass *dateClass) Name() string {
	return "Date"
}

func (dClass *dateClass) NewInstance(arg []jsObject) jsObject {
	if len(arg) != 0 {
		panic("wrong number of args to Date constructor")
	}
	obj := newJSObjGeneric(typeFlagObject, TypeObject, "inst(Date)")
	obj.n = -10101010
	obj.prop = make(map[string]jsObject)
	inst := &dateInstance{
		class:  dClass,
		offset: dClass.offsetFromUTCInMins,
		t:      time.Now(),
	}
	obj.class = dateObj.(*jsObj).class
	obj._this = inst // only instances have _this
	jsObjectMap.put(obj)

	getTimezoneOffset := newJSObjFunc(func(arg []jsObject) jsObject {
		if len(arg) != 1 {
			panic(fmt.Sprintf("wrong number of args to getTimezoneOffset (expected 0 got %d)", len(arg)-1))
		}
		return arg[0].this().(*dateInstance).getTimezoneOffset()
	}, "getTimezoneOffset")

	jsObjectMap.put(getTimezoneOffset)
	obj.prop["getTimezoneOffset"] = getTimezoneOffset

	return obj
}

func (inst *dateInstance) getTimezoneOffset() jsObject {
	_, offset := inst.t.Zone()
	mins := offset / 60
	if mins == 0 {
		return jsObjectMap.get(predefinedZero)
	}
	return newJSObjNum(float64(offset / 60))
}
