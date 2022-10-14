package codegen

import "google.golang.org/protobuf/types/descriptorpb"

// WasmField is like a descriptorpb.FieldDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmField in the templates.
type WasmField struct {
	*descriptorpb.FieldDescriptorProto
	wasmFieldName string
	parent        *WasmMessage
	cgType        *CGType
}

func (w *WasmField) GetGoPackage() string {
	return w.GetParent().GetGoPackage()
}

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

// GetParent returns the parent of this wasm field, which is a wasm message.
func (m *WasmField) GetParent() *WasmMessage {
	return m.parent
}
