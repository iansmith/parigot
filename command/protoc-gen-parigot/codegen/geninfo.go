package codegen

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"log"
)

type GenInfo struct {
	request *pluginpb.CodeGeneratorRequest
	file    *descriptorpb.FileDescriptorProto
	lang    LanguageText
	finder  *SimpleFinder
}

func NewGenInfo() *GenInfo {
	return &GenInfo{
		finder: NewSimpleFinder(),
	}
}

// Service returns all the wasm services of this GenInfo.
// The caller should not change these values, these are pointers in the real
// data structures.
func (g *GenInfo) Service() []*WasmService {
	return g.finder.Service()
}

// Message returns all the wasm messages of this GenInfo.
// The caller should not change these values, these are pointers in the real
// data structures.
func (g *GenInfo) Message() []*WasmMessage {
	return g.finder.Message()
}

// GetFile returns the file data structure that is associated with the currently
// operating code generation.
func (g *GenInfo) GetFile() *descriptorpb.FileDescriptorProto {
	return g.file
}

// GetRequest returns the original request data structure that was delivered via
// stdin.  This request is associated with the currently operating code generation.
func (g *GenInfo) GetRequest() *pluginpb.CodeGeneratorRequest {
	return g.request
}

func (g *GenInfo) IsAbi() bool {
	if g.file.GetOptions() == nil {
		log.Printf("no options on %s\n", g.file.GetName())
		return false
	}
	return IsAbi(g.file.GetOptions().String())
}

type InputParam struct {
	name     string
	typ      *WasmMessage
	paramVar []*ParamVar
	lang     LanguageText
	parent   *WasmMethod
}

func (i *InputParam) GetName() string {
	return i.name
}
func (i *InputParam) GetTyp() *WasmMessage {
	return i.typ
}
func (i *InputParam) GetParamVar() []*ParamVar {
	return i.paramVar
}

func newInputParam(protoPkg string, messageName string, parent *WasmMethod, finder Finder) *InputParam {
	result := &InputParam{
		parent: parent,
	}
	msg := finder.FindMessageByName(protoPkg, messageName, nil)
	if msg == nil {
		log.Fatalf("unable to find input parameter type %s", lastSegmentOfPackage(messageName))
	}
	result.typ = msg
	result.paramVar = make([]*ParamVar, len(msg.GetField()))
	for i, f := range msg.GetField() {
		pVar := newParamVar(f)
		result.paramVar[i] = pVar
	}
	return result
}

func (i *InputParam) IsEmpty() bool {
	return len(i.paramVar) == 0
}

func (i *InputParam) Len() int {
	return len(i.paramVar)
}

type ParamVar struct {
	name  string
	field *WasmField
	lang  LanguageText
}

func (p *ParamVar) GetName() string {
	return p.name
}
func (p *ParamVar) GetField() *WasmField {
	return p.field
}

func newParamVar(f *WasmField) *ParamVar {
	return &ParamVar{
		name:  f.GetWasmFieldName(),
		field: f,
	}
}

func (p *ParamVar) TypeFromProto() string {
	return p.field.GetType().String()
}

func (p *ParamVar) Formal() string {
	return p.name
}

func (p *ParamVar) IsStrictWasmType() bool {
	switch p.field.GetType().String() {
	case "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE":
		return true
	}
	return false
}

func (p *ParamVar) IsWasmType() bool {
	switch p.field.GetType().String() {
	case "TYPE_STRING", "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE", "TYPE_BOOL":
		return true
	}
	return false
}

type OutputParam struct {
	name     string
	typ      *WasmMessage
	paramVar []*ParamVar
	lang     LanguageText
	parent   *WasmMethod
}

func (o *OutputParam) GetName() string {
	return o.name
}
func (o *OutputParam) GetTyp() *WasmMessage {
	return o.typ
}
func (o *OutputParam) GetParamVar() []*ParamVar {
	return o.paramVar
}

func (o *OutputParam) IsEmpty() bool {
	return len(o.paramVar) == 0
}

func (o *OutputParam) IsMultipleReturn() bool {
	return len(o.paramVar) > 1
}

func newOutputParam(protoName string, name string, parent *WasmMethod, finder Finder) *OutputParam {
	result := &OutputParam{
		parent: parent,
	}
	m := finder.FindMessageByName(protoName, name, nil)
	if m == nil {
		log.Fatalf("unable to find type %s for output parameter", name)
	}
	result.typ = m
	result.paramVar = make([]*ParamVar, len(m.GetField()))
	for i, f := range m.GetField() {
		result.paramVar[i] = newParamVar(f)
	}
	return result
}

func removeQuotes(s string) string {
	result := s
	l := len(s)
	if l > 2 && s[0:1] == "\"" && s[l-1:l] == "\"" {
		result = s[1 : l-1]
	}
	return result
}

//func (g *GenInfo) findMessageByName(n string) *WasmMessage {
//	for _, m := range g.wasmMessage {
//		log.Printf("comparing %s to %s", m.GetFullName(), n)
//		if m.GetFullName() == n {
//			return m
//		}
//	}
//	// why do they do this SOME of the time?
//	if len(n) > 0 && n[0] == '.' {
//		for _, m := range g.wasmMessage {
//			log.Printf("comparing %s to %s", m.GetFullName(), n[1:])
//			if m.GetFullName() == n[1:] {
//				return m
//			}
//		}
//	}
//	return nil
//}

type ServiceRecord struct {
	wasmName     string
	protoPackage string
	goPackage    string
}

type MessageRecord struct {
	wasmName     string
	protoPackage string
	goPackage    string
}

func NewMessageRecord(wasmName, protoPackage, goPackage string) *MessageRecord {
	return &MessageRecord{
		wasmName:     wasmName,
		protoPackage: protoPackage,
		goPackage:    goPackage,
	}
}
func NewServiceRecord(wasmName, protoPackage, goPackage string) *ServiceRecord {
	return &ServiceRecord{
		wasmName:     wasmName,
		protoPackage: protoPackage,
		goPackage:    goPackage,
	}
}

func (m *MessageRecord) String() string {
	return fmt.Sprintf("MessageRec(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (m *ServiceRecord) String() string {
	return fmt.Sprintf("ServiceRec(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (g *GenInfo) RegisterService(w *WasmService) {
	g.finder.AddServiceType(w.GetWasmServiceName(), w.GetProtoPackage(), w.GetGoPackage(), w)
}

func (g *GenInfo) RegisterMessage(w *WasmMessage) {
	g.finder.AddMessageType(w.GetWasmMessageName(), w.GetProtoPackage(), w.GetGoPackage(), w)
}

func (g *GenInfo) FindServiceByName(protoPackage string, name string, next Finder) *WasmService {
	return g.finder.FindServiceByName(protoPackage, name, next)
}
func (g *GenInfo) FindMessageByName(protoPackage string, name string, next Finder) *WasmMessage {
	return g.finder.FindMessageByName(protoPackage, name, next)
}

func (g *GenInfo) AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage) {
	g.finder.AddMessageType(wasmName, protoPackage, goPackage, message)

}
func (g *GenInfo) AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService) {
	g.finder.AddServiceType(wasmName, protoPackage, goPackage, service)
}

func (g *GenInfo) SetReqAndFile(request *pluginpb.CodeGeneratorRequest, proto *descriptorpb.FileDescriptorProto) {
	g.request = request
	g.file = proto
}
