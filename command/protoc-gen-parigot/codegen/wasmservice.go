package codegen

import (
	"log"

	"google.golang.org/protobuf/types/descriptorpb"
)

// WasmService is like a descriptorpb.ServiceDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmService in the templates.
type WasmService struct {
	*descriptorpb.ServiceDescriptorProto
	wasmServiceName  string
	wasmServiceErrId string
	parent           *descriptorpb.FileDescriptorProto
	method           []*WasmMethod
	lang             LanguageText
	finder           Finder
	errorIdName      string
	// alwaysPullParameters bool
	// alwaysPullOutput     bool
	// kernel               bool
	// isTest               bool
}

func (w *WasmService) Finder() Finder {
	return w.finder
}

func (w *WasmService) GetLanguage() LanguageText {
	return w.lang
}

// func (w *WasmService) HasKernelOption() bool {
// 	return w.kernel
// }
// func (w *WasmService) NoKernelOption() bool {
// 	return !w.kernel
// }

func (w *WasmService) ProtoPackage() string {
	return w.GetParent().GetPackage()
}
func (w *WasmService) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}

// func (w *WasmService) IsTest() bool {
// 	return w.isTest
// }

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
	n := w.ServiceDescriptorProto.GetName() // if they didn't specify, use normal name
	if w.ServiceDescriptorProto.GetOptions() == nil {
		w.wasmServiceName = n
		return n
	}

	cand, ok := isWasmServiceName(w.ServiceDescriptorProto.GetOptions().String())
	if ok {
		w.wasmServiceName = cand
		return cand
	}
	//YYY
	w.wasmServiceName = removeQuotes(n)
	return w.wasmServiceName
}

// GetWasmServiceErrId gets the string that is the type name
// of the return type of this service.  It will be
// <serviceName>Err.  It caches this value after it is
// computed the first time.
func (w *WasmService) GetWasmServiceErrId() string {
	if w.wasmServiceErrId != "" {
		return w.wasmServiceErrId
	}
	if w.errorIdName != "" {
		return removeQuotes(w.errorIdName)
	}

	return w.GetWasmServiceName() + "Err"
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

// func (s *WasmService) AlwaysPullParameters() bool {
// 	return s.alwaysPullParameters
// }

// func (s *WasmService) AlwaysPullOutput() bool {
// 	return s.alwaysPullOutput
// }

func (s *WasmService) AddImportsNeeded(imp map[string]struct{}) {
	for _, m := range s.GetWasmMethod() {
		log.Printf("wasm method %s needs %+v", m.GetName(), imp)
		m.AddImportsNeeded(imp)
	}
}

func (s *WasmService) Collect() {
	s.method = make([]*WasmMethod, len(s.GetMethod()))
	for j, m := range s.GetMethod() {
		s.method[j] = NewWasmMethod(m, s)
		opt, ok := isStringOptionPresent(m.GetOptions().String(), parigotOptionForHostFuncName)
		if ok {
			s.method[j].HostFuncName = opt
		}
	}
}
