package dom

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/iansmith/parigot/g/dom/v1"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

type DOMError error

var DOMNotFound = errors.New("element not found in document")
var DOMInternalError = errors.New("internal error")
var DOMAlreadyPresent = errors.New("already present")

const ParigotIdAttribute = "parigot-id"

// DOMServer makes the browser work like other servers.   Currently a DOMServer is bound to a Document in the HTML sense.
type DOMServer struct {
	doc js.Value
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
		strId: make(map[string]lib.Id),
	}, nil
}

// ElementById finds the element that has the id given in the request.  That element is returned or the error is
// set to non-nil, probably DOMNotFound.  Note that the children of the node are not fully filled out because
// that could lead a VERY large amount DOM interaction if we tried to fill out the subtree.
func (d *DOMServer) ElementById(in *dommsg.ElementByIdRequest) (*dommsg.ElementByIdResponse, error) {
	value := d.doc.Call("getElementById", in.Id)
	if !value.Truthy() {
		if value.IsNull() {
			return nil, DOMNotFound
		}
		return nil, DOMInternalError
	}
	elem, err := d.buildElement(value)
	if err != nil {
		return nil, err
	}
	return &dommsg.ElementByIdResponse{Elem: elem}, nil
}

// buildBasics returns an parigot *Element with the data provided by domElem.  This
// does not fill out the children, only the tag portion of the element.
func (d *DOMServer) buildBasics(domElem js.Value) (*dommsg.Element, error) {
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

// buildElement returns the element corresponding to the domElem given.  This
// object has the parigot id and the children partially filled out.
func (d *DOMServer) buildElement(domElem js.Value) (*dommsg.Element, error) {
	elem, err := d.buildBasics(domElem)
	if err != nil {
		return nil, err
	}
	if elem.ParigotId == nil {
		pId := lib.NewElementId()
		d.strId[pId.String()] = pId
		elem.ParigotId = lib.Marshal[protosupportmsg.ElementId](pId)
	}

	text := domElem.Get("innerText")
	if text.String() != "" && !text.Truthy() {
		return nil, DOMInternalError
	}
	elem.Text = text.String()
	nChild := domElem.Get("childElementCount")
	if nChild.Truthy() {
		n := nChild.Int()
		nodeList := nChild.Get("children")
		for i := 0; i < n; i++ {
			node := nodeList.Index(i)
			if !node.Truthy() {
				print("Warning, found a non-truthy node in a child list\n")
				continue
			}
			c, err := d.buildBasics(node)
			if err != nil {
				return nil, err
			}
			elem.Child = append(elem.Child, c)
		}
	}
	return elem, nil
}

// ElementByParigotId mimics the behavior of the ElementById method, but checks for an attribute called parigotid
// to try to find the element in question.  Be aware that if you, wrongly, have two elements with the same ParigotId in the
// DOM you'll only get the first one.
func (d *DOMServer) ElementByParigotId(in *dommsg.ElementByParigotIdRequest) (*dommsg.ElementByParigotIdResponse, error) {
	pid := lib.Unmarshal(in.ParigotId)

	obj := d.findById(pid.StringRaw(), "")
	if !obj.Truthy() {
		if obj.IsNull() {
			return nil, DOMNotFound
		}
		return nil, DOMInternalError
	}
	elem, err := d.buildElement(obj)
	if err != nil {
		return nil, err
	}
	return &dommsg.ElementByParigotIdResponse{Elem: elem}, nil
}

func (d *DOMServer) SetChild(in *dommsg.SetChildRequest) (*dommsg.SetChildResponse, error) {
	var parigotId string
	if in.ParigotId != nil {
		pid := lib.Unmarshal(in.ParigotId)
		parigotId = pid.StringRaw()
	}
	parent := d.findById(parigotId, in.Id)
	if !parent.Truthy() {
		if parent.IsNull() {
			return nil, DOMNotFound
		}
		return nil, DOMInternalError
	}

	replaced := d.removeAllChildren(parent)
	var buf bytes.Buffer
	for i := 0; i < len(in.Child); i++ {
		element := in.Child[i]
		check := d.findById(elementParigotIdToString(element), element.Tag.Id)
		if check.Truthy() {
			return nil, DOMAlreadyPresent
		}
		buf.WriteString(toHtml(element))
	}
	parent.Set("innerHTML", buf.String())
	return &dommsg.SetChildResponse{
		Replacements: int32(replaced),
	}, nil
}
func (d *DOMServer) removeAllChildren(parent js.Value) int {
	dropped := 0
	for parent.Get("childElementCount").Int() > 0 {
		dropped++
		parent.Call("removeChild", parent.Get("firstChild"))
	}
	return dropped
}

// CreateElement will build a tree of elements, as described by the request.  The return values are only
// the root and first level children's parigot ids, now that they have been created.
func (d *DOMServer) CreateElement(in *dommsg.CreateElementRequest) (*dommsg.CreateElementResponse, error) {
	var parent js.Value
	if in.Parent != nil {
		parent = d.findById(elementParigotIdToString(in.Parent), in.Parent.Tag.Id)
		if !parent.Truthy() {
			return nil, DOMNotFound
		}
	}
	resp, root, err := d.createElementWithValue(in, in.Parent != nil)
	if resp == nil {
		panic("failed to create response to create element:" + err.Error())
	}
	if err != nil {
		return nil, err
	}
	if in.Parent != nil {
		parent.Call("appendChild", root)
	}

	return resp, err
}

func (d *DOMServer) findById(parigotId, id string) js.Value {
	trimmed := strings.TrimSpace(id)
	if trimmed != "" {
		result := d.doc.Call("getElementById", id)
		if result.Truthy() {
			return result
		}
	}
	trimmed = strings.TrimSpace(parigotId)
	if trimmed != "" {
		result := d.doc.Call("querySelector", fmt.Sprintf("[%s=\"%s\"]", ParigotIdAttribute, trimmed))
		if result.Truthy() {
			return result
		}
	}
	return js.Null()
}

// createElementWithValue will build a tree of elements, as described by the request.
func (d *DOMServer) createElementWithValue(in *dommsg.CreateElementRequest, createDOMElem bool) (*dommsg.CreateElementResponse, js.Value, error) {
	resp := &dommsg.CreateElementResponse{}
	_, v, err := d.createSingleElement(in.Root, createDOMElem)
	if err != nil {
		return resp, js.Null(), err
	}
	resp.Root = in.Root
	if len(in.Root.Child) > 0 {
		for i := 0; i < len(in.Root.Child); i++ {
			recurseResp, c, err := d.createElementWithValue(&dommsg.CreateElementRequest{Root: in.Root.Child[i], Parent: in.Root}, createDOMElem)
			if err != nil {
				return nil, js.Null(), err
			}
			resp.Root.Child[i] = recurseResp.Root
			if createDOMElem {
				v.Call("appendChild", c)
			}
		}
	}
	return resp, v, nil
}

// createSingleElement creates the element and returns the newly created id as well as the value from the dom.
// This changes the element passed in to include it's new ParigotId.
func (d *DOMServer) createSingleElement(element *dommsg.Element, createDOM bool) (lib.Id, js.Value, error) {
	pid := lib.NewElementId()
	pidStr := pid.String()
	d.strId[pidStr] = pid
	element.ParigotId = lib.Marshal[protosupportmsg.ElementId](pid)
	if element.Tag == nil {
		element.Tag = &dommsg.Tag{Name: "span"}
		//fmt.Printf("created fake SPAN node, it has %d child\n", len(element.Child))
	}
	result := d.doc.Call("createElement", element.Tag.Name)
	if !result.Truthy() {
		//fmt.Printf("createSingleElement: result.Truthy?? %v,%v\n", result.Truthy(), result)
		return nil, js.Null(), DOMInternalError
	}
	if element.Tag.Id != "" {
		result.Set("id", element.Tag.Id)
	}
	if len(element.Tag.CssClass) > 0 {
		all := make([]string, len(element.Tag.CssClass))
		for i := 0; i < len(element.Tag.CssClass); i++ {
			all[i] = element.Tag.CssClass[i]
		}
		s := strings.Join(all, " ")
		result.Set("cssClass", s)
	}
	result.Call("setAttribute", ParigotIdAttribute, pidStr)
	result.Set("textContent", element.Text)
	return pid, result, nil
}

// toHtml converts an element and all its children to html text.
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
	if len(e.Child) > 0 {
		child := &bytes.Buffer{}
		for _, c := range e.GetChild() {
			child.WriteString(toHtml(c))
		}
		inner = child.String()
	}

	result := fmt.Sprintf("%s%s%s", t, inner, end)
	return result
}

func jsToInt(value js.Value) int {
	return value.Int()
}

func elementParigotIdToString(elem *dommsg.Element) string {
	if elem == nil {
		return ""
	}
	pid := lib.Unmarshal(elem.ParigotId)
	return pid.StringRaw()
}
