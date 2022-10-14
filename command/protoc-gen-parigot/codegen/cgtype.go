package codegen

import (
	"fmt"
)

var danger = &CGType{}

type CGType struct {
	composite *WasmMessage
	basic     string
	lang      LanguageText
	finder    Finder
	protoPkg  string
	hasValue  bool
}

func NewCGTypeFromComposite(m *WasmMessage, l LanguageText, f Finder,
	protoPackage string) *CGType {
	return &CGType{composite: m, lang: l, protoPkg: protoPackage, finder: f, hasValue: true}
}

func (c *CGType) HasValueBeenSet() bool {
	return c.hasValue
}

func (c *CGType) SetEmpty() {
	if c.HasValueBeenSet() {
		panic("attempt to set the emptiness property on CGType that already has a value")
	}
	c.hasValue = true
}

func (c *CGType) GetCompositeType() *WasmMessage {
	return c.composite
}

func (c *CGType) IsEmpty() bool {
	if c.hasValue == true {
		comp := c.composite
		basic := c.basic
		if comp == nil && basic == "" {
			return true
		}
		return false
	}
	panic("attempt to query the empty status of a CGType before any value set")
}

func NewCGTypeFromBasic(tname string, l LanguageText, f Finder,
	protoPkg string) *CGType {
	return &CGType{composite: nil, lang: l, basic: tname, protoPkg: protoPkg, hasValue: true}
}
func (c *CGType) IsBasic() bool {
	return c.basic != ""
}

func (c *CGType) String(from string) string {
	if c.composite == nil {
		return c.lang.ProtoTypeNameToLanguageTypeName(c.basic)
	}
	addr := c.finder.AddressingNameFromMessage(from, c.composite)
	return addr
}

func GetCGTypeForInputParam(i *InputParam) *CGType {
	// input param has either 0 or 1 entry
	inputName := i.GetParent().GetInputType()
	if inputName == "" {
		return &CGType{hasValue: true} //this is an empty CGType
	}
	finder := i.GetParent().GetFinder()
	protoPkg := i.GetParent().GetProtoPackage()
	lang := i.GetLanguage()
	for _, s := range parigotTypeList {
		if inputName == s {
			cgT := NewCGTypeFromBasic(s, lang, finder, protoPkg)
			return cgT
		}
	}
	// not a basic type, so we'll try composites
	msg := finder.FindMessageByName(protoPkg, i.GetTypeName(), nil)
	if msg == nil {
		panic(fmt.Sprintf("attempted to find message in pkg %s with name %s, but failed",
			protoPkg, i.GetTypeName()))
	}
	cgT := NewCGTypeFromComposite(msg, lang, finder, protoPkg)
	return cgT
}

func GetCGTypeForOutputParam(o *OutputParam) *CGType {
	// input param has either 0 or 1 entry
	outputName := o.GetParent().GetOutputType()
	if outputName == "" {
		return &CGType{hasValue: true} //this is an empty CGType
	}
	finder := o.GetParent().GetFinder()
	protoPkg := o.GetParent().GetProtoPackage()
	lang := o.GetLanguage()
	for _, s := range parigotTypeList {
		if outputName == s {
			cgT := NewCGTypeFromBasic(s, lang, finder, protoPkg)
			return cgT
		}
	}
	// not a basic type, so we'll try composites
	msg := finder.FindMessageByName(protoPkg, o.GetTypeName(), nil)
	if msg == nil {
		panic(fmt.Sprintf("attempted to find message in pkg %s with name %s, but failed",
			protoPkg, o.GetTypeName()))
	}
	cgT := NewCGTypeFromComposite(msg, lang, finder, protoPkg)
	return cgT
}

type CGParameter struct {
	name     string
	field    *WasmField
	cgType   *CGType
	noFormal bool
}

func NewCGParameterFromString(cgType *CGType) *CGParameter {
	return &CGParameter{noFormal: true, cgType: cgType}
}

func NewCGParameterFromField(f *WasmField, cgType *CGType) *CGParameter {
	return &CGParameter{field: f, cgType: cgType}
}

func (c *CGParameter) GetCGType() *CGType {
	return c.cgType
}

func (c *CGParameter) GetFormalName() string {
	if c.noFormal {
		panic("should not be asking for a formal from something that doesn't have one")
	}
	if c.field == nil {
		return c.name
	}
	return c.field.GetName()
}

func (c *CGParameter) String(from string) string {
	return fmt.Sprintf("%s %s", c.GetFormalName(), c.GetCGType().String(from))
}
