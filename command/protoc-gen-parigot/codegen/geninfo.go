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
}

// WasmService is like a descriptorpb.ServiceDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmService in the templates.
type WasmService struct {
	*descriptorpb.ServiceDescriptorProto
	wasmServiceName string
	parent          *descriptorpb.FileDescriptorProto
	method          []*WasmMethod
}

// WasmMessage is like a descriptorpb.DescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmMessage in the templates.
type WasmMessage struct {
	*descriptorpb.DescriptorProto
	wasmMessageName string
	parent          *descriptorpb.FileDescriptorProto
	field           []*WasmField
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
}

// WasmField is like a descriptorpb.FieldDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmField in the templates.
type WasmField struct {
	*descriptorpb.FieldDescriptorProto
	wasmFieldName string
	parent        *WasmMessage
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
	result := w.ServiceDescriptorProto.GetName() // if they didn't specify, use normal name
	if w.ServiceDescriptorProto.GetOptions() == nil {
		return w.wasmServiceName
	}

	cand, ok := isWasmServiceName(w.ServiceDescriptorProto.GetOptions().String())
	if ok {
		result = cand
	}
	return result
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
	return w.wasmFieldName
}

func (g *GenInfo) IsAbi() bool {
	if g.file.GetOptions() == nil {
		log.Printf("no options on %s\n", g.file.GetName())
		return false
	}
	return isAbi(g.file.GetOptions().String())
}

// GetWasmMethod returns all the wasm methods contained inside this service.
func (s *WasmService) GetWasmMethod() []*WasmMethod {
	return s.method
}

// GetField returns all the wasm field contained inside this message.
func (m *WasmMessage) GetField() []*WasmField {
	return m.field
}

// GetFullName returns the name of the object but with package prepended.
func (m *WasmMessage) GetFullName() string {
	return m.parent.GetPackage() + "." + m.GetWasmMessageName()
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
}

func newInputParam(m *WasmMessage) *InputParam {
	result := &InputParam{
		typ:  m,
		name: m.GetWasmMessageName(),
	}
	paramVar := make([]*ParamVar, len(m.GetField()))
	for i, f := range m.GetField() {
		pVar := newParamVar(f)
		paramVar[i] = pVar
	}
	return result
}

func (i *InputParam) IsEmpty() bool {
	return len(i.paramVar) == 0
}

func (i *InputParam) GetParamVar() []*ParamVar {
	return i.paramVar
}

func (i *InputParam) Len() int {
	return len(i.paramVar)
}

type ParamVar struct {
	name  string
	field *WasmField
}

func newParamVar(f *WasmField) *ParamVar {
	return &ParamVar{
		name:  f.GetWasmFieldName(),
		field: f,
	}
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
}

func (o *OutputParam) IsEmpty() bool {
	return len(o.paramVar) == 0
}

func (o *OutputParam) IsMultipleReturn() bool {
	return len(o.paramVar) > 1
}

func newOutputParam(w *WasmMessage) *OutputParam {
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

func (m *WasmMethod) InParams() []*ParamVar {
	return m.input.paramVar
}

func (m *WasmMethod) OutParamas() []*ParamVar {
	if m.output.IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time", m.output.name)
	}
	return m.input.paramVar
}

func (m *WasmMethod) OutType() string {
	if m.output.IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time", m.output.name)
	}
	return m.input.paramVar[0].field.GetType().String()
}

func (m *WasmMethod) AllInputsWithFormals(showFormalName bool) string {
	result := ""
	for i, p := range m.input.paramVar {
		if showFormalName {
			result += p.name + " "
		} else {
			result += "_" + " "
		}
		result += toGoType(p.field.GetType().String())
		if i != len(m.input.paramVar)-1 {
			result += ","
		}
	}
	return result
}

func (m *WasmMethod) OutZeroValue() string {
	panic("outzero")
	if m.output.IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time", m.output.name)
	}
	protoT := m.output.paramVar[0].field.GetType().String()
	goT := toGoType(protoT)
	return goZeroValue(goT)
}

// goZeroValue returns the simplest, empty value for the given go type.
func goZeroValue(s string) string {
	switch s {
	case "string":
		return ""
	case "int32":
		return "int32(0)"
	case "int64":
		return "int64(0)"
	case "float32":
		return "float32(0.0)"
	case "float64":
		return "float64(0.0)"
	case "bool":
		return "false"
	}
	panic("unable to get zero value for go type " + s)
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
