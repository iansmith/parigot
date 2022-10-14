package codegen

import "google.golang.org/protobuf/types/descriptorpb"

// WasmMessage is like a descriptorpb.DescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmMessage in the templates.
type WasmMessage struct {
	*descriptorpb.DescriptorProto
	wasmMessageName string
	parent          *descriptorpb.FileDescriptorProto
	field           []*WasmField
	lang            LanguageText
	finder          Finder
	noPackage       bool
}

func (w *WasmMessage) GetAddressableName(from string) string {
	return w.finder.AddressingNameFromMessage(from, w)
}

func (w *WasmMessage) GetFinder() Finder {
	return w.finder
}

func (w *WasmMessage) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}
func (w *WasmMessage) GetProtoPackage() string {
	return w.GetParent().GetPackage()
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

// GetField returns all the wasm field contained inside this message.
func (m *WasmMessage) GetField() []*WasmField {
	field := m.field
	return field
}

// GetParent returns the parent of this wasm message, which is a descriptor
// for the file containing it.
func (m *WasmMessage) GetParent() *descriptorpb.FileDescriptorProto {
	return m.parent
}
func (m *WasmMessage) GetFullName() string {
	return m.GetParent().GetPackage() + "." + m.GetWasmMessageName()
}

func (m *WasmMessage) HasNoPackageOption() bool {
	return m.noPackage
}

func (m *WasmMessage) NotGoogleMessage() bool {
	return m.GetProtoPackage() != "google.protobuf" && m.GetGoPackage() != "google.golang.org/protobuf/types/descriptorpb)"
}
