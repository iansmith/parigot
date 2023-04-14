package tree

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

const anonPrefix = "_anon"

var anonCount int

type DocSectionNode struct {
	DocFunc      []*DocFuncNode
	Scope_       *SectionScope
	AnonymousNum int
	Program      *ProgramNode
}

func (s *DocSectionNode) FinalizeSemantics(path string) error {
	if s == nil {
		return nil // no doc section, no errors
	}
	panic("FinalizeSemantics")
	for _, fn := range s.DocFunc {
		fn.Section = s
		fn.Elem.attachAnonTextFunc(s.Program.TextSection)
	}
	s.Scope_.DocFn = s.DocFunc
	return nil
}

func (s *DocSectionNode) AttachViewToSection(view *ViewDecl) {
	if view == nil {
		panic("nil view")
	}
	if s == nil {
		s = NewDocSectionNode(view.Section.Program, nil)
		if view.Section.Program.DocSection != nil {
			panic("trying to replace existing doc section")
		}
		view.Section.Program.DocSection = s
		s.Scope_.Brother = GProgram.TextSection.Scope_
	}
	view.DocFn.Section = s
	s.DocFunc = append(s.DocFunc, view.DocFn)
}

func (s *DocSectionNode) VarCheck(filename string) bool {
	for _, fn := range s.DocFunc {
		if !fn.VarCheck(filename) {
			return false
		}
		if !fn.CheckDup(filename) {
			return false
		}
		seen := make(map[string]*ErrorLoc)
		for _, fn := range s.DocFunc {
			e := &ErrorLoc{filename, fn.LineNumber, fn.ColumnNumber}
			if _, ok := seen[fn.Name]; ok {
				log.Printf("two instances of doc func '%s' found at %s and %s", fn.Name, seen[fn.Name].String(), e.String())
				return false
			}
			seen[fn.Name] = e
			if s.Program.TextSection != nil && s.Program.TextSection.Func != nil {
				for _, other := range s.Program.TextSection.Func {
					if fn.Name == other.Name {
						eDoc := &ErrorLoc{filename, fn.LineNumber, fn.ColumnNumber}
						eText := &ErrorLoc{Filename: filename, Line: other.LineNumber, Col: other.ColumnNumber}
						log.Printf("two functions with the same name '%s' found %s and %s", fn.Name, eDoc.String(), eText.String())
						return false
					}
				}
			}
		}

	}
	return true
}

func (s *DocSectionNode) SetNumber() {
	for _, fn := range s.DocFunc {
		fn.SetNumber()
	}
}

func NewDocSectionNode(p *ProgramNode, fn []*DocFuncNode) *DocSectionNode {
	return &DocSectionNode{Program: p, DocFunc: fn, Scope_: NewSectionScope(p.Global)}
}

type DocFuncNode struct {
	Name                     string
	Elem                     *DocElement
	Param, Local             []*PFormal
	PreCode, PostCode        []TextItem
	Section                  *DocSectionNode
	LineNumber, ColumnNumber int
}

func (f *DocFuncNode) SetNumber() {
	if f == nil {
		return
	}
	f.Elem.SetNumber(0)
}

func (f *DocElement) attachAnonTextFunc(sect *TextSectionNode) {
	if f != nil && len(f.TextContent) > 0 {
		tf := NewTextFuncNode()
		tf.Section = sect
		n := fmt.Sprintf("%s_%04d", anonPrefix, anonCount)
		id := NewIdent(n, n, 0, 0)
		tf.Name = id.String()
		anonCount++
		tf.Item_ = f.TextContent
		sect.Func = append(sect.Func, tf)
		sect.Scope_.TextFn = append(sect.Scope_.TextFn, tf)
		first := f.TextContent[0]
		fi := NewFuncInvoc(id, nil, first.GetLine(), first.GetCol())
		ref := NewValueRef(nil, fi, "", first.GetLine(), first.GetCol())
		f.TextContent = []TextItem{NewTextValueRef(ref, first.GetLine(), first.GetCol())}
	}
	if f != nil && f.Child != nil {
		for _, c := range f.Child {
			c.attachAnonTextFunc(sect)
		}
	}
}

