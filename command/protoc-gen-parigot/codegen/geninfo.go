package codegen

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type GenInfo struct {
	request         *pluginpb.CodeGeneratorRequest
	file            *descriptorpb.FileDescriptorProto
	lang            LanguageText
	finder          *SimpleFinder
	kernelInterface string
}

func NewGenInfo() *GenInfo {
	return &GenInfo{
		finder: NewSimpleFinder(),
	}
}

// SetKernelInterface is called by the driver to set the blob of text that needs
// to be added to the kernel service to make it interface correctly with the kernel.
// The generator supplies the filename that this is read from in the templates dir.
func (g *GenInfo) SetKernelInterface(s string) {
	g.kernelInterface = s
}

// KernelInterface is called by a template for a particular language when it wants
// to emit the blob of text that is the extra, kernel-specific code.
func (g *GenInfo) KernelInterface() string {
	return g.kernelInterface
}

// Service returns all the wasm services of this GenInfo.
// The caller should not change these values, these are pointers in the real
// data structures.
func (g *GenInfo) Service() []*WasmService {
	return g.finder.Service()
}

// Message returns all the wasm messages of this GenInfo.
// The caller should not change these values, these are pointers in the real
// data structures.
func (g *GenInfo) Message() []*WasmMessage {
	return g.finder.Message()
}

// GetFile returns the file data structure that is associated with the currently
// operating code generation.
func (g *GenInfo) GetFile() *descriptorpb.FileDescriptorProto {
	return g.file
}

// GetRequest returns the original request data structure that was delivered via
// stdin.  This request is associated with the currently operating code generation.
func (g *GenInfo) GetRequest() *pluginpb.CodeGeneratorRequest {
	return g.request
}

type InputParam struct {
	lang   LanguageText
	parent *WasmMethod
	cgType *CGType
}

func (i *InputParam) SetCGType(cg *CGType) {
	if i.cgType != nil {
		panic("attempt to replace existing cgType in input param")
	}
	i.cgType = cg
}

func (i *InputParam) CGType() *CGType {
	return i.cgType
}

func (i *InputParam) GetTypeName() string {
	p := i.GetParent()
	return p.MethodDescriptorProto.GetInputType()
}
func (i *InputParam) GetLanguage() LanguageText {
	return i.lang
}

func (i *InputParam) GetParent() *WasmMethod {
	return i.parent
}

func (i *InputParam) SetEmpty() {
	i.cgType.SetEmpty()
}

func (i *InputParam) IsEmpty() bool {
	return i.cgType.IsEmpty()
}

func newInputParam(parent *WasmMethod) *InputParam {
	result := &InputParam{
		parent: parent,
		lang:   parent.Language(),
	}
	result.cgType = NewCGTypeFromInput(result, parent, parent.ProtoPackage())
	return result
}

func (i *InputParam) Len() int {
	if i.cgType == nil {
		return 0
	}
	return 1
}

func (o *OutputParam) Len() int {
	if o.cgType == nil {
		return 0
	}
	return 1
}

type OutputParam struct {
	lang   LanguageText
	parent *WasmMethod
	cgType *CGType
}

func (o *OutputParam) IsEmpty() bool {
	if o.cgType == nil {
		panic("should not be calling isEmpty on output param that is not initialized")
	}
	return o.cgType.IsEmpty()
}

func (o *OutputParam) GetCGType() *CGType {
	return o.cgType
}

func (o *OutputParam) GetLanguage() LanguageText {
	return o.lang
}

func (o *OutputParam) SetEmpty() {
	o.cgType.SetEmpty()
}

func (o *OutputParam) SetCGType(c *CGType) {
	if o.cgType != nil {
		panic("attempt to setCGType on an output param with it already set")
	}
	o.cgType = c
}

func (o *OutputParam) GetTypeName() string {
	p := o.GetParent()
	return p.MethodDescriptorProto.GetOutputType()
}

func (o *OutputParam) GetParent() *WasmMethod {
	return o.parent
}

func (o *OutputParam) IsMultipleReturn() bool {
	return false
}

func newOutputParam(parent *WasmMethod) *OutputParam {
	result := &OutputParam{
		parent: parent,
	}
	result.cgType = NewCGTypeFromOutput(result, parent, parent.ProtoPackage())
	return result
}

func removeQuotes(s string) string {
	result := s
	l := len(s)
	if l > 2 && s[0:1] == "\"" && s[l-1:l] == "\"" {
		result = s[1 : l-1]
	}
	return result
}

type ServiceRecord struct {
	wasmName     string
	protoPackage string
	goPackage    string
}

type MessageRecord struct {
	wasmName     string
	protoPackage string
	goPackage    string
}

func NewMessageRecord(wasmName, protoPackage, goPackage string) *MessageRecord {
	return &MessageRecord{
		wasmName:     wasmName,
		protoPackage: protoPackage,
		goPackage:    goPackage,
	}
}
func NewServiceRecord(wasmName, protoPackage, goPackage string) *ServiceRecord {
	return &ServiceRecord{
		wasmName:     wasmName,
		protoPackage: protoPackage,
		goPackage:    goPackage,
	}
}

func (m *MessageRecord) String() string {
	return fmt.Sprintf("MessageRec(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (m *ServiceRecord) String() string {
	return fmt.Sprintf("ServiceRec(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (g *GenInfo) RegisterService(w *WasmService) {
	g.finder.AddServiceType(w.GetWasmServiceName(), w.ProtoPackage(), w.GetGoPackage(), w)
}

func (g *GenInfo) RegisterMessage(w *WasmMessage) {
	g.finder.AddMessageType(w.GetWasmMessageName(), w.GetProtoPackage(), w.GetGoPackage(), w)
}

func (g *GenInfo) FindServiceByName(protoPackage string, name string, next Finder) *WasmService {
	return g.finder.FindServiceByName(protoPackage, name, next)
}
func (g *GenInfo) FindMessageByName(protoPackage string, name string, next Finder) *WasmMessage {
	return g.finder.FindMessageByName(protoPackage, name, next)
}

func (g *GenInfo) AddMessageType(wasmName, protoPackage, goPackage string, message *WasmMessage) {
	g.finder.AddMessageType(wasmName, protoPackage, goPackage, message)

}
func (g *GenInfo) AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService) {
	g.finder.AddServiceType(wasmName, protoPackage, goPackage, service)
}
func (g *GenInfo) AddressingNameFromMessage(currentPkg string, message *WasmMessage) string {
	return g.finder.AddressingNameFromMessage(currentPkg, message)
}

func (g *GenInfo) SetReqAndFile(request *pluginpb.CodeGeneratorRequest, proto *descriptorpb.FileDescriptorProto) {
	g.request = request
	g.file = proto
}
