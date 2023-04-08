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
			//			log.Printf("ignoring protobuf version 2 file from google: %s", import_)
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
	pbNode.Filename = p.CurrentFile
	//log.Printf("finished %s (importFile %+v)", pbNode.Filename, pbNode.ImportFile)
	// package?
	rawPkg := ctx.AllPackageStatement()
	if len(rawPkg) > 1 {
		// what would this even MEAN?
		panic("unable to handle multiple package statements in a .proto file")
	}
	if len(rawPkg) == 1 {
		pbNode.Package = rawPkg[0].GetPkg()
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
	if pbNode == nil {
		panic("nil for pbNode")
	}
	//copy values to message
	for _, m := range pbNode.Message {
		m.Package = pbNode.Package
		m.LocalGoPkg = pbNode.LocalGoPkg
	}
	//log.Printf("finished with pbNode: %+v", pbNode)
	ctx.SetPbNode(pbNode)
}

func (p *Pb3Builder) ExitTopLevelDef(ctx *TopLevelDefContext) {
	if ctx.MessageDef() != nil {
		ctx.SetMsg(ctx.MessageDef().GetMsg())
	}
}

func (p *Pb3Builder) ExitMessageDef(ctx *MessageDefContext) {
	field := make(map[string]*tree.ProtobufMsgElem)
	if ctx.MessageBody() != nil {
		body := ctx.MessageBody().GetBody()
		for _, elem := range body.Elem {
			_, ok := field[elem.Name]
			if ok {
				log.Printf("message definiton '%s' failed because field '%s' is defined more than once", ctx.MessageName().GetText(), elem.Name)
				p.failure = true
				return
			}
			field[elem.Name] = elem
		}
		msg := tree.NewProtobufMessage(ctx.MessageName().GetText(), field)
		ctx.SetMsg(msg)
	} else {
		panic("no message body for message def")
	}
}

func (p *Pb3Builder) ExitMessageBody(ctx *MessageBodyContext) {
	raw := ctx.AllMessageElement()
	e := make([]*tree.ProtobufMsgElem, len(raw))
	for i, m := range raw {
		e[i] = m.GetElem()
	}
	ctx.SetBody(tree.NewProtobufMsgBody(e))
}

func (p *Pb3Builder) ExitField(ctx *FieldContext) {
	var f *tree.ProtobufMsgField
	if ctx.FieldLabel() != nil {
		f = ctx.FieldLabel().GetF()
	} else {
		f = tree.NewProtobufMsgField(false, false)
	}
	name := ctx.FieldName().GetText()
	f.Name = name
	f.Type = ctx.Type_().GetText()
	if isBaseTypeInProtobuf(f.Type) {
		f.TypeBase = true
	}
	ctx.SetF(f)
}

func (p *Pb3Builder) ExitFieldLabel(ctx *FieldLabelContext) {
	opt := false
	rep := false
	if ctx.OPTIONAL() != nil {
		opt = true
	}
	if ctx.REPEATED() != nil {
		rep = true
	}
	ctx.SetF(tree.NewProtobufMsgField(opt, rep))
}

func (p *Pb3Builder) ExitMessageElement(ctx *MessageElementContext) {
	if ctx.Field() == nil && ctx.MapField() == nil {
		panic(fmt.Sprintf("unable to understand anything other than fields as protobuf message components right now: %s", ctx.GetText()))
	}
	var f *tree.ProtobufMsgField
	var m *tree.ProtobufMapField
	if ctx.Field() != nil {
		f = ctx.Field().GetF()
	}
	if ctx.MapField() != nil {
		m = ctx.MapField().GetM()
	}
	//log.Printf("create msg element %v,%v", f == nil, m == nil)
	if f == nil && m == nil {
		log.Printf("not a field, but a %+v", ctx.Field())
		panic("found a field that is not a field or a map")
	}
	ctx.SetElem(tree.NewProtobufMsgElem(f, m))
}

func (p *Pb3Builder) ExitMapField(ctx *MapFieldContext) {
	kt := ctx.KeyType().GetText()
	vt := ctx.Type_().GetText()
	mapName := ctx.MapName().GetText()
	f := tree.NewProtobufMapField(kt, vt, mapName)
	if isBaseTypeInProtobuf(vt) {
		f.ValueTypeBase = true
	}
	ctx.SetM(f)

}

func (p *Pb3Builder) ExitFullIdent(ctx *FullIdentContext) {
	raw := ctx.AllIdent()
	id := make([]string, len(raw))
	for i, r := range raw {
		id[i] = r.GetText()
	}
	ctx.SetFullId(tree.NewFullIdent(id))
}

func isBaseTypeInProtobuf(s string) bool {
	switch s {
	case "double",
		"float",
		"int32",
		"int64",
		"unsigned int32",
		"unsigned int64",
		"signed int32",
		"signed int64",
		"fixed32",
		"fixed64",
		"sfixed32",
		"sfixed64",
		"bool",
		"string",
		"bytes":
		return true
	default:
		return false
	}
}
