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
	// this value doesn't matter if your parent service has the "always"
	// pull parameters flag set
	pullParameters bool
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
	if m.PullParameters() {
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
func (m *WasmMethod) OutType() string {
	return m.GetParent().GetLanguage().OutType(m)
}

func (m *WasmMethod) OutTypeDecl() string {
	return m.GetLanguage().OutTypeDecl(m)
}

func (m *WasmMethod) OutZeroValueDecl() string {
	return m.GetLanguage().OutZeroValueDecl(m)
}

func (m *WasmMethod) RequiresDecode() bool {
	x := m.HasComplexParam()
	y := m.HasComplexOutput()
	return x || y
}
func (m *WasmMethod) NoDecodeRequired() bool {
	return !m.RequiresDecode()
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
	p := m.GetInputParam().GetCGType()
	if p == nil {
		return true
	}
	if p.IsEmpty() {
		return true
	}
	if m.PullParameters() {
		param := ExpandParamInfoForInput(m.GetCGInput(), m, m.GetProtoPackage())
		if param == nil || len(param) == 0 {
			return true
		}
		for _, exp := range param {
			if !exp.GetCGType().IsStrictWasmType() {
				return false
			}
		}
		return true
	} else {
		c := m.GetInputParam().GetCGType()
		return c.IsStrictWasmType()
	}
}

func (m *WasmMethod) GetNumberParametersUsed(c *CGType) int {
	return m.GetLanguage().GetNumberParametersUsed(c)
}
func (m *WasmMethod) NoComplexOutput() bool {
	if m.PullParameters() {
		exp := ExpandReturnInfoForOutput(m.GetCGOutput(), m, m.GetProtoPackage())
		return exp == nil || exp.GetCGType().IsStrictWasmType()
	} else {
		return m.GetCGOutput().GetCGType().IsStrictWasmType()
	}
}

func (m *WasmMethod) AllInputWithFormal(showFormalName bool) string {
	return m.GetLanguage().AllInputWithFormal(m, showFormalName)
}

func (m *WasmMethod) GetLanguage() LanguageText {
	return m.GetParent().GetLanguage()
}

func (m *WasmMethod) GetFormalArgSeparator() string {
	return m.GetLanguage().GetFormalArgSeparator()
}

func (m *WasmMethod) PullParameters() bool {
	if m.parent.AlwaysPullParameters() {
		return true
	}
	return m.pullParameters
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
