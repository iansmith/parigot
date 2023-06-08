package codegen

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type GenInfo struct {
	request         *pluginpb.CodeGeneratorRequest
	nameToFile      map[string]*descriptorpb.FileDescriptorProto
	nameToSvc       map[string][]*descriptorpb.ServiceDescriptorProto
	nameToMsg       map[string][]*descriptorpb.DescriptorProto
	nameToEnumType  map[string][]*descriptorpb.EnumDescriptorProto
	nameToEnumValue map[string][]*descriptorpb.EnumValueDescriptorProto
	finder          *SimpleFinder
}

func NewGenInfo() *GenInfo {
	return &GenInfo{
		finder: NewSimpleFinder(),
	}
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

// GetService returns all the services we know about that will need to be generated.
func (g *GenInfo) GetService(name string) []*descriptorpb.ServiceDescriptorProto {
	return g.nameToSvc[name]
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
func (o *OutputParam) CGType() *CGType {
	return o.cgType
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
	msg          string
	protoPackage string
	goPackage    string
}

type EnumValueRecord struct {
	parent                 *EnumTypeRecord
	wasmName               string
	protoPackage           string
	goPackage              string
	range_end, range_start int
}
type EnumTypeRecord struct {
	protoPackage string
	goPackage    string
	wasmName     string
	_            *EnumValueRecord
}

func (m *MessageRecord) WasmName() string {
	return m.msg
}

func NewMessageRecord(name, protoPackage, goPackage string) *MessageRecord {
	return &MessageRecord{
		msg:          name,
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
func NewEnumTypeRecord(wasmEnumName, protoPackage, goPackage string, wasmEnumType *WasmEnumType) *EnumTypeRecord {
	return &EnumTypeRecord{
		protoPackage: protoPackage,
		goPackage:    goPackage,
		wasmName:     wasmEnumName,
	}
}

func NewEnumValueRecord(wasmEnumValue string, packageName string, goPackage string, parent *EnumTypeRecord, start, end int) *EnumValueRecord {
	v := &EnumValueRecord{
		parent:       parent,
		wasmName:     wasmEnumValue,
		protoPackage: packageName,
		goPackage:    goPackage,
		range_end:    end,
		range_start:  start,
	}
	return v
}

func (m *MessageRecord) String() string {
	return fmt.Sprintf("MessageRec(%s,%s,%s)", m.WasmName(), m.protoPackage, m.goPackage)
}
func (m *ServiceRecord) String() string {
	return fmt.Sprintf("ServiceRec(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (m *EnumTypeRecord) String() string {
	return fmt.Sprintf("EnumType(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (m *EnumValueRecord) String() string {
	return fmt.Sprintf("EnumValue(%s,%s,%s)", m.wasmName, m.protoPackage, m.goPackage)
}
func (g *GenInfo) RegisterService(w *WasmService) {
	g.finder.AddServiceType(w.GetWasmServiceName(), w.ProtoPackage(), w.GetGoPackage(), w)
}

func (g *GenInfo) RegisterMessage(w *WasmMessage) {
	g.finder.AddMessageType(w.GetWasmMessageName(), w.GetProtoPackage(), w.GetGoPackage(), w)
}
func (g *GenInfo) GetAllServiceByName(generatedFile string) []*descriptorpb.ServiceDescriptorProto {
	return g.nameToSvc[generatedFile]
}
func (g *GenInfo) GetAllMessageByName(generatedFile string) []*descriptorpb.DescriptorProto {
	return g.nameToMsg[generatedFile]
}
func (g *GenInfo) GetAllEnumByName(generatedFile string) []*descriptorpb.EnumDescriptorProto {
	return g.nameToEnumType[generatedFile]
}

func (g *GenInfo) GoPackageOption(service []*WasmService, message []*WasmMessage) (string, error) {
	return g.finder.GoPackageOption(service, message)
}

func (g *GenInfo) FindServiceByName(protoPackage, name string) *WasmService {
	//xxx fixme this stinks
	hackyName := fmt.Sprintf(".%s.%s", protoPackage, name)
	return g.finder.FindServiceByName(protoPackage, hackyName)
}

func (g *GenInfo) FindEnumTypeByName(protoPackage, name string) *WasmEnumType {
	//xxx fixme this stinks
	return g.finder.FindEnumTypeByName(protoPackage, name)
}

func (g *GenInfo) FindMessageByName(protoPackage string, name string) *WasmMessage {
	return g.finder.FindMessageByName(protoPackage, name)
}

func (g *GenInfo) AddMessageType(name, protoPackage, goPackage string, message *WasmMessage) {
	g.finder.AddMessageType(name, protoPackage, goPackage, message)

}
func (g *GenInfo) AddServiceType(wasmName, protoPackage, goPackage string, service *WasmService) {
	g.finder.AddServiceType(wasmName, protoPackage, goPackage, service)
}
func (g *GenInfo) AddressingNameFromMessage(currentPkg string, message *WasmMessage) string {
	return g.finder.AddressingNameFromMessage(currentPkg, message)
}

func (g *GenInfo) SetReqAndFileMappings(request *pluginpb.CodeGeneratorRequest,
	n map[string]*descriptorpb.FileDescriptorProto,
	s map[string][]*descriptorpb.ServiceDescriptorProto,
	m map[string][]*descriptorpb.DescriptorProto,
	et map[string][]*descriptorpb.EnumDescriptorProto,
	ev map[string][]*descriptorpb.EnumValueDescriptorProto) {
	g.request = request
	g.nameToFile = n
	g.nameToSvc = s
	g.nameToMsg = m
	g.nameToEnumType = et
	g.nameToEnumValue = ev
}

// GetAllFileName returns the list of string (keys in the two maps) that are visible to this genInfo.
func (g *GenInfo) GetAllFileName() []string {
	result := []string{}
	for k := range g.nameToFile {
		result = append(result, k)
	}
	return result
}

func (g *GenInfo) GetFileByName(name string) *descriptorpb.FileDescriptorProto {
	return g.nameToFile[name]
}

func (g *GenInfo) Contains(name string) bool {
	_, ok := g.nameToFile[name]
	return ok
}

// enum info
func (g *GenInfo) AddEnumType(name, pbPkg, goPkg string, w *WasmEnumType) {
	g.finder.AddEnumType(name, pbPkg, goPkg, w)
}

func (g *GenInfo) AddEnumValue(name, pbPkg, goPkg string, parent *EnumTypeRecord, s, e int, w *WasmEnumValue) {
	g.finder.AddEnumValue(name, pbPkg, goPkg,
		parent, s, e, w)
}
func (g *GenInfo) Enum() []*WasmEnumType {
	return g.finder.Enum()
}
