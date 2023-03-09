package dom

import (
	"bytes"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/iansmith/parigot/g/dom/v1"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

type DOMError error

var DOMNotFound = errors.New("element not found in document")
var DOMInternalError = errors.New("internal error")

type DOMServer struct {
	doc js.Value
	// go from an id to a javascript object that is that element
	id map[string]js.Value
	// go from an id to a parigot id
	domId map[string]lib.Id
	// go from a string rep of a parigot id to the actual id
	strId map[string]lib.Id
}

func LocateDOMServer() (dom.DOMService, error) {

	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return nil, DOMInternalError
	}
	return &DOMServer{
		doc:   doc,
		id:    make(map[string]js.Value),
		domId: make(map[string]lib.Id),
		strId: make(map[string]lib.Id),
	}, nil
}

func (d *DOMServer) ElementById(in *dommsg.ElementByIdRequest) (*dommsg.ElementByIdResponse, error) {
	value := d.doc.Call("getElementById", in.Id)
	if !value.Truthy() {
		if value.IsNull() {
			return nil, DOMNotFound
		}
	}
	elem, err := d.buildElement(value)
	if err != nil {
		return nil, err
	}
	return &dommsg.ElementByIdResponse{Elem: elem}, nil

}

func (d *DOMServer) buildBbasics(domElem js.Value) (*dommsg.Element, error) {
	t := domElem.Get("tagName")
	if !t.Truthy() {
		return nil, DOMInternalError
	}
	i := domElem.Get("id")
	if !i.Truthy() {
		return nil, DOMInternalError
	}

	c := domElem.Call("getAttribute", "class")
	if !c.IsNull() && !c.Truthy() {
		return nil, DOMInternalError
	}
	clazz := []string{}
	if !c.IsNull() {
		l := c.Get("length")
		if !l.Truthy() {
			return nil, DOMInternalError
		}
		for num := 0; num < l.Int(); num++ {
			s := c.Index(num)
			clazz = append(clazz, s.String())
		}
	}

	tagName := t.String()
	id := i.String()

	tag := &dommsg.Tag{
		Name:     tagName,
		Id:       id,
		CssClass: clazz,
	}
	return &dommsg.Element{
		Tag: tag,
	}, nil
}

func (d *DOMServer) buildElement(domElem js.Value) (*dommsg.Element, error) {
	elem, err := d.buildBbasics(domElem)
	if err != nil {
		return nil, err
	}
	domId := elem.Tag.Id

	pId, ok := d.domId[domId]
	if ok {
		elem.Id = lib.Marshal[protosupportmsg.ElementId](pId)
	} else {
		elem.Id = lib.Marshal[protosupportmsg.ElementId](lib.NewElementId())
	}

	text := domElem.Get("innerText")
	if text.String() != "" && !text.Truthy() {
		return nil, DOMInternalError
	}
	elem.Text = text.String()
}

func (d *DOMServer) ElementByParigotId(in *dommsg.ElementByIdRequest) (*dommsg.ElementByIdResponse, error) {
	return nil, nil
}

func (d *DOMServer) SetChild(in *dommsg.SetChildRequest) (*dommsg.SetChildResponse, error) {
	return nil, nil
}

func toHtml(e *dommsg.Element) string {
	t := ""
	end := ""
	tag := e.GetTag()
	if tag != nil {
		allClass := &bytes.Buffer{}
		for _, clazz := range tag.GetCssClass() {
			allClass.WriteString(clazz + " ")
		}
		t = fmt.Sprintf("<%s id=\"%s\" class=\"%s\">", tag.GetName(), tag.GetId(), allClass)
		end = fmt.Sprintf("</%s>", tag.GetName())
	}
	inner := e.GetText()
	if inner == "" {
		child := &bytes.Buffer{}
		for _, c := range e.GetChild() {
			child.WriteString(toHtml(c))
		}
		inner = child.String()
	}

	result := fmt.Sprintf("%s%s%s", t, inner, end)
	print("result of toHtml ", result, "\n")
	return result
}

func jsToInt(value js.Value) int {
	return value.Int()
}
