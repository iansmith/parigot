package util

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"google.golang.org/protobuf/types/descriptorpb"
)

func AddFileContentToFinder(f codegen.Finder, pr *descriptorpb.FileDescriptorProto,
	lang codegen.LanguageText) {
	for _, m := range pr.GetMessageType() {
		msg := codegen.NewWasmMessage(pr, m, lang)
		f.AddMessageType(pr.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), msg)
	}
	for _, s := range pr.GetService() {
		svc := codegen.NewWasmService(pr, s, lang)
		f.AddServiceType(pr.GetName(), pr.GetPackage(), pr.GetOptions().GetGoPackage(), svc)
	}
}
