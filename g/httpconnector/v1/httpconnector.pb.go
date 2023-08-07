// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: httpconnector/v1/httpconnector.proto

package httpconnector

import (
	_ "github.com/iansmith/parigot/g/protosupport/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HttpConnectorErr int32

const (
	HttpConnectorErr_NoError        HttpConnectorErr = 0 // mandatory
	HttpConnectorErr_DispatchError  HttpConnectorErr = 1 //mandatory
	HttpConnectorErr_UnmarshalError HttpConnectorErr = 2 // mandatory
	HttpConnectorErr_MarshalError   HttpConnectorErr = 3 // mandatory
	HttpConnectorErr_InternalError  HttpConnectorErr = 4 // There are internal issues with the httpconnector service
)

// Enum value maps for HttpConnectorErr.
var (
	HttpConnectorErr_name = map[int32]string{
		0: "NoError",
		1: "DispatchError",
		2: "UnmarshalError",
		3: "MarshalError",
		4: "InternalError",
	}
	HttpConnectorErr_value = map[string]int32{
		"NoError":        0,
		"DispatchError":  1,
		"UnmarshalError": 2,
		"MarshalError":   3,
		"InternalError":  4,
	}
)

func (x HttpConnectorErr) Enum() *HttpConnectorErr {
	p := new(HttpConnectorErr)
	*p = x
	return p
}

func (x HttpConnectorErr) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HttpConnectorErr) Descriptor() protoreflect.EnumDescriptor {
	return file_httpconnector_v1_httpconnector_proto_enumTypes[0].Descriptor()
}

func (HttpConnectorErr) Type() protoreflect.EnumType {
	return &file_httpconnector_v1_httpconnector_proto_enumTypes[0]
}

func (x HttpConnectorErr) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HttpConnectorErr.Descriptor instead.
func (HttpConnectorErr) EnumDescriptor() ([]byte, []int) {
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{0}
}

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{0}
}

func (x *CheckRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

var File_httpconnector_v1_httpconnector_proto protoreflect.FileDescriptor

var file_httpconnector_v1_httpconnector_proto_rawDesc = []byte{
	0x0a, 0x24, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a, 0x0c,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x22, 0x27, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x2a, 0x72, 0x0a,
	0x10, 0x48, 0x74, 0x74, 0x70, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x45, 0x72,
	0x72, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x6f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x00, 0x12, 0x11,
	0x0a, 0x0d, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10,
	0x01, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x6e, 0x6d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x04, 0x1a, 0x05, 0xd8, 0x9e, 0x89, 0x02,
	0x01, 0x32, 0x59, 0x0a, 0x0d, 0x48, 0x74, 0x74, 0x70, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x12, 0x48, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1e, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3e, 0x5a, 0x3c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x6e, 0x73, 0x6d,
	0x69, 0x74, 0x68, 0x2f, 0x70, 0x61, 0x72, 0x69, 0x67, 0x6f, 0x74, 0x2f, 0x67, 0x2f, 0x68, 0x74,
	0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x68,
	0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_httpconnector_v1_httpconnector_proto_rawDescOnce sync.Once
	file_httpconnector_v1_httpconnector_proto_rawDescData = file_httpconnector_v1_httpconnector_proto_rawDesc
)

func file_httpconnector_v1_httpconnector_proto_rawDescGZIP() []byte {
	file_httpconnector_v1_httpconnector_proto_rawDescOnce.Do(func() {
		file_httpconnector_v1_httpconnector_proto_rawDescData = protoimpl.X.CompressGZIP(file_httpconnector_v1_httpconnector_proto_rawDescData)
	})
	return file_httpconnector_v1_httpconnector_proto_rawDescData
}

var file_httpconnector_v1_httpconnector_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_httpconnector_v1_httpconnector_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_httpconnector_v1_httpconnector_proto_goTypes = []interface{}{
	(HttpConnectorErr)(0), // 0: httpconnector.v1.HttpConnectorErr
	(*CheckRequest)(nil),  // 1: httpconnector.v1.CheckRequest
	(*CheckResponse)(nil), // 2: httpconnector.v1.CheckResponse
}
var file_httpconnector_v1_httpconnector_proto_depIdxs = []int32{
	1, // 0: httpconnector.v1.HttpConnector.Check:input_type -> httpconnector.v1.CheckRequest
	2, // 1: httpconnector.v1.HttpConnector.Check:output_type -> httpconnector.v1.CheckResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_httpconnector_v1_httpconnector_proto_init() }
func file_httpconnector_v1_httpconnector_proto_init() {
	if File_httpconnector_v1_httpconnector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_httpconnector_v1_httpconnector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_httpconnector_v1_httpconnector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_httpconnector_v1_httpconnector_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_httpconnector_v1_httpconnector_proto_goTypes,
		DependencyIndexes: file_httpconnector_v1_httpconnector_proto_depIdxs,
		EnumInfos:         file_httpconnector_v1_httpconnector_proto_enumTypes,
		MessageInfos:      file_httpconnector_v1_httpconnector_proto_msgTypes,
	}.Build()
	File_httpconnector_v1_httpconnector_proto = out.File
	file_httpconnector_v1_httpconnector_proto_rawDesc = nil
	file_httpconnector_v1_httpconnector_proto_goTypes = nil
	file_httpconnector_v1_httpconnector_proto_depIdxs = nil
}
