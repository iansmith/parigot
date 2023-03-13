package builtin

import (
	"fmt"

	"github.com/iansmith/parigot/ui/css"
)

type ParamChecker func(string) (bool, string)

type BuiltinDef struct {
	Name  string
	Param ParamChecker
}

// CheckCSSClass tests that a given string is a css class we have
// seen when scanning the CSS files.
func CheckIsCssClass(s string) (bool, string) {
	_, ok := css.KnownClasses[s]
	if !ok {
		return false, fmt.Sprintf("'%s' is not a known CSS class", s)
	}
	return ok, ""
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
