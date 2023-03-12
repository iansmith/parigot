package main

import (
	"github.com/iansmith/parigot/apiimpl/dom"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
)

const parentId = "paraLoc"

func main() {

	svc, err := dom.LocateDOMServer()
	if err != nil {
		panic("unable to get DOMService: " + err.Error())
	}
	elem, err := svc.ElementById(&dommsg.ElementByIdRequest{Id: parentId})
	if err != nil {
		panic("unable to get element by id:" + err.Error())
	}

	req := &dommsg.CreateElementRequest{
		Root:   example(),
		Parent: elem.Elem,
	}
	_, err = svc.CreateElement(req)
	if err != nil {
		panic("failed to create element: " + err.Error())
	}

	// setReq := &dommsg.SetChildRequest{
	// 	Id:    parentId,
	// 	Child: []*dommsg.Element{resp.Root},
	// }
	// respSet, err := svc.SetChild(setReq)
	// if err != nil {
	// 	panic("unable to set child:" + err.Error())
	// }
	// fmt.Printf("set success: %+v\n", respSet)
}