func checkDupParamAndLocal(p, l []*PFormal, filename, funcName string, isDoc bool) bool {
	fType := "doc"
	if !isDoc {
		fType = "text"
	}
	seen := make(map[string]struct{})
	for _, param := range p {
		if _, ok := seen[param.Name]; ok {
			e := ErrorLoc{Filename: filename, Line: param.LineNumber, Col: param.ColumnNumber}
			log.Printf("%s function '%s' has duplicate paramater '%s' at %s",
				fType, funcName, param.Name, e.String())
			return false
		}
		if param.Name == funcName {
			e := ErrorLoc{Filename: filename, Line: param.LineNumber, Col: param.ColumnNumber}
			log.Printf("%s function '%s' has a parameter of the same name at %s",
				fType, funcName, e.String())
			return false

		}
		seen[param.Name] = struct{}{}
	}
	seen = make(map[string]struct{})
	for _, local := range l {
		if local.Name == funcName {
			e := ErrorLoc{Filename: filename, Line: local.LineNumber, Col: local.ColumnNumber}
			log.Printf("%s function '%s' has a local of the same name at %s",
				fType, funcName, e.String())
			return false

		}
		if _, ok := seen[local.Name]; ok {
			e := ErrorLoc{Filename: filename, Line: local.LineNumber, Col: local.ColumnNumber}
			log.Printf("%s function '%s' has two local variables named '%s' at %s",
				fType, funcName, local.Name, e.String())
			return false
		}
		seen[local.Name] = struct{}{}

	}
	return true
}

func checkParamShadown(p []*PFormal, filename, funcName string, scope Scope, docFunc bool) bool {
	fType := "doc"
	if !docFunc {
		fType = "text"
	}
	for _, param := range p {
		id := NewIdent(param.Name, param.Name, 0, 0)
		if scope.LookupVar(id) != nil {
			e := ErrorLoc{Filename: filename, Line: param.LineNumber, Col: param.ColumnNumber}
			log.Printf("in %s function '%s', parameter '%s' at %s shadows outer definition", fType, funcName, param.Name, e.String())
			return false
		}
		invoc := NewFuncInvoc(id, nil, 0, 0)
		if scope.LookupFunc(invoc) {
			e := ErrorLoc{Filename: filename, Line: param.LineNumber, Col: param.ColumnNumber}
			log.Printf("in %s function '%s', paramater '%s' at %s shadows outer function", fType, funcName, param.Name, e.String())
			return false
		}
	}
	return true
}

func checkLocalShadow(l, p []*PFormal, filename, funcName string, scope Scope, docFunc bool) bool {
	fType := "doc"
	if !docFunc {
		fType = "text"
	}
	for _, param := range p {
		for _, local := range l {
			if param.Name == local.Name {
				e := ErrorLoc{Filename: filename, Line: param.LineNumber, Col: param.ColumnNumber}
				log.Printf("in %s function '%s', parameter '%s' at %s shadows local declaration of '%s'", fType, funcName, param.Name, e.String(), param.Name)
				return false
			}
		}
	}
	return true
}

func (f *DocFuncNode) CheckDup(filename string) bool {
	if !checkDupParamAndLocal(f.Param, f.Local, f.Name, f.Name, true) {
		return false
	}
	if !checkParamShadown(f.Param, filename, f.Name, f.Section.Scope_, true) {
		return false
	}
	if !checkLocalShadow(f.Local, f.Param, filename, f.Name, f.Section.Scope_, true) {
		return false
	}
	return true
}

func (f *DocFuncNode) VarCheck(filename string) bool {
	if f.Section == nil {
		panic(fmt.Sprintf("no section present on doc func node: %s", f.Name))
	}
	if !CheckAllItems(f.Name, f.PreCode, f.Local, f.Param, f.Section.Scope_, filename) {
		return false
	}
	if !CheckAllItems(f.Name, f.PostCode, f.Local, f.Param, f.Section.Scope_, filename) {
		return false
	}
	if f.Elem != nil && f.Elem.Child != nil {
		for _, fn := range f.Elem.Child {
			if !CheckAllItems(f.Name, fn.TextContent, f.Local, f.Param, f.Section.Scope_, filename) {
				return false
			}
		}
	}
	return true
}

