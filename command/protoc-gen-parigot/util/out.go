package util

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"log"
	"path/filepath"
	"strings"
)

const protoSuffix = ".proto"

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
func (o *OutputFile) Write(p []byte) (n int, err error) {
	return o.buf.Write(p)
}
func (o *OutputFile) Close() error {
	return nil
}

func GenerateOutputFilenameBase(proto *descriptorpb.FileDescriptorProto) string {
	descName := proto.GetName()
	descBase, descFile := filepath.Split(descName)
	n := descFile
	if !strings.HasSuffix(n, protoSuffix) {
		log.Printf("unexpeced filename for processing '%s', expected a %s extension",
			descName, protoSuffix)
	} else {
		n = strings.TrimSuffix(n, protoSuffix)
	}
	return filepath.Join(descBase, n)
}
