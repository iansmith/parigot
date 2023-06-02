package tree

import (
	"fmt"
	"log"
	"strings"
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
		log.Printf("invocation of unknown function '%s' at %s", fn.String(), e.String())
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
			//log.Printf("checking a value ref for %s", ref.String())
			e := &ErrorLoc{
				Filename: filename,
				Line:     ref.FuncInvoc.LineNumber,
				Col:      ref.FuncInvoc.ColumnNumber,
			}
			if !CheckFuncName(ref.FuncInvoc, parent, e) {
				return false
			}
			continue
		}
		if ref.Id.String() == "result" {
			continue
		}
		e := &ErrorLoc{
			Filename: filename,
			Line:     ref.Id.LineNumber,
			Col:      ref.Id.ColumnNumber,
		}

		formal := CheckVarName(fname, ref.Id, local, param, parent, e)
		if formal != nil {
			if strings.Contains(ref.Id.String(), ":") {
				e := &ErrorLoc{
					Filename: filename,
					Line:     ref.Id.LineNumber,
					Col:      ref.Id.ColumnNumber,
				}
				// at this point we have the base name checked and found the formal so we need to walk the message fields
				// and qualifiers
				currIdPart := ref.Id.Part.Qual
				currMsg := formal.Message
				//log.Printf("currIdPart xxx %s and %+v", currIdPart.Id, currMsg.Field)
				first := true
				for currIdPart != nil {
					if !first {
						if !currIdPart.ColonSep {
							log.Printf("cannot use dot separators in a qualifier referring to a protobuf, must use colons at %s: %s", currIdPart.Id, e.String())
							return false
						}
					} else {
						first = false
					}
					if currMsg == nil {
						panic(fmt.Sprintf("no message to search for field %s", currIdPart.Id))
					}
					currentField, ok := currMsg.Field[currIdPart.Id]
					if ok {
						if currentField.Field.Message == nil {
							loc, ok := currMsg.Location[currIdPart.Id]
							if !ok {
								log.Printf("unable to find a protobuf field named '%s' in '%s'", currIdPart.Id, currMsg.Name)
								return false
							}
							currMsg = loc.Message
						} else {
							currMsg = currentField.Field.Message
						}
					}
					currIdPart = currIdPart.Qual
				}
			}
		} else {
			return false
		}

		if formal == nil {
			return false
		}
	}
	return true
}
