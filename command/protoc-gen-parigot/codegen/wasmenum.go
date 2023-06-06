package codegen

import "google.golang.org/protobuf/types/descriptorpb"

// WasmEnumType
type WasmEnumType struct {
	*descriptorpb.EnumDescriptorProto
	wasmEnumTypeName string
	parent           *descriptorpb.FileDescriptorProto
	child            []*WasmEnumValue
	finder           Finder
	lang             LanguageText
}

// each value in the enum has one of these
type WasmEnumValue struct {
	*descriptorpb.EnumValueDescriptorProto
	enumType          *WasmEnumType
	wasmEnumValueName string
	parent            *descriptorpb.FileDescriptorProto
	finder            Finder
	start, end        int
	lang              LanguageText
}

func (w *WasmEnumType) GetFinder() Finder {
	return w.finder
}

func (w *WasmEnumType) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}
func (w *WasmEnumType) GetProtoPackage() string {
	return w.GetParent().GetPackage()
}

func (w *WasmEnumType) GetParent() *descriptorpb.FileDescriptorProto {
	return w.parent
}
func (w *WasmEnumType) AddChild(c *WasmEnumValue) {
	w.child = append(w.child, c)
}

func NewWasmEnumType(file *descriptorpb.FileDescriptorProto, enumDesc *descriptorpb.EnumDescriptorProto,
	lang LanguageText, finder Finder) *WasmEnumType {
	et := &WasmEnumType{
		EnumDescriptorProto: enumDesc,
		wasmEnumTypeName:    "",
		parent:              file,
		child:               []*WasmEnumValue{},
		finder:              finder,
		lang:                lang,
	}
	return et
}

func NewWasmEnumValue(file *descriptorpb.FileDescriptorProto, enumDesc *descriptorpb.EnumValueDescriptorProto,
	lang LanguageText, finder Finder, type_ *WasmEnumType) *WasmEnumValue {
	ev := &WasmEnumValue{
		EnumValueDescriptorProto: enumDesc,
		wasmEnumValueName:        "",
		parent:                   file,
		finder:                   nil,
		lang:                     lang,
		enumType:                 type_,
	}
	type_.AddChild(ev)
	return ev

}
