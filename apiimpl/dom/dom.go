package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
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

var allDomServer = make(map[float64]dom.DOMService)

// DOMServer makes the browser work like other servers.   Currently a DOMServer is bound to a Document in the HTML sense.
type DOMServer struct {
	doc js.Value
	// go from a string rep of a parigot id to the actual id
	strId     map[string]lib.Id
	serverId  float64
	startList []func()
}

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	log.Fatalf("failed because dom is still under development")
}

func LocateDOMServer() (dom.DOMService, error) {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return nil, DOMInternalError
	}
	return &DOMServer{
		doc:       doc,
		strId:     make(map[string]lib.Id),
		startList: []func(){},
		serverId:  rand.Float64(),
	}, nil
}

func (d *DOMServer) elemById(id string) (js.Value, error) {
	value := d.doc.Call("getElementById", id)
	if !value.Truthy() {
		if value.IsNull() {
			return js.Null(), DOMNotFound
		}
		return js.Undefined(), DOMInternalError
	}
	return value, nil
}

func (d *DOMServer) elemByEitherId(parigotId, id string) (js.Value, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed != "" {
		result := d.doc.Call("getElementById", id)
		if !result.IsNull() && !result.Truthy() {
			return js.Undefined(), DOMInternalError
		}
		return result, nil
	}

	trimmed = strings.TrimSpace(parigotId)
	if trimmed != "" {
		result := d.doc.Call("querySelector", fmt.Sprintf("[%s=\"%s\"]", ParigotIdAttribute, trimmed))
		if result.Truthy() {
			return result, nil
		}
	}
	return js.Null(), DOMNotFound
}

// ElementById finds the element that has the id given in the request.  That element is returned or the error is
// set to non-nil, probably DOMNotFound.  Note that the children of the node are not fully filled out because
// that could lead a VERY large amount DOM interaction if we tried to fill out the subtree.
func (d *DOMServer) ElementById(in *dommsg.ElementByIdRequest) (*dommsg.ElementByIdResponse, error) {
	value, err := d.elemById(in.Id)
	if err != nil {
		return nil, err // dom puked
	}
	elem, err := d.buildElement(value)
	if err != nil {
		return nil, err
	}
	if elem.Tag.Id != in.Id {
		panic("Wrong Id found")
	}
	return &dommsg.ElementByIdResponse{Elem: elem}, nil
}

func (d *DOMServer) ServerId() float64 {
	sid := d.serverId
	if _, ok := allDomServer[sid]; !ok {
		allDomServer[sid] = d
	}
	return sid
}

