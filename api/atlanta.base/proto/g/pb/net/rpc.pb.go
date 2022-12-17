// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: pb/net/rpc.proto

package net

import (
	protosupport "github.com/iansmith/parigot/api/proto/g/pb/protosupport"
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

type RPCRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pctx      []byte                  `protobuf:"bytes,1,opt,name=pctx,proto3" json:"pctx,omitempty"`
	Param     []byte                  `protobuf:"bytes,2,opt,name=param,proto3" json:"param,omitempty"`
	ServiceId *protosupport.ServiceId `protobuf:"bytes,3,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	CallId    *protosupport.CallId    `protobuf:"bytes,4,opt,name=call_id,json=callId,proto3" json:"call_id,omitempty"`
	// if you supply method_id, it takes precedence over the string
	MethodId   *protosupport.MethodId `protobuf:"bytes,5,opt,name=method_id,json=methodId,proto3" json:"method_id,omitempty"`
	MethodName string                 `protobuf:"bytes,6,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	Sender     string                 `protobuf:"bytes,7,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *RPCRequest) Reset() {
	*x = RPCRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_net_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCRequest) ProtoMessage() {}

func (x *RPCRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_net_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCRequest.ProtoReflect.Descriptor instead.
func (*RPCRequest) Descriptor() ([]byte, []int) {
	return file_pb_net_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *RPCRequest) GetPctx() []byte {
	if x != nil {
		return x.Pctx
	}
	return nil
}

func (x *RPCRequest) GetParam() []byte {
	if x != nil {
		return x.Param
	}
	return nil
}

func (x *RPCRequest) GetServiceId() *protosupport.ServiceId {
	if x != nil {
		return x.ServiceId
	}
	return nil
}

func (x *RPCRequest) GetCallId() *protosupport.CallId {
	if x != nil {
		return x.CallId
	}
	return nil
}

func (x *RPCRequest) GetMethodId() *protosupport.MethodId {
	if x != nil {
		return x.MethodId
	}
	return nil
}

func (x *RPCRequest) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *RPCRequest) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

type RPCResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pctx     []byte                      `protobuf:"bytes,1,opt,name=pctx,proto3" json:"pctx,omitempty"`
	Result   []byte                      `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	CallId   *protosupport.CallId        `protobuf:"bytes,3,opt,name=call_id,json=callId,proto3" json:"call_id,omitempty"`
	MethodId *protosupport.MethodId      `protobuf:"bytes,4,opt,name=method_id,json=methodId,proto3" json:"method_id,omitempty"` // so you can cache it for later
	KerrId   *protosupport.KernelErrorId `protobuf:"bytes,5,opt,name=kerr_id,json=kerrId,proto3" json:"kerr_id,omitempty"`
}

func (x *RPCResponse) Reset() {
	*x = RPCResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_net_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCResponse) ProtoMessage() {}

func (x *RPCResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_net_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCResponse.ProtoReflect.Descriptor instead.
func (*RPCResponse) Descriptor() ([]byte, []int) {
	return file_pb_net_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *RPCResponse) GetPctx() []byte {
	if x != nil {
		return x.Pctx
	}
	return nil
}

func (x *RPCResponse) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *RPCResponse) GetCallId() *protosupport.CallId {
	if x != nil {
		return x.CallId
	}
	return nil
}

func (x *RPCResponse) GetMethodId() *protosupport.MethodId {
	if x != nil {
		return x.MethodId
	}
	return nil
}

func (x *RPCResponse) GetKerrId() *protosupport.KernelErrorId {
	if x != nil {
		return x.KerrId
	}
	return nil
}

var File_pb_net_rpc_proto protoreflect.FileDescriptor

var file_pb_net_rpc_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x62, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x70, 0x62, 0x2e, 0x6e, 0x65, 0x74, 0x1a, 0x22, 0x70, 0x62, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94,
	0x02, 0x0a, 0x0a, 0x52, 0x50, 0x43, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x63, 0x74, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x70, 0x63, 0x74,
	0x78, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x62,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x49, 0x64, 0x52, 0x06, 0x63, 0x61,
	0x6c, 0x6c, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x49, 0x64, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0xdc, 0x01, 0x0a, 0x0b, 0x52, 0x50, 0x43, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x63, 0x74, 0x78, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x70, 0x63, 0x74, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x30, 0x0a, 0x07, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x49, 0x64, 0x52, 0x06, 0x63, 0x61, 0x6c,
	0x6c, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49,
	0x64, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x37, 0x0a, 0x07, 0x6b,
	0x65, 0x72, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x64, 0x52, 0x06, 0x6b, 0x65,
	0x72, 0x72, 0x49, 0x64, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x74, 0x68, 0x2f, 0x70, 0x61, 0x72, 0x69,
	0x67, 0x6f, 0x74, 0x2f, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6e, 0x65, 0x74, 0x3b, 0x6e, 0x65, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_net_rpc_proto_rawDescOnce sync.Once
	file_pb_net_rpc_proto_rawDescData = file_pb_net_rpc_proto_rawDesc
)

func file_pb_net_rpc_proto_rawDescGZIP() []byte {
	file_pb_net_rpc_proto_rawDescOnce.Do(func() {
		file_pb_net_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_net_rpc_proto_rawDescData)
	})
	return file_pb_net_rpc_proto_rawDescData
}

var file_pb_net_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_net_rpc_proto_goTypes = []interface{}{
	(*RPCRequest)(nil),                 // 0: pb.net.RPCRequest
	(*RPCResponse)(nil),                // 1: pb.net.RPCResponse
	(*protosupport.ServiceId)(nil),     // 2: pb.protosupport.ServiceId
	(*protosupport.CallId)(nil),        // 3: pb.protosupport.CallId
	(*protosupport.MethodId)(nil),      // 4: pb.protosupport.MethodId
	(*protosupport.KernelErrorId)(nil), // 5: pb.protosupport.KernelErrorId
}
var file_pb_net_rpc_proto_depIdxs = []int32{
	2, // 0: pb.net.RPCRequest.service_id:type_name -> pb.protosupport.ServiceId
	3, // 1: pb.net.RPCRequest.call_id:type_name -> pb.protosupport.CallId
	4, // 2: pb.net.RPCRequest.method_id:type_name -> pb.protosupport.MethodId
	3, // 3: pb.net.RPCResponse.call_id:type_name -> pb.protosupport.CallId
	4, // 4: pb.net.RPCResponse.method_id:type_name -> pb.protosupport.MethodId
	5, // 5: pb.net.RPCResponse.kerr_id:type_name -> pb.protosupport.KernelErrorId
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pb_net_rpc_proto_init() }
func file_pb_net_rpc_proto_init() {
	if File_pb_net_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_net_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCRequest); i {
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
		file_pb_net_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCResponse); i {
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
			RawDescriptor: file_pb_net_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_net_rpc_proto_goTypes,
		DependencyIndexes: file_pb_net_rpc_proto_depIdxs,
		MessageInfos:      file_pb_net_rpc_proto_msgTypes,
	}.Build()
	File_pb_net_rpc_proto = out.File
	file_pb_net_rpc_proto_rawDesc = nil
	file_pb_net_rpc_proto_goTypes = nil
	file_pb_net_rpc_proto_depIdxs = nil
}