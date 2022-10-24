package codegen

import "google.golang.org/protobuf/types/descriptorpb"

// WasmService is like a descriptorpb.ServiceDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmService in the templates.
type WasmService struct {
	*descriptorpb.ServiceDescriptorProto
	wasmServiceName      string
	parent               *descriptorpb.FileDescriptorProto
	method               []*WasmMethod
	lang                 LanguageText
	alwaysPullParameters bool
	alwaysPullOutput     bool
	noPackage            bool
	finder               Finder
}

func (w *WasmService) Finder() Finder {
	return w.finder
}

func (w *WasmService) GetLanguage() LanguageText {
	return w.lang
}
func (w *WasmService) HasNoPackageOption() bool {
	return w.noPackage
}
func (w *WasmService) ProtoPackage() string {
	return w.GetParent().GetPackage()
}
func (w *WasmService) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
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

// GetWasmMethod returns all the wasm methods contained inside this service.
func (s *WasmService) GetWasmMethod() []*WasmMethod {
	return s.method
}

// GetParent returns the parent of this wasm service, which is a descriptor
// for the file containing it.
func (s *WasmService) GetParent() *descriptorpb.FileDescriptorProto {
	return s.parent
}
func (s *WasmService) AlwaysPullParameters() bool {
	return s.alwaysPullParameters
}

func (s *WasmService) AlwaysPullOutput() bool {
	return s.alwaysPullOutput
}