// buildBasics returns an parigot *Element with the data provided by domElem.  This
// does not fill out the children, only the tag portion of the element.
func (d *DOMServer) buildBasics(domElem js.Value) (*dommsg.Element, error) {
	t := domElem.Get("tagName")
	if !t.Truthy() {
		return nil, DOMInternalError
	}
	i := domElem.Get("id")

	c := domElem.Get("className")
	if !c.IsNull() && !c.Truthy() {
		return nil, DOMInternalError
	}
	if c.IsNull() {
		log.Printf("no classes found on %s", domElem.Get("id"))
	}
	clazz := []string{}
	if !c.IsNull() {
		clazz = strings.Split(c.String(), " ")
	}

	tagName := t.String()
	id := ""
	if i.Truthy() {
		id = i.String()
	}

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

// ElementByEitherId mimics the behavior of the ElementById method, but checks ADDITIONALLY for an attribute called parigotid
// to try to find the element in question.  Be aware that if you, wrongly, have two elements with the same ParigotId in the
// DOM you'll only get the first one.
func (d *DOMServer) ElementByEitherId(in *dommsg.ElementByEitherIdRequest) (*dommsg.ElementByEitherIdResponse, error) {
	pid := lib.Unmarshal(in.ParigotId)

	obj, err := d.elemByEitherId(pid.StringRaw(), "")
	if err != nil {
		return nil, err
	}
	elem, err := d.buildElement(obj)
	if err != nil {
		return nil, err
	}
	return &dommsg.ElementByEitherIdResponse{Elem: elem}, nil
}

// setChildFirst takes the _root_ of a new tree and tries to place it in the DOM.  If returns
// if it was successful.
func (d *DOMServer) SetChild(in *dommsg.SetChildRequest) (*dommsg.SetChildResponse, error) {
	var parigotId string
	if in.ParigotId != nil {
		pid := lib.Unmarshal(in.ParigotId)
		parigotId = pid.StringRaw()
	}
	parent, err := d.elemByEitherId(parigotId, in.Id)
	if err != nil {
		return nil, err
	}

	replaced := d.removeAllChildren(parent)

	var buf bytes.Buffer
	for i := 0; i < len(in.Child); i++ {
		element := in.Child[i]
		id := lib.Unmarshal(element.ParigotId)
		if id == nil {
			panic(fmt.Sprintf("no parigot id for %+v (%d children, %s)", element.ParigotId, len(in.Child), element.Tag.Name))
		}
		_, err := d.elemByEitherId(id.String(), element.Tag.Id)
		// dom not found is what we expect since it should not be there
		if err != nil && err != DOMNotFound {
			// the "normal" bad case here is DOMAlreadyExists
			return nil, err
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

// CreateElement will build a tree of elements, as described by the request.
// This function returns the tree root as part of the response.
func (d *DOMServer) CreateElement(in *dommsg.CreateElementRequest) (*dommsg.CreateElementResponse, error) {
	var parent js.Value
	if in.Parent != nil {
		var err error
		parent, err = d.elemByEitherId(elementParigotIdToString(in.Parent), in.Parent.Tag.Id)
		if err != nil {
			return nil, err
		}
	}
	resp, root, err := d.createElementWithValue(in, in.Parent != nil)
	if resp == nil {
		panic("failed to create response to create element:" + err.Error())
	}
	if err != nil {
		return nil, err
	}
	if resp.Root.ParigotId == nil {
		panic("bad id (not found) in create")
	}
	if in.Parent != nil {
		parent.Call("appendChild", root)
	}

	return resp, err
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
	}
	result := d.doc.Call("createElement", element.Tag.Name)
	if !result.Truthy() {
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
		result.Call("setAttribute", "class", s)
	}
	result.Call("setAttribute", ParigotIdAttribute, pidStr)
	result.Set("textContent", element.Text)
	return pid, result, nil
}

func (d *DOMServer) AddStartHandler(fn func()) {
	d.startList = append(d.startList, fn)
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
		t = fmt.Sprintf("<%s id=\"%s\" class=\"%s\">", tag.GetName(), tag.GetId(), allClass.String())
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

func (d *DOMServer) AddEvent(selectorString string, eventName string, fn func(this js.Value, arg []js.Value) any) {
	if selectorString[0:1] == "#" {
		elem := d.doc.Call("querySelector", selectorString)
		if !elem.Truthy() {
			fmt.Printf("WARNING: selector '%s' did not match any elements", selectorString)
		} else {
			elem.Call("addEventListener", eventName, js.FuncOf(fn), false)
		}
	}
	if selectorString[0:1] == "." {
		elem := d.doc.Call("querySelector", selectorString)
		if !elem.Truthy() {
			fmt.Printf("WARNING: selector '%s' did not match any elements", selectorString)
		} else {
			elem.Call("addEventListener", eventName, js.FuncOf(fn), false)
		}
	}
	//fmt.Printf("xxx selector for adding event %s\n", selectorString[0:1])
}

func FindByServerId(f float64) dom.DOMService {
	if _, ok := allDomServer[f]; !ok {
		panic("failed to find  dom server:" + fmt.Sprint(f))
	}
	return allDomServer[f]
}

// UpdateCSSClass changes the underlying DOM to have these classes
// on the element provided.
func (d *DOMServer) UpdateCssClass(in *dommsg.UpdateCssClassRequest) error {
	p := lib.Unmarshal(in.Elem.ParigotId)
	e, err := d.elemByEitherId(p.StringRaw(), in.Elem.Tag.Id)
	if err != nil {
		return err
	}
	e.Set("className", strings.Join(in.Elem.Tag.CssClass, " "))
	return nil
}

func DumpElementTree(elem *dommsg.Element, indent int) {
	indentSpc := ""
	for i := 0; i < indent; i++ {
		indentSpc += " "
	}
	children := len(elem.Child) > 0
	pid := "((no parigot id))"
	if elem.ParigotId != nil {
		pid = lib.Unmarshal(elem.ParigotId).String()
	}
	if children {
		log.Printf("%s%s (", indentSpc, pid)
		for i := 0; i < len(elem.Child); i++ {
			DumpElementTree(elem.Child[i], indent+2)
		}
		log.Printf("%s)", indentSpc)
	} else {
		log.Printf("%s%s", indentSpc, pid)
	}
}
