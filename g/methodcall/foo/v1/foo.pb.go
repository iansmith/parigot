// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: methodcall/foo/v1/foo.proto

package foo

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

type FooErr int32

const (
	FooErr_NoError         FooErr = 0 // mandatory
	FooErr_DispatchError   FooErr = 1 //mandatory
	FooErr_UnmarshalError  FooErr = 2 // mandatory
	FooErr_BadParamWritePi FooErr = 3
)

// Enum value maps for FooErr.
var (
	FooErr_name = map[int32]string{
		0: "NoError",
		1: "DispatchError",
		2: "UnmarshalError",
		3: "BadParamWritePi",
	}
	FooErr_value = map[string]int32{
		"NoError":         0,
		"DispatchError":   1,
		"UnmarshalError":  2,
		"BadParamWritePi": 3,
	}
)

func (x FooErr) Enum() *FooErr {
	p := new(FooErr)
	*p = x
	return p
}

func (x FooErr) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FooErr) Descriptor() protoreflect.EnumDescriptor {
	return file_methodcall_foo_v1_foo_proto_enumTypes[0].Descriptor()
}

func (FooErr) Type() protoreflect.EnumType {
	return &file_methodcall_foo_v1_foo_proto_enumTypes[0]
}

func (x FooErr) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FooErr.Descriptor instead.
func (FooErr) EnumDescriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{0}
}

// add or multiplication, has both request and response
type AddMultiplyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value0 int32 `protobuf:"varint,1,opt,name=value0,proto3" json:"value0,omitempty"`
	Value1 int32 `protobuf:"varint,2,opt,name=value1,proto3" json:"value1,omitempty"`
	IsAdd  bool  `protobuf:"varint,3,opt,name=is_add,json=isAdd,proto3" json:"is_add,omitempty"`
}

func (x *AddMultiplyRequest) Reset() {
	*x = AddMultiplyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMultiplyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMultiplyRequest) ProtoMessage() {}

func (x *AddMultiplyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMultiplyRequest.ProtoReflect.Descriptor instead.
func (*AddMultiplyRequest) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{0}
}

func (x *AddMultiplyRequest) GetValue0() int32 {
	if x != nil {
		return x.Value0
	}
	return 0
}

func (x *AddMultiplyRequest) GetValue1() int32 {
	if x != nil {
		return x.Value1
	}
	return 0
}

func (x *AddMultiplyRequest) GetIsAdd() bool {
	if x != nil {
		return x.IsAdd
	}
	return false
}

type AddMultiplyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result int32 `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *AddMultiplyResponse) Reset() {
	*x = AddMultiplyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMultiplyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMultiplyResponse) ProtoMessage() {}

func (x *AddMultiplyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMultiplyResponse.ProtoReflect.Descriptor instead.
func (*AddMultiplyResponse) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{1}
}

func (x *AddMultiplyResponse) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

// no input params
type LucasSequenceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LucasSequenceRequest) Reset() {
	*x = LucasSequenceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LucasSequenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LucasSequenceRequest) ProtoMessage() {}

func (x *LucasSequenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LucasSequenceRequest.ProtoReflect.Descriptor instead.
func (*LucasSequenceRequest) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{2}
}

type LucasSequenceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence []int32 `protobuf:"varint,1,rep,packed,name=sequence,proto3" json:"sequence,omitempty"`
}

func (x *LucasSequenceResponse) Reset() {
	*x = LucasSequenceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LucasSequenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LucasSequenceResponse) ProtoMessage() {}

func (x *LucasSequenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LucasSequenceResponse.ProtoReflect.Descriptor instead.
func (*LucasSequenceResponse) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{3}
}

func (x *LucasSequenceResponse) GetSequence() []int32 {
	if x != nil {
		return x.Sequence
	}
	return nil
}

// no output params
type WritePiRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Terms int32 `protobuf:"varint,1,opt,name=terms,proto3" json:"terms,omitempty"`
}

func (x *WritePiRequest) Reset() {
	*x = WritePiRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WritePiRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WritePiRequest) ProtoMessage() {}

func (x *WritePiRequest) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WritePiRequest.ProtoReflect.Descriptor instead.
func (*WritePiRequest) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{4}
}

func (x *WritePiRequest) GetTerms() int32 {
	if x != nil {
		return x.Terms
	}
	return 0
}

type WritePiResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WritePiResponse) Reset() {
	*x = WritePiResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_methodcall_foo_v1_foo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WritePiResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WritePiResponse) ProtoMessage() {}

func (x *WritePiResponse) ProtoReflect() protoreflect.Message {
	mi := &file_methodcall_foo_v1_foo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WritePiResponse.ProtoReflect.Descriptor instead.
func (*WritePiResponse) Descriptor() ([]byte, []int) {
	return file_methodcall_foo_v1_foo_proto_rawDescGZIP(), []int{5}
}

var File_methodcall_foo_v1_foo_proto protoreflect.FileDescriptor

var file_methodcall_foo_v1_foo_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2f, 0x66, 0x6f, 0x6f,
	0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6f, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31,
	0x1a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x76,
	0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69,
	0x70, 0x6c, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x30, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x30, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0x12, 0x15, 0x0a, 0x06, 0x69, 0x73,
	0x5f, 0x61, 0x64, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x73, 0x41, 0x64,
	0x64, 0x22, 0x2d, 0x0a, 0x13, 0x41, 0x64, 0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x16, 0x0a, 0x14, 0x4c, 0x75, 0x63, 0x61, 0x73, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x33, 0x0a, 0x15, 0x4c, 0x75, 0x63, 0x61,
	0x73, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x26, 0x0a,
	0x0e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x50, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x65, 0x72, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x74, 0x65, 0x72, 0x6d, 0x73, 0x22, 0x11, 0x0a, 0x0f, 0x57, 0x72, 0x69, 0x74, 0x65, 0x50, 0x69,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x58, 0x0a, 0x06, 0x46, 0x6f, 0x6f, 0x45,
	0x72, 0x72, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x6f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x00, 0x12,
	0x11, 0x0a, 0x0d, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x6e, 0x6d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x42, 0x61, 0x64, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x57, 0x72, 0x69, 0x74, 0x65, 0x50, 0x69, 0x10, 0x03, 0x1a, 0x05, 0xd8, 0x9e, 0x89,
	0x02, 0x01, 0x32, 0xa6, 0x02, 0x0a, 0x03, 0x46, 0x6f, 0x6f, 0x12, 0x5c, 0x0a, 0x0b, 0x41, 0x64,
	0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x79, 0x12, 0x25, 0x2e, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64,
	0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x26, 0x2e, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x0d, 0x4c, 0x75, 0x63, 0x61,
	0x73, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x27, 0x2e, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x75,
	0x63, 0x61, 0x73, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x28, 0x2e, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e,
	0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x75, 0x63, 0x61, 0x73, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x07,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x50, 0x69, 0x12, 0x21, 0x2e, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x72, 0x69, 0x74,
	0x65, 0x50, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2e, 0x66, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x50, 0x69, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x0b,
	0xea, 0x9e, 0x89, 0x02, 0x06, 0x46, 0x6f, 0x6f, 0x45, 0x72, 0x72, 0x42, 0x35, 0x5a, 0x33, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x6e, 0x73, 0x6d, 0x69,
	0x74, 0x68, 0x2f, 0x70, 0x61, 0x72, 0x69, 0x67, 0x6f, 0x74, 0x2f, 0x67, 0x2f, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x63, 0x61, 0x6c, 0x6c, 0x2f, 0x66, 0x6f, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x66,
	0x6f, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_methodcall_foo_v1_foo_proto_rawDescOnce sync.Once
	file_methodcall_foo_v1_foo_proto_rawDescData = file_methodcall_foo_v1_foo_proto_rawDesc
)

func file_methodcall_foo_v1_foo_proto_rawDescGZIP() []byte {
	file_methodcall_foo_v1_foo_proto_rawDescOnce.Do(func() {
		file_methodcall_foo_v1_foo_proto_rawDescData = protoimpl.X.CompressGZIP(file_methodcall_foo_v1_foo_proto_rawDescData)
	})
	return file_methodcall_foo_v1_foo_proto_rawDescData
}

var file_methodcall_foo_v1_foo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_methodcall_foo_v1_foo_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_methodcall_foo_v1_foo_proto_goTypes = []interface{}{
	(FooErr)(0),                   // 0: methodcall.foo.v1.FooErr
	(*AddMultiplyRequest)(nil),    // 1: methodcall.foo.v1.AddMultiplyRequest
	(*AddMultiplyResponse)(nil),   // 2: methodcall.foo.v1.AddMultiplyResponse
	(*LucasSequenceRequest)(nil),  // 3: methodcall.foo.v1.LucasSequenceRequest
	(*LucasSequenceResponse)(nil), // 4: methodcall.foo.v1.LucasSequenceResponse
	(*WritePiRequest)(nil),        // 5: methodcall.foo.v1.WritePiRequest
	(*WritePiResponse)(nil),       // 6: methodcall.foo.v1.WritePiResponse
}
var file_methodcall_foo_v1_foo_proto_depIdxs = []int32{
	1, // 0: methodcall.foo.v1.Foo.AddMultiply:input_type -> methodcall.foo.v1.AddMultiplyRequest
	3, // 1: methodcall.foo.v1.Foo.LucasSequence:input_type -> methodcall.foo.v1.LucasSequenceRequest
	5, // 2: methodcall.foo.v1.Foo.WritePi:input_type -> methodcall.foo.v1.WritePiRequest
	2, // 3: methodcall.foo.v1.Foo.AddMultiply:output_type -> methodcall.foo.v1.AddMultiplyResponse
	4, // 4: methodcall.foo.v1.Foo.LucasSequence:output_type -> methodcall.foo.v1.LucasSequenceResponse
	6, // 5: methodcall.foo.v1.Foo.WritePi:output_type -> methodcall.foo.v1.WritePiResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_methodcall_foo_v1_foo_proto_init() }
func file_methodcall_foo_v1_foo_proto_init() {
	if File_methodcall_foo_v1_foo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_methodcall_foo_v1_foo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddMultiplyRequest); i {
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
		file_methodcall_foo_v1_foo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddMultiplyResponse); i {
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
		file_methodcall_foo_v1_foo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LucasSequenceRequest); i {
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
		file_methodcall_foo_v1_foo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LucasSequenceResponse); i {
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
		file_methodcall_foo_v1_foo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WritePiRequest); i {
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
		file_methodcall_foo_v1_foo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WritePiResponse); i {
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
			RawDescriptor: file_methodcall_foo_v1_foo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_methodcall_foo_v1_foo_proto_goTypes,
		DependencyIndexes: file_methodcall_foo_v1_foo_proto_depIdxs,
		EnumInfos:         file_methodcall_foo_v1_foo_proto_enumTypes,
		MessageInfos:      file_methodcall_foo_v1_foo_proto_msgTypes,
	}.Build()
	File_methodcall_foo_v1_foo_proto = out.File
	file_methodcall_foo_v1_foo_proto_rawDesc = nil
	file_methodcall_foo_v1_foo_proto_goTypes = nil
	file_methodcall_foo_v1_foo_proto_depIdxs = nil
}
