package pbmodel

import (
	"fmt"
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

var ignoredCache = make(map[string]struct{})

func (p *Pb3Builder) ExitImportStatement(ctx *ImportStatementContext) {
	import_ := ctx.StrLit().GetText()
	import_ = strings.TrimSpace(import_)
	import_ = strings.TrimPrefix(import_, "\"")
	import_ = strings.TrimSuffix(import_, "\"")
	if strings.HasPrefix(import_, "google/protobuf") {
		_, ok := ignoredCache[import_]
		if !ok {
			log.Printf("ignoring protobuf version 2 file from google: %s", import_)
			ignoredCache[import_] = struct{}{}
		}
		p.Proto2Ignored = append(p.Proto2Ignored, import_)
		return
	}
	//p.OutgoingImport = append(p.OutgoingImport, import_)
	p.FQNameToPath[import_] = p.CurrentFile
	//Pb3Dep.AddEdge(p.CurrentFile, import_)
	ctx.SetImp(import_)
}

func (p *Pb3Builder) ExitPackageStatement(ctx *PackageStatementContext) {
	pkg := ctx.FullIdent().GetText()
	p.CurrentPackage = pkg
	p.CurrentPkgPrefix = p.StripEnds(pkg, p.CurrentFile)

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
	// only one we care about now
	if name == "go_package" {
		triple := tree.NewOptionTriple()
		triple.Name = ctx.OptionName().GetName()
		part := strings.Split(value, ";")
		if len(part) == 1 {
			triple.Value = value
			part := strings.Split(value, "/")
			log.Printf("WARNING: '%s' does not have extra portion of the go_package name denoted by a semicolon, will assume addressing name is %s", value, part[len(part)-1])
			ctx.SetTriple(triple)
			return
		}
		if len(part) > 2 {
			panic(fmt.Sprintf("unable to understand option '%s' with value '%s'", name, value))
		}
		part0 := strings.TrimPrefix(part[0], "\"")
		part0 = strings.TrimSuffix(part0, "\"")
		part1 := strings.TrimPrefix(part[1], "\"")
		part1 = strings.TrimSuffix(part1, "\"")
		triple.Value = part0
		triple.Extra = part1
		ctx.SetTriple(triple)
	}

}

func (p *Pb3Builder) ExitOptionName(ctx *OptionNameContext) {
	if ctx.GetSimple() == nil {
		return
	}
	if ctx.GetSimple().GetText() == "" {
		return
	}
	ctx.SetName(ctx.GetSimple().GetFullId())
}

func (p *Pb3Builder) ExitProto(ctx *ProtoContext) {
	pbNode := tree.NewProtobufFileNode()
	//get imports, but it can have "" in the sequence because of ignored proto2 files
	raw := ctx.AllImportStatement()
	impFile := []string{}
	for _, import_ := range raw {
		if import_.GetImp() != "" {
			impFile = append(impFile, import_.GetImp())
		}
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
		pbNode.PackageName = rawPkg[0].GetPkg()
	}

	// message
	rawMsg := ctx.AllTopLevelDef()
	msg := make([]*tree.ProtobufMessage, len(rawMsg))
	for i, m := range rawMsg {
		msg[i] = m.GetMsg()
	}
	pbNode.Message = msg

	stmt := ctx.AllOptionStatement()
	if len(stmt) > 1 {
		panic(fmt.Sprintf("unable to understand a proto file (%s) that has more than 1 go_package option!", p.CurrentFile))
	}
	if len(stmt) == 1 {
		pbNode.GoPkg = stmt[0].GetTriple().Value
		pbNode.LocalGoPkg = stmt[0].GetTriple().Extra
	}
	ctx.SetPbNode(pbNode)
}

func (p *Pb3Builder) ExitTopLevelDef(ctx *TopLevelDefContext) {
	if ctx.MessageDef() != nil {
		ctx.SetMsg(ctx.MessageDef().GetMsg())
	}
}

func (p *Pb3Builder) ExitMessageDef(ctx *MessageDefContext) {
	msg := tree.NewProtobufMessage(ctx.MessageName().GetText())
	ctx.SetMsg(msg)
}

func (p *Pb3Builder) ExitFullIdent(ctx *FullIdentContext) {
	raw := ctx.AllIdent()
	id := make([]string, len(raw))
	for i, r := range raw {
		id[i] = r.GetText()
	}
	ctx.SetFullId(tree.NewFullIdent(id))
}
