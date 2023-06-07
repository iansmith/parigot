package codegen

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// WasmService is like a descriptorpb.ServiceDescriptorProto (which it contains) but
// also adds things that are specific to parigot, notably wasm-specific names.  This
// should be accessed by .GetWasmService in the templates.
type WasmService struct {
	*descriptorpb.ServiceDescriptorProto
	wasmServiceName      string
	wasmServiceErrId     string
	parent               *descriptorpb.FileDescriptorProto
	method               []*WasmMethod
	lang                 LanguageText
	alwaysPullParameters bool
	alwaysPullOutput     bool
	finder               Finder
	kernel               bool
	isTest               bool
}

func (w *WasmService) Finder() Finder {
	return w.finder
}

func (w *WasmService) GetLanguage() LanguageText {
	return w.lang
}
func (w *WasmService) HasKernelOption() bool {
	return w.kernel
}
func (w *WasmService) NoKernelOption() bool {
	return !w.kernel
}

func (w *WasmService) ProtoPackage() string {
	return w.GetParent().GetPackage()
}
func (w *WasmService) GetGoPackage() string {
	return w.GetParent().GetOptions().GetGoPackage()
}
func (w *WasmService) IsTest() bool {
	return w.isTest
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

// GetWasmServiceName looks through the data structure given that represents the
// original protobuf structure trying to fin constructs like this:
//
//	service Foo {
//	   option (parigot.wasm_service_err_id) = "MyErrId";
//
// If no such construction is found, it returns the golang package name.
func (w *WasmService) GetWasmServiceErrId() string {
	if w.wasmServiceErrId != "" {
		return w.wasmServiceErrId
	}

	//	if w.ServiceDescriptorProto.GetOptions() == nil {
	pkg := w.GetGoPackage()
	part := strings.Split(pkg, ";")
	if len(part) != 2 {
		panic(fmt.Sprintf("unable to understand pkg name from protobuf file: %s", pkg))
	}
	cand := strings.ToUpper(part[1][0:1]) + part[1][1:] + "ErrId"
	w.wasmServiceErrId = cand
	return cand
	//	}
	//
	// cand, ok := isWasmServiceErrId(w.ServiceDescriptorProto.GetOptions().String())
	// w.wasmServiceErrId = cand
	// log.Printf("found an option! candidate is %s, ok is %v", cand, ok)
	// return cand
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

func (s *WasmService) AddImportsNeeded(imp map[string]struct{}) {
	for _, m := range s.GetWasmMethod() {
		m.AddImportsNeeded(imp)
	}
	//if s.kernel {
	//	return
	//}
	//imp["github.com/iansmith/parigot/lib"] = struct{}{}
	//imp["google.golang.org/protobuf/proto"] = struct{}{}
	//imp["fmt"] = struct{}{}
}

func (s *WasmService) Collect() {
	s.method = make([]*WasmMethod, len(s.GetMethod()))
	for j, m := range s.GetMethod() {
		s.method[j] = NewWasmMethod(m, s)
	}
}
