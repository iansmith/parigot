package main

import (
	"fmt"
	"github.com/iansmith/parigot/command/transform"
	"log"
	"strings"
)

const lenLimit = 60
const oneK = 1024
const oneM = 1024 * oneK

type unlink struct {
	imports      map[string]int
	newImport    []*transform.ImportDef
	funcTypes    map[int]*transform.FuncSpec
	candidateFT  map[string]int
	candidateSym map[string]string
	packages     map[string]int
	removed      int
	added        int
	remaining    map[string]struct{}
}

var stdlibName = []string{"$runtime.",
	"$fmt.", "$_*fmt",
	"$unicode",
	"$strings.",
	"$_google.golang.org/protobuf", "$_*google.golang.org/protobuf",
	"$google.golang.org/protobuf", "$_struct_google.golang.org/protobuf",
	"$math",
	"$_time", "$time",
	"$strconv", "$_*strconv",
	"$_*reflect", "$reflect", "$_reflect",
	"$bytes", "$_*bytes",
	"$encoding", "$_*encoding",
	"sort",
}

func newUnlink() *unlink {
	result := &unlink{}
	result.imports = make(map[string]int)
	result.newImport = []*transform.ImportDef{}
	result.funcTypes = make(map[int]*transform.FuncSpec)
	result.candidateFT = make(map[string]int)
	result.candidateSym = make(map[string]string)
	result.packages = make(map[string]int)
	result.remaining = make(map[string]struct{})
	return result
}
func (u *unlink) addImportsForCandidates(mod *transform.Module) {
	for cand, funcNum := range u.candidateFT {
		_, ok := u.funcTypes[funcNum]
		if !ok {
			panic("unable to find func type number")
		}
		split := strings.Split(cand, ".")
		if len(split) < 2 {
			panic("unable to understand candidate name")
		}

		impt := &transform.ImportDef{
			ModuleName:  split[0],
			ImportedAs:  strings.Join(split[1:], "."),
			FuncNameRef: &transform.FuncNameRef{Name: new(string), Type: &transform.TypeRef{Num: funcNum}},
		}
		*impt.FuncNameRef.Name = u.candidateSym[cand]
		count, ok := u.packages[split[0]]
		if ok {
			u.packages[split[0]] = count + 1
		} else {
			u.packages[split[0]] = 1
		}
		u.newImport = append(u.newImport, impt)
		u.added++
		mod.AppendTopLevelDef(impt)
	}
}

func (u *unlink) compileFuncTypes(tl transform.TopLevel) {
	tdef := tl.(*transform.TypeDef)
	u.funcTypes[tdef.Annotation] = tdef.Func
}

func (u *unlink) compileImports(tl transform.TopLevel) {
	idef := tl.(*transform.ImportDef)
	for _, pkg := range stdlibName {
		if strings.HasPrefix(*idef.FuncNameRef.Name, pkg) {
			u.imports[*idef.FuncNameRef.Name] = idef.FuncNameRef.Type.Num
		}
	}
}
func (u *unlink) compileCandidates(tl transform.TopLevel) {
	fdef := tl.(*transform.FuncDef)
	for _, pkg := range stdlibName {
		if fdef.Name != nil && strings.HasPrefix(*fdef.Name, pkg) {
			_, ok := u.imports[*fdef.Name]
			if ok {
				continue //already imported
			}
			noDollar := strings.TrimPrefix(*fdef.Name, "$")
			u.candidateFT[noDollar] = fdef.Type.Num
			u.candidateSym[noDollar] = *fdef.Name
		}
	}
}
func (u *unlink) unlinkStdlib(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.FuncDef)
	for _, pkg := range stdlibName {
		if idef.Name != nil && strings.HasPrefix(*idef.Name, pkg) {
			u.removed++
			return nil
		}
	}
	// xxx what if no name only number?
	_, ok := u.remaining[*idef.Name]
	if !ok {
		u.remaining[*idef.Name] = struct{}{}
	}
	return idef
}

func (u *unlink) dumpStats() {
	log.Printf("moved %d functions from defined to imported", len(u.packages))
	for pkg, count := range u.packages {
		log.Printf("\t %20s: %d", pkg, count)
	}
	log.Printf("the following packages are still remaining in the binary:")
	for k := range u.remaining {
		if strings.HasPrefix(k, "$interface:") {
			log.Printf("\t%70s", "[interface type ignored]")
			continue
		}
		if strings.HasPrefix(k, "$_*github.com/iansmith/parigot") {
			continue
		}
		if len(k) > lenLimit {
			k = k[:lenLimit] + fmt.Sprintf("...%d characters elided", len(k)-lenLimit)
		}
		log.Printf("\t%70s", k)
	}
}

func sizeInBytes(n int64) string {
	if n < oneK {
		return fmt.Sprintf("%d bytes", n)
	}
	if n < oneM {
		f := float32(n) / float32(oneK)
		return fmt.Sprintf("%03.1fK bytes", f)
	}
	f := float32(n) / float32(oneM)
	return fmt.Sprintf("%03.1fM bytes", f)

}
