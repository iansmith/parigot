package tree

import (
	"bytes"
	"fmt"
	"strings"
)

type DocSectionNode struct {
	DocFunc      []*DocFuncNode
	AnonymousNum int
	Program      *ProgramNode
}

func (s *DocSectionNode) FinalizeSemantics() {
	if s == nil {
		return
	}
	s.SetNumber()
	// if len(s.DocFunc) > 0 {
	// 	s.Program.NeedBytes = true
	// }
}

func (s *DocSectionNode) SetNumber() {
	for _, fn := range s.DocFunc {
		fn.SetNumber()
	}
}

func NewDocSectionNode(p *ProgramNode, fn []*DocFuncNode) *DocSectionNode {
	return &DocSectionNode{Program: p, DocFunc: fn}
}

type DocFuncNode struct {
	Name              string
	Elem              *DocElement
	Param, Local      []*PFormal
	PreCode, PostCode []TextItem
	Section           *DocSectionNode
}

func (f *DocFuncNode) SetNumber() {
	f.Elem.SetNumber(0)
}

// func (f *DocFuncNode) CheckForBadVariableUse() string {
// 	for _, seq := range [][]TextItem{f.PreCode, f.PostCode} {
// 		for _, item := range seq {
// 			switch varName := item.(type) {
// 			case *TextValueRef:
// 				msg := f.checkAllForNameDecl(varName.String())
// 				if msg != "" {
// 					return msg
// 				}
// 			}
// 		}
// 	}
// 	return ""
// }

// func (f *DocFuncNode) checkVar(name string, formal []*PFormal) bool {
// 	for _, p := range formal {
// 		if p.Name == name {
// 			return true
// 		}
// 	}
// 	return false
// }

//	func (f *DocFuncNode) checkGlobalAndExtern(name string) bool {
//		return f.Section.Program.checkGlobalAndExtern(name)
//	}
// func (f *DocFuncNode) checkAllForNameDecl(name string) string {
// 	if IsSelfVar(name) {
// 		return ""
// 	}
// 	found := f.checkLocal(name)
// 	if found {
// 		return ""
// 	}
// 	found = f.checkParam(name)
// 	if found {
// 		return ""
// 	}
// 	found = f.checkGlobalAndExtern(name)
// 	if found {
// 		return ""
// 	}

// 	return fmt.Sprintf("in doc function '%s', unknown variable '%s'",
// 		f.Name, name)
// }

// func (f *DocFuncNode) checkLocal(name string) bool {
// 	return f.checkVar(name, f.Local)
// }

//	func (f *DocFuncNode) checkParam(name string) bool {
//		return f.checkVar(name, f.Param)
//	}

func NewDocFuncNode(n string, formal []*PFormal, local []*PFormal, s *DocElement, pre, post []TextItem) *DocFuncNode {
	return &DocFuncNode{Name: n, Param: formal, Local: local, Elem: s, PreCode: pre, PostCode: post}
}

type DocElement struct {
	Number      int
	Tag         *DocTag
	TextContent []TextItem
	Child       []*DocElement
}

func NewDocElementWithText(tag *DocTag, content []TextItem) *DocElement {
	return &DocElement{
		Tag: tag, TextContent: content,
	}
}

func NewDocElementWithChild(child []*DocElement) *DocElement {
	return &DocElement{Child: child}
}

func (e *DocElement) SetNumber(n int) int {
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

func NewDocTag(tag *ValueRef, id *ValueRef, clazz []*ValueRef) (*DocTag, error) {
	if tag.Lit != "" {
		if !validTag(tag.Lit) {
			return nil, fmt.Errorf("unknown tag '%s'", tag.Lit)
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

func validTag(tag string) bool {
	switch tag {
	case
		"article", "aside", "details", "figcaption", "figure", "footer", "header", "legend", "main",
		"mark", "nav", "section", "summary", "time",
		"abbr", "address", "base", "blockquote", "body", "col", "head", "hr", "link", "meta", "noscript",
		"object", "param", "progress", "q", "sub", "sup", "track", "var", "video", "wbr",

		"h1", "h2", "h3", "h4", "h5", "title", "br",
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
