package codegen

import (
	"log"

	"google.golang.org/protobuf/types/descriptorpb"
)

var parigotTypeList = []string{"TYPE_STRING", "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE", "TYPE_BOOL",
	"TYPE_BYTES", "TYPE_BYTE"}

func NewWasmMethod(desc *descriptorpb.MethodDescriptorProto, w *WasmService) *WasmMethod {
	meth := &WasmMethod{
		MethodDescriptorProto: desc,
		wasmMethodName:        "", //computed later
		parent:                w,
	}
	if desc.GetName() == "" || desc.GetInputType() == "" || desc.GetOutputType() == "" {
		log.Fatalf("method data missing: name='%s', inputType='%s', outputType='%s'",
			desc.GetName(), desc.GetInputType(), desc.GetOutputType())
	}
	if meth.GetOptions() != nil {
		meth.pullParameters = pullParamsOption(meth.GetOptions().String())
		meth.pullOutput = pullOutputOption(meth.GetOptions().String())
		meth.isTest = hasMethodTestOption(meth.GetOptions().String())
	}
	return meth
}
func NewWasmMessage(file *descriptorpb.FileDescriptorProto, message *descriptorpb.DescriptorProto,
	lang LanguageText, finder Finder) *WasmMessage {
	m := &WasmMessage{
		DescriptorProto: message,
		wasmMessageName: "",
		parent:          file,
		field:           []*WasmField{},
		lang:            lang,
		finder:          finder,
	}
	if message.GetOptions() != nil {
		m.noPackage = hasNoPackageOption(message.GetOptions().String())
	}
	return m
}

func NewWasmService(file *descriptorpb.FileDescriptorProto,
	service *descriptorpb.ServiceDescriptorProto, lang LanguageText, finder Finder) *WasmService {
	s := &WasmService{
		ServiceDescriptorProto: service,
		wasmServiceName:        "",
		parent:                 file,
		method:                 []*WasmMethod{},
		lang:                   lang,
		finder:                 finder,
	}
	if service.GetOptions() != nil {
		s.alwaysPullParameters = alwaysPullParamsOption(service.GetOptions().String())
		s.alwaysPullOutput = alwaysPullOutputOption(service.GetOptions().String())
		s.noPackage = hasNoPackageOption(service.GetOptions().String())
		s.kernel = hasKernelOption(service.GetOptions().String())
		s.isTest = hasServiceTestOption(service.GetOptions().String())
	}
	return s
}

func NewWasmField(proto *descriptorpb.FieldDescriptorProto, w *WasmMessage) *WasmField {
	result := &WasmField{
		FieldDescriptorProto: proto,
		parent:               w,
	}

	return result
}
