package tree

import (
	"fmt"
	"log"
)

type ErrorLoc struct {
	Filename  string
	Line, Col int
}

func (e *ErrorLoc) String() string {
	return fmt.Sprintf("%s:%d:%d", e.Filename, e.Line, e.Col)
}

func CheckFuncName(fn *FuncInvoc, scope Scope, e *ErrorLoc) bool {
	if !scope.LookupFunc(fn) {
		log.Printf("invocation of unknown function '%s' in %s", fn.String(), e.String())
		return false
	}
	return true
}

func ResolveReferenceFormal(formal *PFormal, ident *Ident) *PFormal {
	if formal.Name == ident.Part.Id {
		return formal
	}
	return nil
}

func CheckLocalAndParam(fname string, id *Ident, local, param []*PFormal, parent Scope) *PFormal {
	var f *PFormal
	found := false
	for _, formal := range local {
		f = ResolveReferenceFormal(formal, id)
		if f != nil {
			found = true
			break
		}
	}
	if found {
		return f
	}
	for _, formal := range param {
		f = ResolveReferenceFormal(formal, id)
		if f != nil {
			found = true
			break
		}
	}
	if found {
		return f
	}
	return parent.LookupVar(id)
}

func CheckVarName(fname string, id *Ident, local, param []*PFormal, parent Scope, e *ErrorLoc) *PFormal {
	copy := *e
	copy.Line = id.LineNumber
	copy.Col = id.ColumnNumber
	result := CheckLocalAndParam(fname, id, local, param, parent)
	if result == nil {
		log.Printf("use of unknown variable '%s' at %s", id.String(), e.String())
	}
	return result
}

func CheckAllItems(fname string, item []TextItem, local, param []*PFormal, parent Scope, filename string) bool {
	for _, i := range item {
		if i.SubTemplate() != ValueRefTemplate {
			continue
		}
		ref := i.(*TextValueRef).Ref
		if ref.Lit != "" {
			continue
		}
		if ref.FuncInvoc != nil {
			e := &ErrorLoc{
				Filename: filename,
				Line:     ref.FuncInvoc.LineNumber,
				Col:      ref.FuncInvoc.ColumnNumber,
			}
			return CheckFuncName(ref.FuncInvoc, parent, e)
		}
		if ref.Id.String() == "result" {
			continue
		}
		e := &ErrorLoc{
			Filename: filename,
			Line:     ref.Id.LineNumber,
			Col:      ref.Id.ColumnNumber,
		}

		return CheckVarName(fname, ref.Id, local, param, parent, e) != nil

	}
	return true
}