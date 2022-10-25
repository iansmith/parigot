package codegen

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
	"strings"
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
	pullParameters       bool
	pullOutput           bool
	abiCall              bool
	protoPackageOverride *string
}

func (w *WasmMethod) InputFields() []*WasmField {
	if w.CGInput().Len() == 0 {
		log.Fatalf("attempt to use InputFields but no input fields present in %s",
			w.WasmMethodName())
	}
	return w.CGInput().CGType().CompositeType().GetField()
}

func (w *WasmMethod) GetOutputFields() []*WasmField {
	if w.CGOutput().Len() == 0 {
		log.Fatalf("attempt to use GetOutputFields but no out fields present in %s",
			w.WasmMethodName())
	}
	return w.CGOutput().GetCGType().CompositeType().GetField()
}

func (w *WasmMethod) importForMessage(m *WasmMessage) string {
	fullName := m.GetFullName()
	parts := strings.Split(fullName, ".")
	formattedName := w.Finder().AddressingNameFromMessage(w.ProtoPackage(), m)
	if len(parts) > 2 {
		return fmt.Sprintf("github.com/iansmith/parigot/g/%s", strings.Join(parts[:len(parts)-1], "/"))
	}
	log.Fatalf("method: %s: full name of input: %s, formatted name: %s [%d]",
		w.WasmMethodName(), fullName, formattedName, len(parts))
	return ""
}
func (w *WasmMethod) addImportForInput(comp *WasmMessage, imp map[string]struct{}) {
	imp[w.importForMessage(comp)] = struct{}{}
}
func (w *WasmMethod) addImportForOutput(comp *WasmMessage, imp map[string]struct{}) {
	imp[w.importForMessage(comp)] = struct{}{}
}
func (w *WasmMethod) AddImportsNeeded(imp map[string]struct{}) {
	w.addImportForInput(w.CGInput().CGType().CompositeType(), imp)
	w.addImportForOutput(w.CGOutput().GetCGType().CompositeType(), imp)
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

func (w *WasmMethod) ProtoPackage() string {
	if w.protoPackageOverride != nil {
		return *w.protoPackageOverride
	}
	return w.Parent().ProtoPackage()
}

func (w *WasmMethod) Finder() Finder {
	return w.Parent().Finder()
}

func (w *WasmMethod) SetProtoPackage(p string) string {
	w.protoPackageOverride = new(string)
	*w.protoPackageOverride = p
	return ""
}

func (w *WasmMethod) GoPackage() string {
	return w.Parent().GetGoPackage()
}

func (w *WasmMethod) CGInput() *InputParam {
	return w.input
}
func (w *WasmMethod) CGOutput() *OutputParam {
	return w.output
}

// WasmMethodName looks through the data structure given that represents the
// original protobuf structure trying to find constructs like this:
//
//		service Foo {
//	  rpc Baz(BazRequest) returns (BazResponse) {
//		   option (parigot.wasm_method_name) = "heffalumph";
//	}
//
// If no such construction is found, it returns the simple name.
func (w *WasmMethod) WasmMethodName() string {
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

// Parent returns the parent of this wasm method, which is a wasm service.
func (m *WasmMethod) Parent() *WasmService {
	return m.parent
}
func (m *WasmMethod) EmtpyInput() bool {
	return m.CGInput().CGType() == nil
}
func (m *WasmMethod) NotEmtpyInput() bool {
	return m.CGInput().CGType() != nil
}
func (m *WasmMethod) EmtpyOutput() bool {
	return m.CGOutput().GetCGType() == nil
}
func (m *WasmMethod) NotEmptyOutput() bool {
	if m.PullOutput() {
		exp := ExpandReturnInfoForOutput(m.CGOutput(), m, m.ProtoPackage())
		if exp == nil {
			return false
		}
		return !exp.GetCGType().IsEmpty()
	}
	t := m.CGOutput().GetCGType()
	if t == nil {
		return false
	}
	return !t.IsEmpty()
}
func (m *WasmMethod) InputParam() *InputParam {
	return m.input
}
func (m *WasmMethod) OutputParam() *OutputParam {
	return m.output
}
func (m *WasmMethod) InputCodeNeeded() bool {
	if m.CGInput() == nil {
		panic("should have input param")
	}
	if m.CGInput().CGType() == nil {
		panic("should have input param's type")
	}
	if m.CGInput().CGType().IsBasic() {
		panic("should not have a simple type at top level of proto definition")
	}
	if m.CGInput().CGType().IsCompositeNoFields() {
		return false
	}
	return true
}

func (m *WasmMethod) OutputCodeNeeded() bool {
	if m.CGOutput() == nil {
		panic("should have an output param")
	}
	if m.CGOutput().GetCGType() == nil {
		panic("should have param's type")
	}
	if m.CGOutput().GetCGType().IsBasic() {
		panic("should not have a simple type at top level of a proto definition")
	}
	cgt := m.CGOutput().GetCGType()
	if cgt.IsCompositeNoFields() {
		return false
	}
	//if m.PullOutput() {
	//	field := cgt.CompositeType().GetField()
	//	if len(field) > 1 {
	//		log.Fatalf("Currently, it is not allowed to pull up an output with more than 1 field:%s",
	//			cgt.CompositeType().GetName())
	//	}
	//	// we checked for zero before, so we know len(field)==1
	//	only := cgt.CompositeType().GetField()[0]
	//	cgt = NewCGTypeFromField(only, m, m.ProtoPackage())
	//	// we ban the pull up to a composite, because that's just confusing
	//	// should we recurse on that?
	//	if !cgt.IsBasic() {
	//		log.Fatalf("Pulling up an output to another composite result is ambiguous:",
	//			cgt.CompositeType().GetName())
	//	}
	//	if cgt.Basic() == "TYPE_INT64" {
	//		log.Fatalf("Due to a javascript limitation, you cannot return a single int64 from a function "+
	//			" and your pull up of the output of type %s would result in that", m.CGOutput().GetTypeName())
	//	}
	//}
	return true
}

func recurseCompositesToFindTypeSequence(t *CGType, m *WasmMethod, seq []string) []string {
	if t.IsCompositeNoFields() {
		return []string{}
	}
	if t.IsBasic() {
		return m.Language().BasicTypeToWasm(t.Basic())
	}
	//it's a composite with 1 or more fields
	comp := t.CompositeType()
	field := comp.GetField()
	result := seq
	for _, f := range field {
		cgt := NewCGTypeFromField(f, m, m.ProtoPackage())
		if !cgt.IsBasic() {
			result = append(result, recurseCompositesToFindTypeSequence(t, m, seq)...)
		} else {
			result = append(result, m.Language().BasicTypeToWasm(cgt.Basic())...)
		}
	}
	return result
}

func (m *WasmMethod) FuncChoice() *FuncChooser {
	return m.Language().FuncChoice()
}

func (m *WasmMethod) MarkInputOutputMessages() {
	if !m.input.IsEmpty() {
		if !m.input.CGType().IsBasic() {
			m.input.CGType().CompositeType().MarkSource(true, m)
		}
	}
	if !m.output.IsEmpty() {
		if !m.output.GetCGType().IsBasic() {
			m.output.GetCGType().CompositeType().MarkSource(false, m)
		}
	}
}

func (m *WasmMethod) Language() LanguageText {
	return m.Parent().GetLanguage()
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
