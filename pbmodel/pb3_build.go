package pbmodel

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
}
func (b *Pb3Builder) Failed() bool {
	return b.failure
}

func (b *Pb3Builder) StripEnds(pkg, path string) string {
	pkgPart := strings.Split(pkg, ".")
	dir := filepath.Dir(path)
	dir = filepath.Clean(dir)
	elem := strings.Split(dir, string(os.PathSeparator))
	if elem == nil {
		panic("empty path given to strip ends")
	}
	for len(pkgPart) > 0 {
		lastPart := pkgPart[len(pkgPart)-1]
		lastElem := elem[len(elem)-1]
		if lastElem != lastPart {
			break
		}
		pkgPart = pkgPart[:len(pkgPart)-1]
		elem = elem[:len(elem)-1]
	}
	if len(pkgPart) != 0 {
		panic(fmt.Sprintf("had package parts left: %s", filepath.Join(pkgPart...)))
	}
	return filepath.Join(elem...)
}

func (p *Pb3Builder) ExitOptionStatement(ctx *OptionStatementContext) {
	name, value := ctx.OptionName().GetText(), ctx.Constant().GetText()
	if name == "go_package" {
		p.CurrentGoPackage = value
	}
}
