package pbmodel

import (
	"log"
	"strings"

	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/ui/parser/tree"
)

type Pb3Builder struct {
	*Baseprotobuf3Listener
	CurrentFile      string
	CurrentPkgPrefix string
	OutgoingImport   []string
	failure          bool
	Proto2Ignored    []string
	CurrentGoPackage string
	FQNameToPath     map[string]string
	CurrentPackage   string
}

var _ protobuf3Listener = &Pb3Builder{}

func NewPb3Builder() *Pb3Builder {
	return &Pb3Builder{
		FQNameToPath: make(map[string]string),
	}
}

func (p *Pb3Builder) Reset(path string) {
	p.CurrentFile = path
	p.CurrentPackage = ""
	p.CurrentPkgPrefix = ""
	p.failure = false //shouldn't be needed
	p.OutgoingImport = nil

}

func (p *Pb3Builder) ExitImportStatement(ctx *ImportStatementContext) {
	import_ := ctx.StrLit().GetText()
	import_ = strings.TrimSpace(import_)
	import_ = strings.TrimPrefix(import_, "\"")
	import_ = strings.TrimSuffix(import_, "\"")
	if strings.HasPrefix(import_, "google/protobuf") {
		log.Printf("ignoring protobuf version 2 file from google: %s", import_)
		p.Proto2Ignored = append(p.Proto2Ignored, import_)
		return
	}
	p.OutgoingImport = append(p.OutgoingImport, import_)
	p.FQNameToPath[import_] = p.CurrentFile
	pb3Import.AddEdge(p.CurrentFile, import_)
}

func (p *Pb3Builder) ExitPackageStatement(ctx *PackageStatementContext) {
	pkg := ctx.FullIdent().GetText()
	p.CurrentPackage = pkg
	p.CurrentPkgPrefix = p.StripEnds(pkg, p.CurrentFile)
	log.Printf("inside source %s, set current package to %s (with package prefix %s)", p.CurrentFile, p.CurrentPackage, p.CurrentPkgPrefix)

	ctx.SetPkg(ctx.FullIdent().GetText())
}
func (b *Pb3Builder) Failed() bool {
	return b.failure
}

func (b *Pb3Builder) StripEnds(pkg, path string) string {
	return helper.StripEndsOfPathForPkg(pkg, path)
}

func (p *Pb3Builder) ExitOptionStatement(ctx *OptionStatementContext) {
	name, value := ctx.OptionName().GetText(), ctx.Constant().GetText()
	if name == "go_package" {
		p.CurrentGoPackage = value
	}
}

func (p *Pb3Builder) ExitProto(ctx *ProtoContext) {
	pbNode := tree.NewProtobufFileNode()
	//get imports
	raw := ctx.AllImportStatement()
	impFile := make([]string, len(raw))
	for i, import_ := range raw {
		impFile[i] = import_.GetImp()
	}
	pbNode.ImportFile = impFile
	pbNode.FileName = p.CurrentFile

	// package?
	rawPkg := ctx.AllPackageStatement()
	if len(rawPkg) > 1 {
		// what would this even MEAN?
		panic("unable to handle multiple package statements in a .proto file")
	}
	if len(rawPkg) == 1 {
		log.Printf("got the package name %s", rawPkg[0].GetPkg())
		pbNode.PackageName = rawPkg[0].GetPkg()
	}

	// message
	rawMsg := ctx.AllTopLevelDef()
	msg := make([]*tree.ProtobufMessage, len(rawMsg))
	for i, m := range rawMsg {
		msg[i] = m.GetMsg()
	}
	pbNode.Message = msg
}

func (p *Pb3Builder) ExitTopLevelDef(ctx *TopLevelDefContext) {
	ctx.SetMsg(ctx.MessageDef().GetMsg())
}

func (p *Pb3Builder) ExitMessageDef(ctx *MessageDefContext) {
	msg := tree.NewProtobufMessage(ctx.MessageName().GetText())
	log.Printf("create protobuf message with %s", msg.Name)
	ctx.SetMsg(msg)
}
