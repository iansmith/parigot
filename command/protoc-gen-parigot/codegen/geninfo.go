package codegen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"log"
)

type GenInfo struct {
	request     *pluginpb.CodeGeneratorRequest
	file        *descriptorpb.FileDescriptorProto
	wasmService []*WasmService
	wasmMessage []*WasmMessage
	lang        LanguageText
}

// WasmService is like a descriptorpb.ServiceDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmService in the templates.
type WasmService struct {
	*descriptorpb.ServiceDescriptorProto
	wasmServiceName string
	parent          *descriptorpb.FileDescriptorProto
	method          []*WasmMethod
	lang            LanguageText
}

// WasmMessage is like a descriptorpb.DescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmMessage in the templates.
type WasmMessage struct {
	*descriptorpb.DescriptorProto
	wasmMessageName string
	parent          *descriptorpb.FileDescriptorProto
	field           []*WasmField
	lang            LanguageText
}

// WasmMethod is like a descriptorpb.MethodDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmMethod in the templates.
type WasmMethod struct {
	*descriptorpb.MethodDescriptorProto
	wasmMethodName string
	parent         *WasmService
	input          *InputParam
	output         *OutputParam
	lang           LanguageText
	abilang        ABIText
}

func (w *WasmMethod) GetInput() *InputParam {
	return w.input
}
func (w *WasmMethod) GetOutput() *OutputParam {
	return w.output
}

// WasmField is like a descriptorpb.FieldDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmField in the templates.
type WasmField struct {
	*descriptorpb.FieldDescriptorProto
	wasmFieldName string
	parent        *WasmMessage
	lang          LanguageText
}

// GetWasmService returns all the wasm services of this GenInfo.
func (g *GenInfo) GetWasmService() []*WasmService {
	return g.wasmService
}

