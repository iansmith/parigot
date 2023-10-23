// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: httpconnector/v1/httpconnector.proto

package httpconnector

import (
	v1 "github.com/iansmith/parigot/g/protosupport/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
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
	HttpConnectorErr_NoReceiver     HttpConnectorErr = 5 // We could not find any service that implements HttpConnector
	HttpConnectorErr_ReceiverFailed HttpConnectorErr = 6 // We attempted to call the receiver, but got a failure back
)

// Enum value maps for HttpConnectorErr.
var (
	HttpConnectorErr_name = map[int32]string{
		0: "NoError",
		1: "DispatchError",
		2: "UnmarshalError",
		3: "MarshalError",
		4: "InternalError",
		5: "NoReceiver",
		6: "ReceiverFailed",
	}
	HttpConnectorErr_value = map[string]int32{
		"NoError":        0,
		"DispatchError":  1,
		"UnmarshalError": 2,
		"MarshalError":   3,
		"InternalError":  4,
		"NoReceiver":     5,
		"ReceiverFailed": 6,
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

type HandleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HttpMethod string     `protobuf:"bytes,1,opt,name=http_method,json=httpMethod,proto3" json:"http_method,omitempty"`
	Url        string     `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	ServiceId  *v1.IdRaw  `protobuf:"bytes,3,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	MethodId   *v1.IdRaw  `protobuf:"bytes,4,opt,name=method_id,json=methodId,proto3" json:"method_id,omitempty"`
	ReqAny     *anypb.Any `protobuf:"bytes,5,opt,name=req_any,json=reqAny,proto3" json:"req_any,omitempty"`
}

func (x *HandleRequest) Reset() {
	*x = HandleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleRequest) ProtoMessage() {}

func (x *HandleRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use HandleRequest.ProtoReflect.Descriptor instead.
func (*HandleRequest) Descriptor() ([]byte, []int) {
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{0}
}

func (x *HandleRequest) GetHttpMethod() string {
	if x != nil {
		return x.HttpMethod
	}
	return ""
}

func (x *HandleRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *HandleRequest) GetServiceId() *v1.IdRaw {
	if x != nil {
		return x.ServiceId
	}
	return nil
}

func (x *HandleRequest) GetMethodId() *v1.IdRaw {
	if x != nil {
		return x.MethodId
	}
	return nil
}

func (x *HandleRequest) GetReqAny() *anypb.Any {
	if x != nil {
		return x.ReqAny
	}
	return nil
}

type HandleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HttpStatus   int32             `protobuf:"varint,1,opt,name=http_status,json=httpStatus,proto3" json:"http_status,omitempty"`
	HttpResponse []byte            `protobuf:"bytes,2,opt,name=http_response,json=httpResponse,proto3" json:"http_response,omitempty"`
	Header       map[string]string `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *HandleResponse) Reset() {
	*x = HandleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleResponse) ProtoMessage() {}

func (x *HandleResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use HandleResponse.ProtoReflect.Descriptor instead.
func (*HandleResponse) Descriptor() ([]byte, []int) {
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{1}
}

func (x *HandleResponse) GetHttpStatus() int32 {
	if x != nil {
		return x.HttpStatus
	}
	return 0
}

func (x *HandleResponse) GetHttpResponse() []byte {
	if x != nil {
		return x.HttpResponse
	}
	return nil
}

func (x *HandleResponse) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[2]
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
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{2}
}

func (x *CheckRequest) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Whether receive a http request from outside
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_httpconnector_v1_httpconnector_proto_msgTypes[3]
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
	return file_httpconnector_v1_httpconnector_proto_rawDescGZIP(), []int{3}
}

func (x *CheckResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_httpconnector_v1_httpconnector_proto protoreflect.FileDescriptor

var file_httpconnector_v1_httpconnector_proto_rawDesc = []byte{
	0x0a, 0x24, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdd, 0x01, 0x0a, 0x0d, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x74, 0x74,
	0x70, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x68, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x35, 0x0a, 0x0a,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x64, 0x52, 0x61, 0x77, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x64, 0x52, 0x61, 0x77, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x5f,
	0x61, 0x6e, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x06, 0x72, 0x65, 0x71, 0x41, 0x6e, 0x79, 0x22, 0xd7, 0x01, 0x0a, 0x0e, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x74,
	0x74, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x68, 0x74, 0x74, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x68,
	0x74, 0x74, 0x70, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x44, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x28, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x29, 0x0a, 0x0d, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2a, 0x96, 0x01, 0x0a, 0x10, 0x48, 0x74, 0x74, 0x70, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x45, 0x72, 0x72, 0x12, 0x0b, 0x0a, 0x07, 0x4e,
	0x6f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x44, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x55,
	0x6e, 0x6d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x02, 0x12,
	0x10, 0x0a, 0x0c, 0x4d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10,
	0x03, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x6f, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x06, 0x1a, 0x05, 0xd8, 0x9e, 0x89, 0x02, 0x01, 0x32,
	0x5c, 0x0a, 0x0d, 0x48, 0x74, 0x74, 0x70, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x12, 0x4b, 0x0a, 0x06, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x1f, 0x2e, 0x68, 0x74, 0x74,
	0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x68, 0x74,
	0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3e, 0x5a,
	0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x6e, 0x73,
	0x6d, 0x69, 0x74, 0x68, 0x2f, 0x70, 0x61, 0x72, 0x69, 0x67, 0x6f, 0x74, 0x2f, 0x67, 0x2f, 0x68,
	0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x3b,
	0x68, 0x74, 0x74, 0x70, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_httpconnector_v1_httpconnector_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_httpconnector_v1_httpconnector_proto_goTypes = []interface{}{
	(HttpConnectorErr)(0),  // 0: httpconnector.v1.HttpConnectorErr
	(*HandleRequest)(nil),  // 1: httpconnector.v1.HandleRequest
	(*HandleResponse)(nil), // 2: httpconnector.v1.HandleResponse
	(*CheckRequest)(nil),   // 3: httpconnector.v1.CheckRequest
	(*CheckResponse)(nil),  // 4: httpconnector.v1.CheckResponse
	nil,                    // 5: httpconnector.v1.HandleResponse.HeaderEntry
	(*v1.IdRaw)(nil),       // 6: protosupport.v1.IdRaw
	(*anypb.Any)(nil),      // 7: google.protobuf.Any
}
var file_httpconnector_v1_httpconnector_proto_depIdxs = []int32{
	6, // 0: httpconnector.v1.HandleRequest.service_id:type_name -> protosupport.v1.IdRaw
	6, // 1: httpconnector.v1.HandleRequest.method_id:type_name -> protosupport.v1.IdRaw
	7, // 2: httpconnector.v1.HandleRequest.req_any:type_name -> google.protobuf.Any
	5, // 3: httpconnector.v1.HandleResponse.header:type_name -> httpconnector.v1.HandleResponse.HeaderEntry
	1, // 4: httpconnector.v1.HttpConnector.Handle:input_type -> httpconnector.v1.HandleRequest
	2, // 5: httpconnector.v1.HttpConnector.Handle:output_type -> httpconnector.v1.HandleResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_httpconnector_v1_httpconnector_proto_init() }
func file_httpconnector_v1_httpconnector_proto_init() {
	if File_httpconnector_v1_httpconnector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_httpconnector_v1_httpconnector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleRequest); i {
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
			switch v := v.(*HandleResponse); i {
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
		file_httpconnector_v1_httpconnector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_httpconnector_v1_httpconnector_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
			NumMessages:   5,
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
