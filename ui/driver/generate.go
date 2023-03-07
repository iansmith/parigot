package driver

import (
	"fmt"
	"io"
	"text/template"

	"github.com/iansmith/parigot-ui/parser"
)

const golang = "template/go.tmpl"

var allLanguageTemplates = []string{golang}

type generateContext struct {
	templateName string
	program      *parser.ProgramNode
	global       map[string]any
	scope        *parser.ScopeStack
}

func newGenerateContext(languageName string) *generateContext {
	return &generateContext{
		templateName: languageName,
		global:       make(map[string]any),
		scope:        parser.NewScopeStack(),
	}
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
	funcMap := template.FuncMap{}
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