// GetWasmService returns all the wasm messages of this GenInfo.
func (g *GenInfo) GetWasmMessage() []*WasmMessage {
	return g.wasmMessage
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

// GetWasmServiceName looks through the data structure given that represents the
// original protobuf structure trying to find constructs like this:
//
//	service Foo {
//	   option (parigot.wasm_service_name) = "foozie";
//
// If no such construction is found, it returns the simple name.
func (w *WasmService) GetWasmServiceName() string {
	if w.wasmServiceName != "" {
		return w.wasmServiceName
	}
	// look for it
	w.wasmServiceName = w.ServiceDescriptorProto.GetName() // if they didn't specify, use normal name
	if w.ServiceDescriptorProto.GetOptions() == nil {
		return w.wasmServiceName
	}

	cand, ok := isWasmServiceName(w.ServiceDescriptorProto.GetOptions().String())
	if ok {
		w.wasmServiceName = cand
	}
	w.wasmServiceName = removeQuotes(w.wasmServiceName)
	return w.wasmServiceName
}

// GetWasmMessageName looks through the data structure given that represents the
// original protobuf structure trying to find constructs like this:
//
//	message Foo {
//	   option (parigot.wasm_message_name) = "yeech";
//
// If no such construction is found, it returns the simple name.
func (w *WasmMessage) GetWasmMessageName() string {
	if w.wasmMessageName != "" {
		return w.wasmMessageName
	}
	// look for it
	w.wasmMessageName = w.DescriptorProto.GetName() // if they didn't specify, use normal name
	if w.DescriptorProto.GetOptions() == nil {
		return w.wasmMessageName
	}

	cand, ok := isWasmMessageName(w.DescriptorProto.GetOptions().String())
	if ok {
		w.wasmMessageName = cand
	}
	w.wasmMessageName = removeQuotes(w.wasmMessageName)
	return w.wasmMessageName
}

// GetWasmMethodName looks through the data structure given that represents the
// original protobuf structure trying to find constructs like this:
//
//		service Foo {
//	  rpc Baz(BazRequest) returns (BazResponse) {
//		   option (parigot.wasm_method_name) = "heffalumph";
//	}
//
// If no such construction is found, it returns the simple name.
func (w *WasmMethod) GetWasmMethodName() string {
	if w.wasmMethodName != "" {
		return w.wasmMethodName
	}
	// look for it
	w.wasmMethodName = w.MethodDescriptorProto.GetName() // if they didn't specify, use normal name
	if w.MethodDescriptorProto.GetOptions() == nil {
		return w.wasmMethodName
	}

	cand, ok := isWasmMessageName(w.MethodDescriptorProto.GetOptions().String())
	if ok {
		w.wasmMethodName = cand
	}
	w.wasmMethodName = removeQuotes(w.wasmMethodName)
	return w.wasmMethodName
}

// GetWasmFieldName looks through the data structure given that represents the
// original protobuf structure trying to find constructs like this:
//
//		Message Bar {
//	  rpc Grill(GrillRequest) returns (GrillResponse) {
//		   option (parigot.wasm_field_name) = "tgimcfunsters";
//	}
//
// If no such construction is found, it returns the simple name.
func (w *WasmField) GetWasmFieldName() string {
	if w.wasmFieldName != "" {
		return w.wasmFieldName
	}

	// look for it
	w.wasmFieldName = w.FieldDescriptorProto.GetName() // if they didn't specify, use normal name
	if w.FieldDescriptorProto.GetOptions() == nil {
		return w.wasmFieldName
	}
	cand, ok := isWasmMessageName(w.FieldDescriptorProto.GetOptions().String())
	if ok {
		w.wasmFieldName = cand
	}
	w.wasmFieldName = removeQuotes(w.wasmFieldName)
	return w.wasmFieldName
}

func (g *GenInfo) IsAbi() bool {
	if g.file.GetOptions() == nil {
		log.Printf("no options on %s\n", g.file.GetName())
		return false
	}
	return IsAbi(g.file.GetOptions().String())
}

// GetWasmMethod returns all the wasm methods contained inside this service.
func (s *WasmService) GetWasmMethod() []*WasmMethod {
	return s.method
}

// GetField returns all the wasm field contained inside this message.
func (m *WasmMessage) GetField() []*WasmField {
	return m.field
}

// GetParent returns the parent of this wasm service, which is a descriptor
// for the file containing it.
func (s *WasmService) GetParent() *descriptorpb.FileDescriptorProto {
	return s.parent
}

// GetParent returns the parent of this wasm message, which is a descriptor
// for the file containing it.
func (m *WasmMessage) GetParent() *descriptorpb.FileDescriptorProto {
	return m.parent
}

// GetParent returns the parent of this wasm method, which is a wasm service.
func (m *WasmMethod) GetParent() *WasmService {
	return m.parent
}

// GetParent returns the parent of this wasm field, which is a wasm message.
func (m *WasmField) GetParent() *WasmMessage {
	return m.parent
}

type InputParam struct {
	name     string
	typ      *WasmMessage
	paramVar []*ParamVar
	lang     LanguageText
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

func newInputParam(g *GenInfo, messageName string) *InputParam {
	msg := g.findMessageByName(messageName)
	if msg == nil {
		log.Fatalf("unable to find input parameter type %s", lastSegmentOfPackage(messageName))
	}

	result := &InputParam{
		typ:  msg,
		name: msg.GetWasmMessageName(),
	}
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

func newOutputParam(g *GenInfo, name string) *OutputParam {
	w := g.findMessageByName(name)
	if w == nil {
		log.Fatalf("unable to find type %s for output parameter", name)
	}
	result := &OutputParam{
		name: w.GetWasmMessageName(),
		typ:  w,
	}
	result.paramVar = make([]*ParamVar, len(w.GetField()))
	for i, f := range w.GetField() {
		result.paramVar[i] = newParamVar(f)
	}
	return result
}

func (m *WasmMethod) EmtpyInput() bool {
	return m.input.IsEmpty()
}
func (m *WasmMethod) NotEmtpyInput() bool {
	return !m.input.IsEmpty()
}
func (m *WasmMethod) EmtpyOutput() bool {
	return m.output.IsEmpty()
}
func (m *WasmMethod) NotEmptyOutput() bool {
	return !m.output.IsEmpty()
}
func (m *WasmMessage) GetFullName() string {
	return m.GetParent().GetPackage() + "." + m.GetWasmMessageName()
}
func (m *WasmMethod) InParams() []*ParamVar {
	return m.input.paramVar
}

func (m *WasmMethod) OutType() string {
	return m.lang.OutType(m)
}

func (m *WasmMethod) AllInputParamWithFormalWasmLevel(showFormalName bool) string {
	return m.abilang.AllInputParamWithFormalWasmLevel(m, showFormalName)
}

func (m *WasmMethod) HasComplexParam() bool {
	return !m.NoComplexParam()
}

func (m *WasmMethod) HasComplexOutput() bool {
	return !m.NoComplexOutput()
}

func (m *WasmMethod) NoComplexParam() bool {
	for _, p := range m.input.paramVar {
		if !p.IsStrictWasmType() {
			return false
		}
	}
	return true
}

func (m *WasmMethod) NoComplexOutput() bool {
	for _, p := range m.output.paramVar {
		if !p.IsStrictWasmType() {
			return false
		}
	}
	return true
}

func (m *WasmMethod) AllInputWithFormal(showFormalName bool) string {
	return m.lang.AllInputWithFormal(m, showFormalName)
}

func (m *WasmMethod) AllInputFormal() string {
	return m.lang.AllInputFormal(m)
}

func (m *WasmMethod) OutZeroValue() string {
	return m.lang.OutZeroValue(m)
}

// toGoType returns the string that is the equivalent of the given string, a
// protobuf type.
func toGoType(s string) string {
	switch s {
	case "TYPE_STRING":
		return "string"
	case "TYPE_INT32":
		return "int32"
	case "TYPE_INT64":
		return "int64"
	case "TYPE_FLOAT":
		return "float32"
	case "TYPE_DOUBLE":
		return "float64"
	case "TYPE_BOOL":
		return "bool"
	}
	panic("unable to convert " + s + " to go type")
}

func (m *WasmMethod) AllInputNumberedParam() string {
	return m.lang.AllInputNumberedParam(m)
}

func (m *WasmMethod) AllInputParamWasmToGoImpl() string {
	return m.abilang.AllInputParamWasmToGoImpl(m)
}

func removeQuotes(s string) string {
	result := s
	l := len(s)
	if l > 2 && s[0:1] == "\"" && s[l-1:l] == "\"" {
		result = s[1 : l-1]
	}
	return result
}

func (g *GenInfo) findMessageByName(n string) *WasmMessage {
	for _, m := range g.wasmMessage {
		log.Printf("comparing %s to %s", m.GetFullName(), n)
		if m.GetFullName() == n {
			return m
		}
	}
	// why do they do this SOME of the time?
	if len(n) > 0 && n[0] == '.' {
		for _, m := range g.wasmMessage {
			log.Printf("comparing %s to %s", m.GetFullName(), n[1:])
			if m.GetFullName() == n[1:] {
				return m
			}
		}
	}
	return nil
}
