package util

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type OutputFile struct {
	name string
	buf  bytes.Buffer
}

func NewOutputFile(name string) *OutputFile {
	return &OutputFile{
		name: name,
	}
}

func (o *OutputFile) Printf(spec string, args ...interface{}) {
	s := fmt.Sprintf(spec, args)
	o.buf.WriteString(s)
}

func (o *OutputFile) ToGoogleCGResponseFile() *pluginpb.CodeGeneratorResponse_File {
	content := o.buf.String()
	return &pluginpb.CodeGeneratorResponse_File{Name: &o.name, Content: &content}
}

func GenerateOutputFilenameBase(pb *descriptorpb.FileDescriptorProto) string {
	return pb.GetPackage() + "///" + pb.GetName()
}
