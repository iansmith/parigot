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
	// these values doesn't matter if your parent service has the "always"
	// version of it set
	pullParameters bool
	pullOutput     bool
	abiCall        bool
}

func (w *WasmMethod) GetInputFields() []*WasmField {
	if w.GetCGInput().Len() == 0 {
		log.Fatalf("attempt to use GetInputFields but no input fields present in %s",
			w.GetWasmMethodName())
	}
	return w.GetCGInput().GetCGType().GetCompositeType().GetField()
}

func (w *WasmMethod) GetOutputFields() []*WasmField {
	if w.GetCGOutput().Len() == 0 {
		log.Fatalf("attempt to use GetOutputFields but no out fields present in %s",
			w.GetWasmMethodName())
	}
	return w.GetCGOutput().GetCGType().GetCompositeType().GetField()
}

func (w *WasmMethod) HasNoPackageOption() bool {
	return w.parent.HasNoPackageOption()
}

func (w *WasmMethod) HasAbiCallOption() bool {
	return w.abiCall
}
func (w *WasmMethod) NoAbiCallOption() bool {
	return !w.abiCall
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
	if m.PullOutput() {
		exp := ExpandReturnInfoForOutput(m.GetCGOutput(), m, m.GetProtoPackage())
		if exp == nil {
			return false
		}
		return !exp.GetCGType().IsEmpty()
	}
	t := m.GetCGOutput().GetCGType()
	if t == nil {
		return false
	}
	return !t.IsEmpty()
}
func (m *WasmMethod) GetInputParam() *InputParam {
	return m.input
}
func (m *WasmMethod) GetOutputParam() *OutputParam {
	return m.output
}
func (m *WasmMethod) InputCodeNeeded() bool {
	if m.GetCGInput() == nil {
		panic("should have input param")
	}
	if m.GetCGInput().GetCGType() == nil {
		panic("should have input param's type")
	}
	if m.GetCGInput().GetCGType().IsBasic() {
		panic("should not have a simple type at top level of method")
	}
	if m.GetCGInput().GetCGType().IsEmpty() {
		return false
	}
	if m.GetCGInput().GetCGType().IsCompositeNoFields() {
		return false
	}
	return true
}
func (m *WasmMethod) OutputCodeNeeded() bool {
	if m.GetCGOutput() == nil {
		panic("should have an output param")
	}
	if m.GetCGOutput().GetCGType() == nil {
		panic("should have input param's type")
	}
	if m.GetCGOutput().GetCGType().IsBasic() {
		panic("should not have a simple type at top level of method")
	}
	if m.GetCGOutput().GetCGType().IsEmpty() {
		return false
	}
	if m.GetCGOutput().GetCGType().IsCompositeNoFields() {
		return false
	}
	return true
}

func (m *WasmMethod) FuncChoice() *FuncChooser {
	return m.GetLanguage().FuncChoice()
}

func (m *WasmMethod) MarkInputOutputMessages() {
	if !m.input.IsEmpty() {
		if !m.input.GetCGType().IsBasic() {
			m.input.GetCGType().GetCompositeType().MarkSource(true, m)
		}
	}
	if !m.output.IsEmpty() {
		if !m.output.GetCGType().IsBasic() {
			m.output.GetCGType().GetCompositeType().MarkSource(false, m)
		}
	}
}

func (m *WasmMethod) GetLanguage() LanguageText {
	return m.GetParent().GetLanguage()
}

func (m *WasmMethod) PullParameters() bool {
	if m.parent.AlwaysPullParameters() {
		return true
	}
	return m.pullParameters
}

func (m *WasmMethod) PullOutput() bool {
	if m.parent.AlwaysPullOutput() {
		return true
	}
	return m.pullOutput
}
