package driver

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/ui/parser/tree"
)

const golang = "template/go.tmpl"

var allLanguageTemplates = []string{golang}

type generateContext struct {
	templateName string
	program      *tree.ProgramNode
	global       map[string]any
}

func newGenerateContext(languageName string) *generateContext {
	return &generateContext{
		templateName: languageName,
		global:       make(map[string]any),
	}
}

func zerothElem(in []*tree.FuncActual) *tree.ValueRef {
	return in[0].Ref
}

func transformFormalType(formal *tree.PFormal) string {
	s := formal.Type.String()
	if !strings.Contains(s, ":") {
		return s
	}
	if !strings.HasPrefix(s, ":") {
		panic(fmt.Sprintf("unable to understand type %s (no model?)", formal.Type.String()))
	}
	s = s[1:]
	if !strings.Contains(s, ":") {
		panic(fmt.Sprintf("unable to understand type %s (no qualifier?)", formal.Type.String()))
	}
	if formal.Message == nil {
		panic(fmt.Sprintf("we got a formal that is a model:message ('%s'), but no message registered on formal!", formal.Type.String()))
	}
	pkg := formal.Message.Package
	part := strings.Split(pkg, ".")

	if strings.HasPrefix(pkg, "msg") {
		if len(part) < 2 {
			panic(fmt.Sprintf("unable to understand package name: %s", pkg))
		}
		p := part[len(part)-1]
		if strings.HasPrefix(p, "v") {
			p = part[len(part)-2]
			return "msg" + p + "." + formal.Message.Name
		}
	}
	return part[len(part)-1] + "." + formal.Message.Name
}

func runTemplate(ctx *generateContext, out io.Writer) error {
	root, err := loadTemplates()
	if err != nil {
		return err
	}
	t := root.Lookup(ctx.templateName)
	if t == nil {
		return err
	}
	return t.Execute(out, ctx.global)
}

// loadTemplates not only loads the templates proper from the embedded FS but
// also sets up the default functions from codegen.FuncMap.  Generators have
// a chance later to add functions if te want.  The list of files to load
// is from generator.TemplateName() and the list of extra functions is
// generator.FuncMap().
func loadTemplates() (*template.Template, error) {
	// create root template
	root := template.New("root")
	// add functions
	funcMap := template.FuncMap{
		"zerothElem":          zerothElem,
		"transformFormalType": transformFormalType,
	}
	root = root.Funcs(funcMap)

	// these calls are meant to be "chained" so this construction is needed
	// to capture the "new" value of root.
	t := root
	// template loading
	for _, f := range allLanguageTemplates {
		all, err := templateFS.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s from embedded fs:%v", f, err)
		}
		t, err = root.New(f).Parse(string(all))
		if err != nil {
			return nil, fmt.Errorf("unable to parse template %s:%v", f, err)
		}
	}
	return t, nil
}
