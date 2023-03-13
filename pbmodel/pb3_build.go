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
	currentFile      string
	OutgoingImport   []string
	failure          bool
	currentPkgPrefix string
	Proto2Ignored    []string
}

var _ protobuf3Listener = &Pb3Builder{}

func NewPb3Builder() *Pb3Builder {
	return &Pb3Builder{}
}

func (p *Pb3Builder) Reset(path string) {
	p.currentFile = path
	p.currentPkgPrefix = ""
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
	pb3Import.AddEdge(p.currentFile, import_)
}

func (p *Pb3Builder) ExitPackageStatement(ctx *PackageStatementContext) {
	p.currentPkgPrefix = p.StripEnds(ctx.FullIdent().GetText(), p.currentFile)
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
