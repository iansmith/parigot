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
	packageName := proto.GetPackage()
	descName := proto.GetName()
	descBase, descFile := filepath.Split(descName)
	parts := strings.Split(descBase, "/")
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	pkgParts := strings.Split(packageName, ".")
	///////////
	for {
		if len(parts) == 0 || len(pkgParts) == 0 {
			break
		}
		last := parts[len(parts)-1]
		lastPkg := pkgParts[len(pkgParts)-1]
		if last == lastPkg {
			parts = parts[:len(parts)-1]
			pkgParts = pkgParts[:len(pkgParts)-1]
		}
	}
	descBase = filepath.Join(parts...)
	descName = filepath.Join(descBase, descFile)
	n := descName
	if !strings.HasSuffix(descName, protoSuffix) {
		log.Printf("unexpeced filename for processing '%s', expected a %s extension",
			descName, protoSuffix)
	} else {
		n = strings.TrimSuffix(n, protoSuffix)
	}
	return strings.Replace(packageName, ".", "/", -1) + "/" + n
}
