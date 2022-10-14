package codegen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
)

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

func (w *WasmMethod) HasNoPackageOption() bool {
	return w.parent.HasNoPackageOption()
}

func (w *WasmMethod) GetProtoPackage() string {
	return w.GetParent().GetProtoPackage()
}

func (w *WasmMethod) GetFinder() Finder {
	return w.GetParent().GetFinder()
}

func (w *WasmMethod) GetGoPackage() string {
	return w.GetParent().GetGoPackage()
}

func (w *WasmField) GetProtoPackage() string {
	return w.GetParent().GetProtoPackage()
}

func (w *WasmMethod) GetCGInput() *InputParam {
	return w.input
}
func (w *WasmMethod) GetCGOutput() *OutputParam {
	return w.output
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

// GetParent returns the parent of this wasm method, which is a wasm service.
func (m *WasmMethod) GetParent() *WasmService {
	return m.parent
}
func (m *WasmMethod) EmtpyInput() bool {
	return m.GetCGInput().GetCGType() == nil
}
func (m *WasmMethod) NotEmtpyInput() bool {
	return m.GetCGInput().GetCGType() != nil
}
func (m *WasmMethod) EmtpyOutput() bool {
	return m.GetCGOutput().GetCGType() == nil
}
func (m *WasmMethod) NotEmptyOutput() bool {
	return m.GetCGOutput().GetCGType() == nil
}
func (m *WasmMethod) GetInputParam() *InputParam {
	return m.input
}
func (m *WasmMethod) GetOutputParam() *OutputParam {
	return m.output
}
func (m *WasmMethod) OutType() string {
	return m.GetParent().GetLanguage().OutType(m)
}
func (m *WasmMethod) AllInputWithFormalWasmLevel(showFormalName bool) string {
	return m.GetParent().GetLanguage().(AbiLanguageText).AllInputWithFormalWasmLevel(m, showFormalName)
}
func (m *WasmMethod) HasComplexParam() bool {
	return !m.NoComplexParam()
}

func (m *WasmMethod) HasComplexOutput() bool {
	return !m.NoComplexOutput()
}

func (m *WasmMethod) NoComplexParam() bool {
	return true
}

func (m *WasmMethod) NoComplexOutput() bool {
	return true
}

func (m *WasmMethod) AllInputWithFormal(showFormalName bool) string {
	if m.PullParameters() {
		log.Printf("trying to get to pulling params: %s,%s", m.GetName(),
			m.GetCGInput().GetCGType().String(m.GetProtoPackage()))
	}
	return m.GetLanguage().AllInputWithFormal(m, showFormalName)
}

func (m *WasmMethod) GetLanguage() LanguageText {
	return m.GetParent().GetLanguage()
}
func (m *WasmMethod) PullParameters() bool {
	return m.parent.AlwaysPullParameters()
}

func (m *WasmMethod) AllInputFormal() string {
	return m.GetLanguage().AllInputFormal(m)
}

func (m *WasmMethod) OutZeroValue() string {
	return m.GetLanguage().OutZeroValue(m)
}

func (m *WasmMethod) AllInputNumberedParam() string {
	return m.GetLanguage().AllInputNumberedParam(m)
}

func (m *WasmMethod) AllInputWasmToGoImpl() string {
	return m.GetLanguage().(AbiLanguageText).AllInputWasmToGoImpl(m)
}
