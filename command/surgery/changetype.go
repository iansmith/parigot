package main

import (
	"github.com/iansmith/parigot/command/transform"
	"log"
)

type changeType struct {
	name string
}

func newChangeType(name string) *changeType {
	return &changeType{
		name: name,
	}
}

func (c *changeType) Run(mod *transform.Module) {
	typeCount := 0
	findToplevelInModule(mod, transform.TypeDefT, func(tl transform.TopLevel) {
		typeCount++
	})
	log.Printf("\tsurgery: found %d types in wasm", typeCount)
	findToplevelInModule(mod, transform.ImportDefT, func(tl transform.TopLevel) {
		idef := tl.(*transform.ImportDef)
		if idef.ImportedAs == c.name {
			c.modify(mod, typeCount)
			typeCount++
		}
	})
}

func (c *changeType) modify(mod *transform.Module, typeNumber int) {
	typeDef := &transform.TypeDef{
		Annotation: typeNumber,
		Func: &transform.FuncSpec{
			Param: &transform.ParamDef{
				Type: &transform.TypeNameSeq{
					Name: []string{"i32", "i32", "i32", "i32", "i32", "i32"},
				},
			},
			Result: nil,
		},
	}
	mod.AppendTopLevelDef(typeDef)

	log.Printf("\tsurgery: appended type def %s", typeDef.IndentedString(0))
	changeToplevelInModule(mod, transform.ImportDefT, func(tl transform.TopLevel) transform.TopLevel {
		iDef := tl.(*transform.ImportDef)
		if iDef.ImportedAs == c.name {
			iDef.FuncNameRef.Type.Num = typeNumber
			log.Printf("\tsurgery:changed import %s", iDef.IndentedString(0))
		}
		return iDef
	})
	log.Printf("\tsurgery: operation completed,putting the wasm back together")
}
