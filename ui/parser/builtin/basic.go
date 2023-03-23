package builtin

import (
	"fmt"

	"github.com/iansmith/parigot/ui/css"
	"github.com/iansmith/parigot/ui/parser/tree"
)

type ParamChecker func(*tree.FuncInvoc) (bool, string)

type BuiltinDef struct {
	Name  string
	Param ParamChecker
}

// This is the type checker for any builtin func that has one parameter that
// is a CSS class.  We can only verify the CSS class if the literal is in the code.
func CheckIsCssClass(f *tree.FuncInvoc) (bool, string) {
	if len(f.Actual) != 1 {
		return false, fmt.Sprintf("wrong number of arguments to %s, expected 1 but got %d", f.Name, len(f.Actual))
	}
	arg0 := f.Actual[0].Ref.Lit
	if arg0 != "" {
		_, ok := css.KnownClasses[arg0]
		if !ok {
			return false, fmt.Sprintf("in call to '%s', '%s' is not a known CSS class", f.Name, arg0)
		}
	}
	return true, ""
}

var all map[string]*BuiltinDef = map[string]*BuiltinDef{
	"ToggleSingle": {Name: "ToggleSingle", Param: CheckIsCssClass},
}

func GetBuiltinChecker(name string) (ParamChecker, error) {
	bdef, ok := all[name]
	if !ok {
		return nil, fmt.Errorf("builtin '%s' not found", name)
	}
	return bdef.Param, nil

}
