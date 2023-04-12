package driver

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

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
	s := formal.TypeName.Type.String()
	if !strings.Contains(s, ":") {
		return formal.TypeName.String()
	}
	if !strings.HasPrefix(s, ":") {
		panic(fmt.Sprintf("unable to understand type %s (no model?)", formal.TypeName.String()))
	}
	if !strings.Contains(s, ":") {
		panic(fmt.Sprintf("unable to understand type %s (no qualifier?)", formal.TypeName.String()))
	}
	if formal.Message == nil {
		panic(fmt.Sprintf("we got a formal that is a model:message ('%s'), but no message registered on formal!", formal.TypeName.String()))
	}
	if formal.Message.LocalGoPkg != "" {
		ts := formal.TypeName.TypeStarter
		switch ts {
		case "":
			ts = "*"
		case "[]":
			ts = "[]*"
		}
		return ts + formal.Message.LocalGoPkg + "." + formal.Message.Name
	}
	pkg := formal.Message.Package
	part := strings.Split(pkg, ".")

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
		"idOutputGo":          idOutputGo,
		"fqProtobufNameGo":    fqProtobufNameGo,
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
func fqProtobufNameGo(msg *tree.ProtobufMessage) string {
	if len(msg.Package) == 0 && msg.LocalGoPkg == "" {
		return "*" + msg.Name
	}
	if msg.LocalGoPkg != "" {
		return "*" + msg.LocalGoPkg + "." + msg.Name
	}

	part := strings.Split(msg.Package, ".")
	if len(part) == 1 {
		return "*" + msg.Package + "." + msg.Name
	}
	candidate := part[len(part)-1]
	if len(candidate) < 2 {
		return "*" + candidate
	}
	allDigit := true
	if candidate[0] == 'v' {
		for i := 1; i < len(candidate); i++ {
			if !unicode.IsDigit(rune(candidate[i])) {
				allDigit = false
				break
			}
		}
	}
	if allDigit {
		return "*" + part[len(part)-2] + "." + msg.Name
	}
	return "*" + part[len(part)-1] + "." + msg.Name
}

func idOutputGo(id string) string {
	hasDot := strings.Contains(id, ".")
	hasColon := strings.Contains(id, ":")

	if !hasDot && !hasColon {
		return id
		// r := rune(id[0])
		// return string(unicode.ToUpper(r)) + id[1:]
	}
	var part []string
	if hasDot {
		part = strings.Split(id, ".")
	} else {
		part = strings.Split(id, ":")
	}
	result := make([]string, len(part))
	for i, segment := range part {
		if i == 0 {
			result[0] = part[0]
			continue
		}
		var b = ""
		for _, u := range strings.Split(segment, "_") {
			cap := string(unicode.ToUpper(rune(u[0]))) + u[1:]
			b += cap
		}
		result[i] = b
	}
	return strings.Join(result, ".")
}
