// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: log/log.proto

package log

import (
	log "github.com/iansmith/parigot/api/proto/g/pb/log"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_log_log_proto protoreflect.FileDescriptor

var file_log_log_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x6c, 0x6f, 0x67, 0x1a, 0x10, 0x70, 0x62, 0x2f, 0x6c, 0x6f, 0x67, 0x2f, 0x6c, 0x6f, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x33, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x2c, 0x0a,
	0x03, 0x4c, 0x6f, 0x67, 0x12, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x6f,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x6c, 0x6f,
	0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x31, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x61, 0x6e, 0x73, 0x6d, 0x69,
	0x74, 0x68, 0x2f, 0x70, 0x61, 0x72, 0x69, 0x67, 0x6f, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x3b, 0x6c, 0x6f, 0x67, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_log_log_proto_goTypes = []interface{}{
	(*log.LogRequest)(nil), // 0: pb.log.LogRequest
	(*log.LogResult)(nil),  // 1: pb.log.LogResult
}
var file_log_log_proto_depIdxs = []int32{
	0, // 0: log.Log.Log:input_type -> pb.log.LogRequest
	1, // 1: log.Log.Log:output_type -> pb.log.LogResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_log_log_proto_init() }
func file_log_log_proto_init() {
	if File_log_log_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_log_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_log_log_proto_goTypes,
		DependencyIndexes: file_log_log_proto_depIdxs,
	}.Build()
	File_log_log_proto = out.File
	file_log_log_proto_rawDesc = nil
	file_log_log_proto_goTypes = nil
	file_log_log_proto_depIdxs = nil
}