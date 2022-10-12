package codegen

import "google.golang.org/protobuf/types/descriptorpb"

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
}

func (w *WasmMessage) GetProtoPackage() string {
	return w.GetParent().GetPackage()
}

func (w *WasmMessage) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}

func (w *WasmService) GetProtoPackage() string {
	return w.GetParent().GetPackage()
}

func (w *WasmService) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}

func (w *WasmMethod) GetProtoPackage() string {
	return w.GetParent().GetProtoPackage()
}

func (w *WasmMethod) GetGoPackage() string {
	return w.GetParent().GetGoPackage()
}

func (w *WasmField) GetProtoPackage() string {
	return w.GetParent().GetProtoPackage()
}

func (w *WasmField) GetGoPackage() string {
	return w.GetParent().GetGoPackage()
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

// GetWasmServiceName looks through the data structure given that represents the
// original protobuf structure trying to fin constructs like this:
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
	return m.lang.(AbiLanguageText).AllInputParamWithFormalWasmLevel(m, showFormalName)
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

func (m *WasmMethod) AllInputNumberedParam() string {
	return m.lang.AllInputNumberedParam(m)
}

func (m *WasmMethod) AllInputParamWasmToGoImpl() string {
	return m.lang.(AbiLanguageText).AllInputParamWasmToGoImpl(m)
}

func NewWasmMessage(file *descriptorpb.FileDescriptorProto, message *descriptorpb.DescriptorProto,
	lang LanguageText) *WasmMessage {
	m := &WasmMessage{
		DescriptorProto: message,
		wasmMessageName: "",
		parent:          file,
		field:           []*WasmField{},
		lang:            lang,
	}
	return m
}

func NewWasmService(file *descriptorpb.FileDescriptorProto,
	service *descriptorpb.ServiceDescriptorProto, lang LanguageText) *WasmService {
	s := &WasmService{
		ServiceDescriptorProto: service,
		wasmServiceName:        "",
		parent:                 file,
		method:                 []*WasmMethod{},
		lang:                   lang,
	}
	return s
}