func NewDocFuncNode(n string, formal []*PFormal, local []*PFormal, s *DocElement, pre, post []TextItem) *DocFuncNode {
	return &DocFuncNode{Name: n, Param: formal, Local: local, Elem: s, PreCode: pre, PostCode: post}
}

type DocElement struct {
	Number      int
	Tag         *DocTag
	TextContent []TextItem
	Child       []*DocElement
	Parent      *DocElement
}

/*
func NewDocElementWithText(tag *DocTag, content []TextItem) *DocElement {
	return &DocElement{
		Tag: tag, TextContent: content,
	}
}

func NewDocElementWithChild(child []*DocElement) *DocElement {
	return &DocElement{Child: child}
}
*/

func (e *DocElement) SetNumber(n int) int {
	if e == nil {
		return n
	}
	if e.TextContent == nil && len(e.Child) == 0 {
		e.Number = n
		return n + 1
	}
	e.Number = n
	n++
	for _, c := range e.Child {
		n = c.SetNumber(n)
	}
	return n
}

type DocTag struct {
	Tag   *ValueRef
	Id    *ValueRef
	Class []*ValueRef
}

func NewDocTag(tag *ValueRef, id *ValueRef, clazz []*ValueRef, existing map[string]struct{}) (*DocTag, error) {
	if tag.Lit != "" {
		if !ValidTag(tag.Lit) {
			return nil, fmt.Errorf("unknown tag '%s'", tag.Lit)
		}
	}
	for _, class := range clazz {
		if class.Lit == "" {
			continue
		}
		_, ok := existing["."+class.Lit]
		if !ok {
			return nil, fmt.Errorf("unknown css class '%s'", class.Lit)
		}
	}
	return &DocTag{Tag: tag, Id: id, Class: clazz}, nil
}

func (d *DocTag) String() string {
	if d.Tag == nil {
		return `&dommsg.Tag{Name:"span"}`
	}
	t := d.Tag.String()
	id := ""
	class := []string{}
	if d.Id != nil {
		id = "#" + d.Id.String()
	}
	if len(d.Class) != 0 {
		for _, c := range d.Class {
			str := c.String()
			str = strings.TrimPrefix(str, "\"")
			str = strings.TrimSuffix(str, "\"")
			class = append(class, "."+str)
		}
	}
	buf := &bytes.Buffer{}
	buf.WriteString("&dommsg.Tag{Name:" + t + ",\n")
	if id != "" {
		buf.WriteString("Id:" + id + ",\n")
	}
	if len(class) != 0 {
		buf.WriteString("Class: []string{\n")
		for _, c := range class {
			buf.WriteString(c + ",\n")
		}
	} else {
		buf.WriteString("Class: []string{},\n")
	}
	buf.WriteString("}")
	result := buf.String()
	result = strings.TrimSpace(result)
	return result
}

func ValidTag(tag string) bool {
	switch tag {
	case
		"article", "aside", "details", "figcaption", "figure", "footer", "header", "legend", "main",
		"mark", "nav", "section", "summary", "time",
		"abbr", "address", "base", "blockquote", "body", "col", "head", "hr", "link", "meta", "noscript",
		"object", "param", "progress", "q", "sub", "sup", "track", "var", "video", "wbr",

		"h1", "h2", "h3", "h4", "h5", "h6", "title", "br",
		"strong", "em",
		"a", "p", "span", "div",
		"form", "input", "fieldset", "label", "keygen", "optgroup", "option", "textarea", "button",
		"ul", "ol", "dl", "dd", "dt", "li",
		"img",
		"code", "kbd", "pre", "samp",
		"script",
		"table", "tbody", "td", "tr":
		return true
	}
	return false
}
